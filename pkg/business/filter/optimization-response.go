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
	CurrStep        int64         `json:"currStep"`
	MaxSteps        int64         `json:"maxSteps"`
	StartTime       time.Time     `json:"startTime"`
	EndTime         time.Time     `json:"endTime"`
	Status          string        `json:"status"`
	Error           string        `json:"error"`
	Runs            []*Run        `json:"runs"`
	BaseValue       float64       `json:"baseValue"`
	BestValue       float64       `json:"bestValue"`
	FieldToOptimize string        `json:"fieldToOptimize"`
	Duration        int64         `json:"duration"`
	Filter         *Flags         `json:"filter"`
}

//=============================================================================

func NewOptimizationResponse(info *OptimizationInfo) *OptimizationResponse {
	or := &OptimizationResponse{}
	or.StartDate       = info.StartDate
	or.CurrStep        = info.CurrStep
	or.MaxSteps        = info.MaxSteps
	or.StartTime       = info.StartTime
	or.EndTime         = info.EndTime
	or.Status          = info.Status
	or.Runs            = info.GetRuns()
	or.BaseValue       = info.BaseValue
	or.BestValue       = info.BestValue
	or.FieldToOptimize = info.FieldToOptimize
	or.Duration        = int64(time.Now().Sub(info.StartTime).Seconds())
	or.Filter          = &info.Filter

	if info.Error != nil {
		or.Error = info.Error.Error()
	}

	return or
}

//=============================================================================
