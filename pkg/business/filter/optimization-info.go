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
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== Run
//===
//=============================================================================

type Run struct {
	Filter       *model.TradingFilter `json:"filter"`
	FitnessValue  float64             `json:"fitnessValue"`
	NetReturn     float64             `json:"netReturn"`
	AvgTrade      float64             `json:"avgTrade"`
	MaxDrawdown   float64             `json:"maxDrawdown"`
	random        int
}

//=============================================================================
//===
//=== OptimizationInfo
//===
//=============================================================================

const OptimStatusIdle     = "idle"
const OptimStatusRunning  = "running"
const OptimStatusComplete = "complete"
const OptimStatusError    = "error"

//-----------------------------------------------------------------------------

type OptimizationInfo struct {
	sync.RWMutex
	CurrStep  int64
	MaxSteps  int64
	StartTime time.Time
	EndTime   time.Time
	Status    string
	Error     error
	results   *core.SortedResults[*Run]

	StartDate       *time.Time
	BaseValue       float64
	BestValue       float64
	FieldToOptimize string
	Filter          Flags
}

//=============================================================================

type Flags struct {
	Activation   EnableSet   `json:"activation"`
	Deactivation EnableSet   `json:"deactivation"`
}

//=============================================================================

func (f *Flags) fromConfig(fc *optimization.FilterConfig) {
	f.Activation  .fromConfig(&fc.Activation)
	f.Deactivation.fromConfig(&fc.Activation)
}

//=============================================================================

type EnableSet struct {
	NetReturn         bool  `json:"netReturn"`
	AverageTrade      bool  `json:"averageTrade"`
	EquityVsAverage   bool  `json:"equityVsAverage"`
	WinningPercentage bool  `json:"winningPercentage"`
	MaxDrawdown       bool  `json:"maxDrawdown"`
	TrendLine         bool  `json:"trendLine"`
	ProfitFactor      bool  `json:"profitFactor"`
	SharpeRatio       bool  `json:"sharpeRatio"`
	Sqn100            bool  `json:"sqn100"`
	ReturnMaxDDRatio  bool  `json:"returnMaxDDRatio"`
}

//=============================================================================

func (s *EnableSet) fromConfig(c *optimization.Config) {
	s.NetReturn         = c.NetReturn        .Enabled
	s.AverageTrade      = c.AverageTrade     .Enabled
	s.EquityVsAverage   = c.EquityVsAverage  .Enabled
	s.WinningPercentage = c.WinningPercentage.Enabled
	s.MaxDrawdown       = c.MaxDrawdown      .Enabled
	s.TrendLine         = c.TrendLine        .Enabled
	s.ProfitFactor      = c.ProfitFactor     .Enabled
	s.SharpeRatio       = c.SharpeRatio      .Enabled
	s.Sqn100            = c.Sqn100           .Enabled
	s.ReturnMaxDDRatio  = c.ReturnMaxDDRatio .Enabled
}

//=============================================================================

func NewOptimizationInfo(maxResultSize int, field string, fc *optimization.FilterConfig, steps int64,
						 baseValue float64, startDate *time.Time) *OptimizationInfo {
	oi := &OptimizationInfo{}
	oi.CurrStep        = 0
	oi.MaxSteps        = steps
	oi.StartTime       = time.Now()
	oi.Status          = OptimStatusRunning
	oi.results         = core.NewSortedResults[*Run](maxResultSize, runComparator)
	oi.BaseValue       = baseValue
	oi.BestValue       = baseValue
	oi.FieldToOptimize = field
	oi.StartDate       = startDate

	oi.Filter.fromConfig(fc)

	return oi
}

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func (oi *OptimizationInfo) GetRuns() []*Run {
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
