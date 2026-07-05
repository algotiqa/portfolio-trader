//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package model

import (
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Specs
//===
//=============================================================================

var DefAverageLength = 20
var DefMaxVolatility = 1.5

var SpecAverageLength = core.NewNumberParamSpec[int    ]( "averageLength", true, 1, 1000, &DefAverageLength)
var SpecMaxVolatility = core.NewNumberParamSpec[float64]( "maxVolatility", true, 0.1, 50, &DefMaxVolatility)

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type PercentVolatilityConfig struct {
	averageLength int
	maxVolatility float64
}

//=============================================================================
//===
//=== Model
//===
//=============================================================================

type PercentVolatilityModel struct {
	config *PercentVolatilityConfig
}

//=============================================================================

func NewPercentVolatilityModel() *PercentVolatilityModel {
	return &PercentVolatilityModel{
		config: &PercentVolatilityConfig{
			averageLength: DefAverageLength,
			maxVolatility: DefMaxVolatility,
		},
	}
}

//=============================================================================

func (m *PercentVolatilityModel) Name() db.ModelName {
	return db.ModelPercentVolatility
}

//=============================================================================

func (m *PercentVolatilityModel) Init(config map[string]any) error {
	avgLen,err1 := core.MapNumber[int    ](config, SpecAverageLength)
	maxVol,err2 := core.MapNumber[float64](config, SpecMaxVolatility)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	m.config.averageLength = *avgLen
	m.config.maxVolatility = *maxVol
	return nil
}

//=============================================================================

func (m *PercentVolatilityModel) Config() map[string]any {
	cfg := make(map[string]any)
	cfg[SpecAverageLength.Name] = m.config.averageLength
	cfg[SpecMaxVolatility.Name] = m.config.maxVolatility

	return cfg
}

//=============================================================================

func (m *PercentVolatilityModel) PositionInit(ts *TradingSnapshot) {}

//=============================================================================

func (m *PercentVolatilityModel) PositionFor(ts *TradingSnapshot) int {
	capAtRisk := ts.CurrentCapital * m.config.maxVolatility / 100
	units     := int(capAtRisk / ts.AtrValue)

	return units
}

//=============================================================================
