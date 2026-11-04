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

	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

func dailyTimes(start time.Time, n int) []time.Time {
	res := make([]time.Time, n)

	for i := range res {
		res[i] = start.AddDate(0, 0, i)
	}

	return res
}

//=============================================================================

func TestFirstSaturday(t *testing.T) {
	for d := 0; d < 7; d++ {
		tm   := time.Date(2026, 1, 4+d, 12, 0, 0, 0, time.UTC)
		got  := firstSaturday(tm)

		expDay := 10
		if d == 6 {
			expDay = 17
		}

		exp := time.Date(2026, 1, expDay, 12, 0, 0, 0, time.UTC)

		if !got.Equal(exp) {
			t.Errorf("firstSaturday(%v): expected %v but got %v", tm.Weekday(), exp, got)
		}
	}
}

//=============================================================================

func TestFirstMonth(t *testing.T) {
	midMonth := time.Date(2026, 1, 15, 14, 30, 0, 0, time.UTC)
	firstOfNext := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	if !firstMonth(midMonth).Equal(firstOfNext) {
		t.Errorf("firstMonth(%v): expected %v but got %v", midMonth, firstOfNext, firstMonth(midMonth))
	}

	firstDay := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)

	if !firstMonth(firstDay).Equal(firstOfNext) {
		t.Errorf("firstMonth(%v): expected %v but got %v", firstDay, firstOfNext, firstMonth(firstDay))
	}
}

//=============================================================================

func TestWeeklyPeriodSampler(t *testing.T) {
	start := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	times := dailyTimes(start, 21)

	activations := make([]int, 21)
	activations[4]  = 1
	activations[18] = 1

	tests := []struct {
		runPeriod int
		expected  []int
	}{
		{1, []int{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1}},
		{2, []int{0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
	}

	for _, tt := range tests {
		s := NewWeeklyPeriodSampler(tt.runPeriod)
		got,_ := s.Sample(times, activations)

		if len(got) != len(activations) {
			t.Fatalf("Bad sample length. Expected %v but got %v", len(activations), len(got))
		}

		if !slices.Equal(got, tt.expected) {
			t.Errorf("Bad weekly sampling (runPeriod %v).\nExpected %v\nbut got  %v", tt.runPeriod, tt.expected, got)
		}
	}
}

//=============================================================================

func TestMonthlyPeriodSampler(t *testing.T) {
	start := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	times := dailyTimes(start, 55)

	activations := make([]int, 55)
	activations[16] = 1

	tests := []struct {
		runPeriod int
		lastOne   int
	}{
		{1, 45},
		{2, 55},
	}

	for _, tt := range tests {
		s := NewMonthlyPeriodSampler(tt.runPeriod)
		got,_ := s.Sample(times, activations)

		if len(got) != len(activations) {
			t.Fatalf("Bad sample length. Expected %v but got %v", len(activations), len(got))
		}

		expected := make([]int, 55)

		for i := 17; i < tt.lastOne; i++ {
			expected[i] = 1
		}

		if !slices.Equal(got, expected) {
			t.Errorf("Bad monthly sampling (runPeriod %v).\nExpected %v\nbut got  %v", tt.runPeriod, expected, got)
		}
	}
}

//=============================================================================

func TestNewPeriodicSampler(t *testing.T) {
	p := db.Portfolio{PeriodType: db.PeriodTypeWeekly, RunPeriod: 2}

	if _, ok := NewPeriodicSampler(&p).(*WeeklyPeriodSampler); !ok {
		t.Errorf("Expected a weekly sampler")
	}

	p = db.Portfolio{PeriodType: db.PeriodTypeMonthly, RunPeriod: 3}

	if _, ok := NewPeriodicSampler(&p).(*MonthlyPeriodSampler); !ok {
		t.Errorf("Expected a monthly sampler")
	}

	defer func() {
		if recover() == nil {
			t.Errorf("Expected a panic for an unknown period type")
		}
	}()

	p = db.Portfolio{PeriodType: db.PeriodType("XX")}
	NewPeriodicSampler(&p)
}

//=============================================================================
