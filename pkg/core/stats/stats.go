//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package stats

import (
	"math"
	"slices"
	"time"
)

//=============================================================================

func Mean[T float64|int](data []T) float64 {
	if data == nil || len(data) == 0 {
		return math.NaN()
	}

	sum := 0.0

	for _, v := range data {
		sum += float64(v)
	}

	return sum/float64(len(data))
}

//=============================================================================

func Median(data []float64) float64 {
	if data == nil || len(data) == 0 {
		return math.NaN()
	}

	aux := slices.Clone(data)
	slices.Sort(aux)
	size := len(aux)
	if size % 2 == 1 {
		return aux[size/2]
	} else {
		return (aux[size/2 -1] + aux[size/2]) / 2
	}
}

//=============================================================================

func StdDev(data []float64, mean float64) float64 {
	if data == nil || len(data) == 0 {
		return math.NaN()
	}

	sum  := 0.0
	diff := 0.0

	for _, v := range data {
		diff = v - mean
		sum += diff * diff
	}

	return math.Sqrt(sum/float64(len(data)))
}

//=============================================================================

func SharpeRatio(average, stdDev float64) float64{
	if average == math.NaN() {
		return math.NaN()
	}

	if stdDev == math.NaN() {
		return math.NaN()
	}

	if stdDev == 0 {
		return math.Inf(1)
	}

	return average/stdDev
}

//=============================================================================

func Skewness(mean, median, stdDev float64) float64 {
	if stdDev == 0 {
		return 0
	}

	return 3*(mean - median)/stdDev
}

//=============================================================================

func LinearRegression(x []time.Time, y []float64) float64 {
	xAxis := calcXAxis(x)
	xMean := Mean(xAxis)
	yMean := Mean(y)

	num := 0.0
	den := 0.0
	aux := 0.0

	for i,_ := range y {
		aux = xAxis[i] - xMean
		num += aux * (y[i] - yMean)
		den += aux * aux
	}

	return num / den
}

//=============================================================================

func Min(data []float64) float64 {
	if len(data) == 0 {
		return math.NaN()
	}

	minValue := data[0]

	for _, v := range data {
		if v < minValue {
			minValue = v
		}
	}

	return minValue
}

//=============================================================================

func Max(data []float64) float64 {
	if len(data) == 0 {
		return math.NaN()
	}

	maxValue := data[0]

	for _, v := range data {
		if v > maxValue {
			maxValue = v
		}
	}

	return maxValue
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcXAxis(time []time.Time) []float64 {
	var x []float64

	for _, t := range time {
		v := calcHours(t, time[0])
		x = append(x, v)
	}

	return x
}

//=============================================================================

func calcHours(t,t0 time.Time) float64 {
	delta := t.UnixMilli() - t0.UnixMilli()
	return float64(delta / 1000 / 3600)
}

//=============================================================================
