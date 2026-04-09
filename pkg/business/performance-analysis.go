//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/portfolio-trader/pkg/business/performance"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func RunPerformanceAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, req *performance.AnalysisRequest) (*performance.AnalysisResponse, error) {

	//--- Get trading system

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	//--- Get location of timezone to shift dates

	loc, err := core.GetLocation(req.Timezone, ts)
	if err != nil {
		c.Log.Error("RunPerformanceAnalysis: Bad timezone", "timezone", req.Timezone, "error", err)
		return nil, err
	}

	fromTime, toTime, err := core.CalcSelectedPeriod(&req.SelectedPeriod, loc)
	if err != nil {
		c.Log.Error("RunPerformanceAnalysis: Bad fromDate or toDate", "fromDate", req.FromDate, "toDate", req.ToDate, "error", err)
		return nil, err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return nil, err
	}
	shiftTradesTimezone(trades, loc)

	returns, err := db.FindDailyReturnsByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return nil, err
	}

	livePeriods, err := db.FindLivePeriodsByTradingSystemId(tx, ts.Id)
	if err != nil {
		return nil, err
	}
	shiftLiveTimezone(livePeriods, loc)

	res := performance.GetPerformanceAnalysis(ts, trades, returns, livePeriods)

	return res, nil
}

//=============================================================================

func shiftTradesTimezone(trades *[]db.Trade, loc *time.Location) {
	for i := 0; i < len(*trades); i++ {
		tr := &(*trades)[i]
		tr.EntryDate         = shiftLocation(tr.EntryDate,         loc)
		tr.ExitDate          = shiftLocation(tr.ExitDate,          loc)
		tr.EntryDateAtBroker = shiftLocation(tr.EntryDateAtBroker, loc)
		tr.ExitDateAtBroker  = shiftLocation(tr.ExitDateAtBroker,  loc)
	}
}

//=============================================================================

func shiftLiveTimezone(periods *[]db.LivePeriod, loc *time.Location) {
	for i := 0; i < len(*periods); i++ {
		lp := &(*periods)[i]
		aux := shiftLocation(&lp.Period, loc)
		lp.Period = *aux
	}
}

//=============================================================================

func shiftLocation(t *time.Time, loc *time.Location) *time.Time {
	if t == nil {
		return nil
	}

	out := t.In(loc)

	return &out
}

//=============================================================================
