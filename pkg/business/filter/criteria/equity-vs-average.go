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

const EquityVsAverage = "equityVsAverage"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type EquityVsAverageConfig struct {
	Enabled bool    `json:"enabled"`
	Length  int     `json:"length"`
	Value   float64 `json:"value"`
}

//=============================================================================

func (c *EquityVsAverageConfig) validate() error {
	if c.Length < 1 || c.Length > 1000 {
		return req.NewBadRequestError("equityVsAverage.length must be between %v and %v", 1, 1000)
	}

	return nil
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type EquityVsAverageCriteria struct {
	BaseCriteria
	config *EquityVsAverageConfig
	equity []float64
}

//=============================================================================

func NewEquityVsAverageCriteria(c *EquityVsAverageConfig, runMode RunMode) (*EquityVsAverageCriteria, error) {
	cr := &EquityVsAverageCriteria{
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

func (c *EquityVsAverageCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	c.equity = s.Equity

	length := c.config.Length
	maSum  := 0.0

	a := NewActivation(EquityVsAverage, length)

	for i,t := range s.Time {
		maSum += s.Equity[i]

		if i >= length-1 {
			if i > length-1 {
				maSum -= s.Equity[i-length]
			}
			a.AddMetric(t, maSum/float64(length))
		}
	}

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate()
	}
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func (c *EquityVsAverageCriteria) evaluate() {
	length := c.config.Length
	value  := c.config.Value

	for i, m := range c.activation.Metrics {
		switch c.runMode {
			case RunModeActivation:
				if c.equity[i+length-1] >= m + value {
					c.activation.Flags[i] = 1
				}

			case RunModeDeactivation:
				if c.equity[i+length-1] < m - value {
					c.activation.Flags[i] = 1
				}
		}
	}
}

//=============================================================================
