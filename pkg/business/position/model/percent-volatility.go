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
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Model
//===
//=============================================================================

type PercentVolatilityModel struct {
	averageLength int
	maxVolatility float64
}

//=============================================================================

func NewPercentVolatilityModel() *PercentVolatilityModel {
	return &PercentVolatilityModel{
		averageLength: 20,
		maxVolatility: 1.5,
	}
}

//=============================================================================

func (m *PercentVolatilityModel) Name() db.ModelName {
	return db.ModelPercentVolatility
}

//=============================================================================

func (m *PercentVolatilityModel) Init(config map[string]any) error {
	return nil
}

//=============================================================================

func (m *PercentVolatilityModel) Config() map[string]any {
	cfg := make(map[string]any)
	return cfg
}

//=============================================================================

func (m *PercentVolatilityModel) PositionInit(ts *TradingSnapshot) {}

//=============================================================================

func (m *PercentVolatilityModel) PositionFor(ts *TradingSnapshot) int {
	return 1
}

//=============================================================================
