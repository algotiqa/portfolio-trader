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
	"math"

	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== Run comparator
//===
//=== Notes:
//===  - in reverse order: max to min
//=============================================================================

func runComparator(a any, b any) int {
	r1 := a.(*Run)
	r2 := b.(*Run)
	v1 := r1.FitnessValue
	v2 := r2.FitnessValue

	if v1 < v2 {
		return +1
	}
	if v1 > v2 {
		return -1
	}

	if r1.random == r2.random {
		return 0
	}

	if r1.random < r2.random {
		return 1
	}

	return -1
}

//=============================================================================
//===
//=== Fitness functions
//===
//=============================================================================

type FitnessFunction func(r *Run) float64

//=============================================================================

func GetFitnessFunction(field string) FitnessFunction {
	switch field {
	case FieldToOptimizeNetProfit:
		return ffNetProfit

	case FieldToOptimizeAvgTrade:
		return ffAvgTrade

	case FieldToOptimizeDrawDown:
		return ffMaxDrawdown

	case FieldToOptimizeNetProfitAvgTrade:
		return ffNetProfitAvgTrade

	case FieldToOptimizeNetProfitAvgTradeMaxDD:
		return ffNetProfitAvgTradeMaxDD

	default:
		panic("Unknown field to optimize: " + field)
	}
}

//=============================================================================

func ffNetProfit(r *Run) float64 {
	return r.NetProfit
}

//=============================================================================

func ffAvgTrade(r *Run) float64 {
	return r.AvgTrade
}

//=============================================================================

func ffMaxDrawdown(r *Run) float64 {
	return r.MaxDrawdown
}

//=============================================================================

func ffNetProfitAvgTrade(r *Run) float64 {
	fv := r.NetProfit * r.AvgTrade

	//--- We have to push down results where both netProfit and avgTrade are negative
	//--- because their product is positive

	if r.NetProfit < 0 && r.AvgTrade < 0 {
		fv *= -1
	}

	if fv <= -1000 || fv >= 1000 {
		fv = math.Round(fv)
	}

	return fv
}

//=============================================================================

func ffNetProfitAvgTradeMaxDD(r *Run) float64 {
	fv := r.NetProfit * r.AvgTrade
	dd := math.Abs(r.MaxDrawdown)

	//--- We have to push down results where both netProfit and avgTrade are negative
	//--- because their product is positive

	if r.NetProfit < 0 && r.AvgTrade < 0 {
		fv *= -1
	}

	if dd == 0 {
		dd = 1
	}

	fv = core.Trunc2d(fv / dd)

	if fv <= -1000 || fv >= 1000 {
		fv = math.Round(fv)
	}

	return fv
}

//=============================================================================
