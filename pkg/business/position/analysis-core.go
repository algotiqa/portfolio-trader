//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package position

import (
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== AnalysisResponse building
//===
//=============================================================================

func RunAnalysis(ts *db.TradingSystem, model model.PositionModel, list *[]db.Trade) *AnalysisResponse {
	res := &AnalysisResponse{}
	//res.TradingSystem.Id = ts.Id
	//res.TradingSystem.Name = ts.Name
	//res.Filter = filter
	//
	////--- Creates slices
	//
	//if len(*list) == 0 {
	//	return res
	//}
	//
	//e := &res.Equities
	//
	////--- Calc unfiltered equity and days
	//calcUnfilteredEquityAndProfit(e, ts, list)
	//
	//if res.Filter.EquAvgEnabled {
	//	e.Average = calcAverageEquity(e.Time, e.UnfilteredEquity, res.Filter.EquAvgLen)
	//}
	//
	//res.Activations = calcActivations(e, filter)
	//calcFilterActivation(e, res.Activations, filter)
	//calcFilteredEquity(res)
	//
	//unfilteredDrawdown, maxUnfDD := core.BuildDrawDown(&e.UnfilteredEquity)
	//filteredDrawdown, maxFilDD := core.BuildDrawDown(&e.FilteredEquity)
	//
	//e.UnfilteredDrawdown = *unfilteredDrawdown
	//e.FilteredDrawdown = *filteredDrawdown
	//
	//calcSummary(res, maxUnfDD, maxFilDD)

	return res
}

//=============================================================================

//=============================================================================
