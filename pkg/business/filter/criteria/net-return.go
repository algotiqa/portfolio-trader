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

const NetReturn = "netReturn"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type NetReturnConfig struct {
	Config
}

//=============================================================================

func (c *NetReturnConfig) validate() error {
	return c.Config.validate()
}

//=============================================================================
//===
//=== Criteria
//===
//=============================================================================

type NetReturnCriteria struct {
	BaseCriteria
	config *NetReturnConfig
}

//=============================================================================

func NewNetReturnCriteria(c *NetReturnConfig, runMode RunMode) (*NetReturnCriteria, error) {
	cr := &NetReturnCriteria{
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

func (c *NetReturnCriteria) Eval(s *Snapshot) {
	if !c.config.Enabled {
		return
	}

	retSum := 0.0
	length := c.config.Length

	a := NewActivation(NetReturn, length)

	for i,t := range s.Time {
		if i >= length-1 {
			retSum = s.Equity[i]

			if i > length-1 {
				retSum -= s.Equity[i-length]
			}

			a.AddMetric(t, retSum)
		}
	}

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate(c.config.Value, nil)
	}
}

//=============================================================================
