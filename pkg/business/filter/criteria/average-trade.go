//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package criteria

//=============================================================================

const AverageTrade = "averageTrade"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type AverageTradeConfig struct {
	Config
}

//=============================================================================

func (c *AverageTradeConfig) validate() error {
	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type AverageTradeCriteria struct {
	BaseCriteria
	config *AverageTradeConfig
}

//=============================================================================

func NewAverageTradeCriteria(c *AverageTradeConfig, runMode RunMode) (*AverageTradeCriteria, error) {
	cr := &AverageTradeCriteria{
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

func (c *AverageTradeCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	retSum := 0.0
	count  := 0.0
	length := c.config.Length

	a := NewActivation(AverageTrade, length)

	for i,t := range s.Time {
		if i >= length-1 {
			retSum = s.Equity[i]
			count++

			if i > length-1 {
				retSum -= s.Equity[i-length]
				count--
			}

			a.AddMetric(t,retSum/count)
		}
	}

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate(c.config.Value, nil)
	}
}

//=============================================================================
