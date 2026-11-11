//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package model

import "github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"

//=============================================================================
//===
//=== Trading filter
//===
//=============================================================================

type TradingFilter struct {
	Activation    FilterConfig  `json:"activation"`
	Deactivation  FilterConfig  `json:"deactivation"`
}

//=============================================================================

func NewTradingFilter() *TradingFilter {
	return &TradingFilter{
		Activation  : FilterConfig{
			CombinationMode: CombinationModeAll,
		},
		Deactivation: FilterConfig{
			CombinationMode: CombinationModeOne,
		},
	}
}

//=============================================================================

type CombinationMode string

//-----------------------------------------------------------------------------

const (
	CombinationModeAll  CombinationMode = "A" //--- This is an AND
	CombinationModeMost CombinationMode = "M" //--- Majority of conditions are satisfied (50% +1)
	CombinationModeOne  CombinationMode = "O" //--- This is an OR
)

//-----------------------------------------------------------------------------

type FilterConfig struct {
	CombinationMode   CombinationMode                  `json:"combinationMode"`
	NetReturn         criteria.NetReturnConfig         `json:"netReturn"`
	AverageTrade      criteria.AverageTradeConfig      `json:"averageTrade"`
	EquityVsAverage   criteria.EquityVsAverageConfig   `json:"equityVsAverage"`
	WinningPercentage criteria.WinningPercentageConfig `json:"winningPercentage"`
	MaxDrawdown       criteria.MaxDrawdownConfig       `json:"maxDrawdown"`
	TrendLine         criteria.TrendLineConfig         `json:"trendLine"`
	ProfitFactor      criteria.ProfitFactorConfig      `json:"profitFactor"`
	SharpeRatio       criteria.SharpeRatioConfig       `json:"sharpeRatio"`
	Sqn100            criteria.Sqn100Config            `json:"sqn100"`
	ReturnMaxDDRatio  criteria.ReturnMaxDDRatioConfig  `json:"returnMaxDDRatio"`
}

//=============================================================================


