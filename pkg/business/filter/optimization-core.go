//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import (
	"log/slog"
	"sync"
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

var jobs = struct {
	sync.RWMutex
	m map[uint]*OptimizationProcess
}{m: make(map[uint]*OptimizationProcess)}

//-----------------------------------------------------------------------------

var workers = core.WorkerPool{}

//=============================================================================
//===
//=== Init
//===
//=============================================================================

func init() {
	num := 1
	workers.Init(num, 100)
	go periodicCleanup()
}

//=============================================================================
//===
//=== API methods
//===
//=============================================================================

func StartOptimization(ts *db.TradingSystem, trades *[]db.Trade, or *OptimizationRequest) {
	jobs.Lock()
	defer jobs.Unlock()

	fop, ok := jobs.m[ts.Id]
	if ok {
		slog.Error("Stopping a previous optimization process", "tsId", ts.Id)
		fop.Stop()
		delete(jobs.m, ts.Id)
	}

	fop = &OptimizationProcess{
		ts:     ts,
		trades: trades,
		optReq: or,
	}

	fop.Start()
	jobs.m[ts.Id] = fop
}

//=============================================================================

func StopOptimization(tsId uint) error {
	jobs.Lock()
	defer jobs.Unlock()

	fop, ok := jobs.m[tsId]
	if ok {
		fop.Stop()
	}

	return nil
}

//=============================================================================

func GetOptimizationInfo(tsId uint) *OptimizationInfo {
	jobs.Lock()
	defer jobs.Unlock()

	fop, ok := jobs.m[tsId]
	if !ok {
		return &OptimizationInfo{
			Status: OptimStatusIdle,
		}
	}

	return fop.GetInfo()
}

//=============================================================================
//===
//=== Cleanup process
//===
//=============================================================================

func periodicCleanup() {
	for {
		time.Sleep(time.Minute * 5)
		purge()
	}
}

//=============================================================================

func purge() {
	jobs.Lock()
	defer jobs.Unlock()

	for tsId, op := range jobs.m {
		if op.info.Status == OptimStatusComplete {
			delta := time.Now().Sub(op.info.EndTime)
			if delta.Minutes() >= 30 {
				slog.Info("purge: Purging optimization process entry for trading system", "tsId", tsId)
				delete(jobs.m, tsId)
			}
		}
	}
}

//=============================================================================
