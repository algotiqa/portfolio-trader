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
	"strconv"
	"time"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
	"gorm.io/gorm"
)

//=============================================================================

func calcAllocation(a *db.Allocation) error {
	job,err := retrieveInfo(a)
	if err != nil {
		return err
	}

	if job.tradingSystems == nil {
		job.Log(db.LogLevelInfo, "No trading systems assigned to portfolio")
	} else if len(job.tradingSystems) < 2 {
		job.Log(db.LogLevelInfo, "Only 1 trading system assigned to portfolio")
	} else {
		calcFilterActivation(job)
		calcSystemCorrelation(job)
	}

	return saveAllocationResults(job)
}

//=============================================================================

func retrieveInfo(a *db.Allocation) (*AllocationJob,error) {
	var job = AllocationJob{
		allocation: a,
	}

	//--- Start taking trades from 5 years ago
	fromTrade  := time.Now().AddDate(-5,0,0)
	fromReturn := time.Now().AddDate(0, 0, -a.CorrelationPeriod)
	fromRetDate:= types.ToDate(&fromReturn)

	var tsList *[]db.TradingSystem
	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		filt := map[string]any{}
		filt["portfolio_id"] = a.PortfolioId

		var errx error
		tsList,errx = db.GetTradingSystems(tx, filt, 0, 5000)
		if errx == nil {
			for _,ts := range *tsList {
				trades,errt := db.FindTradesByTsIdFromTime(tx, ts.Id, &fromTrade, nil)
				if errt != nil {
					return errt
				}

				returns,errr := db.FindDailyReturnsByTsIdFromTime(tx, ts.Id, &fromRetDate, nil)
				if errr != nil {
					return errr
				}

				tsf,errf := db.GetTradingFilterByTsId(tx, ts.Id)
				if errf != nil {
					return errf
				}

				tsi := &TradingSystemInfo{
					system : &ts,
					filter : tsf,
					trades : trades,
					returns: returns,
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

func saveAllocationResults(job *AllocationJob) error {
	return dbms.RunInTransaction(func(tx *gorm.DB) error {
		job.allocation.Status = calcAllocationStatus(job.logs)

		err := db.UpdateAllocation(tx, job.allocation)
		if err == nil {
			err = saveAllocationFilters(tx, job.filters)
			if err == nil {
				err = saveAllocationLogs(tx, job.logs)
				if err == nil {
					err = saveSystemCorrelations(tx, job.correlations)
				}
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

func saveSystemCorrelations(tx *gorm.DB, correlations []*db.SystemCorrelation) error {
	for _, corr := range correlations {
		err := db.AddSystemCorrelation(tx, corr)
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func calcFilterActivation(job *AllocationJob) {
	for _, tsi := range job.tradingSystems {
		tsi.filterPassed = true
		if tsi.filter != nil {
			tsi.filterPassed = filter.CalcActivation(tsi.system, tsi.filter, *tsi.trades)
		}

		job.AddFilter(tsi.system.Id, tsi.filterPassed)
	}
}

//=============================================================================

func calcSystemCorrelation(job *AllocationJob) {
	var list []*TradingSystemInfo
	for _, tsi := range job.tradingSystems {
		if tsi.filterPassed {
			list = append(list, tsi)
		}
	}

	job.Log(db.LogLevelInfo, "Trading systems that passed the filter: "+ strconv.Itoa(len(list)))

	for i := 0; i < len(list) -1; i++ {
		for j := i + 1; j < len(list); j++ {
			ts1 := list[i].system
			ts2 := list[j].system

			corr,err := core.CalcCorrelation(list[i].returns, list[j].returns)
			message := ""
			if err != nil {
				message = err.Error()
			}

			sc := &db.SystemCorrelation{
				AllocationId    : job.allocation.Id,
				TradingSystem1Id: ts1.Id,
				TradingSystem2Id: ts2.Id,
				Correlation     : core.Trunc2d(corr),
				Message         : message,
			}

			job.correlations = append(job.correlations, sc)
		}
	}

	for _,tsi := range list {
		if len(*tsi.returns) < 2 {
			job.Log(db.LogLevelError, "Missing daily returns to calculate correlations for '"+tsi.system.Name+"'")
		} else if len(*tsi.returns) < MinDailyReturns {
			job.Log(db.LogLevelWarning, "Insufficient daily returns to calculate good correlations for '"+tsi.system.Name+"'")
		}
	}
}

//=============================================================================
