//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package portfolio

import (
	"time"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

type AllocationJob struct {
	allocation     *db.Allocation
	tradingSystems []*TradingSystemInfo
	filters        []*db.AllocationFilter
	logs           []*db.AllocationLog
}

//=============================================================================

type TradingSystemInfo struct {
	system *db.TradingSystem
	filter *db.TradingFilter
	trades *[]db.Trade
}

//=============================================================================

func calcAllocation(a *db.Allocation) error {
	job,err := retrieveInfo(a)
	if err != nil {
		return err
	}

	if job.tradingSystems == nil {
		addError(job, "No trading systems assigned to portfolio")
	} else {
		calcFilterActivation(job)
	}

	return saveAllocationResults(job)
}

//=============================================================================

func retrieveInfo(a *db.Allocation) (*AllocationJob,error) {
	var job = AllocationJob{
		allocation: a,
	}

	//--- Start taking trades from 5 years ago
	from := time.Now().AddDate(-5,0,0)

	var tsList *[]db.TradingSystem
	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		filter := map[string]any{}
		filter["portfolio_id"] = a.PortfolioId

		var errx error
		tsList,errx = db.GetTradingSystems(tx, filter, 0, 5000)
		if errx == nil {
			for _,ts := range *tsList {
				trades,errt := db.FindTradesByTsIdFromTime(tx, ts.Id, &from, nil)
				if errt != nil {
					return errt
				}

				tsf,errf := db.GetTradingFilterByTsId(tx, ts.Id)
				if errf != nil {
					return errf
				}

				tsi := &TradingSystemInfo{
					system: &ts,
					filter: tsf,
					trades: trades,
				}
				job.tradingSystems = append(job.tradingSystems, tsi)
			}

			a.Status = db.AllocStatusRunning
			errx = db.UpdateAllocation(tx, a)
		}

		return errx
	})

	return &job,err
}

//=============================================================================

func addError(job *AllocationJob, message string) {
	log := &db.AllocationLog{
		AllocationId: job.allocation.Id,
		Level       : db.LogLevelError,
		Message     : message,
	}

	job.logs = append(job.logs, log)
}

//=============================================================================

func saveAllocationResults(job *AllocationJob) error {
	return dbms.RunInTransaction(func(tx *gorm.DB) error {
		job.allocation.Status = calcAllocationStatus(job.logs)

		err := db.UpdateAllocation(tx, job.allocation)
		if err == nil {
			err = saveAllocationFilters(tx, job.filters)
			if err == nil {
				err = saveAllocationLogs(tx, job.logs)
			}
		}
		return err
	})
}

//=============================================================================

func calcAllocationStatus(list []*db.AllocationLog) db.AllocStatus {
	for _, log := range list {
		if log.Level == db.LogLevelError {
			return db.AllocStatusErrors
		}
	}

	for _, log := range list {
		if log.Level == db.LogLevelWarning {
			return db.AllocStatusWarnings
		}
	}

	return db.AllocStatusDone
}

//=============================================================================

func saveAllocationFilters(tx *gorm.DB, filters []*db.AllocationFilter) error {
	for _, filter := range filters {
		err := db.AddAllocationFilter(tx, filter)
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func saveAllocationLogs(tx *gorm.DB, logs []*db.AllocationLog) error {
	for _, log := range logs {
		err := db.AddAllocationLog(tx, log)
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func calcFilterActivation(job *AllocationJob) {
	for _, tsi := range job.tradingSystems {
		activation := true
		if tsi.filter != nil {
			activation = filter.CalcActivation(tsi.system, tsi.filter, *tsi.trades)
		}

		af := &db.AllocationFilter{
			AllocationId   : job.allocation.Id,
			TradingSystemId: tsi.system.Id,
			FilterPassed   : activation,
			Comment        : "",
		}

		job.filters = append(job.filters, af)
	}
}

//=============================================================================
