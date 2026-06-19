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
	"sync"
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Run
//===
//=============================================================================

type Run struct {
	Filter *db.TradingFilter `json:"filter"`

	FitnessValue float64 `json:"fitnessValue"`
	NetProfit    float64 `json:"netProfit"`
	AvgTrade     float64 `json:"avgTrade"`
	MaxDrawdown  float64 `json:"maxDrawdown"`
	random       int
}

//=============================================================================
//===
//=== OptimizationInfo
//===
//=============================================================================

const OptimStatusIdle = "idle"
const OptimStatusRunning = "running"
const OptimStatusComplete = "complete"

type OptimizationInfo struct {
	sync.RWMutex
	CurrStep  uint
	MaxSteps  uint
	StartTime time.Time
	EndTime   time.Time
	Status    string
	results   *core.SortedResults

	StartDate       *time.Time
	BaseValue       float64
	BestValue       float64
	FieldToOptimize string

	Filter struct {
		PosProfit bool
		OldVsNew  bool
		WinPerc   bool
		EquVsAvg  bool
		Trendline bool
		Drawdown  bool
	}
}

//=============================================================================

func NewOptimizationInfo(maxResultSize int, field string, fc *optimization.FilterConfig,
	steps uint, baseValue float64, startDate *time.Time) *OptimizationInfo {
	oi := &OptimizationInfo{}
	oi.CurrStep = 0
	oi.StartTime = time.Now()
	oi.Status = OptimStatusRunning
	oi.results = core.NewSortedResults(maxResultSize, runComparator)
	oi.BaseValue = baseValue
	oi.BestValue = baseValue
	oi.MaxSteps = steps
	oi.FieldToOptimize = field
	oi.StartDate = startDate

	oi.Filter.PosProfit = fc.EnablePosProfit
	oi.Filter.OldVsNew = fc.EnableOldNew
	oi.Filter.WinPerc = fc.EnableWinPerc
	oi.Filter.EquVsAvg = fc.EnableEquAvg
	oi.Filter.Trendline = fc.EnableTrendline
	oi.Filter.Drawdown = fc.EnableDrawdown

	return oi
}

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func (oi *OptimizationInfo) GetRuns() []any {
	oi.Lock()
	defer oi.Unlock()

	if oi.results != nil {
		return oi.results.ToList()
	}

	return nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func (oi *OptimizationInfo) addResult(r *Run) {
	oi.Lock()
	defer oi.Unlock()

	oi.CurrStep++
	oi.results.Add(r)

	fv := r.FitnessValue

	if oi.BestValue < fv {
		oi.BestValue = fv
	}
}

//=============================================================================

func (oi *OptimizationInfo) isStatusComplete() bool {
	oi.Lock()
	defer oi.Unlock()

	if oi.CurrStep != oi.MaxSteps {
		return false
	}

	oi.EndTime = time.Now()
	oi.Status = OptimStatusComplete

	return true
}

//=============================================================================
