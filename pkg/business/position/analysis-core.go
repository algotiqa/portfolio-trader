//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
