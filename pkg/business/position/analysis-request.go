//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package position

import (
	"fmt"

	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== TradingPosition
//===
//=============================================================================

type Model struct {
	Name    db.ModelName    `json:"name"`
	Config  map[string]any  `json:"config"`
}

//=============================================================================

func NewModel(name db.ModelName, config map[string]any) *Model {
	return &Model{
		Name  : name,
		Config: config,
	}
}

//=============================================================================
//===
//=== TradingPosition
//===
//=============================================================================

type TradingPosition struct {
	Params  Parameters  `json:"params"`
	Model   Model       `json:"model"`
}

//=============================================================================
//===
//=== AnalysisRequest
//===
//=============================================================================

type AnalysisRequest struct {
	Params  *Parameters         `json:"params"`
	Model   *Model              `json:"model"`
	Period  core.SelectedPeriod `json:"period"`
}

//=============================================================================
//===
//=== Position sizing general parameters
//===
//=============================================================================

var DefInitialCapital = 15000.0
var DefRuinPercentage = 25.0
var DefMaxUnits       = 1
var DefRiskValue      = 150.0
var DefRiskPerUnit    = db.RpuStopLoss

var SpecInitialCapital = core.NewNumberParamSpec[float64]   ("initialCapital",  true, 1, 1000000000, &DefInitialCapital)
var SpecRuinParcentage = core.NewNumberParamSpec[float64]   ("ruinPercentage",  true, 0,        100, &DefRuinPercentage)
var SpecMarginOverride = core.NewNumberParamSpec[float64]   ("marginOverride", false, 0,    1000000, nil)
var SpecMaxUnits       = core.NewNumberParamSpec[int]       ("maxUnits",        true, 1,     100000, &DefMaxUnits)
var SpecRiskPerUnit    = core.NewListParamSpec  [db.RpuType]("riskPerUnit",     true, db.RpuDomain, DefRiskPerUnit)
var SpecRiskValue      = core.NewNumberParamSpec[float64]   ("riskValue",      false, 1,    50000.0, &DefRiskValue)

//=============================================================================

type Parameters struct {
	InitialCapital  *float64       `json:"initialCapital"`
	RuinPercentage  *float64       `json:"ruinPercentage"`
	MarginOverride  *float64       `json:"marginOverride"`
	MaxUnits        *int           `json:"maxUnits"`
	RiskPerUnit     db.RpuType     `json:"riskPerUnit"`
	RiskValue       *float64       `json:"riskValue"`
}

//=============================================================================

func (p *Parameters) Validate() error {
	var err error

	p.InitialCapital,err = SpecInitialCapital.Validate(p.InitialCapital)
	if err != nil {
		return err
	}

	p.RuinPercentage,err = SpecRuinParcentage.Validate(p.RuinPercentage)
	if err != nil {
		return err
	}

	p.MarginOverride,err = SpecMarginOverride.Validate(p.MarginOverride)
	if err != nil {
		return err
	}

	p.MaxUnits,err = SpecMaxUnits.Validate(p.MaxUnits)
	if err != nil {
		return err
	}

	p.RiskPerUnit,err = SpecRiskPerUnit.Validate(p.RiskPerUnit)
	if err != nil {
		return err
	}

	p.RiskValue,err = SpecRiskValue.Validate(p.RiskValue)
	if err != nil {
		return err
	}

	if p.RiskPerUnit == db.RpuFixedValue && p.RiskValue == nil {
		return fmt.Errorf("'riskValue' is mandatory when riskPerUnit='fixedValue'")
	}

	return nil
}

//=============================================================================
