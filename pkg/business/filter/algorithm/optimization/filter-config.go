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
	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== FilterConfig
//===
//=============================================================================

const MaxTradesLength     = 300
const MaxTrendlineValue   = 200
const MaxOldNewPercentage = 200
const MaxWinningPercentage= 100
const MaxDrawdown         = 50000

//=============================================================================

type FilterConfig struct {
	EnablePosProfit bool   `json:"enablePosProfit"`
	EnableOldNew    bool   `json:"enableOldNew"`
	EnableWinPerc   bool   `json:"enableWinPerc"`
	EnableEquAvg    bool   `json:"enableEquAvg"`
	EnableTrendline bool   `json:"enableTrendline"`
	EnableDrawdown  bool   `json:"enableDrawdown"`

	PosProLen      core.FieldOptimization[int]  `json:"posProLen"`
	OldNewOldLen   core.FieldOptimization[int]  `json:"oldNewOldLen"`
	OldNewNewLen   core.FieldOptimization[int]  `json:"oldNewNewLen"`
	OldNewOldPerc  core.FieldOptimization[int]  `json:"oldNewOldPerc"`
	WinPercLen     core.FieldOptimization[int]  `json:"winPercLen"`
	WinPercPerc    core.FieldOptimization[int]  `json:"winPercPerc"`
	EquAvgLen      core.FieldOptimization[int]  `json:"equAvgLen"`
	TrendlineLen   core.FieldOptimization[int]  `json:"trendlineLen"`
	TrendlineValue core.FieldOptimization[int]  `json:"trendlineValue"`
	DrawdownMin    core.FieldOptimization[int]  `json:"drawdownMin"`
	DrawdownMax    core.FieldOptimization[int]  `json:"drawdownMax"`
}

//=============================================================================

func (fc *FilterConfig) Validate() error {
	if err := fc.PosProLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.OldNewOldLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.OldNewNewLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.OldNewOldPerc.Validate(1, MaxOldNewPercentage); err != nil {
		return err
	}

	if err := fc.WinPercLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.WinPercPerc.Validate(1, MaxWinningPercentage); err != nil {
		return err
	}

	if err := fc.EquAvgLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.TrendlineLen.Validate(1, MaxTradesLength); err != nil {
		return err
	}

	if err := fc.TrendlineValue.Validate(1, MaxTrendlineValue); err != nil {
		return err
	}

	if err := fc.DrawdownMin.Validate(1, MaxDrawdown); err != nil {
		return err
	}

	if err := fc.DrawdownMax.Validate(1, MaxDrawdown); err != nil {
		return err
	}

	return nil
}

//=============================================================================
