//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Filter activation calculation
//===
//=============================================================================

func CalcActivation(ts *db.TradingSystem, filter *TradingFilter, list *[]db.Trade, p *db.Portfolio) bool {
	if len(*list) == 0 {
		return true
	}

	eval,err := NewEvaluator(filter)
	if err != nil {
		panic("cannot create evaluator:" +err.Error())
	}

	//--- Calc unfiltered equity and net profits
	s := calcSnapshot(ts, list)

	//--- Calc unfiltered drawdown
	unfilteredDrawdown, _ := core.BuildDrawDown(&s.Equity)
	s.Drawdown = *unfilteredDrawdown

	eval.Eval(s, p)

	return eval.sampledFinal[len(eval.sampledFinal) -1] == 1
}

//=============================================================================
//===
//=== AnalysisResponse building
//===
//=============================================================================

func RunAnalysis(ts *db.TradingSystem, filter *TradingFilter, list *[]db.Trade, p *db.Portfolio, eval *Evaluator) *AnalysisResponse {
	res := &AnalysisResponse{}
	res.TradingSystem.Id   = ts.Id
	res.TradingSystem.Name = ts.Name
	res.Filter             = filter
	res.Portfolio          = p

	if len(*list) == 0 {
		return res
	}

	//--- Calc unfiltered equity and net profits
	s := calcSnapshot(ts, list)

	//--- Calc unfiltered drawdown
	unfilteredDrawdown, maxUnfDD := core.BuildDrawDown(&s.Equity)
	s.Drawdown = *unfilteredDrawdown

	eval.Eval(s, p)

	e := &res.Equities
	e.Time               = s.Time
	e.UnfilteredEquity   = s.Equity
	e.UnfilteredDrawdown = s.Drawdown
	e.FilteredEquity     = calcFilteredEquity(s.NetReturn, eval.sampledFinal)

	filteredDrawdown, maxFilDD := core.BuildDrawDown(&e.FilteredEquity)
	e.FilteredDrawdown   = *filteredDrawdown

	res.Outcomes = eval.Outcomes()
	core.Trunc2dSlice(e.UnfilteredDrawdown)
	core.Trunc2dSlice(e.FilteredDrawdown)

	calcSummary(res, maxUnfDD, maxFilDD, s.NetReturn)

	return res
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcSnapshot(ts *db.TradingSystem, tradeList *[]db.Trade) *criteria.Snapshot {
	netEquity     := 0.0
	costPerOperat := float64(ts.CostPerOperation)
	size          := len(*tradeList)

	s := &criteria.Snapshot{
		Time      : make([]time.Time, size),
		NetReturn : make([]float64  , size),
		Equity    : make([]float64  , size),
	}

	for i, t := range *tradeList {
		netProfit := core.NetReturn(t.GrossReturn, costPerOperat)
		netEquity += netProfit

		s.Time     [i] = *t.ExitDate
		s.Equity   [i] = core.Trunc2d(netEquity)
		s.NetReturn[i] = core.Trunc2d(netProfit)
	}

	return s
}

//=============================================================================

func calcFilteredEquity(netProfit []float64, activations []int) []float64 {
	sum := 0.0

	filteredEquity := make([]float64, len(netProfit))

	for i, value := range netProfit {
		if i > 0 {
			//--- We have to use the activation at time [i-1]
			if activations[i-1] == 0 {
				value = 0.0
			}
		}

		sum += value
		filteredEquity[i] = core.Trunc2d(sum)
	}

	return filteredEquity
}

//=============================================================================

func calcSummary(res *AnalysisResponse, maxUnfDD, maxFilDD float64, netReturn []float64) {
	sum := &res.Summary
	last := len(res.Equities.Time) - 1

	sum.UnfProfit       = core.Trunc2d(res.Equities.UnfilteredEquity[last])
	sum.FilProfit       = core.Trunc2d(res.Equities.FilteredEquity[last])
	sum.UnfMaxDrawdown  = core.Trunc2d(maxUnfDD)
	sum.FilMaxDrawdown  = core.Trunc2d(maxFilDD)
	sum.UnfWinningPerc  = core.CalcWinningPercentage(netReturn, nil)
	sum.FilWinningPerc  = core.CalcWinningPercentage(netReturn, res.Outcomes.SampledFinal)
	sum.UnfAverageTrade = core.CalcAverageTrade(netReturn, nil)
	sum.FilAverageTrade = core.CalcAverageTrade(netReturn, res.Outcomes.SampledFinal)
}

//=============================================================================
