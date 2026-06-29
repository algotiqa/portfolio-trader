//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/importexport"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingSystems(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int, details bool) (*[]db.TradingSystem, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	return db.GetTradingSystems(tx, filter, offset, limit)
}

//=============================================================================

func GetTradingSystem(tx *gorm.DB, c *auth.Context, id uint) (*db.TradingSystem, error) {
	ts, err := db.GetTradingSystemById(tx, id)
	if err != nil {
		return nil, err
	}

	if ts == nil {
		return nil, req.NewNotFoundError("trading system not found : %v", id)
	}

	if !c.Session.IsAdmin() {
		if ts.Username != c.Session.Username {
			return nil, req.NewForbiddenError("user not allowed : %v", ts.Username)
		}
	}

	return ts, nil
}

//=============================================================================

func DeleteTradingSystem(tx *gorm.DB, id uint) error {
	err := db.DeleteAllTradesByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteAllEquityBarsByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteAllLivePeriodsByTradingSystemId(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteTradingFilter(tx, id)
	if err != nil {
		return err
	}

	err = db.DeleteTradingPosition(tx, id)
	if err != nil {
		return err
	}

	return db.DeleteTradingSystem(tx, id)
}

//=============================================================================

func GetTrades(tx *gorm.DB, c *auth.Context, id uint) (*[]db.Trade, error) {
	_, err := getTradingSystemAndCheckAccess(tx, c, id)
	if err != nil {
		return nil, err
	}

	return db.FindTradesByTradingSystemId(tx, id)
}

//=============================================================================

func ExportTradingSystems(tx *gorm.DB, c *auth.Context, ids []uint) (*importexport.ExportedData, error){
	c.Log.Info("ExportTradingSystems: Exporting trading systems", "count", len(ids))

	systems, err1 := db.GetTradingSystemsById(tx, c.Session.Username, ids)
	if err1 != nil {
		return nil, err1
	}

	filters, err2 := db.GetTradingFiltersByTsIds(tx, ids)
	if err2 != nil {
		return nil, err2
	}

	trades,  err3 := db.FindTradesByTradingSystemsId(tx, ids)
	if err3 != nil {
		return nil, err3
	}

	periods, err5 := db.FindLivePeriodsByTradingSystemsId(tx, ids)
	if err5 != nil {
		return nil, err5
	}

	positions, err6 := db.GetTradingPositionsByTsIds(tx, ids)
	if err6 != nil {
		return nil, err6
	}

	tss := importexport.BuildTradingSystems(systems, filters, trades, periods, positions)

	return importexport.EncodeTradingSystems(tss)
}

//=============================================================================
