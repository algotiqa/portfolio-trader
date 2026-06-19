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
	"time"
)

//=============================================================================
//===
//=== OptimizationResponse
//===
//=============================================================================

type OptimizationResponse struct {
	StartDate       *time.Time    `json:"startDate"`
	CurrStep        uint          `json:"currStep"`
	MaxSteps        uint          `json:"maxSteps"`
	StartTime       time.Time     `json:"startTime"`
	EndTime         time.Time     `json:"endTime"`
	Status          string        `json:"status"`
	Runs            []any         `json:"runs"`
	BaseValue       float64       `json:"baseValue"`
	BestValue       float64       `json:"bestValue"`
	FieldToOptimize string        `json:"fieldToOptimize"`
	Duration        int64         `json:"duration"`
	Filter struct {
		PosProfit bool `json:"posProfit"`
		OldVsNew  bool `json:"oldVsNew"`
		WinPerc   bool `json:"winPerc"`
		EquVsAvg  bool `json:"equVsAvg"`
		Trendline bool `json:"trendline"`
		Drawdown  bool `json:"drawdown"`
	} `json:"filter"`
}

//=============================================================================

func NewOptimizationResponse(info *OptimizationInfo) *OptimizationResponse {
	or := &OptimizationResponse{}
	or.StartDate = info.StartDate
	or.CurrStep  = info.CurrStep
	or.MaxSteps  = info.MaxSteps
	or.StartTime = info.StartTime
	or.EndTime   = info.EndTime
	or.Status    = info.Status
	or.BaseValue = info.BaseValue
	or.BestValue = info.BestValue

	or.FieldToOptimize   = info.FieldToOptimize
	or.Filter.PosProfit = info.Filter.PosProfit
	or.Filter.OldVsNew  = info.Filter.OldVsNew
	or.Filter.WinPerc   = info.Filter.WinPerc
	or.Filter.EquVsAvg  = info.Filter.EquVsAvg
	or.Filter.Trendline = info.Filter.Trendline
	or.Filter.Drawdown  = info.Filter.Drawdown

	or.Runs     = info.GetRuns()
	or.Duration = int64(time.Now().Sub(info.StartTime).Seconds())

	return or
}

//=============================================================================
