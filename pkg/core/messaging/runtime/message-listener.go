//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package runtime

import (
	"encoding/json"
	"log/slog"
	"sort"
	"time"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/portfolio-trader/pkg/consts"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/portfolio-trader/pkg/platform"
	"gorm.io/gorm"
)

//=============================================================================

func HandleMessage(m *msg.Message) bool {

	slog.Info("New message received", "source", m.Source, "type", m.Type)

	if m.Source == msg.SourceTrade {
		tm := TradeListMessage{}
		err := json.Unmarshal(m.Entity, &tm)
		if err != nil {
			slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
			return true
		}

		if m.Type == msg.TypeCreate {
			return handleNewTrades(&tm)
		}
	}

	slog.Error("Dropping message with unknown source/type!", "source", m.Source, "type", m.Type)
	return true
}

//=============================================================================

func handleNewTrades(tm *TradeListMessage) bool {
	tsId := tm.TradingSystemId

	slog.Info("handleNewTrades: Processing new trades for trading systems", "id", tsId)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		ts, err := db.GetTradingSystemById(tx, tsId)
		if err != nil {
			slog.Error("handleNewTrades: Cannot retrieve trading system", "id", tsId, "error", err.Error())
			return err
		}

		if ts == nil {
			slog.Error("handleNewTrades: Trading system not found. Discarding trades", "id", tsId)
			return nil
		}

		if !ts.Trading {
			slog.Warn("handleNewTrades: Trading system is not in TRADING mode. Discarding trades", "id", tsId)
			return nil
		}

		var trades = &[]db.Trade{}

		if tm.Reload {
			err = deleteTrades(tx, ts)
		} else {
			trades, err = db.FindTradesByTradingSystemId(tx, tsId)
		}

		if err == nil {
			trades, err = addNewTrades(tx, ts, trades, tm.Trades)
			if err == nil {
				err = updateTradingSystem(tx, ts)
			}
		}

		return err
	})

	if err == nil {
		slog.Info("handleNewTrades: Ending processing of new trades", "id", tsId)
	} else {
		slog.Error("handleNewTrades: Ending process with error", "id", tsId, "error", err)
	}
	return err == nil
}

//=============================================================================

func deleteTrades(tx *gorm.DB, ts *db.TradingSystem) error {
	slog.Info("deleteTrades: Deleting all trades on trading system", "id", ts.Id, "name", ts.Name)

	err := db.DeleteAllTradesByTradingSystemId(tx, ts.Id)
	if err != nil {
		return err
	}

	//--- When this will get heavy, we can move this task into a queue

	err = db.DeleteAllEquityBarsByTradingSystemId(tx, ts.Id)
	if err != nil {
		return err
	}

	err = db.DeleteAllLivePeriodsByTradingSystemId(tx, ts.Id)
	if err != nil {
		return err
	}

	ts.FirstTrade      = nil
	ts.LastTrade       = nil
	ts.LastNetProfit   = 0
	ts.LastNetAvgTrade = 0
	ts.LastNumTrades   = 0

	err = platform.DeleteEquityChart(ts.Username, ts.Id)
	if err != nil {
		return err
	}

	slog.Info("deleteTrades: Operation ended", "id", ts.Id)
	return nil
}

//=============================================================================

func addNewTrades(tx *gorm.DB, ts *db.TradingSystem, trades *[]db.Trade, newTrades []*TradeItem) (*[]db.Trade, error) {
	list := *trades

	//--- Keep a map of the already saved trades. Some tools (like MultiCharts) may create duplicates

	tradeSet := map[string]bool{}
	for _, dbt := range *trades {
		tradeSet[dbt.String()] = true
	}

	for _, tr := range newTrades {
		dbTr := toDbTrade(ts.Id, tr)
		_, exists := tradeSet[dbTr.String()]
		if exists {
			continue
		}

		//--- We need to add trades that are outside of [firstTrade .. lastTrade]
		//--- because we will have duplicates when importing from external strategies
		//--- Example: we have @NQ and we run the strategy on the full period to get lots of data.
		//--- Then, when switching to live, the instrument will switch to something like @NQM25 for roughly
		//--- 180 days. @NQ and @NQM25 are slightly different and there will be a jump in the continuous
		//--- contract causing duplicates between @NQ and @NQM25 during the last 180 days

		if isTheTradeOutOfRange(ts, dbTr) {
			tradeSet[dbTr.String()] = true
			err := addTrade(tx, dbTr, tr)
			if err != nil {
				return nil, err
			}

			list = append(list, *dbTr)

			//--- Update information on trading system
			//--- It is better to use the exit date for first/last trade because a trade could last
			//--- for 7+ days and the IDLE flag is impacted

			if ts.FirstTrade == nil || ts.FirstTrade.After(*tr.ExitDate) {
				ts.FirstTrade = tr.ExitDate
			}

			if ts.LastTrade == nil || ts.LastTrade.Before(*tr.ExitDate) {
				ts.LastTrade = tr.ExitDate
			}
		}
	}

	//--- Sort final list as new trades could be in the past

	sort.Slice(list, func(i, j int) bool {
		return list[i].ExitDate.Before(*list[j].ExitDate)
	})

	return &list, nil
}

//=============================================================================

func isTheTradeOutOfRange(ts *db.TradingSystem, t *db.Trade) bool {
	if ts.FirstTrade == nil || ts.LastTrade == nil {
		return true
	}

	//--- Better to use the exit date (please, see above)
	return t.ExitDate.Before(*ts.FirstTrade) || t.ExitDate.After(*ts.LastTrade)
}

//=============================================================================

func toDbTrade(tsId uint, t *TradeItem) *db.Trade {
	return &db.Trade{
		TradingSystemId: tsId,
		TradeType      : t.TradeType,
		EntryDate      : t.EntryDate,
		EntryPrice     : t.EntryPrice,
		EntryLabel     : t.EntryLabel,
		ExitDate       : t.ExitDate,
		ExitPrice      : t.ExitPrice,
		ExitLabel      : t.ExitLabel,
		GrossReturn    : t.GrossReturn,
		MaxContracts   : t.MaxContracts,
	}
}

//=============================================================================

func addTrade(tx *gorm.DB, tr *db.Trade, mtr *TradeItem) error {
	err := db.AddTrade(tx, tr)
	if err != nil {
		return err
	}

	for _, eb := range mtr.Equity {
		dbEb := toDbEquityBar(tr.Id, eb)
		err = db.AddEquityBar(tx, dbEb)
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func toDbEquityBar(trId int64, e *EquityBarItem) *db.EquityBar {
	return &db.EquityBar{
		TradeId    : trId,
		Date       : e.Date,
		GrossReturn: e.GrossReturn,
		Contracts  : e.Contracts,
	}
}

//=============================================================================

func updateTradingSystem(tx *gorm.DB, ts *db.TradingSystem) error {

	//--- If we got new trades, probably we have to set an idle/broken state to running

	if ts.Status == db.TsStatusIdle || ts.Status == db.TsStatusBroken {
		idleStart := time.Now().Add(-time.Hour * 24 * time.Duration(consts.IdleDays))

		if ts.LastTrade.After(idleStart) {
			ts.Status = db.TsStatusRunning
		}
	}

	return db.UpdateTradingSystem(tx, ts)
}

//=============================================================================

//func updateActivationStatus(ts *db.TradingSystem, trades *[]db.Trade, f *db.TradingFilter) {
//	if !ts.Running {
//		ts.SuggestedAction = db.TsActionNone
//		ts.Status = db.TsStatusOff
//		return
//	}
//
//	//--- The trading system is running (i.e. live)
//
//	activValue := false
//	if f != nil {
//		activValue = filter.CalcActivation(ts, f, *trades)
//	}
//
//	if ts.AutoActivation {
//		handleAutomaticActivation(ts, activValue)
//	} else {
//		handleManualActivation(ts, activValue)
//	}
//}

//=============================================================================

//func handleManualActivation(ts *db.TradingSystem, activValue bool) {
//	if !ts.Active {
//		if !activValue {
//			ts.SuggestedAction = db.TsActionNone
//		} else {
//			ts.SuggestedAction = db.TsActionTurnOn
//		}
//	} else {
//		if !activValue {
//			ts.SuggestedAction = db.TsActionTurnOff
//		} else {
//			ts.SuggestedAction = db.TsActionNone
//		}
//	}
//}

//=============================================================================

//func handleAutomaticActivation(ts *db.TradingSystem, activValue bool) {
//	ts.SuggestedAction = db.TsActionNone
//
//	if !ts.Active {
//		if activValue {
//			ts.Status = db.TsStatusRunning
//			ts.Active = true
//			activate(ts)
//			notifyRuntime(ts)
//		}
//	} else {
//		if !activValue {
//			ts.Status = db.TsStatusPaused
//			ts.Active = false
//			activate(ts)
//			notifyRuntime(ts)
//		}
//	}
//}

//=============================================================================
