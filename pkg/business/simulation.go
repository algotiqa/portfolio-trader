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
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/simulation"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func StartSimulation(tx *gorm.DB, c *auth.Context, tsId uint, rq *simulation.Request) error {
	//--- Get trading system

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	c.Log.Info("StartSimulation: Starting", "id", tsId, "name", ts.Name, "runs", rq.Runs)

	fromDate,toDate,err := core.CalcSelectedPeriod(&rq.SelectedPeriod, time.UTC)
	if err != nil {
		return err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromDate, toDate)
	if err != nil {
		return err
	}

	if len(*trades) == 0 {
		return req.NewUnprocessableEntityError("no trades found for given time")
	}

	simulation.Start(rq, ts, trades)
	c.Log.Info("StartSimulation: Ending", "id", tsId, "name", ts.Name, "runs", rq.Runs)
	return nil
}

//=============================================================================

func StopSimulation(c *auth.Context, tsId uint) bool {
	return simulation.Stop(tsId)
}

//=============================================================================

func GetSimulationResult(c *auth.Context, tsId uint) *simulation.Result {
	return simulation.GetResult(tsId)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================
