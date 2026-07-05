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
	"math"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== AnalysisResponse building
//===
//=============================================================================

func RunAnalysis(ts *db.TradingSystem, curModel, selModel model.PositionModel, trades *[]db.Trade, pos *db.TradingPosition) (*AnalysisResponse,error) {
	baseline    := model.NewFixedUnitModel()
	ruinCapital := core.Trunc2d(pos.InitialCapital * (100 - pos.RuinPercentage) / 100)

	res := &AnalysisResponse{}
	res.TradingSystem = buildTradingSystem(ts)
	res.Params        = extractParameters(pos)
	res.ParamSpecs    = buildParamSpecs()
	res.ModelSpecs    = buildModelSpecs()
	res.UsedMargin    = ts.MarginValue

	if pos.MarginOverride != nil {
		res.UsedMargin = *pos.MarginOverride
	}

	grossRisk,errG := calcRiskValue(pos.RiskPerUnit, pos.RiskValue, trades, 0)
	netRisk,errN   := calcRiskValue(pos.RiskPerUnit, pos.RiskValue, trades, ts.CostPerOperation)

	res.NoLosses         = (errG != nil) || (errN != nil)
	res.GrossRisk        = grossRisk
	res.NetRisk          = netRisk
	res.RuinCapital      = ruinCapital
	res.CostPerOperation = ts.CostPerOperation

	res.Baseline      = calcAnalysisResult(baseline, trades, res)
	res.Current       = calcAnalysisResult(curModel, trades, res)
	res.Selected      = calcAnalysisResult(selModel, trades, res)

	return res,nil
}

//=============================================================================

func buildTradingSystem(ts *db.TradingSystem) *TradingSystem {
	return &TradingSystem{
		Id  : ts.Id,
		Name: ts.Name,
	}
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

func calcAnalysisResult(model model.PositionModel, trades *[]db.Trade, res *AnalysisResponse) *AnalysisResult {
	if model == nil {
		return nil
	}

	ar := &AnalysisResult{
		Model : NewModel(model.Name(), model.Config()),
	}

	if ! (trades == nil || len(*trades)==0 || res.NoLosses) {
		ar.Gross = calcModelPerformance(model, trades, res, 0, res.GrossRisk)
		ar.Net   = calcModelPerformance(model, trades, res, res.CostPerOperation, res.NetRisk)
	}

	return ar
}

//=============================================================================

func calcModelPerformance(mod model.PositionModel, trades *[]db.Trade, res *AnalysisResponse, costPerOper float64, risk float64) *ModelPerformance {
	var equity    []float64
	var positions []int

	params := res.Params

	snapshot := &model.TradingSnapshot{
		InitialCapital: *params.InitialCapital,
		CurrentCapital: *params.InitialCapital,
		RiskValue     : risk,
		AtrValue      : 0,
	}

	mod.PositionInit(snapshot)

	atrMap := make(map[int]float64) //TODO
	margin := res.UsedMargin
	ruined := false

	for _, trade := range *trades {
		snapshot.AtrValue = calcAtr(&trade, atrMap)
		position := mod.PositionFor(snapshot)

		//--- Check if we go above the max allowed position
		if position > *params.MaxUnits {
			position = *params.MaxUnits
		}

		//--- Check if we have enough margin to trade
		if margin * float64(position) >= snapshot.CurrentCapital {
			position = int(snapshot.CurrentCapital / margin)
		}

		currReturn := (trade.GrossReturn - 2 * costPerOper) * float64(position)
		snapshot.CurrentCapital += currReturn

		equity    = append(equity,    snapshot.CurrentCapital)
		positions = append(positions, position)

		if snapshot.CurrentCapital < res.RuinCapital  {
			ruined = true
			break
		}
	}

	drawdown, maxDrawdown := core.BuildDrawDown(&equity)

	ratio := 0.0
	if maxDrawdown != 0 {
		ratio = -((snapshot.CurrentCapital - snapshot.InitialCapital) / maxDrawdown)
	}

	return &ModelPerformance{
		Equity          : equity,
		Drawdown        : *drawdown,
		Positions       : positions,
		Return          : core.Trunc2d(snapshot.CurrentCapital),
		MaxDrawdown     : core.Trunc2d(maxDrawdown),
		ReturnDrawdRatio: core.Trunc2d(ratio),
		Ruined          : ruined,
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

func buildModelSpecs() map[string]any {
	specs := make(map[string]any)

	for k,v := range model.NewFixedUnitModel().Spec() {
		specs[k] = v
	}

	for k,v := range model.NewPercentRiskModel().Spec() {
		specs[k] = v
	}

	for k,v := range model.NewPercentVolatilityModel().Spec() {
		specs[k] = v
	}
	for k,v := range model.NewMarketMoneyModel().Spec() {
		specs[k] = v
	}

	return specs
}

//=============================================================================

func calcRiskValue(riskPerUnit db.RpuType, riskValue *float64, trades *[]db.Trade, costPerOper float64) (float64,error) {
	if riskPerUnit == db.RpuFixedValue {
		return *riskValue +2 * costPerOper,nil
	}

	returns := core.GetReturns(trades, db.TradeTypeAll, costPerOper)

	if riskPerUnit == db.RpuStopLoss {
		return core.CalcRisk(returns, costPerOper)
	}

	sum    := 0.0
	maxVal := 0.0

	var losses []float64
	for _, ret := range returns {
		if ret < 0 {
			ret    = -ret
			maxVal = math.Max(maxVal, ret)
			sum   += ret
			losses = append(losses, ret)
		}
	}

	if losses == nil || len(losses) == 0 {
		return 0, req.NewUnprocessableEntityError("no losses found")
	}

	if riskPerUnit == db.RpuMaxLoss {
		return maxVal, nil
	}

	if riskPerUnit == db.RpuAvgLoss {
		return sum / float64(len(losses)), nil
	}

	panic("Unknown riskPerUnit type: "+ riskPerUnit)
}

//=============================================================================

func calcAtr(t *db.Trade, atrMap map[int]float64) float64 {
	return 0 //TODO
}

//=============================================================================
