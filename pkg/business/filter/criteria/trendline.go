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
	Enabled bool    `json:"enabled"`
	Length  int     `json:"length"`
	Value   float64 `json:"value"`
}

//=============================================================================

func (c *TrendLineConfig) validate() error {
	if c.Length < 1 || c.Length > 1000 {
		return req.NewBadRequestError("trendLine.length must be between %v and %v", 1, 1000)
	}

	if c.Value < 1 || c.Value > 300 {
		return req.NewBadRequestError("trendLine.value must be between %v and %v", 1, 300)
	}

	return nil
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

	for i, _ := range s.Time {
		if i >= length-1 {
			slope := stats.LinearRegression(s.Time[i-length+1:i], s.Equity[i-length+1:i])
			a.AddMetric(slope)
		}
	}

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate(c.config.Value, nil)
	}
}

//=============================================================================
