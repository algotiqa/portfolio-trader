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

var DefRiskPerTradeOnCap   = 1.5
var DefRiskPerTradeOnEarn  = 4.0
var DefPercentageOnCapital = 50.0
var DefMoneyConversion     = McTypePercOnCapital

var SpecRiskPerTradeOnCap  = core.NewNumberParamSpec[float64]("riskPerTradeOnCap",   true, 0.1,  50, &DefRiskPerTradeOnCap)
var SpecRiskPerTradeOnEarn = core.NewNumberParamSpec[float64]("riskPerTradeOnEarn",  true, 0.1,  50, &DefRiskPerTradeOnEarn)
var SpecPercentageOnCapital= core.NewNumberParamSpec[float64]("percentageOnCapital", true,   1, 200, &DefPercentageOnCapital)

var SpecMoneyConversion    = core.NewListParamSpec[MoneyConversionType]("moneyConversion", true, MoneyConversionDomain, DefMoneyConversion)

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type MarketMoneyConfig struct {
	riskPerTradeOnCap   float64
	riskPerTradeOnEarn  float64
	moneyConversion     MoneyConversionType
	percentageOnCapital float64
}

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
	config *MarketMoneyConfig
	baseCapital  float64
}

//=============================================================================

func NewMarketMoneyModel() *MarketMoneyModel {
	return &MarketMoneyModel{
		config: &MarketMoneyConfig{
			riskPerTradeOnCap  : DefRiskPerTradeOnCap,
			riskPerTradeOnEarn : DefRiskPerTradeOnEarn,
			moneyConversion    : DefMoneyConversion,
			percentageOnCapital: DefPercentageOnCapital,
		},
	}
}

//=============================================================================

func (m *MarketMoneyModel) Name() db.ModelName {
	return db.ModelMarketMoney
}

//=============================================================================

func (m *MarketMoneyModel) Init(config map[string]any) error {
	riskOnCap ,err1 := core.MapNumber[float64](config, SpecRiskPerTradeOnCap)
	riskOnEarn,err2 := core.MapNumber[float64](config, SpecRiskPerTradeOnEarn)
	monConv,   err3 := core.MapList[MoneyConversionType](config, SpecMoneyConversion)

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}

	m.config.riskPerTradeOnCap  = *riskOnCap
	m.config.riskPerTradeOnEarn = *riskOnEarn
	m.config.moneyConversion    = monConv

	if monConv == McTypePercOnCapital {
		percOnCap, err4 := core.MapNumber[float64](config, SpecPercentageOnCapital)
		if err4 != nil {
			return err4
		}

		m.config.percentageOnCapital = *percOnCap
	}

	return nil
}

//=============================================================================

func (m *MarketMoneyModel) Config() map[string]any {
	cfg := make(map[string]any)
	cfg[SpecRiskPerTradeOnCap  .Name] = m.config.riskPerTradeOnCap
	cfg[SpecRiskPerTradeOnEarn .Name] = m.config.riskPerTradeOnEarn
	cfg[SpecMoneyConversion    .Name] = m.config.moneyConversion
	cfg[SpecPercentageOnCapital.Name] = m.config.percentageOnCapital

	return cfg
}

//=============================================================================

func (m *MarketMoneyModel) Spec() map[string]any {
	specs := make(map[string]any)

	specs[SpecRiskPerTradeOnCap  .Name] = SpecRiskPerTradeOnCap
	specs[SpecRiskPerTradeOnEarn .Name] = SpecRiskPerTradeOnEarn
	specs[SpecMoneyConversion    .Name] = SpecMoneyConversion
	specs[SpecPercentageOnCapital.Name] = SpecPercentageOnCapital

	return specs
}

//=============================================================================

func (m *MarketMoneyModel) PositionInit(ts *TradingSnapshot) {
	m.baseCapital = ts.InitialCapital
}

//=============================================================================

func (m *MarketMoneyModel) PositionFor(ts *TradingSnapshot) int {
	capAtRisk := m.baseCapital                       * m.config.riskPerTradeOnCap  / 100
	earAtRisk := (ts.CurrentCapital - m.baseCapital) * m.config.riskPerTradeOnEarn / 100

	if earAtRisk < 0 {
		//--- When we are losing money, which back to the PercentRisk model
		capAtRisk = ts.CurrentCapital * m.config.riskPerTradeOnCap  / 100
		earAtRisk = 0
	}

	if m.config.moneyConversion == McTypePercOnCapital {
		if ts.CurrentCapital > m.baseCapital * (1 + m.config.percentageOnCapital/100) {
			m.baseCapital = ts.CurrentCapital
		}
	}

	units := int( (capAtRisk + earAtRisk) / ts.RiskValue)

	return units
}

//=============================================================================
