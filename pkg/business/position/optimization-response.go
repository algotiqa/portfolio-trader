//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package position

import "time"

//=============================================================================

type OptimizationResponse struct{
	CurrStep    uint               `json:"currStep"`
	TotalSteps  uint               `json:"totalSteps"`
	Duration    int64              `json:"duration"`
	StartTime   time.Time          `json:"startTime"`
	EndTime     time.Time          `json:"endTime"`
	Status      string             `json:"status"`
	Results     []*ExecutionResult `json:"results"`
	BestRun     *ExecutionResult   `json:"bestRun"`
	Models struct {
		FixedUnit    bool `json:"fixedUnit"`
		PercentRisk  bool `json:"percentRisk"`
		PercentVol   bool `json:"percentVol"`
		MarketMoney  bool `json:"marketMoney"`
	} `json:"models"`
	Targets struct {
		DesReturnOnAccount  float64 `json:"desReturnOnAccount"`
		MinRetMaxDrawdRatio float64 `json:"minRetMaxDrawdRatio"`
		MaxTolDrawdPerc     float64 `json:"maxTolDrawdPerc"`
	} `json:"targets"`
}

//=============================================================================

func NewOptimizationResponse(info *OptimizationInfo) *OptimizationResponse {
	or := &OptimizationResponse{}

	or.CurrStep   = info.CurrStep
	or.TotalSteps = info.TotalSteps
	or.Duration   = int64(time.Now().Sub(info.StartTime).Seconds())
	or.StartTime  = info.StartTime
	or.EndTime    = info.EndTime
	or.Status     = info.Status

	if info.Status != OptimStatusIdle {
		or.Models.FixedUnit   = info.config.EnableFixedUnit
		or.Models.PercentRisk = info.config.EnablePercentRisk
		or.Models.PercentVol  = info.config.EnablePercentVol
		or.Models.MarketMoney = info.config.EnableMarketMoney

		or.Targets.DesReturnOnAccount  = *info.targets.DesReturnOnAccount
		or.Targets.MaxTolDrawdPerc     = info.spec.MaxTolDrawdPerc
		or.Targets.MinRetMaxDrawdRatio = *info.targets.MinRetMaxDrawdRatio

		or.BestRun = info.bestRun

		if info.Status != OptimStatusRunning {
			//--- Status is either completed or aborted
			or.Results = info.results.ToList()
		}
	}

	return or
}

//=============================================================================
