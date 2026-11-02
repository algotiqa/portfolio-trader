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

const Sqn100 = "sqn100"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type Sqn100Config struct {
	Enabled bool    `json:"enabled"`
	Length  int     `json:"length"`
	Value   float64 `json:"value"`
}

//=============================================================================

func (c *Sqn100Config) validate() error {
	if c.Length < 1 || c.Length > 1000 {
		return req.NewBadRequestError("sqn100.length must be between %v and %v", 1, 1000)
	}

	return nil
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

	a := NewActivation(Sqn100, 0)

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate(0,nil)
	}
}

//=============================================================================
