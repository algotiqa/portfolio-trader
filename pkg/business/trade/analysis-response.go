//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package trade

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type AnalysisResponse struct {
	TradingSystem *db.TradingSystem `json:"tradingSystem"`
	Trades        []*Entry          `json:"trades"`
}

//=============================================================================

func NewAnalysisResponse() *AnalysisResponse {
	return &AnalysisResponse{
	}
}

//=============================================================================

type Entry struct {
	TradeType    string     `json:"tradeType"`
	EntryDate    *time.Time `json:"entryDate"`
	EntryLabel   string     `json:"entryLabel"`
	ExitDate     *time.Time `json:"exitDate"`
	ExitLabel    string     `json:"exitLabel"`
	GrossReturn  float64    `json:"grossReturn"`
	MaxContracts int        `json:"maxContracts"`

	GrossEquity *EquityInfo  `json:"grossEquity"`
	NetEquity   *EquityInfo  `json:"netEquity"`
	Contracts   []int        `json:"contracts"`
}

//=============================================================================

type EquityInfo struct {
	Equity   []float64  `json:"equity"`
	Return   float64    `json:"return"`
	RunUp    float64    `json:"runUp"`
	Drawdown float64    `json:"drawdown"`
}

//=============================================================================
