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
	"github.com/algotiqa/portfolio-trader/pkg/core/stats"
)

//=============================================================================

const TrendLine = "trendLine"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type TrendLineConfig struct {
	Config
}

//=============================================================================

func (c *TrendLineConfig) validate() error {
	if c.Value < MinTrendlineValue || c.Value > MaxTrendlineValue {
		return req.NewBadRequestError("trendLine.value must be between %v and %v", MinTrendlineValue, MaxTrendlineValue)
	}

	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type TrendLineCriteria struct {
	BaseCriteria
	config *TrendLineConfig
}

//=============================================================================

func NewTrendLineCriteria(c *TrendLineConfig, runMode RunMode) (*TrendLineCriteria, error) {
	cr := &TrendLineCriteria{
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

func (c *TrendLineCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	length := c.config.Length

	a := NewActivation(TrendLine, length)

	for i,t := range s.Time {
		if i >= length-1 {
			slope := stats.LinearRegression(s.Time[i-length+1:i+1], s.Equity[i-length+1:i+1])
			a.AddMetric(t, slope)
		}
	}

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate(c.config.Value, nil)
	}
}

//=============================================================================
