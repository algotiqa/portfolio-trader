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
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type PeriodicSampler interface {
	Sample(time []time.Time, activations []int) ([]int, *criteria.Activation)
}

//=============================================================================

func NewPeriodicSampler(p *db.Portfolio) PeriodicSampler {
	if p.PeriodType == db.PeriodTypeWeekly {
		return NewWeeklyPeriodSampler(p.RunPeriod)
	}

	if p.PeriodType == db.PeriodTypeMonthly {
		return NewMonthlyPeriodSampler(p.RunPeriod)
	}

	panic("Unknown period type: " + p.PeriodType)
}

//=============================================================================
//===
//=== Weekly sampler
//===
//=============================================================================

type WeeklyPeriodSampler struct {
	runPeriod int
}

//=============================================================================

func NewWeeklyPeriodSampler(runPeriod int) *WeeklyPeriodSampler {
	return &WeeklyPeriodSampler{
		runPeriod: runPeriod,
	}
}

//=============================================================================

func (s *WeeklyPeriodSampler) Sample(time []time.Time, activations []int) ([]int, *criteria.Activation) {
	date := firstSaturday(time[0])

	res  := make([]int, len(activations))
	prev := activations[0]

	a := criteria.NewActivation("finalActivation", 0)
	a.AddFlag(time[0], activations[0])

	for i, t := range time {
		if t.Before(date) {
			res[i] = prev
		} else {
			prev   = activations[i-1]
			res[i] = prev
			a.AddFlag(date, prev)
			date = date.AddDate(0, 0, 7*s.runPeriod)
			for date.Before(t) {
				a.AddFlag(date, prev)
				date = date.AddDate(0, 0, 7*s.runPeriod)
			}
		}
	}

	return res, a
}

//=============================================================================
//===
//=== Monthly sampler
//===
//=============================================================================

type MonthlyPeriodSampler struct {
	runPeriod int
}

//=============================================================================

func NewMonthlyPeriodSampler(runPeriod int) *MonthlyPeriodSampler {
	return &MonthlyPeriodSampler{
		runPeriod: runPeriod,
	}
}

//=============================================================================

func (s *MonthlyPeriodSampler) Sample(time []time.Time, activations []int) ([]int, *criteria.Activation) {
	date := firstMonth(time[0])

	res := make([]int, len(activations))
	prev := activations[0]

	a := criteria.NewActivation("finalActivation", 0)
	a.AddFlag(time[0], activations[0])

	for i, t := range time {
		if t.Before(date) {
			res[i] = prev
		} else {
			prev   = activations[i-1]
			res[i] = prev
			a.AddFlag(date, prev)
			date = date.AddDate(0, s.runPeriod, 0)
			for date.Before(t) {
				a.AddFlag(date, prev)
				date = date.AddDate(0, s.runPeriod, 0)
			}
		}
	}

	return res, a
}

//=============================================================================
//===
//=== Private utility functions
//===
//=============================================================================

func firstSaturday(t time.Time) time.Time {
	days := int((time.Saturday - t.Weekday() + 7) % 7)
	if days == 0 {
		days = 7
	}
	return t.AddDate(0, 0, days)
}

//=============================================================================

func firstMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
}

//=============================================================================
