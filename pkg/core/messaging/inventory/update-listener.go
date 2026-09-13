//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package inventory

import (
	"encoding/json"
	"log/slog"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/portfolio-trader/pkg/business"
	"github.com/algotiqa/portfolio-trader/pkg/business/importexport"
	"github.com/algotiqa/portfolio-trader/pkg/business/position"
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func HandleMessage(m *msg.Message) bool {

	slog.Info("New message received", "source", m.Source, "type", m.Type)

	if m.Source == msg.SourceTradingSystem {
		return handleTradingSystem(m)
	} else if m.Source == msg.SourceDataProduct {
		return handleDataProduct(m)
	} else if m.Source == msg.SourceBrokerProduct {
		return handleBrokerProduct(m)
	} else if m.Source == msg.SourceAccount {
		return handleAccount(m)
	} else if m.Source == msg.SourcePortfolio {
		return handlePortfolio(m)
	}

	slog.Error("Dropping message with unknown source!", "source", m.Source)
	return true
}

//=============================================================================
//=== Trading systems
//=============================================================================

func handleTradingSystem(m *msg.Message) bool {
	tsm := TradingSystemMessage{}
	err := json.Unmarshal(m.Entity, &tsm)
	if err != nil {
		slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
		return true
	}

	if m.Type == msg.TypeCreate {
		return setTradingSystem(&tsm, true)
	}
	if m.Type == msg.TypeUpdate {
		return setTradingSystem(&tsm, false)
	}
	if m.Type == msg.TypeDelete {
		return deleteTradingSystem(&tsm)
	}

	slog.Error("Dropping message with unknown type!", "type", m.Type)
	return true
}

//=============================================================================

func setTradingSystem(tsm *TradingSystemMessage, create bool) bool {
	slog.Info("setTradingSystem: Trading system change received", "create", create, "id", tsm.TradingSystem.Id)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		isNew := true

		ts, err := db.GetTradingSystemById(tx, tsm.TradingSystem.Id)

		if err != nil {
			return err
		}

		if ts == nil {
			ts = &db.TradingSystem{}
			ts.Running         = false
			ts.AutoActivation  = false
			ts.Status          = db.TsStatusOff
			ts.Active          = false
			ts.SuggestedAction = db.TsActionNone
		} else {
			isNew = false

			if ts.Username != tsm.TradingSystem.Username {
				slog.Error("Trading system '%v' not owned by user '%v'! Dropping message", tsm.TradingSystem.Id, tsm.TradingSystem.Username)
				return nil
			}
		}

		ts.Id                 = tsm.TradingSystem.Id
		ts.Username           = tsm.TradingSystem.Username
		ts.Name               = tsm.TradingSystem.Name
		ts.Timeframe          = tsm.TradingSystem.Timeframe
		ts.DataProductId      = tsm.TradingSystem.DataProductId
		ts.DataSymbol         = tsm.DataProduct.Symbol
		ts.BrokerProductId    = tsm.TradingSystem.BrokerProductId
		ts.BrokerSymbol       = tsm.BrokerProduct.Symbol
		ts.BrokerConnectionId = tsm.BrokerProduct.ConnectionId
		ts.PointValue         = tsm.BrokerProduct.PointValue
		ts.CostPerOperation   = tsm.BrokerProduct.CostPerOperation
		ts.MarginValue        = tsm.BrokerProduct.MarginValue
		ts.Increment          = tsm.BrokerProduct.Increment
		ts.MarketType         = tsm.BrokerProduct.MarketType
		ts.CurrencyId         = tsm.Currency.Id
		ts.CurrencyCode       = tsm.Currency.Code
		ts.CurrencySymbol     = tsm.Currency.Symbol
		ts.TradingSessionId   = tsm.TradingSession.Id
		ts.SessionName        = tsm.TradingSession.Name
		ts.SessionConfig      = tsm.TradingSession.Config
		ts.StrategyType       = tsm.TradingSystem.StrategyType
		ts.Overnight          = tsm.TradingSystem.Overnight
		ts.Tags               = tsm.TradingSystem.Tags
		ts.Finalized          = tsm.TradingSystem.Finalized
		ts.Timezone           = tsm.Exchange.Timezone
		ts.AgentProfileId     = tsm.TradingSystem.AgentProfileId
		ts.ExternalRef        = tsm.TradingSystem.ExternalRef
		ts.InSampleFrom       = tsm.TradingSystem.InSampleFrom
		ts.InSampleTo         = tsm.TradingSystem.InSampleTo
		ts.EngineCode         = tsm.TradingSystem.EngineCode
		ts.Trading            = isNew && tsm.TradingSystem.Finalized

		err = db.UpdateTradingSystem(tx, ts)

		if err == nil && isNew {
			if len(tsm.PortfolioPack) == 0 {
				err = db.SetTradingFilter(tx, &db.TradingFilter{
					TradingSystemId: ts.Id,
				})
				if err == nil {
					var p *db.TradingPosition
					p,err = getDefaultTradingPosition()
					if err == nil {
						p.TradingSystemId = ts.Id
						err = db.SetTradingPosition(tx, p)
					}
				}
			} else {
				err = importexport.ImportTradingSystem(tx, ts, tsm.PortfolioPack)
			}
		}

		return err
	})

	if err != nil {
		slog.Error("Raised error while processing message")
	} else {
		slog.Info("setTradingSystem: Operation complete")
	}

	return err == nil
}

//=============================================================================

func getDefaultTradingPosition() (*db.TradingPosition,error) {
	mod := model.NewFixedUnitModel()

	data,err := json.Marshal(mod.Config())
	if err != nil {
		return nil,err
	}

	ps := &db.TradingPosition{
		InitialCapital : position.DefInitialCapital,
		MaxTolDrawdPerc: position.DefMaxTolDrawdPerc,
		MarginOverride : nil,
		MaxUnits       : position.DefMaxUnits,
		RiskPerUnit    : position.DefRiskPerUnit,
		RiskValue      : &position.DefRiskValue,
		Model          : mod.Name(),
		Config         : string(data),
	}

	return ps,nil
}

//=============================================================================

func deleteTradingSystem(tsm *TradingSystemMessage) bool {
	id := tsm.TradingSystem.Id
	slog.Info("deleteTradingSystem: Trading system deletion received", "id", id)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		return business.DeleteTradingSystem(tx, id)
	})

	if err != nil {
		slog.Error("deleteTradingSystem: Raised error while deleting trading system", "error", err.Error())
	} else {
		slog.Info("deleteTradingSystem: Operation complete", "id", id)
	}

	return err == nil
}

//=============================================================================
//=== Data products
//=============================================================================

func handleDataProduct(m *msg.Message) bool {
	dpm := DataProductMessage{}
	err := json.Unmarshal(m.Entity, &dpm)
	if err != nil {
		slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
		return true
	}

	if m.Type == msg.TypeCreate {
		//--- If the broker product is new, there are no trading systems to update. Just return 'true'
		return true
	}
	if m.Type == msg.TypeUpdate {
		return updateDataProduct(&dpm)
	}

	slog.Error("Dropping message with unknown type!", "type", m.Type)
	return true
}

//=============================================================================

func updateDataProduct(dpm *DataProductMessage) bool {
	slog.Info("updateDataProduct: Data product change received", "sourceId", dpm.DataProduct.Id)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		values := map[string]interface{}{
			//--- data_symbol cannot be changed, anyway
			"data_symbol": dpm.DataProduct.Symbol,
		}

		return db.UpdateDataProductInfo(tx, dpm.DataProduct.Id, values)
	})

	if err != nil {
		slog.Error("Raised error while processing message")
	} else {
		slog.Info("updateDataProduct: Operation complete")
	}

	return err == nil
}

//=============================================================================
//=== Broker products
//=============================================================================

func handleBrokerProduct(m *msg.Message) bool {
	bpm := BrokerProductMessage{}
	err := json.Unmarshal(m.Entity, &bpm)
	if err != nil {
		slog.Error("Dropping badly formatted message!", "entity", string(m.Entity))
		return true
	}

	if m.Type == msg.TypeCreate {
		//--- If the broker product is new, there are no trading systems to update. Just return 'true'
		return true
	}
	if m.Type == msg.TypeUpdate {
		return updateBrokerProduct(&bpm)
	}
	if m.Type == msg.TypeDelete { return true }

	slog.Error("Dropping message with unknown type!", "type", m.Type)
	return true
}

//=============================================================================

func updateBrokerProduct(bpm *BrokerProductMessage) bool {
	slog.Info("updateBrokerProduct: Broker product change received", "sourceId", bpm.BrokerProduct.Id)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		values := map[string]interface{}{
			"broker_symbol":      bpm.BrokerProduct.Symbol,
			"point_value":        bpm.BrokerProduct.PointValue,
			"cost_per_operation": bpm.BrokerProduct.CostPerOperation,
			"margin_value":       bpm.BrokerProduct.MarginValue,
			"increment":          bpm.BrokerProduct.Increment,
		}

		return db.UpdateBrokerProductInfo(tx, bpm.BrokerProduct.Id, values)
	})

	if err != nil {
		slog.Error("Raised error while processing message")
	} else {
		slog.Info("updateBrokerProduct: Operation complete")
	}

	return err == nil
}

//=============================================================================
//=== Accounts
//=============================================================================

func handleAccount(m *msg.Message) bool {
	am := AccountMessage{}
	err := json.Unmarshal(m.Entity, &am)
	if err != nil {
		slog.Error("handleAccount: Dropping badly formatted message!", "entity", string(m.Entity))
		return true
	}

	if m.Type == msg.TypeCreate {
		//--- If the broker product is new, there are no trading systems to update. Just return 'true'
		return true
	}
	if m.Type == msg.TypeUpdate {
		return updateAccount(&am)
	}
	if m.Type == msg.TypeDelete { return true }

	slog.Error("handleAccount: Dropping message with unknown type!", "type", m.Type)
	return true
}

//=============================================================================

func updateAccount(am *AccountMessage) bool {
	slog.Info("updateAccount: Account change received", "sourceId", am.Account.Id)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		values := map[string]interface{}{
			"account_code"           : am.Account.Code,
			"account_name"           : am.Account.Name,
			"account_active"         : am.Account.StatusMessage == "",
			"account_current_capital": am.Account.CurrentCapital,
			"account_currency_id"    : am.Currency.Id,
			"account_currency_code"  : am.Currency.Code,
			"account_currency_symbol": am.Currency.Symbol,
		}

		return db.UpdateAccountInfo(tx, am.Account.Id, values)
	})

	if err != nil {
		slog.Error("updateAccount: Raised error while processing message")
	} else {
		slog.Info("updateAccount: Operation complete")
	}

	return err == nil
}

//=============================================================================
//=== Portfolios
//=============================================================================

func handlePortfolio(m *msg.Message) bool {
	pm  := PortfolioMessage{}
	err := json.Unmarshal(m.Entity, &pm)
	if err != nil {
		slog.Error("Dropping badly formatted message for a portfolio!", "entity", string(m.Entity))
		return true
	}

	if m.Type == msg.TypeCreate {
		return setPortfolio(&pm, true)
	}
	if m.Type == msg.TypeUpdate {
		return setPortfolio(&pm, false)
	}
	if m.Type == msg.TypeDelete {
		return deletePortfolio(&pm)
	}

	slog.Error("Dropping portfolio message with unknown type!", "type", m.Type)
	return true
}

//=============================================================================

func setPortfolio(pm *PortfolioMessage, create bool) bool {
	slog.Info("setPortfolio: Portfolio change received", "create", create, "id", pm.Portfolio.Id)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		p, err := db.GetPortfolioById(tx, pm.Portfolio.Id)
		if err != nil {
			return err
		}

		if p == nil {
			p = &db.Portfolio{}
		} else {
			if p.Username != pm.Portfolio.Username {
				slog.Error("Portfolio '%v' not owned by user '%v'! Dropping message", pm.Portfolio.Id, pm.Portfolio.Username)
				return nil
			}
		}

		p.Id                    = pm.Portfolio.Id
		p.Username              = pm.Portfolio.Username
		p.Name                  = pm.Portfolio.Name
		p.AccountPerc           = pm.Portfolio.AccountPerc
		p.MaxMarginPerc         = pm.Portfolio.MaxMarginPerc
		p.AccountId             = pm.Portfolio.AccountId
		p.AccountCode           = pm.Account.Code
		p.AccountName           = pm.Account.Name
		p.AccountCurrentCapital = pm.Account.CurrentCapital
		p.AccountActive         = pm.Account.StatusMessage == ""
		p.AccountCurrencyId     = pm.Account.CurrencyId
		p.AccountCurrencyCode   = pm.Currency.Code
		p.AccountCurrencySymbol = pm.Currency.Symbol
		p.SupportsAccounting    = pm.Account.SupportsAccounting
		p.ConnectionId          = pm.Account.ConnectionId

		return db.UpdatePortfolio(tx, p)
	})

	if err != nil {
		slog.Error("setPortfolio: Raised error while processing message")
	} else {
		slog.Info("setPortfolio: Operation complete")
	}

	return err == nil
}

//=============================================================================

func deletePortfolio(pm *PortfolioMessage) bool {
	id := pm.Portfolio.Id
	slog.Info("deletePortfolio: Portfolio deletion received", "id", id)

	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		return business.DeletePortfolio(tx, id)
	})

	if err != nil {
		slog.Error("deletePortfolio: Raised error while deleting portfolio", "error", err.Error())
	} else {
		slog.Info("deletePortfolio: Operation complete", "id", id)
	}

	return err == nil

}

//=============================================================================
