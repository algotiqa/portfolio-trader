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

	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================

type BarRange struct {
	MinValue float64 `json:"minValue"`
	MaxValue float64 `json:"maxValue"`
}

//=============================================================================

type Histogram struct {
	Bars     []int       `json:"bars"`
	Gaussian []float64   `json:"gaussian"`
	Ranges   []*BarRange `json:"ranges"`
	slots    int
}

//=============================================================================

func NewHistogram(data []float64) *Histogram {
	slots := int(math.Sqrt(float64(len(data))))
	if slots%2 == 0 {
		slots++
	}

	h := &Histogram{
		slots: slots,
	}

	h.Bars = make([]int, slots)
	h.Gaussian = make([]float64, slots)
	h.Ranges = make([]*BarRange, slots)

	h.initRanges(data)
	h.populate(data)
	h.createGaussian(data)

	return h
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func (h *Histogram) initRanges(data []float64) {
	minV := Min(data)
	maxV := Max(data)
	delta := (maxV - minV) / float64(h.slots)

	for i := 0; i < h.slots; i++ {
		h.Ranges[i] = &BarRange{
			MinValue: core.Trunc2d(minV),
			MaxValue: core.Trunc2d(minV + delta),
		}
		minV = minV + delta
	}
}

//=============================================================================

func (h *Histogram) populate(data []float64) {
	for _, value := range data {
		for i := 0; i < h.slots; i++ {
			barRange := h.Ranges[i]
			if barRange.MinValue <= value && value < barRange.MaxValue {
				h.Bars[i]++
				break
			} else if i == h.slots-1 {
				h.Bars[i]++
			}
		}
	}
}

//=============================================================================

func (h *Histogram) createGaussian(data []float64) {
	mean := Mean(data)
	stdDev := StdDev(data, mean)
	meanIdx := 0

	for i, br := range h.Ranges {
		medValue := (br.MinValue + br.MaxValue) / 2
		gasValue := gaussian(medValue, mean, stdDev)
		h.Gaussian[i] = gasValue

		if br.MinValue <= mean && mean < br.MaxValue {
			meanIdx = i
		}
	}

	//--- The gaussian has been calculated, now we have to scale it to match the histogram

	topVal := gaussian(mean, mean, stdDev)
	meanCount := float64(h.Bars[meanIdx])

	for i, v := range h.Gaussian {
		h.Gaussian[i] = core.Trunc2d(meanCount * v / topVal)
	}
}

//=============================================================================

func gaussian(x, mean, stdDev float64) float64 {
	v := (x - mean) / stdDev
	exponent := -v * v / 2
	return math.Exp(exponent) / (stdDev * math.Sqrt(2*math.Pi))
}

//=============================================================================
