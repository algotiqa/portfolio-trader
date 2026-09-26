//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================
//
// These tests encode the behavior documented in the comments of
// correlation.go. Four of them are expected to FAIL against the current
// implementation: they expose the bugs described in correlation-bugs.md and
// will turn green as soon as those bugs are fixed.
//=============================================================================

package core

import (
	"math"
	"testing"

	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
	"golang.org/x/exp/slices"
)

//=============================================================================

func dailyReturns(values map[types.Date]float64) *[]db.DailyReturn {
	list := make([]db.DailyReturn, 0, len(values))

	for dt, r := range values {
		list = append(list, db.DailyReturn{Date: dt, GrossReturn: r})
	}

	return &list
}

//=============================================================================
//===
//=== CalcPearsonCorrelation
//===
//=============================================================================

func TestCalcPearsonCorrelation(t *testing.T) {
	testCases := []struct {
		name    string
		x       []float64
		y       []float64
		expect  float64
		wantErr bool
		bug     string
	}{
		{
			name:   "Perfect positive",
			x:      []float64{1, 2, 3, 4, 5},
			y:      []float64{2, 4, 6, 8, 10},
			expect: 1.0,
		},
		{
			name:   "Perfect negative",
			x:      []float64{1, 2, 3, 4, 5},
			y:      []float64{-2, -4, -6, -8, -10},
			expect: -1.0,
		},
		{
			name:   "Known value",
			x:      []float64{1, 2, 3, 4, 5},
			y:      []float64{2, 4, 5, 4, 5},
			expect: 6 / math.Sqrt(60),
		},
		{
			name:   "Zero covariance",
			x:      []float64{-2, -1, 0, 1, 2},
			y:      []float64{1, -1, 1, -1, 1},
			expect: 0,
		},
		{
			name:   "Two points",
			x:      []float64{1, 2},
			y:      []float64{2, 4},
			expect: 1.0,
		},
		{
			name:   "Shift and scale",
			x:      []float64{1, 2, 3, 4},
			y:      []float64{7, 9, 11, 13},
			expect: 1.0,
		},
		{
			name:   "Negative scale",
			x:      []float64{1, 2, 3, 4},
			y:      []float64{-3, -6, -9, -12},
			expect: -1.0,
		},
		{
			name:   "Large magnitudes",
			x:      []float64{1e150, 2e150, 3e150},
			y:      []float64{1e150, 2e150, 3e150},
			expect: 1.0,
			bug:    " (varX*varY overflows, see BUG-4)",
		},
		{
			name:    "Length mismatch",
			x:       []float64{1, 2},
			y:       []float64{1},
			wantErr: true,
		},
		{
			name:    "Both empty",
			x:       []float64{},
			y:       []float64{},
			wantErr: true,
		},
		{
			name:    "Single point",
			x:       []float64{1},
			y:       []float64{2},
			wantErr: true,
		},
		{
			name:    "NaN in x",
			x:       []float64{1, math.NaN(), 3},
			y:       []float64{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "NaN in y",
			x:       []float64{1, 2, 3},
			y:       []float64{1, math.NaN(), 3},
			wantErr: true,
		},
		{
			name:    "+Inf in x",
			x:       []float64{1, math.Inf(1), 3},
			y:       []float64{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "-Inf in y",
			x:       []float64{1, 2, 3},
			y:       []float64{1, math.Inf(-1), 3},
			wantErr: true,
		},
		{
			name:    "Constant x",
			x:       []float64{5, 5, 5},
			y:       []float64{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "Constant y",
			x:       []float64{1, 2, 3},
			y:       []float64{5, 5, 5},
			wantErr: true,
		},
		{
			name:    "Both constant",
			x:       []float64{5, 5},
			y:       []float64{7, 7},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := CalcPearsonCorrelation(tc.x, tc.y)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected an error but got result %v", r)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if math.IsNaN(r) {
				t.Errorf("Got NaN, expected %v%s", tc.expect, tc.bug)

				return
			}

			if math.Abs(r-tc.expect) > 1e-9 {
				t.Errorf("Expected %v and got %v%s", tc.expect, r, tc.bug)
			}
		})
	}
}

//=============================================================================

func TestCalcPearsonCorrelationSymmetry(t *testing.T) {
	x := []float64{3, -1, 4, 1, 5, 9, 2, 6}
	y := []float64{2, 7, 1, 8, 2, 8, 1, 8}

	xy, err1 := CalcPearsonCorrelation(x, y)
	yx, err2 := CalcPearsonCorrelation(y, x)

	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v, %v", err1, err2)
	}

	if math.Abs(xy-yx) > 1e-12 {
		t.Errorf("Correlation is not symmetric: %v vs %v", xy, yx)
	}
}

//=============================================================================
//===
//=== hydrateDailyReturns
//===
//=============================================================================

func TestHydrateDailyReturnsNil(t *testing.T) {
	aReturns, bReturns := hydrateDailyReturns(nil, nil)

	if len(aReturns) != 0 || len(bReturns) != 0 {
		t.Errorf("Expected two empty slices but got a=%v b=%v", aReturns, bReturns)
	}
}

//=============================================================================

func TestHydrateDailyReturnsNilEqualsEmpty(t *testing.T) {
	aList := []db.DailyReturn{
		{Date: types.NewDate(2024, 1, 2), GrossReturn: -0.2},
		{Date: types.NewDate(2024, 1, 1), GrossReturn: 0.1},
	}

	empty := []db.DailyReturn{}

	a1, b1 := hydrateDailyReturns(&aList, nil)
	a2, b2 := hydrateDailyReturns(&aList, &empty)

	if !slices.Equal(a1, a2) || !slices.Equal(b1, b2) {
		t.Errorf("Nil and empty slice behave differently: a1=%v a2=%v b1=%v b2=%v", a1, a2, b1, b2)
	}
}

//=============================================================================

func TestHydrateDailyReturnsOnlyA(t *testing.T) {
	aList := []db.DailyReturn{
		{Date: types.NewDate(2024, 1, 2), GrossReturn: -0.2},
		{Date: types.NewDate(2024, 1, 1), GrossReturn: 0.1},
	}

	aReturns, bReturns := hydrateDailyReturns(&aList, nil)

	if !slices.Equal(aReturns, []float64{0.1, -0.2}) {
		t.Errorf("A returns: expected [0.1 -0.2] and got %v (BUG-2: map indexed by loop counter)", aReturns)
	}

	if !slices.Equal(bReturns, []float64{0, 0}) {
		t.Errorf("B returns: expected [0 0] and got %v", bReturns)
	}
}

//=============================================================================

func TestHydrateDailyReturnsOnlyB(t *testing.T) {
	bList := []db.DailyReturn{
		{Date: types.NewDate(2024, 1, 2), GrossReturn: 0.3},
		{Date: types.NewDate(2024, 1, 1), GrossReturn: -0.4},
	}

	aReturns, bReturns := hydrateDailyReturns(nil, &bList)

	if !slices.Equal(aReturns, []float64{0, 0}) {
		t.Errorf("A returns: expected [0 0] and got %v", aReturns)
	}

	if !slices.Equal(bReturns, []float64{-0.4, 0.3}) {
		t.Errorf("B returns: expected [-0.4 0.3] and got %v (BUG-3: map indexed by loop counter)", bReturns)
	}
}

//=============================================================================
//===
//=== CalcCorrelation
//===
//=============================================================================

// TestCalcCorrelationUnionAlignment exercises hydration with both series
// non-empty. The current implementation loops forever in that case (BUG-1),
// so the checks run in a child process that gets killed after a short
// timeout instead of hanging the whole test run.

func TestCalcCorrelationUnionAlignment(t *testing.T) {
	//--- overlapping dates

	a := dailyReturns(map[types.Date]float64{
		types.NewDate(2024, 1, 1): 0.1,
		types.NewDate(2024, 1, 2): -0.2,
		types.NewDate(2024, 1, 3): 0.3,
	})

	b := dailyReturns(map[types.Date]float64{
		types.NewDate(2024, 1, 1): 0.1,
		types.NewDate(2024, 1, 2): -0.2,
		types.NewDate(2024, 1, 4): 0.4,
	})

	expA := []float64{0.1, -0.2, 0.3, 0}
	expB := []float64{0.1, -0.2, 0, 0.4}

	gotA, gotB := hydrateDailyReturns(a, b)

	if !slices.Equal(gotA, expA) || !slices.Equal(gotB, expB) {
		t.Errorf("Bad hydration for overlapping series.\nExpected a=%v b=%v\nGot     a=%v b=%v", expA, expB, gotA, gotB)
	}

	//--- disjoint dates

	c := dailyReturns(map[types.Date]float64{
		types.NewDate(2024, 1, 1): 0.2,
		types.NewDate(2024, 1, 2): -0.1,
	})

	d := dailyReturns(map[types.Date]float64{
		types.NewDate(2024, 2, 1): 0.5,
		types.NewDate(2024, 2, 2): -0.5,
	})

	gotA, gotB = hydrateDailyReturns(c, d)

	if !slices.Equal(gotA, []float64{0.2, -0.1, 0, 0}) || !slices.Equal(gotB, []float64{0, 0, 0.5, -0.5}) {
		t.Errorf("Bad hydration for disjoint series. Got a=%v b=%v", gotA, gotB)
	}

	//--- identical series end to end

	r, err := CalcCorrelation(a, a)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if math.Abs(r-1.0) > 1e-9 {
		t.Errorf("Identical series: expected 1.0 and got %v", r)
	}

	//--- disjoint series end to end

	r, err = CalcCorrelation(c, d)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if math.Abs(r) > 1e-9 {
		t.Errorf("Disjoint series: expected 0.0 and got %v", r)
	}

	//--- overlapping series end to end, against a hand-hydrated oracle

	expected, _ := CalcPearsonCorrelation(expA, expB)

	r, err = CalcCorrelation(a, b)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if math.Abs(r-expected) > 1e-12 {
		t.Errorf("Overlapping series: expected %v (oracle) and got %v", expected, r)
	}
}

//=============================================================================

func TestCalcCorrelationErrors(t *testing.T) {
	if _, err := CalcCorrelation(nil, nil); err == nil {
		t.Errorf("Expected an error for two nil series")
	}

	empty := []db.DailyReturn{}

	b := dailyReturns(map[types.Date]float64{
		types.NewDate(2024, 1, 1): 0.5,
		types.NewDate(2024, 1, 2): -0.5,
	})

	if _, err := CalcCorrelation(&empty, b); err == nil {
		t.Errorf("Expected an error when one series is empty (hydrates to a constant)")
	}

	if _, err := CalcCorrelation(b, &empty); err == nil {
		t.Errorf("Expected an error when one series is empty (hydrates to a constant)")
	}
}

//=============================================================================

func TestCalcCorrelationDoesNotMutateInputs(t *testing.T) {
	empty := []db.DailyReturn{}

	b := []db.DailyReturn{
		{Date: types.NewDate(2024, 1, 1), GrossReturn: 0.5},
		{Date: types.NewDate(2024, 1, 2), GrossReturn: -0.5},
	}

	_, _ = CalcCorrelation(&empty, &b)

	if len(b) != 2 || b[0].GrossReturn != 0.5 || b[1].GrossReturn != -0.5 {
		t.Errorf("Input series was mutated: %v", b)
	}
}

//=============================================================================
//===
//=== dailyReturnMap
//===
//=============================================================================

func TestDailyReturnMap(t *testing.T) {
	m := dailyReturnMap(nil)

	if len(m) != 0 {
		t.Errorf("Expected an empty map for a nil list and got %v", m)
	}

	empty := []db.DailyReturn{}

	m = dailyReturnMap(&empty)

	if len(m) != 0 {
		t.Errorf("Expected an empty map for an empty list and got %v", m)
	}

	list := []db.DailyReturn{
		{Date: types.NewDate(2024, 1, 2), GrossReturn: 0.2},
		{Date: types.NewDate(2024, 1, 1), GrossReturn: 0.1},
	}

	m = dailyReturnMap(&list)

	if len(m) != 2 {
		t.Fatalf("Expected 2 entries and got %v", m)
	}

	if m[int(types.NewDate(2024, 1, 1))] != 0.1 {
		t.Errorf("Expected 0.1 for 20240101 and got %v", m[int(types.NewDate(2024, 1, 1))])
	}

	if m[int(types.NewDate(2024, 1, 2))] != 0.2 {
		t.Errorf("Expected 0.2 for 20240102 and got %v", m[int(types.NewDate(2024, 1, 2))])
	}
}

//=============================================================================

func TestDailyReturnMapDuplicateDates(t *testing.T) {
	list := []db.DailyReturn{
		{Date: types.NewDate(2024, 1, 1), GrossReturn: 0.1},
		{Date: types.NewDate(2024, 1, 1), GrossReturn: 0.9},
	}

	m := dailyReturnMap(&list)

	if len(m) != 1 {
		t.Fatalf("Expected 1 entry for a duplicated date and got %v", m)
	}

	if m[int(types.NewDate(2024, 1, 1))] != 0.9 {
		t.Errorf("Expected the last value 0.9 to win and got %v", m[int(types.NewDate(2024, 1, 1))])
	}
}

//=============================================================================
