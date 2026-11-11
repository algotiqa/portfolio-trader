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
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/model"
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
	fromTime        *time.Time
	portfolio       *db.Portfolio
	stopping        bool
}

//=============================================================================

func (op *OptimizationProcess) Start() {
	field := op.optReq.FieldToOptimize

	op.fitnessFunction = GetFitnessFunction(field)

	algo := algorithm.New(op.optReq.Algorithm.Type)
	ctx  := NewContext(op)
	algo.Init(ctx)

	fc := op.optReq.FilterConfig
	base,err := op.calcBaseValue()
	op.info = NewOptimizationInfo(MaxResultSize, field, fc, algo.Steps(), base, op.fromTime)

	if err != nil {
		op.ExitWithError(err)
	} else {
		slog.Info("generate: Started", "tsId", op.ts.Id, "tsName", op.ts.Name, "algorithm", op.optReq.Algorithm)
		go algo.Optimize()
	}
}

//=============================================================================

func (op *OptimizationProcess) Stop() {
	slog.Info("Stop: Stopping filter optimization process", "tsId", op.ts.Id)
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

func (op *OptimizationProcess) runAnalysis(filter *model.TradingFilter) (float64,error) {
	eval,err := NewEvaluator(filter)
	if err != nil {
		return 0, err
	}

	sum := RunAnalysis(op.ts, filter, op.trades, op.portfolio, eval).Summary
	run := op.createRun(filter, &sum)
	op.info.addResult(run)

	return run.FitnessValue, nil
}

//=============================================================================

func (op *OptimizationProcess) createRun(filter *model.TradingFilter, sum *Summary) *Run {
	r := &Run{
		Filter     : filter,
		NetReturn  : sum.FilProfit,
		AvgTrade   : sum.FilAverageTrade,
		MaxDrawdown: sum.FilMaxDrawdown,
		random     : rand.Int(),
	}

	r.FitnessValue = op.fitnessFunction(r)

	return r
}

//=============================================================================

func (op *OptimizationProcess) calcBaseValue() (float64,error) {
	baseline := op.optReq.Baseline

	eval,err := NewEvaluator(baseline)
	if err != nil {
		return 0, err
	}

	sum := RunAnalysis(op.ts, baseline, op.trades, op.portfolio, eval).Summary
	run := op.createRun(baseline, &sum)

	return run.FitnessValue, nil
}

//=============================================================================

func (op *OptimizationProcess) EndProcess() {
	op.info.EndTime = time.Now()
	op.info.Status  = OptimStatusComplete

	slog.Info("EndProcess: Completed.")
}

//=============================================================================

func (op *OptimizationProcess) ExitWithError(err error) {
	op.info.EndTime = time.Now()
	op.info.Status  = OptimStatusError
	op.info.Error   = err
	op.stopping     = true

	slog.Error(err.Error(), "tsId", op.ts.Id, "tsName", op.ts.Name)
}

//=============================================================================
