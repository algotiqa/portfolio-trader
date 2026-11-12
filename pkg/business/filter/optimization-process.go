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
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
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
	workers         core.WorkerPool
}

//=============================================================================

func (op *OptimizationProcess) Start() {
	field := op.optReq.FieldToOptimize
	op.fitnessFunction = GetFitnessFunction(field)
	fc := op.optReq.FilterConfig
	fc.Deactivation.WarmUp()

	algo := algorithm.New(op.optReq.Algorithm.Type)
	algo.Init(op)

	base,err := op.calcBaseValue()

	if err != nil {
		op.exitWithError(err)
	} else {
		op.info = NewOptimizationInfo(MaxResultSize, field, fc, algo.StepsCount(), base, op.fromTime)
		go op.generate(algo)
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
//=== Context implementation
//===
//=============================================================================

func (op *OptimizationProcess) FilterConfig() *optimization.FilterConfig {
	return op.optReq.FilterConfig
}

//=============================================================================

func (op *OptimizationProcess) AlgorithmConfig() *optimization.AlgorithmConfig {
	return &op.optReq.Algorithm.Config
}

//=============================================================================
//--- To protect

func (op *OptimizationProcess) IsStopping() bool {
	return op.stopping
}

//=============================================================================

func (op *OptimizationProcess) Baseline() model.TradingFilter {
	return *op.optReq.Baseline
}

//=============================================================================

func (op *OptimizationProcess) LogInfo(message string) {
	slog.Info(message, "tsId", op.ts.Id, "tsName", op.ts.Name)
}

//=============================================================================
//--- To protect

func (op *OptimizationProcess) RunAnalysis(filter *model.TradingFilter) (float64,bool) {
	eval,err := NewEvaluator(filter)
	if err != nil {
		op.exitWithError(err)
		return 0, false
	}

	sum := RunAnalysis(op.ts, filter, op.trades, op.portfolio, eval).Summary
	run := op.createRun(filter, &sum)
	op.info.addResult(run)						//	<-- protected

	return run.FitnessValue, true
}

//=============================================================================

func (op *OptimizationProcess) Submit(task func()) {
	op.workers.Submit(task)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

//--- Running in a Go routine

func (op *OptimizationProcess) generate(algo optimization.Algorithm) {
	slog.Info("generate: Started", "tsId", op.ts.Id, "tsName", op.ts.Name, "algorithm", op.optReq.Algorithm)
	op.workers.Init(8, 500)
	algo.Optimize()
	op.workers.Wait()

	if op.info.Status == OptimStatusRunning {
		op.info.EndTime = time.Now()
		op.info.Status  = OptimStatusComplete
	}

	slog.Info("EndProcess: Completed.")
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

func (op *OptimizationProcess) exitWithError(err error) {
	op.info.EndTime = time.Now()
	op.info.Status  = OptimStatusError
	op.info.Error   = err
	op.stopping     = true

	slog.Error(err.Error(), "tsId", op.ts.Id, "tsName", op.ts.Name)
}

//=============================================================================
