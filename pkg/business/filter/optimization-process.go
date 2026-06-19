//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

const MaxResultSize = 1000

//=============================================================================
//===
//=== OptimizationProcess
//===
//=============================================================================

type OptimizationProcess struct {
	ts              *db.TradingSystem
	trades          *[]db.Trade
	optReq          *OptimizationRequest
	info            *OptimizationInfo
	fitnessFunction FitnessFunction
	stopping        bool
}

//=============================================================================

func (op *OptimizationProcess) Start() {
	field := op.optReq.FieldToOptimize
	startDate := op.optReq.StartDate

	op.fitnessFunction = GetFitnessFunction(field)

	algo := algorithm.New(op.optReq.Algorithm.Type)
	ctx := NewContext(op)
	algo.Init(ctx)

	fc := op.optReq.FilterConfig
	op.info = NewOptimizationInfo(MaxResultSize, field, fc, algo.StepsCount(), op.calcBaseValue(), startDate)

	go op.generate(algo)
}

//=============================================================================

func (op *OptimizationProcess) Stop() {
	slog.Info("Stop: Stopping optimization process", "tsId", op.ts.Id)
	op.stopping = true
}

//=============================================================================

func (op *OptimizationProcess) GetInfo() *OptimizationInfo {
	return op.info
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

//--- GoRoutine

func (op *OptimizationProcess) generate(algo optimization.Algorithm) {
	slog.Info("generate: Started", "tsId", op.ts.Id, "tsName", op.ts.Name, "algorithm", op.optReq.Algorithm)

	algo.Optimize()

	for !op.info.isStatusComplete() {
		time.Sleep(time.Second * 1)
	}

	slog.Info("generate: Complete.")
}

//=============================================================================

func (op *OptimizationProcess) runAnalysis(filter *db.TradingFilter) float64 {
	sum := RunAnalysis(op.ts, filter, op.trades).Summary
	run := op.createRun(filter, &sum)
	op.info.addResult(run)

	return run.FitnessValue
}

//=============================================================================

func (op *OptimizationProcess) createRun(filter *db.TradingFilter, sum *Summary) *Run {
	r := &Run{
		Filter:      filter,
		NetProfit:   sum.FilProfit,
		AvgTrade:    sum.FilAverageTrade,
		MaxDrawdown: sum.FilMaxDrawdown,
		random:      rand.Int(),
	}

	r.FitnessValue = op.fitnessFunction(r)

	return r
}

//=============================================================================

func (op *OptimizationProcess) calcBaseValue() float64 {
	baseline := op.optReq.Baseline

	sum := RunAnalysis(op.ts, baseline, op.trades).Summary
	run := op.createRun(baseline, &sum)

	return run.FitnessValue
}

//=============================================================================
