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

type MoneyConversionType string

//-----------------------------------------------------------------------------

const (
	McTypePercOnCapital MoneyConversionType = "percOnCapital"
)

//-----------------------------------------------------------------------------

var MoneyConversionDomain = []MoneyConversionType{
	McTypePercOnCapital,
}

//-----------------------------------------------------------------------------

type MarketMoneyModel struct {
	riskPerTradeOnCap   float64
	riskPerTradeOnEarn  float64
	moneyConversion     MoneyConversionType
	percentageOnCapital int
}

//=============================================================================

func NewMarketMoneyModel() *MarketMoneyModel {
	return &MarketMoneyModel{
		riskPerTradeOnCap  : 1.5,
		riskPerTradeOnEarn : 4.0,
		moneyConversion    : McTypePercOnCapital,
		percentageOnCapital: 50.0,
	}
}

//=============================================================================

func (m *MarketMoneyModel) Name() db.ModelName {
	return db.ModelMarketMoney
}

//=============================================================================

func (m *MarketMoneyModel) Init(config map[string]any) error {
	return nil
}

//=============================================================================

func (m *MarketMoneyModel) Config() map[string]any {
	cfg := make(map[string]any)
	return cfg
}

//=============================================================================

func (m *MarketMoneyModel) PositionInit(ts *TradingSnapshot) {}

//=============================================================================

func (m *MarketMoneyModel) PositionFor(ts *TradingSnapshot) int {
	return 1
}

//=============================================================================
