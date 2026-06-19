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
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

func GetTradeAnalysis(ts *db.TradingSystem, trades *[]db.Trade, barMap map[int64][]*db.EquityBar) (*AnalysisResponse, error) {
	res := NewAnalysisResponse()
	res.TradingSystem = ts

	var ti *Entry

	for _, trade := range *trades {
		bars, found := barMap[trade.Id]
		if found {
			ti = buildTradeInfo(&trade, bars, ts.CostPerOperation)
		} else {
			ti = buildTradeInfo(&trade, nil, ts.CostPerOperation)
		}
		res.Trades = append(res.Trades, ti)
	}

	return res, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func buildTradeInfo(tr *db.Trade, bars []*db.EquityBar, costPerOper float64) *Entry {

	var grossEq   []float64
	var netEq     []float64
	var contracts []int

	for _, eb := range bars {
		grossEq   = append(grossEq,   eb.GrossReturn)
		netEq     = append(netEq,     eb.GrossReturn - 2*costPerOper)
		contracts = append(contracts, eb.Contracts)
	}

	return &Entry{
		TradeType   : tr.TradeType,
		EntryDate   : tr.EntryDate,
		EntryLabel  : tr.EntryLabel,
		ExitDate    : tr.ExitDate,
		ExitLabel   : tr.ExitLabel,
		GrossReturn : tr.GrossReturn,
		MaxContracts: tr.MaxContracts,
		GrossEquity : buildEquity(grossEq, tr.GrossReturn, 0),
		NetEquity   : buildEquity(  netEq, tr.GrossReturn,  costPerOper),
		Contracts   : contracts,
	}
}

//=============================================================================

func buildEquity(eq []float64, ret float64, costPerOper float64) *EquityInfo {
	runUp, drawd := core.CalcRunUpAndDrawdown(eq)

	return &EquityInfo{
		Equity  : eq,
		Return  : ret - 2*costPerOper,
		RunUp   : runUp,
		Drawdown: drawd,
	}
}

//=============================================================================
