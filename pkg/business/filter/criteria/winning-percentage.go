//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package criteria

import "github.com/algotiqa/core/req"

//=============================================================================

const WinningPercentage = "winningPercentage"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type WinningPercentageConfig struct {
	Config
}

//=============================================================================

func (c *WinningPercentageConfig) validate() error {
	if c.Value < MinWinningPercentage || c.Value > MaxWinningPercentage {
		return req.NewBadRequestError("profitFactor.value must be between %v and %v", MinWinningPercentage, MaxWinningPercentage)
	}

	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type WinningPercentageCriteria struct {
	BaseCriteria
	config *WinningPercentageConfig
}

//=============================================================================

func NewWinningPercentageCriteria(c *WinningPercentageConfig, runMode RunMode) (*WinningPercentageCriteria, error) {
	cr := &WinningPercentageCriteria{
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

func (c *WinningPercentageCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	posCount := 0.0
	totCount := 0.0
	length   := c.config.Length

	a := NewActivation(WinningPercentage, length)

	for i,t := range s.Time {
		if s.NetReturn[i] != 0 {
			totCount++

			if s.NetReturn[i] > 0 {
				posCount++
			}
		}

		if i >= length-1 {
			if i > length-1 {
				if s.NetReturn[i-length] != 0 {
					totCount--

					if s.NetReturn[i-length] > 0 {
						posCount--
					}
				}
			}

			metric := -1.0

			if totCount > 0 {
				metric = posCount * 100 / totCount
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
