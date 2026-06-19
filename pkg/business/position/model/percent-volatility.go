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

const PercentVolatility = "PV"

//=============================================================================

type PercentVolatilityConfig struct {
	AverageLength int     `json:"averageLength"`
	MaxVolatility float64 `json:"maxVolatility"`
}

//=============================================================================

type PercentVolatilityModel struct {
	config *PercentVolatilityConfig
}

//=============================================================================

func newPercentVolatilityModel(config string) (*PercentVolatilityModel,error) {
	c := &PercentVolatilityConfig{}
	err := json.Unmarshal([]byte(config), c)
	if err != nil {
		return nil, err
	}

	return &PercentVolatilityModel{
		config: c,
	}, nil
}

//=============================================================================

func (fm *PercentVolatilityModel) Calc() {}

//=============================================================================
