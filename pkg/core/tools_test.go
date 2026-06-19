//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

//=============================================================================

var days = []time.Time{
	time.Date(2024, 1 , 1, 10, 11, 12, 0, time.UTC),
	time.Date(2024, 2 ,21,  3, 11, 12, 0, time.UTC),
	time.Date(2024, 7 ,11, 14, 11, 12, 0, time.UTC),
	time.Date(2025, 1 ,23, 22, 11, 12, 0, time.UTC),
}

var xAxis = []float64{ 0, 1217, 4612, 9324 }

var yAxis1= []float64{ 125,  87,  90, 130 }
var yAxis2= []float64{ 250, -30, 120, -12 }

//=============================================================================

func TestMean(t *testing.T) {
	mean := calcMean(xAxis)

	if mean != 3788.25 {
		t.Errorf("Bad mean: Expected 3788.25 and got %v", mean)
	}

	mean = calcMean(yAxis1)

	if mean != 108 {
		t.Errorf("Bad mean: Expected 108 and got %v", mean)
	}

	mean = calcMean(yAxis2)

	if mean != 82 {
		t.Errorf("Bad mean: Expected 82 and got %v", mean)
	}
}

//=============================================================================

func TestXAxisCalculation(t *testing.T) {
	axis := calcXAxis(days)

	if !slices.Equal(axis, xAxis) {
		t.Errorf("Bad xAxis result. Got %v", axis)
	}
}

//=============================================================================

func TestLinearRegression(t *testing.T) {
	slope := LinearRegression(days, yAxis1)

	if slope < 0.00184 || slope > 0.00185 {
		t.Errorf("Bad linear regression: Expected ~0.00184 and got %v", slope)
	}

	slope = LinearRegression(days, yAxis2)

	if slope < -0.01602 || slope > -0.01601 {
		t.Errorf("Bad linear regression: Expected ~-0.01601 and got %v", slope)
	}
}

//=============================================================================
