//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package model

//=============================================================================

import "github.com/algotiqa/portfolio-trader/pkg/db"

//=============================================================================
//===
//=== PositionModel
//===
//=============================================================================

type PositionModel interface {
	Name() db.ModelName
	Init(config map[string]any) error
	Config() map[string]any
	PositionInit(ts *TradingSnapshot)
	PositionFor(ts *TradingSnapshot) int
}

//=============================================================================

type TradingSnapshot struct {
	InitialCapital float64
	CurrentCapital float64
	RiskValue      float64
	AtrValue       float64
}

//=============================================================================

func New(model db.ModelName, config map[string]any) (PositionModel, error) {
	var m PositionModel
	switch model {
		case db.ModelFixedUnit:
			m = NewFixedUnitModel()

		case db.ModelPercentRisk:
			m = NewPercentRiskModel()

		case db.ModelPercentVolatility:
			m = NewPercentVolatilityModel()

		case db.ModelMarketMoney:
			m = NewMarketMoneyModel()

		default:
			panic("Unknown position sizing model: " + model)
	}

	return m, m.Init(config)
}

//=============================================================================
