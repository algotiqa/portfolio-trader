//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package simulation

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
	m map[uint]*Process
}{m: make(map[uint]*Process)}

//-----------------------------------------------------------------------------

var workers = core.WorkerPool{}

//=============================================================================
//===
//=== Init
//===
//=============================================================================

func init() {
	num := 4
	workers.Init(num, 100)
	go periodicCleanup()
}

//=============================================================================
//===
//=== API methods
//===
//=============================================================================

func Start(req *Request, ts *db.TradingSystem, trades *[]db.Trade) {
	jobs.Lock()
	defer jobs.Unlock()

	sp, ok := jobs.m[ts.Id]
	if ok {
		slog.Error("Stopping a previous simulation process", "tsId", ts.Id)
		sp.Stop()
		delete(jobs.m, ts.Id)
	}

	sp = NewProcess(ts, trades, req)
	workers.Submit(sp.Start)

	jobs.m[ts.Id] = sp
}

//=============================================================================

func Stop(tsId uint) bool {
	jobs.Lock()
	defer jobs.Unlock()

	sp, ok := jobs.m[tsId]
	if ok {
		sp.Stop()
		delete(jobs.m, tsId)
	}

	return ok
}

//=============================================================================

func GetResult(tsId uint) *Result {
	jobs.Lock()
	defer jobs.Unlock()

	sp, ok := jobs.m[tsId]
	if !ok {
		return &Result{
			Status: SimStatusIdle,
		}
	}

	return sp.GetResult()
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

	for tsId, sp := range jobs.m {
		if sp.result.Status == SimStatusComplete {
			delta := time.Now().Sub(sp.result.EndTime)
			if delta.Minutes() >= 30 {
				slog.Info("purge: Purging simulation process entry for trading system", "tsId", tsId)
				delete(jobs.m, tsId)
			}
		}
	}
}

//=============================================================================
