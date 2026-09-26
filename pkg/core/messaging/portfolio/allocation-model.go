//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package portfolio

import "github.com/algotiqa/portfolio-trader/pkg/db"

//=============================================================================

const MinDailyReturns = 15

//=============================================================================

type AllocationJob struct {
	allocation     *db.Allocation
	tradingSystems []*TradingSystemInfo
	filters        []*db.AllocationFilter
	logs           []*db.AllocationLog
	correlations   []*db.SystemCorrelation
}

//=============================================================================

func (j *AllocationJob) Log(level db.LogLevel, message string) {
	log := &db.AllocationLog{
		AllocationId: j.allocation.Id,
		Level       : level,
		Message     : message,
	}

	j.logs = append(j.logs, log)
}

//=============================================================================

func (j *AllocationJob) AddFilter(tsId uint, filterPassed bool) {
	af := &db.AllocationFilter{
		AllocationId   : j.allocation.Id,
		TradingSystemId: tsId,
		FilterPassed   : filterPassed,
		Comment        : "",
	}

	j.filters = append(j.filters, af)
}

//=============================================================================

type TradingSystemInfo struct {
	system        *db.TradingSystem
	filter        *db.TradingFilter
	filterPassed  bool
	trades        *[]db.Trade
	returns       *[]db.DailyReturn
}

//=============================================================================
