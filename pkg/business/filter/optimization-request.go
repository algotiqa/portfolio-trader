//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import (
	"errors"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== OptimizationRequest
//===
//=============================================================================

const FieldToOptimizeNetReturn              = "netReturn"
const FieldToOptimizeAvgTrade               = "avgTrade"
const FieldToOptimizeDrawDown               = "maxDD"
const FieldToOptimizeNetReturnAvgTrade      = "netReturn*avgTrade"
const FieldToOptimizeNetReturnAvgTradeMaxDD = "netReturn*avgTrade/maxDD"

//=============================================================================

type AlgorithmSpec struct {
	Type   string                       `json:"type"`
	Config optimization.AlgorithmConfig `json:"config"`
}

//=============================================================================

type OptimizationRequest struct {
	Period          *core.SelectedPeriod       `json:"period,omitempty"`
	FieldToOptimize string                     `json:"fieldToOptimize"`
	FilterConfig    *optimization.FilterConfig `json:"filterConfig"`
	Algorithm       *AlgorithmSpec             `json:"algorithm"`
	Baseline        *model.TradingFilter       `json:"baseline"`
}

//=============================================================================

func (r *OptimizationRequest) Validate() error {
	if  r.FieldToOptimize != FieldToOptimizeNetReturn         &&
		r.FieldToOptimize != FieldToOptimizeAvgTrade          &&
		r.FieldToOptimize != FieldToOptimizeDrawDown          &&
		r.FieldToOptimize != FieldToOptimizeNetReturnAvgTrade &&
		r.FieldToOptimize != FieldToOptimizeNetReturnAvgTradeMaxDD {
		return errors.New("Invalid field to optimize: " + r.FieldToOptimize)
	}

	algoType := r.Algorithm.Type

	if algoType != algorithm.Simple && algoType != algorithm.Genetic {
		return errors.New("Invalid optimization algorithm: " + algoType)
	}

	if err := r.FilterConfig.Validate(); err != nil {
		return err
	}

	return nil
}

//=============================================================================
