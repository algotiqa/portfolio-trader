//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"errors"
	"math"
	"slices"

	"github.com/algotiqa/portfolio-trader/pkg/db"
	"golang.org/x/exp/maps"
)

//=============================================================================
//===
//=== Computes the Pearson product-moment correlation coefficient
//=== between two equally sized series. It uses the two-pass algorithm: means
//=== are computed first and deviations are accumulated afterwards. This avoids
//=== the catastrophic cancellation of the naive single-pass formula and is the
//=== right trade-off for daily return series, whose values are small and close
//=== to their mean.
//===
//=============================================================================

func CalcPearsonCorrelation(x []float64, y []float64) (float64, error) {
	if len(x) != len(y) {
		return 0, errors.New("Cannot compute correlation: series have different lengths")
	}

	n := len(x)

	if n < 2 {
		return 0, errors.New("Cannot compute correlation: at least two points are required")
	}

	var sumX, sumY float64

	for i := range n {
		if math.IsNaN(x[i]) || math.IsNaN(y[i]) || math.IsInf(x[i], 0) || math.IsInf(y[i], 0) {
			return 0, errors.New("Cannot compute correlation: series contain non finite values")
		}

		sumX += x[i]
		sumY += y[i]
	}

	meanX := sumX / float64(n)
	meanY := sumY / float64(n)

	var cov, varX, varY float64

	for i := range n {
		dx := x[i] - meanX
		dy := y[i] - meanY

		cov  += dx * dy
		varX += dx * dx
		varY += dy * dy
	}

	if varX == 0 || varY == 0 {
		return 0, errors.New("Cannot compute correlation: one of the series is constant")
	}

	return cov / (math.Sqrt(varX) * math.Sqrt(varY)), nil
}

//=============================================================================
//===
//=== CorrelationDailyReturns computes the Pearson correlation between the
//=== daily returns of two trading systems. The two series are aligned on the
//=== union of their trading dates: a date present in only one system is
//=== treated as a flat day (zero return) for the other one, because systems
//=== trading different instruments or sessions do not share the same calendar.
//===
//=============================================================================

func CalcCorrelation(a *[]db.DailyReturn, b *[]db.DailyReturn) (float64, error) {
	x, y := hydrateDailyReturns(a, b)
	return CalcPearsonCorrelation(x, y)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func hydrateDailyReturns(a *[]db.DailyReturn, b *[]db.DailyReturn) ([]float64, []float64) {
	aMap := dailyReturnMap(a)
	bMap := dailyReturnMap(b)

	aDates := maps.Keys(aMap)
	bDates := maps.Keys(bMap)

	slices.Sort(aDates)
	slices.Sort(bDates)

	var aIndex, bIndex int
	var aReturns,bReturns []float64

	for aIndex < len(aDates) && bIndex < len(bDates) {
		aDate := aDates[aIndex]
		bDate := bDates[bIndex]

		switch {
			case aDate < bDate:
				aReturns = append(aReturns, aMap[aDate])
				bReturns = append(bReturns, 0)
				aIndex++
			case aDate > bDate:
				aReturns = append(aReturns, 0)
				bReturns = append(bReturns, bMap[bDate])
				bIndex++
			default:
				aReturns = append(aReturns, aMap[aDate])
				bReturns = append(bReturns, bMap[bDate])
				aIndex++
				bIndex++
		}
	}

	//--- Add trailing values from A

	for ; aIndex < len(aDates); aIndex++ {
		aDate := aDates[aIndex]
		aReturns = append(aReturns, aMap[aDate])
		bReturns = append(bReturns, 0)
	}

	//--- Add trailing values from B

	for ; bIndex < len(bDates); bIndex++ {
		bDate := bDates[bIndex]
		bReturns = append(bReturns, bMap[bDate])
		aReturns = append(aReturns, 0)
	}

	return aReturns, bReturns
}

//=============================================================================

func dailyReturnMap(list *[]db.DailyReturn) map[int]float64 {
	m := map[int]float64{}

	if list == nil {
		return m
	}

	for _, dr := range *list {
		m[int(dr.Date)] = dr.GrossReturn
	}

	return m
}

//=============================================================================
