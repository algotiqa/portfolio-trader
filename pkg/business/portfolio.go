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
	"log/slog"

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

	return db.GetPortfolios(tx, filter, offset, limit)
}

//=============================================================================

func GetPortfolioTree(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]*PortfolioTree, error) {

	//--- The only valid filter can be the username

	//--- Get all portfolios

	poList, err := GetPortfolios(tx, c, filter, offset, limit)
	if err != nil {
		return nil, req.NewServerErrorByError(err)
	}

	//--- Get all trading systems

	tsList, err := GetTradingSystems(tx, c, filter, offset, limit, false)
	if err != nil {
		return nil, req.NewServerErrorByError(err)
	}

	return buildPortfolioTree(c.Log, poList, tsList), nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func buildPortfolioTree(log *slog.Logger, poList *[]db.Portfolio, tsList *[]db.TradingSystem) *[]*PortfolioTree {

	//--- Step 1: Collect all nodes into a map

	nodeMap := map[uint]*PortfolioTree{}
	fullMap := map[uint]*PortfolioTree{}

	for _, p := range *poList {
		pt := &PortfolioTree{
			Portfolio:      p,
			Children:       []*PortfolioTree{},
			TradingSystems: []*db.TradingSystem{},
		}
		nodeMap[p.Id] = pt
		fullMap[p.Id] = pt
	}

	//--- Step 2: Build the tree

	for key, p := range fullMap {
		if p.ParentId != 0 {
			parent := fullMap[p.ParentId]
			parent.AddChild(p)
			delete(nodeMap, key)
		}
	}

	//--- Step 2: Add trading system information

	for _, ts := range *tsList {
		aux := ts
		portfolio := fullMap[*ts.PortfolioId]
		portfolio.AddTradingSystem(&aux)
	}

	//--- Step 3: Return tree

	if len(*poList) > 0 && len(nodeMap) == 0 {
		log.Error("Portfolios have circular loops (!)")
	}

	var result []*PortfolioTree

	for _, p := range nodeMap {
		result = append(result, p)
	}

	return &result
}

//=============================================================================
