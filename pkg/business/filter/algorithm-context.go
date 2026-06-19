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

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type OptimizationContext struct {
	op *OptimizationProcess
}

//=============================================================================
//=== Constructor
//=============================================================================

func NewContext(op *OptimizationProcess) optimization.Context {
	return &OptimizationContext{
		op: op,
	}
}

//=============================================================================
//=== Methods
//=============================================================================

func (oc *OptimizationContext) FilterConfig() *optimization.FilterConfig {
	return oc.op.optReq.FilterConfig
}

//=============================================================================

func (oc *OptimizationContext) AlgorithmConfig() *optimization.AlgorithmConfig {
	return &oc.op.optReq.Algorithm.Config
}

//=============================================================================

func (oc *OptimizationContext) IsStopping() bool {
	return oc.op.stopping
}

//=============================================================================

func (oc *OptimizationContext) Baseline() db.TradingFilter {
	return *oc.op.optReq.Baseline
}

//=============================================================================

func (oc *OptimizationContext) RunAnalysis(filter *db.TradingFilter) float64 {
	return oc.op.runAnalysis(filter)
}

//=============================================================================

func (oc *OptimizationContext) LogInfo(message string) {
	slog.Info(message, "tsId", oc.op.ts.Id, "tsName", oc.op.ts.Name)
}

//=============================================================================
