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
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== AnalysisResponse building
//===
//=============================================================================

func RunAnalysis(ts *db.TradingSystem, curModel, selModel model.PositionModel, trades *[]db.Trade, pos *db.TradingPosition) *AnalysisResponse {
	baseline := model.NewFixedUnitModel()

	res := &AnalysisResponse{}
	res.TradingSystem = buildTradingSystem(ts)
	res.Params        = extractParameters(pos)
	res.Baseline      = calcAnalysisResult(baseline, trades, ts.CostPerOperation)
	res.Current       = calcAnalysisResult(curModel, trades, ts.CostPerOperation)
	res.Selected      = calcAnalysisResult(selModel, trades, ts.CostPerOperation)
	res.ParamSpecs    = buildParamSpecs()
	return res
}

//=============================================================================

func buildTradingSystem(ts *db.TradingSystem) *TradingSystem {
	return &TradingSystem{
		Id  : ts.Id,
		Name: ts.Name,
	}
}

//=============================================================================

func calcAnalysisResult(model model.PositionModel, trades *[]db.Trade, costPerOper float64) *AnalysisResult {
	if model == nil {
		return nil
	}

	gross := calcModelPerformance()
	net   := calcModelPerformance()

	return &AnalysisResult{
		Model : NewModel(model.Name(), model.Config()),
		Gross : gross,
		Net   : net,
	}
}

//=============================================================================

func calcModelPerformance() *ModelPerformance {
	return nil
}

//=============================================================================

func extractParameters(p *db.TradingPosition) *Parameters {
	return &Parameters{
		InitialCapital: &p.InitialCapital,
		RuinPercentage: &p.RuinPercentage,
		MarginOverride: p.MarginOverride,
		MaxUnits      : &p.MaxUnits,
		RiskPerUnit   : p.RiskPerUnit,
		RiskValue     : p.RiskValue,
	}
}

//=============================================================================

func buildParamSpecs() map[string]any {
	specs := make(map[string]any)

	specs[SpecInitialCapital.Name] = SpecInitialCapital
	specs[SpecRuinParcentage.Name] = SpecRuinParcentage
	specs[SpecMarginOverride.Name] = SpecMarginOverride
	specs[SpecMaxUnits.Name]       = SpecMaxUnits
	specs[SpecRiskPerUnit.Name]    = SpecRiskPerUnit
	specs[SpecRiskValue.Name]      = SpecRiskValue

	return specs
}

//=============================================================================
