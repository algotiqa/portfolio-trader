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

type PercentRiskModel struct {
	riskPerTrade float64
}

//=============================================================================

func NewPercentRiskModel() *PercentRiskModel {
	return &PercentRiskModel{
		riskPerTrade : 1.75,
	}
}

//=============================================================================

func (m *PercentRiskModel) Name() db.ModelName {
	return db.ModelPercentRisk
}

//=============================================================================

func (m *PercentRiskModel) Init(config map[string]any) error {
	return nil
}

//=============================================================================

func (m *PercentRiskModel) Config() map[string]any {
	cfg := make(map[string]any)
	return cfg
}

//=============================================================================

func (m *PercentRiskModel) PositionInit(ts *TradingSnapshot) {}

//=============================================================================

func (m *PercentRiskModel) PositionFor(ts *TradingSnapshot) int {
	capAtRisk := ts.CurrentCapital * m.riskPerTrade / 100
	units     := int(capAtRisk / ts.RiskValue)

	return units
}

//=============================================================================
