//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func GetPortfolios(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]db.Portfolio, error) {
	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	return nil,nil //db.GetPortfolios(tx, filter, offset, limit)
}

//=============================================================================

func DeletePortfolio(tx *gorm.DB, id uint) error {
	return db.DeletePortfolio(tx, id)
}

//=============================================================================

func GetAssignableTradingSystems(tx *gorm.DB, c *auth.Context, id uint) (*[]db.TradingSystemAssignable,error) {
	p,err := getPortfolio(tx, c, id, "GetAssignableSystems")
	if err != nil {
		return nil, err
	}

	return db.GetAssignableTradingSystems(tx, id, p.ConnectionId, !p.SupportsAccounting)
}

//=============================================================================

func GetAssignedTradingSystems(tx *gorm.DB, c *auth.Context, id uint) (*[]db.TradingSystem,error) {
	_,err := getPortfolio(tx, c, id, "GetAssignedTradingSystems")
	if err != nil {
		return nil, err
	}

	filter := map[string]any{}
	filter["portfolio_id"] = id

	return db.GetTradingSystems(tx, filter, 0, 5000)
}

//=============================================================================

func AssignTradingSystemsToPortfolio(tx *gorm.DB, c *auth.Context, id uint, list []uint) error {
	p,err := getPortfolio(tx, c, id, "AssignTradingSystemsToPortfolio")
	if err != nil {
		return err
	}

	assignables,errX := db.GetAssignableTradingSystems(tx, id, p.ConnectionId, !p.SupportsAccounting)
	if errX != nil {
		return errX
	}

	assMap := buildAssignableMap(assignables)
	for _, tsId := range list {
		if _,ok := assMap[tsId]; !ok {
			return req.NewBadRequestError("Trading system with id=%v is not assignable to portfolio", tsId)
		}
	}

	err = db.UpdatePortfolioForTradingSystems(tx, list, &id)

	if err == nil {
		c.Log.Info("AssignTradingSystemsToPortfolio: Assigned portfolio to trading systems", "portfolioId", id, "tradingSystems", list)
	}

	return err
}

//=============================================================================

func UnassignTradingSystemsFromPortfolio(tx *gorm.DB, c *auth.Context, id uint, list []uint) error {
	_,err := getPortfolio(tx, c, id, "UnassignTradingSystemsFromPortfolio")
	if err != nil {
		return err
	}

	filter := map[string]any{}
	filter["portfolio_id"] = id

	assigned,errX := db.GetTradingSystems(tx, filter, 0, 5000)
	if errX != nil {
		return errX
	}

	assMap := buildAssignedMap(assigned)
	for _, tsId := range list {
		if _,ok := assMap[tsId]; !ok {
			return req.NewBadRequestError("Trading system with id=%v is not assignable to portfolio", tsId)
		}
	}

	err = db.UpdatePortfolioForTradingSystems(tx, list, nil)

	if err == nil {
		c.Log.Info("UnassignTradingSystemsFromPortfolio: Unassigned portfolio to trading systems", "portfolioId", id, "tradingSystems", list)
	}

	return err
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func getPortfolio(tx *gorm.DB, c *auth.Context, id uint, function string) (*db.Portfolio, error) {
	p, err := db.GetPortfolioById(tx, id)

	if err != nil {
		c.Log.Error(function+": Could not retrieve portfolio", "error", err.Error())
		return nil, err
	}

	if p == nil {
		c.Log.Error(function+": Portfolio was not found", "id", id)
		return nil, req.NewNotFoundError("Portfolio was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if p.Username != c.Session.Username {
			c.Log.Error(function+": Portfolio not owned by user", "id", id)
			return nil, req.NewForbiddenError("Portfolio is not owned by user: %v", id)
		}
	}

	return p, nil
}

//=============================================================================

func buildAssignableMap(list *[]db.TradingSystemAssignable) map[uint]bool {
	var result map[uint]bool = make(map[uint]bool)
	for _, ts := range *list {
		result[ts.Id] = true
	}

	return result
}

//=============================================================================

func buildAssignedMap(list *[]db.TradingSystem) map[uint]bool {
	var result map[uint]bool = make(map[uint]bool)
	for _, ts := range *list {
		result[ts.Id] = true
	}

	return result
}

//=============================================================================
