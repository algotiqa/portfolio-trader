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

const SharpeRatio = "sharpeRatio"

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type SharpeRatioConfig struct {
	Enabled bool    `json:"enabled"`
	Length  int     `json:"length"`
	Value   float64 `json:"value"`
}

//=============================================================================

func (c *SharpeRatioConfig) validate() error {
	if c.Length < 1 || c.Length > 1000 {
		return req.NewBadRequestError("sharpeRatio.length must be between %v and %v", 1, 1000)
	}

	return nil
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

	a := NewActivation(SharpeRatio, 0)

	//--- We may not be able to calculate if length is too high

	if a.Flags != nil {
		c.activation = a
		c.evaluate(0,nil)
	}
}

//=============================================================================
