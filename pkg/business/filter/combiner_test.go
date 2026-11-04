//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import (
	"testing"
	"time"

	"golang.org/x/exp/slices"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
)

//=============================================================================

func testSnapshot(equity []float64, netReturn []float64, drawdown []float64) *criteria.Snapshot {
	size := len(equity)

	s := &criteria.Snapshot{
		Time      : make([]time.Time, size),
		Equity    : equity,
		NetReturn : netReturn,
		Drawdown  : drawdown,
	}

	for i := range s.Time {
		s.Time[i] = time.Date(2026, 1, 5+i, 0, 0, 0, 0, time.UTC)
	}

	return s
}

//=============================================================================

func TestNewCombinerValidation(t *testing.T) {
	unknown := &Config{CombinationMode: "X"}

	if _, err := NewCombiner(unknown, criteria.RunModeActivation); err == nil {
		t.Errorf("Expected an error for an unknown combination mode")
	}

	invalid := &Config{
		CombinationMode: CombinationModeAll,
		NetReturn:       criteria.NetReturnConfig{Enabled: true, Length: 0, Value: 0},
	}

	if _, err := NewCombiner(invalid, criteria.RunModeActivation); err == nil {
		t.Errorf("Expected an error for an invalid enabled criterion")
	}

	disabled := &Config{
		CombinationMode: CombinationModeAll,
		NetReturn:       criteria.NetReturnConfig{Enabled: false, Length: 0, Value: 0},
	}

	if _, err := NewCombiner(disabled, criteria.RunModeActivation); err != nil {
		t.Errorf("Disabled criteria must not be validated, got %v", err)
	}
}

//=============================================================================

func TestCombinerCombinationModes(t *testing.T) {
	equity    := []float64{100, 110,  90, 95}
	netProfit := []float64{ 10,  -5,  20, -3}
	drawdown  := []float64{ -1,  -6,  -6, -1}

	s := testSnapshot(equity, netProfit, drawdown)

	tests := []struct {
		mode     CombinationMode
		expected []int
	}{
		{CombinationModeOne , []int{1, 1, 1, 1}},
		{CombinationModeAll , []int{0, 0, 0, 0}},
		{CombinationModeMost, []int{1, 0, 0, 1}},
	}

	for _, tt := range tests {
		cfg := &Config{
			CombinationMode:   tt.mode,
			NetReturn:         criteria.NetReturnConfig        {Enabled: true, Length: 1, Value: 0},
			WinningPercentage: criteria.WinningPercentageConfig{Enabled: true, Length: 3, Value: 50},
			MaxDrawdown:       criteria.MaxDrawdownConfig      {Enabled: true, Value: -5},
		}

		c, err := NewCombiner(cfg, criteria.RunModeActivation)
		if err != nil {
			t.Fatalf("cannot create combiner: %v", err)
		}

		c.Eval(s)

		if !slices.Equal(c.Outcome(), tt.expected) {
			t.Errorf("Bad %v outcome. Expected %v but got %v", tt.mode, tt.expected, c.Outcome())
		}
	}
}

//=============================================================================

func TestCombinerActivations(t *testing.T) {
	equity    := []float64{100, 110,  90, 95}
	netProfit := []float64{ 10,  -5,  20, -3}
	drawdown  := []float64{ -1,  -6,  -6, -1}

	s := testSnapshot(equity, netProfit, drawdown)

	cfg := &Config{
		CombinationMode:   CombinationModeOne,
		NetReturn:         criteria.NetReturnConfig        {Enabled: true, Length: 1, Value: 0},
		WinningPercentage: criteria.WinningPercentageConfig{Enabled: true, Length: 3, Value: 50},
		MaxDrawdown:       criteria.MaxDrawdownConfig      {Enabled: true, Value: -5},
	}

	c, err := NewCombiner(cfg, criteria.RunModeActivation)
	if err != nil {
		t.Fatalf("cannot create combiner: %v", err)
	}

	c.Eval(s)

	acts := c.Activations()

	if len(acts) != 3 {
		t.Fatalf("Expected 3 activations but got %v", len(acts))
	}

	var names []string

	for _, a := range acts {
		names = append(names, a.Name)
	}

	slices.Sort(names)

	expected := []string{criteria.MaxDrawdown, criteria.NetReturn, criteria.WinningPercentage}

	if !slices.Equal(names, expected) {
		t.Errorf("Bad activation names. Expected %v but got %v", expected, names)
	}
}

//=============================================================================

func TestCombinerNoCriteria(t *testing.T) {
	s := testSnapshot([]float64{1, 2, 3, 4}, []float64{0, 0, 0, 0}, []float64{0, 0, 0, 0})

	for _, mode := range []CombinationMode{CombinationModeOne, CombinationModeAll, CombinationModeMost} {
		cfg := &Config{CombinationMode: mode}

		c, err := NewCombiner(cfg, criteria.RunModeActivation)
		if err != nil {
			t.Fatalf("cannot create combiner: %v", err)
		}

		c.Eval(s)

		if !slices.Equal(c.Outcome(), []int{0, 0, 0, 0}) {
			t.Errorf("A combiner with no active criteria must return all zeros, got %v for mode %v", c.Outcome(), mode)
		}

		if len(c.Activations()) != 0 {
			t.Errorf("Expected no activations but got %v", len(c.Activations()))
		}
	}
}

//=============================================================================
