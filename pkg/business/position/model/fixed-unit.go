//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package model

import "encoding/json"

//=============================================================================

const FixedUnit = "FU"

//=============================================================================

type FixedUnitConfig struct {
	Units float64 `json:"units"`
}

//=============================================================================

type FixedUnitModel struct {
	config *FixedUnitConfig
}

//=============================================================================

func newFixedUnitModel(config string) (*FixedUnitModel, error) {
	c := &FixedUnitConfig{}
	err := json.Unmarshal([]byte(config), c)
	if err != nil {
		return nil, err
	}

	return &FixedUnitModel{
		config : c,
	},nil
}

//=============================================================================

func (fm *FixedUnitModel) Calc() {
}

//=============================================================================
