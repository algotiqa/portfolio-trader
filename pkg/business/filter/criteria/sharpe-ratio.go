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
	"golang.org/x/exp/stats"
)

//=============================================================================

const SharpeRatio = "sharpeRatio"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type SharpeRatioConfig struct {
	Config
}

//=============================================================================

func (c *SharpeRatioConfig) validate() error {
	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type SharpeRatioCriteria struct {
	BaseCriteria
	config *SharpeRatioConfig
}

//=============================================================================

func NewSharpeRatioCriteria(c *SharpeRatioConfig, runMode RunMode) (*SharpeRatioCriteria, error) {
	cr := &SharpeRatioCriteria{
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

func (c *SharpeRatioCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	length := c.config.Length

	a := NewActivation(SharpeRatio,length)

	for i,t := range s.Time {
		if i >= length-1 {
			mean, stdd := stats.MeanAndStdDev(s.NetReturn[i-length+1:i+1])

			sr := -1.0

			if stdd > 0 {
				sr = mean/stdd
			}

			a.AddMetric(t, sr)
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
