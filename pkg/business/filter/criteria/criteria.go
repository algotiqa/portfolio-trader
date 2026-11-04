//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package criteria

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== FilteringCriteria
//===
//=============================================================================

type RunMode string

const (
	RunModeActivation   RunMode = "A"
	RunModeDeactivation RunMode = "D"
)

//=============================================================================

type FilteringCriteria interface {
	Eval(s *Snapshot)
	Activation() *Activation
}

//=============================================================================

type Snapshot struct {
	Time      []time.Time
	NetReturn []float64
	Equity    []float64
	Drawdown  []float64
}

//=============================================================================
//===
//=== Activation
//===
//=============================================================================

type Activation struct {
	Name    string      `json:"name"`
	Times   []time.Time `json:"times"`
	Flags   []int       `json:"flags"`
	Metrics []float64   `json:"metrics"`
	length  int
}

//=============================================================================

func NewActivation(name string, length int) *Activation {
	return &Activation{
		Name  : name,
		length: length,
	}
}

//=============================================================================

func (a *Activation) AddMetric(time time.Time, metric float64) {
	a.Times   = append(a.Times,   time)
	a.Flags   = append(a.Flags,   0)
	a.Metrics = append(a.Metrics, core.Trunc2d(metric))
}

//=============================================================================

func (a *Activation) AddFlag(time time.Time, flag int) {
	a.Times = append(a.Times,   time)
	a.Flags = append(a.Flags,   flag)
}

//=============================================================================

func (a *Activation) LastActive() int {
	return a.Flags[len(a.Flags)-1]
}

//=============================================================================

func (a *Activation) Active(i int) int {
	length := a.length

	if i >= length-1 {
		return a.Flags[i-length+1]
	}

	return 0
}

//=============================================================================
