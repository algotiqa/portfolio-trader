//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package stats

import (
	"math"
	"testing"
)

//=============================================================================

var prices = []float64{ 1.64, 5.85, 9.22, 3.51, -0.88, 1.07, 13.03, 9.4, 10.49, -5.08, 0, 0 }

//=============================================================================

func TestSharpeRatio(t *testing.T) {
	t.Logf("Avg: %v", Mean(prices))
	t.Logf("Std: %v", StdDev(prices, Mean(prices)))
	t.Logf("SR : %v", SharpeRatio(10, 0.05))
	t.Logf("SRA: %v", SharpeRatio(10, 0.05)*math.Sqrt(12))

	//if !slices.Equal(*equ, equity) {
	//	t.Errorf("Bad gross equity calculation. Expected %v but got %v", equity, equ)
	//}
}

//=============================================================================
