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

	livePeriods, err := db.FindLivePeriodsByTradingSystemId(tx, ts.Id)
	if err != nil {
		return nil, err
	}
	shiftLiveTimezone(livePeriods, loc)

	res := performance.GetPerformanceAnalysis(ts, trades, livePeriods)

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
