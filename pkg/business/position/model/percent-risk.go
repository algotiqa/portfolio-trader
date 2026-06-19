//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package model

import "encoding/json"

//=============================================================================

const PercentRisk = "PR"

//=============================================================================

type RpuType string

//-----------------------------------------------------------------------------

const (
	RpuStopLoss   RpuType = "stopLoss"
	RpuMaxLoss    RpuType = "maxLoss"
	RpuAvgLoss    RpuType = "avgLoss"
	RpuFixedValue RpuType = "fixedValue"
)

//-----------------------------------------------------------------------------

type PercentRiskConfig struct {
	RiskPerTrade float64 `json:"riskPerTrade"`
	RiskPerUnit  RpuType `json:"riskPerUnit"`
	RiskValue    float64 `json:"riskValue"`
}

//=============================================================================

type PercentRiskModel struct {
	config *PercentRiskConfig
}

//=============================================================================

func newPercentRiskModel(config string) (*PercentRiskModel,error) {
	c := &PercentRiskConfig{}
	err := json.Unmarshal([]byte(config), c)
	if err != nil {
		return nil, err
	}

	return &PercentRiskModel{
		config : c,
	},nil
}

//=============================================================================

func (fm *PercentRiskModel) Calc() {}

//=============================================================================
