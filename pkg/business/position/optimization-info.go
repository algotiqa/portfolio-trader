//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package position

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================

const MaxResultSize = 5000

//=============================================================================

const OptimStatusIdle      = "idle"
const OptimStatusRunning   = "running"
const OptimStatusCompleted = "completed"
const OptimStatusAborted   = "aborted"

//=============================================================================
//===
//=== OptimizationInfo
//===
//=============================================================================

type OptimizationInfo struct {
	CurrStep    uint
	TotalSteps  uint
	StartTime   time.Time
	EndTime     time.Time
	Status      string
	spec       *RunSpec
	config     *ModelConfig
	targets    *Targets
	results    *core.SortedResults[*ExecutionResult]
	bestRun    *ExecutionResult
}

//=============================================================================

func NewOptimizationInfo(spec *RunSpec, config *ModelConfig, targets *Targets, totalSteps uint, maxResultSize int) *OptimizationInfo {
	oi := &OptimizationInfo{}
	oi.CurrStep   = 0
	oi.TotalSteps = totalSteps
	oi.StartTime  = time.Now()
	oi.Status     = OptimStatusRunning
	oi.spec       = spec
	oi.config     = config
	oi.targets    = targets
	oi.results    = core.NewSortedResults[*ExecutionResult](maxResultSize, runComparator)

	return oi
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func (oi *OptimizationInfo) addResult(er *ExecutionResult, targets *Targets) {
	oi.CurrStep++

	if isAboveTargets(er, targets) {
		oi.results.Add(er)
		fv := er.FitnessValue()

		if oi.bestRun == nil || oi.bestRun.FitnessValue()  < fv {
			oi.bestRun = er
		}
	}
}

//=============================================================================

func (oi *OptimizationInfo) finish(stopped bool) {
	oi.EndTime = time.Now()

	if stopped {
		oi.Status = OptimStatusAborted
	} else {
		oi.Status = OptimStatusCompleted
	}
}

//=============================================================================
//===
//=== Run comparator
//===
//=== Notes:
//===  - in reverse order: max to min
//=============================================================================

func runComparator(a any, b any) int {
	r1 := a.(*ExecutionResult)
	r2 := b.(*ExecutionResult)
	v1 := r1.FitnessValue()
	v2 := r2.FitnessValue()

	if v1 < v2 {
		return +1
	}

	return -1
}

//=============================================================================

func isAboveTargets(er *ExecutionResult, t *Targets) bool {
	return	er.MedianRetDrawdRatio >= *t.MinRetMaxDrawdRatio
}

//=============================================================================
