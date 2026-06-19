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
)

//=============================================================================

type Percentile struct {
	data []float64
}

//=============================================================================

func NewPercentile(data []float64) *Percentile {
	perc := slices.Clone(data)
	slices.Sort(perc)

	return &Percentile{
		data: perc,
	}
}

//=============================================================================

func (p *Percentile) Get(percentile float64) float64 {
	size := len(p.data)
	idx  := int(math.Floor(percentile * float64(size -1) / 100))
	if idx >= size {
		idx = size - 1
	}

	return p.data[idx]
}

//=============================================================================
