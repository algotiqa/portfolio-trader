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

const ReturnMaxDDRatio = "returnMaxDDRatio"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type ReturnMaxDDRatioConfig struct {
	Config
}

//=============================================================================

func (c *ReturnMaxDDRatioConfig) validate() error {
	if c.Value < MinReturnMaxDDValue || c.Value > MaxReturnMaxDDValue {
		return req.NewBadRequestError("returnMaxDDRatio.value must be between %v and %v", MinReturnMaxDDValue, MaxReturnMaxDDValue)
	}

	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type ReturnMaxDDRatioCriteria struct {
	BaseCriteria
	config *ReturnMaxDDRatioConfig
}

//=============================================================================

func NewReturnMaxDDRatioCriteria(c *ReturnMaxDDRatioConfig, runMode RunMode) (*ReturnMaxDDRatioCriteria, error) {
	cr := &ReturnMaxDDRatioCriteria{
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

func (c *ReturnMaxDDRatioCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	//retSum := 0.0
	length := c.config.Length

	a := NewActivation(ReturnMaxDDRatio, length)

	for i, t := range s.Time {
		if i >= length-1 {
			returns := s.NetReturn[i-length+1:i+1]
			equity  := core.BuildEquity(&returns)
			_,maxDD := core.BuildDrawDown(equity)

			ret := (*equity)[len(*equity) -1]

			if maxDD == 0 {
				ret = -1
			} else {
				ret = ret / -maxDD
			}

			a.AddMetric(t, ret)
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
