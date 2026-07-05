//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package model

import (
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Specs
//===
//=============================================================================

var DefUnits = 1

var SpecUnits = core.NewNumberParamSpec[int]( "units", true, 1, 10000, &DefUnits)

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type FixedUnitConfig struct {
	units int
}

//=============================================================================
//===
//=== Model
//===
//=============================================================================

type FixedUnitModel struct {
	config *FixedUnitConfig
}

//=============================================================================

func NewFixedUnitModel() *FixedUnitModel {
	return &FixedUnitModel{
		config: &FixedUnitConfig{
			units: DefUnits,
		},
	}
}

//=============================================================================

func (m *FixedUnitModel) Name() db.ModelName {
	return db.ModelFixedUnit
}

//=============================================================================

func (m *FixedUnitModel) Init(config map[string]any) error {
	units,err := core.MapNumber[int](config, SpecUnits)
	if err != nil {
		return err
	}

	m.config.units = *units
	return nil
}

//=============================================================================

func (m *FixedUnitModel) Config() map[string]any {
	cfg := make(map[string]any)
	cfg[SpecUnits.Name] = m.config.units

	return cfg
}

//=============================================================================

func (m *FixedUnitModel) Spec() map[string]any {
	specs := make(map[string]any)

	specs[SpecUnits.Name] = SpecUnits

	return specs
}

//=============================================================================

func (m *FixedUnitModel) PositionInit(ts *TradingSnapshot) {}

//=============================================================================

func (m *FixedUnitModel) PositionFor(ts *TradingSnapshot) int {
	return m.config.units
}

//=============================================================================
