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

type PositionModel interface {
	Calc()
}

//=============================================================================

func New(model, config string) (PositionModel, error) {
	switch model {
		case FixedUnit:
			return newFixedUnitModel(config)

		case PercentRisk:
			return newPercentRiskModel(config)

		case PercentVolatility:
			return newPercentVolatilityModel(config)

		case MarketMoney:
			return newMarketMoneyModel(config)

		default:
		panic("Unknown position sizing model: " + model)
	}
}

//=============================================================================

func NewFixedUnitDefaultConfig() *FixedUnitConfig {
	return &FixedUnitConfig{
		Units: 1,
	}
}

//=============================================================================
