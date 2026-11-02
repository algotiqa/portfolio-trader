//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func SetTradingFilters(tx *gorm.DB, c *auth.Context, tsId uint, f *filter.TradingFilter) error {
	_, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	//--- Validate input

	_,err = filter.NewEvaluator(f)
	if err != nil {
		return err
	}

	tf,err := filter.ConvertToDbFilter(f)
	if err != nil {
		return err
	}

	tf.TradingSystemId = tsId

	return db.SetTradingFilter(tx, tf)
}

//=============================================================================

func RunFilterAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, far *filter.AnalysisRequest) (*filter.AnalysisResponse, error) {
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	//--- Get filtering config (if not provided)

	var filters *filter.TradingFilter

	if far.Filter == nil {
		var dbf *db.TradingFilter
		dbf, err = db.GetTradingFilterByTsId(tx, tsId)
		if err != nil {
			return nil, err
		}
		filters,err = filter.ConvertFromDbFilter(dbf)
		if err != nil {
			return nil, err
		}
	} else {
		filters = far.Filter
	}

	//--- Validate input, creating an evaluator

	eval,err := filter.NewEvaluator(filters)
	if err != nil {
		return nil, err
	}

	//--- Get trades

	fromTime, toTime, err := core.CalcSelectedPeriod(&far.Period, time.UTC)
	if err != nil {
		return nil, err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return nil, err
	}

	//--- Get portfolio

	if ts.PortfolioId == nil || *ts.PortfolioId == 0 {
		return nil, req.NewBadRequestError("missing portfolio id")
	}

	p,err := getPortfolio(tx, c, *ts.PortfolioId, "RunFilterAnalysis")
	if err != nil {
		return nil, err
	}

	//--- Run analysis

	core.NormalizeContracts(trades)
	res := filter.RunAnalysis(ts, filters, trades, p, eval)

	return res, err
}

//=============================================================================

func StartFilterOptimization(tx *gorm.DB, c *auth.Context, tsId uint, oreq *filter.OptimizationRequest) error {
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, oreq.StartDate, nil)
	if err != nil {
		return err
	}

	err = oreq.Validate()
	if err != nil {
		return err
	}

	c.Log.Info("StartFilterOptimization: Starting optimization", "tsId", ts.Id, "tsName", ts.Name)
	core.NormalizeContracts(trades)
	filter.StartOptimization(ts, trades, oreq)

	return nil
}

//=============================================================================

func StopFilterOptimization(c *auth.Context, tsId uint) error {
	c.Log.Info("StopFilterOptimization: Stopping optimization", "tsId", tsId)
	err := filter.StopOptimization(tsId)

	return err
}

//=============================================================================

func GetFilterOptimizationInfo(c *auth.Context, tsId uint) (*filter.OptimizationResponse, error) {
	info := filter.GetOptimizationInfo(tsId)
	return filter.NewOptimizationResponse(info), nil
}

//=============================================================================
