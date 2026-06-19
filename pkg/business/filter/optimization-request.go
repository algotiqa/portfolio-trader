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
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== OptimizationRequest
//===
//=============================================================================

const FieldToOptimizeNetProfit = "netProfit"
const FieldToOptimizeAvgTrade = "avgTrade"
const FieldToOptimizeDrawDown = "maxDD"
const FieldToOptimizeNetProfitAvgTrade = "netProfit*avgTrade"
const FieldToOptimizeNetProfitAvgTradeMaxDD = "netProfit*avgTrade/maxDD"

//=============================================================================

type AlgorithmSpec struct {
	Type   string                       `json:"type"`
	Config optimization.AlgorithmConfig `json:"config"`
}

//=============================================================================

type OptimizationRequest struct {
	StartDate       *time.Time                 `json:"startDate,omitempty"`
	FieldToOptimize string                     `json:"fieldToOptimize"`
	FilterConfig    *optimization.FilterConfig `json:"filterConfig"`
	Algorithm       *AlgorithmSpec             `json:"algorithm"`
	Baseline        *db.TradingFilter          `json:"baseline"`
}

//=============================================================================

func (r *OptimizationRequest) Validate() error {
	if r.FieldToOptimize != FieldToOptimizeNetProfit &&
		r.FieldToOptimize != FieldToOptimizeAvgTrade &&
		r.FieldToOptimize != FieldToOptimizeDrawDown &&
		r.FieldToOptimize != FieldToOptimizeNetProfitAvgTrade &&
		r.FieldToOptimize != FieldToOptimizeNetProfitAvgTradeMaxDD {
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
