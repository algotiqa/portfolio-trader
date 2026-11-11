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
	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================

const MaxDrawdown = "maxDrawdown"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type MaxDrawdownConfig struct {
	Config
}

//=============================================================================

func (c *MaxDrawdownConfig) validate() error {
	if c.Value < MinDrawdownValue || c.Value > MaxDrawdownValue {
		return req.NewBadRequestError("maxDrawdown.value must be between %v and %v", MinDrawdownValue, MaxDrawdownValue)
	}

	return c.Config.validate()
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

	length := c.config.Length

	a := NewActivation(MaxDrawdown, length)

	for i,t := range s.Time {
		if i >= length-1 {
			returns  := s.NetReturn[i-length+1:i+1]
			equity   := core.BuildEquity(&returns)
			_, maxDD := core.BuildDrawDown(equity)

			a.AddMetric(t, maxDD)
		}
	}

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate(c.config.Value, nil)
	}
}

//=============================================================================
