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
	"github.com/algotiqa/portfolio-trader/pkg/business/position"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func RunPositionAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, par *position.AnalysisRequest) (*position.AnalysisResponse, error) {
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	fromTime, toTime, err := core.CalcSelectedPeriod(&par.SelectedPeriod, time.UTC)
	if err != nil {
		return nil, err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return nil, err
	}

	res := position.RunAnalysis(ts, nil, trades)

	return res, err
}

//=============================================================================

//func StartFilterOptimization(tx *gorm.DB, c *auth.Context, tsId uint, oreq *filter.OptimizationRequest) error {
//	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
//	if err != nil {
//		return err
//	}
//
//	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, oreq.StartDate, nil)
//	if err != nil {
//		return err
//	}
//
//	err = oreq.Validate()
//	if err != nil {
//		return err
//	}
//
//	c.Log.Info("StartFilterOptimization: Starting optimization", "tsId", ts.Id, "tsName", ts.Name)
//	filter.StartOptimization(ts, trades, oreq)
//
//	return nil
//}

//=============================================================================

//func StopFilterOptimization(c *auth.Context, tsId uint) error {
//	c.Log.Info("StopFilterOptimization: Stopping optimization", "tsId", tsId)
//	err := filter.StopOptimization(tsId)
//
//	return err
//}

//=============================================================================

//func GetFilterOptimizationInfo(c *auth.Context, tsId uint) (*filter.OptimizationResponse, error) {
//	info := filter.GetOptimizationInfo(tsId)
//	return filter.NewOptimizationResponse(info), nil
//}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================



//=============================================================================
