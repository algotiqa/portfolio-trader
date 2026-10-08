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
	"encoding/json"
	"time"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/portfolio-trader/pkg/business/allocation"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
	"gorm.io/gorm"
)

//=============================================================================

// TradeYears Start taking trades from 5 years ago
const TradeYears = 5

//=============================================================================

func calcAllocation(a *db.Allocation) error {
	job,err := retrieveInfo(a)
	if err != nil {
		return err
	}

	report := job.BuildPortfolioAllocation()
	return saveAllocationResults(a, report)
}

//=============================================================================

func retrieveInfo(a *db.Allocation) (*allocation.Job,error) {
	var job = allocation.NewJob(buildSpec(a))

	fromTrade  := time.Now().AddDate(-TradeYears,0,0)
	fromReturn := time.Now().AddDate(0, 0, -a.CorrelationPeriod)
	fromRetDate:= types.ToDate(&fromReturn)

	var tsList *[]db.TradingSystem
	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		p,err := db.GetPortfolioById(tx, a.PortfolioId)
		if err != nil {
			return err
		}
		job.SetPortfolio(p)

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

				tsp,errp := db.GetTradingPositionByTsId(tx, ts.Id)
				if errp != nil {
					return errp
				}

				core.NormalizeContracts(trades)
				tsi := allocation.NewTradingSystemInfo(&ts, tsf, tsp, trades, returns)
				job.AddTradingSystem(tsi)
			}

			a.Status = db.AllocStatusRunning
			errx = db.UpdateAllocation(tx, a)
		}

		return errx
	})

	return job,err
}

//=============================================================================

func buildSpec(a *db.Allocation) *allocation.JobSpec {
	return &allocation.JobSpec{
		AccountCapital: a.AccountCapital,
		AccountPerc   : a.AccountPerc,
		MaxMarginPerc : a.MaxMarginPerc,
		TradeYears    : TradeYears,
	}
}

//=============================================================================

func saveAllocationResults(a *db.Allocation, report *allocation.Report) error {
	data,err := json.Marshal(report)
	if err != nil {
		return err
	}

	return dbms.RunInTransaction(func(tx *gorm.DB) error {
		a.Status = calcAllocationStatus(report)
		a.Report = string(data)
		return db.UpdateAllocation(tx, a)
	})
}

//=============================================================================

func calcAllocationStatus(r *allocation.Report) db.AllocStatus {

	//--- Errors

	for _, log := range r.Logs {
		if log.Level == allocation.LogLevelError {
			return db.AllocStatusErrors
		}
	}

	//--- Warnings

	for _, ts := range r.TradingSystems {
		if ts.HealthLevel != allocation.HealthLevelOk {
			return db.AllocStatusWarnings
		}
	}

	for _, log := range r.Logs {
		if log.Level == allocation.LogLevelWarning {
			return db.AllocStatusWarnings
		}
	}

	//--- Ok

	return db.AllocStatusDone
}

//=============================================================================
