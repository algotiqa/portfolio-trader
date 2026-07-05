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

var DefRiskPerTrade = 1.5

var SpecRiskPerTrade = core.NewNumberParamSpec[float64]( "riskPerTrade", true, 0.1, 50, &DefRiskPerTrade)

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type PercentRiskConfig struct {
	riskPerTrade float64
}

//=============================================================================
//===
//=== Model
//===
//=============================================================================

type PercentRiskModel struct {
	config *PercentRiskConfig
}

//=============================================================================

func NewPercentRiskModel() *PercentRiskModel {
	return &PercentRiskModel{
		config: &PercentRiskConfig{
			riskPerTrade : DefRiskPerTrade,
		},
	}
}

//=============================================================================

func (m *PercentRiskModel) Name() db.ModelName {
	return db.ModelPercentRisk
}

//=============================================================================

func (m *PercentRiskModel) Init(config map[string]any) error {
	risk,err := core.MapNumber[float64](config, SpecRiskPerTrade)
	if err != nil {
		return err
	}

	m.config.riskPerTrade = *risk
	return nil
}

//=============================================================================

func (m *PercentRiskModel) Config() map[string]any {
	cfg := make(map[string]any)
	cfg[SpecRiskPerTrade.Name] = m.config.riskPerTrade

	return cfg
}

//=============================================================================

func (m *PercentRiskModel) PositionInit(ts *TradingSnapshot) {}

//=============================================================================

func (m *PercentRiskModel) PositionFor(ts *TradingSnapshot) int {
	capAtRisk := ts.CurrentCapital * m.config.riskPerTrade / 100
	units     := int(capAtRisk / ts.RiskValue)

	return units
}

//=============================================================================
