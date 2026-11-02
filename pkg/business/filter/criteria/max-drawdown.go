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
	"github.com/algotiqa/core/req"
)

//=============================================================================

const MaxDrawdown = "maxDrawdown"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type MaxDrawdownConfig struct {
	Enabled bool    `json:"enabled"`
	Value   float64 `json:"value"`
}

//=============================================================================

func (c *MaxDrawdownConfig) validate() error {
	if c.Value > 0 {
		return req.NewBadRequestError("maxDrawdown.value must be <= 0")
	}

	return nil
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type MaxDrawdownCriteria struct {
	BaseCriteria
	config *MaxDrawdownConfig
}

//=============================================================================

func NewMaxDrawdownCriteria(c *MaxDrawdownConfig, runMode RunMode) (*MaxDrawdownCriteria, error) {
	cr := &MaxDrawdownCriteria{
		config : c,
	}

	cr.runMode = runMode

	var err error
	if c.Enabled {
		err = c.validate()
	}

	return cr, err
}

//=============================================================================

func (c *MaxDrawdownCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	a := NewActivation(MaxDrawdown, 1)

	for i := range s.Time {
		a.AddMetric(s.Drawdown[i])
	}

	c.activation = a
	c.evaluate(c.config.Value, nil)
}

//=============================================================================
