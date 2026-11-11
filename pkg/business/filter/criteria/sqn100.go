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
	"math"

	"golang.org/x/exp/stats"
)

//=============================================================================

const Sqn100 = "sqn100"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type Sqn100Config struct {
	Config
}

//=============================================================================

func (c *Sqn100Config) validate() error {
	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type Sqn100Criteria struct {
	BaseCriteria
	config *Sqn100Config
}

//=============================================================================

func NewSqn100Criteria(c *Sqn100Config, runMode RunMode) (*Sqn100Criteria, error) {
	cr := &Sqn100Criteria{
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

func (c *Sqn100Criteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	length := c.config.Length

	a := NewActivation(Sqn100,length)

	for i,t := range s.Time {
		if i >= length-1 {
			mean, stdd := stats.MeanAndStdDev(s.NetReturn[i-length+1:i+1])

			sqn := -1.0

			if stdd > 0 {
				trades := math.Min(float64(length), 100)
				sqn = (mean/stdd)*math.Sqrt(trades)
			}

			a.AddMetric(t, sqn)
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
