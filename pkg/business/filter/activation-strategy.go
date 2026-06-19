//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import "time"

//=============================================================================
//===
//=== ActivationStrategy
//===
//=============================================================================

type ActivationStrategy struct {
	activation *Activation
	enabled    bool
	index      int
}

//=============================================================================

func (as *ActivationStrategy) IsActive(t time.Time) bool {
	//--- Strategy not enabled: skip it returning always 1
	if !as.enabled {
		return true
	}

	//--- Strategy not computable: return true because we must align with unfiltered equity
	if as.activation == nil {
		return true
	}

	if t.Before(as.activation.Time[as.index]) {
		return true
	}

	if t != as.activation.Time[as.index] {
		panic("Help!")
	}

	as.index++
	return as.activation.Values[as.index -1] != 0
}

//=============================================================================

func NewActivationStrategy(a *Activation, enabled bool) *ActivationStrategy {
	return &ActivationStrategy{
		activation: a,
		enabled: enabled,
		index: 0,
	}
}

//=============================================================================
