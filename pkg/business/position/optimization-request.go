//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package position

import (
	"errors"

	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== Specs
//===
//=============================================================================

//--- RunConfig

var SpecRunsPerSimul    = core.NewNumberParamSpec[int]( "runsPerSimul",    true, 1, 50000, nil)
var SpecProjectionYears = core.NewNumberParamSpec[int]( "projectionYears", true, 1,     5, nil)
var SpecTradesPerYear   = core.NewNumberParamSpec[int]( "tradesPerYear",  false, 1,  1000, nil)

//--- Targets

var SpecDesReturnOnAccount = core.NewNumberParamSpec[float64]("desReturnOnAccount",  true, 0, 10000, nil)
var SpecMinRetMaxDrawdRatio= core.NewNumberParamSpec[float64]("minRetMaxDrawdRatio", true, 0,  1000, nil)

//=============================================================================
//===
//=== OptimizationRequest
//===
//=============================================================================

type OptimizationRequest struct{
	Params      *Parameters     `json:"params"`
	RunConfig   *RunConfig      `json:"runConfig"`
	Targets     *Targets        `json:"targets"`
	ModelConfig *ModelConfig    `json:"modelConfig"`
}

//=============================================================================

func (o *OptimizationRequest) Validate() error {
	if o.Params == nil {
		return errors.New("params are required")
	}

	err := o.Params.Validate()
	if err != nil {
		return err
	}

	if o.RunConfig == nil {
		return errors.New("runConfig is required")
	}

	err = o.RunConfig.Validate()
	if err != nil {
		return err
	}

	if o.Targets == nil {
		return errors.New("targets are required")
	}

	err = o.Targets.Validate()
	if err != nil {
		return err
	}

	if o.ModelConfig == nil {
		return errors.New("modelConfig is required")
	}

	return o.ModelConfig.Validate()
}

//=============================================================================
//===
//=== RunConfig
//===
//=============================================================================

type RunConfig struct {
	Period            core.SelectedPeriod `json:"period"`
	RunsPerSimul     *int                 `json:"runsPerSimul"`
	ProjectionYears  *int                 `json:"projectionYears"`
	TradesPerYear    *int                 `json:"tradesPerYear"`
}

//=============================================================================

func (c *RunConfig) Validate() error {
	var err error

	c.RunsPerSimul,err = SpecRunsPerSimul.Validate(c.RunsPerSimul)
	if err != nil {
		return err
	}

	c.ProjectionYears,err = SpecProjectionYears.Validate(c.ProjectionYears)
	if err != nil {
		return err
	}

	c.TradesPerYear,err = SpecTradesPerYear.Validate(c.TradesPerYear)
	return err
}

//=============================================================================
//===
//=== Targets
//===
//=============================================================================

type Targets struct {
	DesReturnOnAccount  *float64 `json:"desReturnOnAccount"`
	MinRetMaxDrawdRatio *float64 `json:"minRetMaxDrawdRatio"`
}

//=============================================================================

func (t *Targets) Validate() error {
	var err error

	t.DesReturnOnAccount,err = SpecDesReturnOnAccount.Validate(t.DesReturnOnAccount)
	if err != nil {
		return err
	}

	t.MinRetMaxDrawdRatio,err = SpecMinRetMaxDrawdRatio.Validate(t.MinRetMaxDrawdRatio)
	return err
}

//=============================================================================
//===
//=== ModelConfig
//===
//=============================================================================

type ModelConfig struct {
	EnableFixedUnit   bool   `json:"enableFixedUnit"`
	EnablePercentRisk bool   `json:"enablePercentRisk"`
	EnablePercentVol  bool   `json:"enablePercentVol"`
	EnableMarketMoney bool   `json:"enableMarketMoney"`

	Units                core.FieldOptimization[int]     `json:"units"`
	RiskPerTrade         core.FieldOptimization[float64] `json:"riskPerTrade"`
	AverageLength        core.FieldOptimization[int]     `json:"averageLength"`
	MaxVolatility        core.FieldOptimization[float64] `json:"maxVolatility"`
	RiskPerTradeOnCap    core.FieldOptimization[float64] `json:"riskPerTradeOnCap"`
	RiskPerTradeOnEarn   core.FieldOptimization[float64] `json:"riskPerTradeOnEarn"`
	PercentageOnCapital  core.FieldOptimization[float64] `json:"percentageOnCapital"`
}

//=============================================================================

func (c *ModelConfig) Validate() error {
	if c.EnableFixedUnit {
		if err := c.Units.ValidateWithSpec(model.SpecUnits); err != nil {
			return err
		}
	}

	if c.EnablePercentRisk {
		if err := c.RiskPerTrade.ValidateWithSpec(model.SpecRiskPerTrade); err != nil {
			return err
		}
	}

	if c.EnablePercentVol {
		if err := c.AverageLength.ValidateWithSpec(model.SpecAverageLength); err != nil {
			return err
		}

		if err := c.MaxVolatility.ValidateWithSpec(model.SpecMaxVolatility); err != nil {
			return err
		}
	}

	if c.EnableMarketMoney {
		if err := c.RiskPerTradeOnCap.ValidateWithSpec(model.SpecRiskPerTradeOnCap); err != nil {
			return err
		}

		if err := c.RiskPerTradeOnEarn.ValidateWithSpec(model.SpecRiskPerTradeOnEarn); err != nil {
			return err
		}

		if err := c.PercentageOnCapital.ValidateWithSpec(model.SpecPercentageOnCapital); err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================
