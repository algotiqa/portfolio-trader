//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package importexport

import (
	"encoding/json"

	"github.com/algotiqa/portfolio-trader/pkg/db"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

//=============================================================================

func BuildTradingSystems(systems *[]db.TradingSystem, filters *[]db.TradingFilter, trades *[]db.Trade,
						 periods *[]db.LivePeriod, positions *[]db.TradingPosition) []*TradingSystem {

	tsMap := map[uint]*TradingSystem{}

	for _,ts := range *systems {
		tsMap[ts.Id] = NewTradingSystem(&ts)
	}

	for _, f := range *filters {
		ts,ok := tsMap[f.TradingSystemId]
		if ok {
			ts.TradingFilter = &f
		}
	}

	for _, p := range *positions {
		ts,ok := tsMap[p.TradingSystemId]
		if ok {
			ts.TradingPosition = &p
		}
	}

	for _, tr := range *trades {
		ts,ok := tsMap[tr.TradingSystemId]
		if ok {
			ts.Trades = append(ts.Trades, NewTrade(&tr))
		}
	}

	for _, lp := range *periods {
		ts,ok := tsMap[lp.TradingSystemId]
		if ok {
			ts.LivePeriods = append(ts.LivePeriods, NewLivePeriod(&lp))
		}
	}

	return maps.Values(tsMap)
}

//=============================================================================

func EncodeTradingSystems(list []*TradingSystem) (*ExportedData, error) {
	ed := &ExportedData{}

	for _, ts := range list {
		data, err := json.MarshalIndent(ts, "", "\t")
		if err != nil {
			return nil, err
		}

		es := &EncodedSystem{
			Id      : ts.Id,
			JsonData: data,
		}

		ed.TradingSystems = append(ed.TradingSystems, es)
	}

	return ed, nil
}

//=============================================================================

func ImportTradingSystem(tx *gorm.DB, ts *db.TradingSystem, data []byte) error {
	its := TradingSystem{}
	err := json.Unmarshal(data, &its)
	if err == nil {
		err = updateTradingSystem(tx, ts, &its)
		if err == nil {
			err = setTradingFilter(tx, ts.Id, its.TradingFilter)
			if err == nil {
				err = setTradingPosition(tx, ts.Id, its.TradingPosition)
				if err == nil {
					err = addTrades(tx, ts.Id, its.Trades)
					if err == nil {
						err = addLivePeriods(tx, ts.Id, its.LivePeriods)
					}
				}
			}
		}
	}

	return err
}

//=============================================================================

func updateTradingSystem(tx *gorm.DB, ts *db.TradingSystem, its *TradingSystem) error {
	ts.Trading         = its.Trading
	ts.AutoActivation  = its.AutoActivation
	ts.Active          = its.Active
	ts.FirstTrade      = its.FirstTrade
	ts.LastTrade       = its.LastTrade
	ts.LastNetProfit   = its.LastNetProfit
	ts.LastNetAvgTrade = its.LastNetAvgTrade
	ts.LastNumTrades   = its.LastNumTrades

	//--- These cannot be set (an imported trading system must be off)
	ts.Running         = false
	ts.Status          = db.TsStatusOff

	return db.UpdateTradingSystem(tx, ts)
}

//=============================================================================

func setTradingFilter(tx *gorm.DB, id uint, f *db.TradingFilter) error {
	f.TradingSystemId = id
	return db.SetTradingFilter(tx, f)
}

//=============================================================================

func setTradingPosition(tx *gorm.DB, id uint, p *db.TradingPosition) error {
	p.TradingSystemId = id
	return db.SetTradingPosition(tx, p)
}

//=============================================================================

func addTrades(tx *gorm.DB, id uint, list []*Trade) error {
	for _, t := range list {
		err := db.AddTrade(tx, convertTrade(id,t))
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func convertTrade(id uint, t *Trade) *db.Trade {
	return &db.Trade{
		TradingSystemId   : id,
		TradeType         : t.TradeType,
		EntryDate         : t.EntryDate,
		EntryPrice        : t.EntryPrice,
		EntryLabel        : t.EntryLabel,
		ExitDate          : t.ExitDate,
		ExitPrice         : t.ExitPrice,
		ExitLabel         : t.ExitLabel,
		GrossReturn       : t.GrossReturn,
		MaxContracts      : t.MaxContracts,
		EntryDateAtBroker : t.EntryDateAtBroker,
		EntryPriceAtBroker: t.EntryPriceAtBroker,
		ExitDateAtBroker  : t.ExitDateAtBroker,
		ExitPriceAtBroker : t.ExitPriceAtBroker,
	}
}

//=============================================================================

func addLivePeriods(tx *gorm.DB, id uint, list []*LivePeriod) error {
	for _, lp := range list {
		err := db.AddLivePeriod(tx, convertLivePeriod(id,lp))
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func convertLivePeriod(id uint, l *LivePeriod) *db.LivePeriod {
	return &db.LivePeriod{
		TradingSystemId: id,
		Period         : l.Period,
		Active         : l.Active,
	}
}

//=============================================================================
