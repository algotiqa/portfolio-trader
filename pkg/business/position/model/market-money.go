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

const MarketMoney = "MM"

//=============================================================================

type MoneyConversionType string

//-----------------------------------------------------------------------------

const (
	McTypePercOnCapital MoneyConversionType = "percOnCapital"
)

//-----------------------------------------------------------------------------

type MarketMoneyConfig struct {
	RiskPerTradeOnCap   float64             `json:"riskPerTradeOnCap"`
	RiskPerTradeOnEarn  float64             `json:"riskPerTradeOnEarn"`
	MoneyConversion     MoneyConversionType `json:"moneyConversion"`
	PercentageOnCapital float64             `json:"percentageOnCapital"`
}

//=============================================================================

type MarketMoneyModel struct {
	config *MarketMoneyConfig
}

//=============================================================================

func newMarketMoneyModel(config string) (*MarketMoneyModel,error) {
	c := &MarketMoneyConfig{}
	err := json.Unmarshal([]byte(config), c)
	if err != nil {
		return nil, err
	}

	return &MarketMoneyModel{
		config : c,
	},nil
}

//=============================================================================

func (fm *MarketMoneyModel) Calc() {}

//=============================================================================
