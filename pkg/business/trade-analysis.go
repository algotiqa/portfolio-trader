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

package business

import (
	"slices"
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/portfolio-trader/pkg/business/trade"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func RunTradeAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, req *trade.AnalysisRequest) (*trade.AnalysisResponse, error) {
	//--- Get trading system

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	//--- Get location of timezone to shift dates

	loc, err := time.LoadLocation(ts.Timezone)
	if err != nil {
		c.Log.Error("RunTradeAnalysis: Bad timezone", "timezone", ts.Timezone, "error", err)
		return nil, err
	}

	fromDate,toDate,err := core.CalcSelectedPeriod(&req.SelectedPeriod, loc)
	if err != nil {
		return nil, err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	shiftTradesTimezone(trades, loc)

	barMap, err := getEquityBars(tx, trades)
	if err != nil {
		return nil, err
	}

	return trade.GetTradeAnalysis(ts, trades, barMap)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func getEquityBars(tx *gorm.DB, trades *[]db.Trade) (map[int64][]*db.EquityBar, error) {

	//--- Find all bars associated to given trades

	var ids []int64
	for _, tr := range *trades {
		ids = append(ids, tr.Id)
	}

	list,err := db.FindEquityBarsByTradesId(tx, ids)
	if err != nil {
		return nil, err
	}

	//--- Build a map to group bars by trade

	barMap := make(map[int64][]*db.EquityBar)

	for _, eb := range *list {
		barList,found := barMap[eb.TradeId]
		if !found {
			barList = make([]*db.EquityBar, 0)
		}

		barList = append(barList, &eb)
		barMap[eb.TradeId] = barList
	}

	for _,bars := range barMap {
		slices.SortFunc(bars, func(a, b *db.EquityBar) int {
			return int(a.Date.Unix() - b.Date.Unix())
		})
	}

	return barMap, nil
}

//=============================================================================
