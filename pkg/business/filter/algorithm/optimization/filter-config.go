//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package optimization

import (
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== FilterConfig
//===
//=============================================================================

type FilterConfig struct {
	Activation   Config `json:"activation"`
	Deactivation Config `json:"deactivation"`
}

//=============================================================================

func (fc *FilterConfig) Validate() error {
	err := fc.Activation.Validate()
	if err != nil {
		return err
	}

	err = fc.Deactivation.Validate()
	if err != nil {
		return err
	}

	return nil
}

//=============================================================================

type Config struct {
	NetReturn         NetReturnOptimiz
	AverageTrade      AverageTradeOptimiz
	EquityVsAverage   EquityVsAverageOptimiz
	WinningPercentage WinningPercentageOptimiz
	MaxDrawdown       MaxDrawdownOptimiz
	TrendLine         TrendLineOptimiz
	ProfitFactor      ProfitFactorOptimiz
	SharpeRatio       SharpeRatioOptimiz
	Sqn100            Sqn100Optimiz
	ReturnMaxDDRatio  ReturnMaxDDRatioOptimiz
}

//=============================================================================

func (c *Config) Validate() error {
	if err := c.NetReturn.Validate(); err != nil {
		return err
	}

	if err := c.AverageTrade.Validate(); err != nil {
		return err
	}

	if err := c.EquityVsAverage.Validate(); err != nil {
		return err
	}

	if err := c.WinningPercentage.Validate(); err != nil {
		return err
	}

	if err := c.MaxDrawdown.Validate(); err != nil {
		return err
	}

	if err := c.TrendLine.Validate(); err != nil {
		return err
	}

	if err := c.ProfitFactor.Validate(); err != nil {
		return err
	}

	if err := c.SharpeRatio.Validate(); err != nil {
		return err
	}

	if err := c.Sqn100.Validate(); err != nil {
		return err
	}

	if err := c.ReturnMaxDDRatio.Validate(); err != nil {
		return err
	}

	return nil
}

//=============================================================================

func (c *Config) StepsCount() int64 {
	return	c.NetReturn        .StepsCount() +
			c.AverageTrade     .StepsCount()+
			c.EquityVsAverage  .StepsCount() +
			c.WinningPercentage.StepsCount() +
			c.MaxDrawdown      .StepsCount() +
			c.TrendLine        .StepsCount() +
			c.ProfitFactor     .StepsCount() +
			c.SharpeRatio      .StepsCount() +
			c.Sqn100           .StepsCount() +
			c.ReturnMaxDDRatio .StepsCount()
}

//=============================================================================

func (c *Config) AllOff() bool {
	return	!c.NetReturn.Enabled         &&
			!c.AverageTrade.Enabled      &&
			!c.EquityVsAverage.Enabled   &&
			!c.WinningPercentage.Enabled &&
			!c.MaxDrawdown.Enabled       &&
			!c.TrendLine.Enabled         &&
			!c.ProfitFactor.Enabled      &&
			!c.SharpeRatio.Enabled       &&
			!c.Sqn100.Enabled            &&
			!c.ReturnMaxDDRatio.Enabled
}

//=============================================================================

func (c *Config) WarmUp() {
	c.NetReturn.WarmUp()
	c.AverageTrade.WarmUp()
	c.EquityVsAverage.WarmUp()
	c.WinningPercentage.WarmUp()
	c.MaxDrawdown.WarmUp()
	c.TrendLine.WarmUp()
	c.ProfitFactor.WarmUp()
	c.SharpeRatio.WarmUp()
	c.Sqn100.WarmUp()
	c.ReturnMaxDDRatio.WarmUp()
}

//=============================================================================
//===
//=== Base optimization config
//===
//=============================================================================

type Optimization struct {
	Enabled bool                            `json:"enabled"`
	Length  core.FieldOptimization[int]     `json:"length"`
	Value   core.FieldOptimization[float64] `json:"value"`
}

//=============================================================================

func (c *Optimization) Validate() error {
	if c.Enabled {
		if err := c.Length.Validate(criteria.MinTradesLength, criteria.MaxTradesLength); err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func (c *Optimization) StepsCount() int64 {
	if c.Enabled {
		return c.Length.StepsCount() * c.Value.StepsCount()
	}

	return 0
}

//=============================================================================

func (c *Optimization) WarmUp() {
	c.Length.Steps()
	c.Value.Steps()
}

//=============================================================================
//===
//=== Net return
//===
//=============================================================================

type NetReturnOptimiz struct {
	Optimization
}

//=============================================================================

func (c *NetReturnOptimiz) Validate() error {
	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== AverageTrade
//===
//=============================================================================

type AverageTradeOptimiz struct {
	Optimization
}

//=============================================================================

func (c *AverageTradeOptimiz) Validate() error {
	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Equity vs average
//===
//=============================================================================

type EquityVsAverageOptimiz struct {
	Optimization
}

//=============================================================================

func (c *EquityVsAverageOptimiz) Validate() error {
	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Winning percentage
//===
//=============================================================================

type WinningPercentageOptimiz struct {
	Optimization
}

//=============================================================================

func (c *WinningPercentageOptimiz) Validate() error {
	if c.Enabled {
		if err := c.Value.Validate(criteria.MinWinningPercentage, criteria.MaxWinningPercentage); err != nil {
			return err
		}
	}

	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Max drawdown
//===
//=============================================================================

type MaxDrawdownOptimiz struct {
	Optimization
}

//=============================================================================

func (c *MaxDrawdownOptimiz) Validate() error {
	if c.Enabled {
		if err := c.Value.Validate(criteria.MinDrawdownValue, criteria.MaxDrawdownValue); err != nil {
			return err
		}
	}

	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Trendline
//===
//=============================================================================

type TrendLineOptimiz struct {
	Optimization
}

//=============================================================================

func (c *TrendLineOptimiz) Validate() error {
	if c.Enabled {
		if err := c.Value.Validate(criteria.MinTrendlineValue, criteria.MaxTrendlineValue); err != nil {
			return err
		}
	}

	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Profit factor
//===
//=============================================================================

type ProfitFactorOptimiz struct {
	Optimization
}

//=============================================================================

func (c *ProfitFactorOptimiz) Validate() error {
	if c.Enabled {
		if err := c.Value.Validate(criteria.MinProfitFactor, criteria.MaxProfitFactor); err != nil {
			return err
		}
	}

	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Sharpe ratio
//===
//=============================================================================

type SharpeRatioOptimiz struct {
	Optimization
}

//=============================================================================

func (c *SharpeRatioOptimiz) Validate() error {
	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Sqn 100
//===
//=============================================================================

type Sqn100Optimiz struct {
	Optimization
}

//=============================================================================

func (c *Sqn100Optimiz) Validate() error {
	return c.Optimization.Validate()
}

//=============================================================================
//===
//=== Return / max drawdown
//===
//=============================================================================

type ReturnMaxDDRatioOptimiz struct {
	Optimization
}

//=============================================================================

func (c *ReturnMaxDDRatioOptimiz) Validate() error {
	if c.Enabled {
		if err := c.Value.Validate(criteria.MinReturnMaxDDValue, criteria.MaxReturnMaxDDValue); err != nil {
			return err
		}
	}

	return c.Optimization.Validate()
}

//=============================================================================
