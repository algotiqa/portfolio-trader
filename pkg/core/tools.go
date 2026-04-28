//=============================================================================
/*
Copyright © 2024 Andrea Carboni andrea.carboni71@gmail.com

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

package core

import (
	"math"
	"slices"
	"time"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

func Trunc2d(value float64) float64 {
	return float64(int(math.Floor(value*100))) / 100
}

//=============================================================================

func GetLocation(timezone string, ts *db.TradingSystem) (*time.Location, error) {
	if timezone == "exchange" {
		timezone = ts.Timezone
	}

	return time.LoadLocation(timezone)
}

//=============================================================================

func ToNonZeroDailyReturnSlice(list *[]db.DailyReturn) []float64 {
	var res []float64

	for _, dr := range *list {
		if dr.GrossProfit != 0 {
			res = append(res, dr.GrossProfit)
		}
	}

	return res
}

//=============================================================================

func CalcRisk(returns []float64, costPerOper float64) (float64, error) {
	if returns == nil || len(returns) == 0 {
		return 0, req.NewUnprocessableEntityError("no losses found")
	}

	//--- Step 1: calculate distribution of losses

	counts := map[float64]int{}

	for _, ret := range returns {
		if ret < -2*costPerOper {
			count, ok := counts[ret]
			if !ok {
				counts[ret] = 1
			} else {
				counts[ret] = count + 1
			}
		}
	}

	if len(counts) == 0 {
		return 0, req.NewUnprocessableEntityError("no losses found")
	}

	//--- Step 2: Order all losses by their count

	var losses []*riskInfo

	for loss, count := range counts {
		losses = append(losses, &riskInfo{ loss, count })
	}

	slices.SortFunc(losses, func(a, b *riskInfo) int {
		return b.count - a.count
	})

	//--- Step 3: find the stop loss (the loss with most hits)

	if len(losses) == 1 {
		return math.Abs(losses[0].loss), nil
	}

	if losses[0].count > losses[1].count * 3 {
		return math.Abs(losses[0].loss), nil
	}

	// Step 4: Stop loss not identifies (maybe varies depending on ATR). Calc average

	sum := 0.0

	for _, ri := range losses {
		sum += ri.loss
	}

	loss := math.Abs(sum / float64(len(losses)))
	return Trunc2d(loss), nil
}

//-----------------------------------------------------------------------------

type riskInfo struct {
	loss  float64
	count int
}

//=============================================================================

func GetReturns(trades *[]db.Trade, tradeType string, costPerOper float64) []float64 {
	var list []float64

	for _, t := range *trades {
		if tradeType == db.TradeTypeAll || tradeType == t.TradeType {
			returns := t.GrossProfit - 2*costPerOper
			list = append(list, returns)
		}
	}

	return list
}

//=============================================================================

func CalcRMultiple(returns []float64, risk float64) []float64 {
	var list []float64

	for _, r := range returns {
		list = append(list, r/risk)
	}

	return list
}

//=============================================================================
