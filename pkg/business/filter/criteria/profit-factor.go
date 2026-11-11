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

const ProfitFactor = "profitFactor"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type ProfitFactorConfig struct {
	Config
}

//=============================================================================

func (c *ProfitFactorConfig) validate() error {
	if c.Value < MinProfitFactor || c.Value > MaxProfitFactor {
		return req.NewBadRequestError("profitFactor.value must be between %v and %v", MinProfitFactor, MaxProfitFactor)
	}

	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type ProfitFactorCriteria struct {
	BaseCriteria
	config *ProfitFactorConfig
}

//=============================================================================

func NewProfitFactorCriteria(c *ProfitFactorConfig, runMode RunMode) (*ProfitFactorCriteria, error) {
	cr := &ProfitFactorCriteria{
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

func (c *ProfitFactorCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	profit := 0.0
	loss   := 0.0
	length := c.config.Length

	a := NewActivation(ProfitFactor, length)

	for i,t := range s.Time {
		if s.NetReturn[i] > 0 {
			profit += s.NetReturn[i]
		}

		if s.NetReturn[i] < 0 {
			loss += - s.NetReturn[i]
		}

		if i >= length-1 {
			if i > length - 1 {
				oldVal := s.NetReturn[i-length]

				if oldVal > 0 {
					profit -= oldVal
				}

				if oldVal < 0 {
					loss += oldVal
				}
			}

			metric := -1.0

			if loss > 0 {
				metric = profit / loss
			}

			a.AddMetric(t, metric)
		}
	}

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		skipValue := -1.0
		c.evaluate(c.config.Value, &skipValue)
	}
}

//=============================================================================
