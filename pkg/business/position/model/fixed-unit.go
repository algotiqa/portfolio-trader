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

var specUnits = core.NewNumberParamSpec( "units", true, 1, 10000, &DefUnits)

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
			units: 1,
		},
	}
}

//=============================================================================

func (m *FixedUnitModel) Name() db.ModelName {
	return db.ModelFixedUnit
}

//=============================================================================

func (m *FixedUnitModel) Init(config map[string]any) error {
	//m.units = config.FuUnits
	//
	//if m.units < 1 {
	//	return errors.New("units must be positive: " + strconv.Itoa(m.units))
	//}
	//
	//if m.units > 10000 {
	//	return errors.New("units too big: " + strconv.Itoa(m.units))
	//}

	return nil
}

//=============================================================================

func (m *FixedUnitModel) Config() map[string]any {
	cfg := make(map[string]any)
	cfg[specUnits.Name] = m.config.units

	return cfg
}

//=============================================================================

func (m *FixedUnitModel) PositionInit(ts *TradingSnapshot) {}

//=============================================================================

func (m *FixedUnitModel) PositionFor(ts *TradingSnapshot) int {
	return m.config.units
}

//=============================================================================
