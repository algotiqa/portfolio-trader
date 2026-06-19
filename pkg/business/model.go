//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Portfolio tree
//===
//=============================================================================

type PortfolioTree struct {
	db.Portfolio
	Children       []*PortfolioTree    `json:"children"`
	TradingSystems []*db.TradingSystem `json:"tradingSystems"`
}

//-----------------------------------------------------------------------------

func (pt *PortfolioTree) AddChild(p *PortfolioTree) {
	pt.Children = append(pt.Children, p)
}

//-----------------------------------------------------------------------------

func (pt *PortfolioTree) AddTradingSystem(ts *db.TradingSystem) {
	pt.TradingSystems = append(pt.TradingSystems, ts)
}

//=============================================================================
