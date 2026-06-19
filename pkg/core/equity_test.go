//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"testing"

	"golang.org/x/exp/slices"
)

//=============================================================================

var trades = []float64{ 4, 2, -3, -1, 2 }
var equity = []float64{ 4, 6,  3,  2, 4 }

//=============================================================================

func TestBuildGrossEquity(t *testing.T) {
	equ := BuildEquity(&trades)

	if !slices.Equal(*equ, equity) {
		t.Errorf("Bad gross equity calculation. Expected %v but got %v", equity, equ)
	}
}

//=============================================================================

var equity1 = []float64{ 3, 5,  4, 7, 1 }
var ddown1  = []float64{ 0, 0, -1, 0, -6 }
var maxDd1  = -6.0

//=============================================================================

func TestBuildDrawDown(t *testing.T) {
	dd1, max1 := BuildDrawDown(&equity1)

	if !slices.Equal(*dd1, ddown1) {
		t.Errorf("Bad drawdown result. Got %v", dd1)
	}

	if max1 != maxDd1 {
		t.Errorf("Bad max drawdown result. Expected %v but got %v", maxDd1, max1)
	}
}

//=============================================================================
