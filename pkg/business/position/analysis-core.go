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
	"log/slog"
	"time"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/portfolio-trader/pkg/platform"
	"github.com/algotiqa/types"
)

//=============================================================================

const AtrMaxBackDays = 5

//=============================================================================
//===
//=== AnalysisResponse building
//===
//=============================================================================

func RunAnalysis(ts *db.TradingSystem, curModel, selModel model.PositionModel, trades *[]db.Trade, pos *db.TradingPosition) (*AnalysisResponse,error) {
	res := &AnalysisResponse{}
	res.TradingSystem = buildTradingSystem(ts)
	res.Params        = extractParameters(pos)
	res.ParamSpecs    = buildParamSpecs()
	res.ModelSpecs    = buildModelSpecs()
	res.UsedMargin    = ts.MarginValue

	if pos.MarginOverride != nil {
		res.UsedMargin = *pos.MarginOverride
	}

	grossRisk,errG := core.CalcRiskForPosition(pos.RiskPerUnit, pos.RiskValue, trades, 0)
	netRisk,errN   := core.CalcRiskForPosition(pos.RiskPerUnit, pos.RiskValue, trades, ts.CostPerOperation)

	baseline := model.NewFixedUnitModel()

	res.NoLosses         = (errG != nil) || (errN != nil)
	res.GrossRisk        = grossRisk
	res.NetRisk          = netRisk
	res.CostPerOperation = ts.CostPerOperation

	var errC, errS error

	res.Baseline,_    = calcAnalysisResult(ts, baseline, trades, res)
	res.Current, errC = calcAnalysisResult(ts, curModel, trades, res)
	res.Selected,errS = calcAnalysisResult(ts, selModel, trades, res)
	res.Time          = calcTime(trades)

	if errC != nil {
		return nil, errC
	}

	if errS != nil {
		return nil, errS
	}

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
		InitialCapital : &p.InitialCapital,
		MaxTolDrawdPerc: &p.MaxTolDrawdPerc,
		MarginOverride : p.MarginOverride,
		MaxUnits       : &p.MaxUnits,
		RiskPerUnit    : p.RiskPerUnit,
		RiskValue      : p.RiskValue,
	}
}

//=============================================================================

func calcAnalysisResult(ts *db.TradingSystem, model model.PositionModel, trades *[]db.Trade, res *AnalysisResponse) (*AnalysisResult,error) {
	if model == nil {
		return nil,nil
	}

	ar := &AnalysisResult{
		Model : NewModel(model.Name(), model.Config()),
	}

	var atrMap map[types.Date]float64
	if model.AtrLen() > 0 {
		data,err := platform.GetLatestDataProductAnalysis(ts.Username, ts.DataProductId,365*100, model.AtrLen(), 1440)
		if err != nil {
			return nil,req.NewServerErrorByError(err)
		}
		atrMap = buildAtrMap(data.BarResults)
	}

	var errG,errN error

	if ! (trades == nil || len(*trades)==0 || res.NoLosses) {
		ar.Gross,errG = calcModelPerformance(model, trades, res, 0,                    res.GrossRisk, atrMap, ts.PointValue)
		ar.Net  ,errN = calcModelPerformance(model, trades, res, res.CostPerOperation, res.NetRisk,   atrMap, ts.PointValue)

		if errG != nil || errN != nil {
			//--- ErrG and errN are the same, we just need to check one
			return nil,errG
		}
	}

	return ar,nil
}

//=============================================================================

func calcModelPerformance(mod model.PositionModel, trades *[]db.Trade, res *AnalysisResponse, costPerOper float64,
						  risk float64, atrMap map[types.Date]float64, pointValue float64) (*ModelPerformance,error) {
	var equity    []float64
	var positions []int
	var err       error

	params := res.Params

	snapshot := &model.TradingSnapshot{
		InitialCapital: *params.InitialCapital,
		CurrentCapital: *params.InitialCapital,
		RiskValue     : risk,
		AtrValue      : 0,
		PointValue    : pointValue,
	}

	mod.PositionInit(snapshot)

	margin := res.UsedMargin
	ruined := false

	for _, trade := range *trades {
		snapshot.AtrValue,err = calcAtr(&trade, atrMap)
		if err != nil {
			return nil,err
		}

		position := int(mod.PositionFor(snapshot))

		//--- Check if we go above the max allowed position
		if position > *params.MaxUnits {
			position = *params.MaxUnits
		}

		//--- Check if we have enough margin to trade
		if margin * float64(position) >= snapshot.CurrentCapital {
			position = int(snapshot.CurrentCapital / margin)
		}

		currReturn := core.NetReturn(trade.GrossReturn, costPerOper) * float64(position)
		snapshot.CurrentCapital += currReturn

		equity    = append(equity,    snapshot.CurrentCapital)
		positions = append(positions, position)

		if snapshot.CurrentCapital < margin {
			ruined = true
			break
		}
	}

	drawdownPerc, maxDrawdownPerc, maxDrawdown := core.BuildDrawDownPercentage(&equity)

	finalReturn := snapshot.CurrentCapital - snapshot.InitialCapital
	retOnAcc    := finalReturn * 100 / snapshot.InitialCapital
	ratio       := 0.0

	if maxDrawdownPerc != 0 {
		ratio = -(retOnAcc / maxDrawdownPerc)
	}

	return &ModelPerformance{
		Equity          : equity,
		DrawdownPerc    : *drawdownPerc,
		Positions       : positions,
		Return          : finalReturn,
		MaxDrawdown     : maxDrawdown,
		MaxDrawdownPerc : maxDrawdownPerc,
		ReturnOnAccount : core.Trunc2d(retOnAcc),
		ReturnDrawdRatio: core.Trunc2d(ratio),
		Ruined          : ruined,
	},nil
}

//=============================================================================

func buildParamSpecs() map[string]any {
	specs := make(map[string]any)

	specs[SpecInitialCapital.Name]  = SpecInitialCapital
	specs[SpecMaxTolDrawdPerc.Name] = SpecMaxTolDrawdPerc
	specs[SpecMarginOverride.Name]  = SpecMarginOverride
	specs[SpecMaxUnits.Name]        = SpecMaxUnits
	specs[SpecRiskPerUnit.Name]     = SpecRiskPerUnit
	specs[SpecRiskValue.Name]       = SpecRiskValue

	return specs
}

//=============================================================================

func buildModelSpecs() map[string]any {
	specs := make(map[string]any)

	for k,v := range model.NewFixedUnitModel().Specs() {
		specs[k] = v
	}

	for k,v := range model.NewPercentRiskModel().Specs() {
		specs[k] = v
	}

	for k,v := range model.NewPercentVolatilityModel().Specs() {
		specs[k] = v
	}
	for k,v := range model.NewMarketMoneyModel().Specs() {
		specs[k] = v
	}

	return specs
}

//=============================================================================

func calcAtr(t *db.Trade, atrMap map[types.Date]float64) (float64,error) {
	if atrMap == nil {
		return 0,nil
	}

	//--- If we need to calculate the position for today, we have to use the ATR value
	//--- of yesterday, because today is not closed yet. This is why we are subtracting 1

	origDate := types.ToDate(t.EntryDate)
	posDate  := origDate.AddDays(-1)

	for i := 0; i<AtrMaxBackDays; i++ {
		atr, ok := atrMap[posDate]
		if ok {
			return atr,nil
		}

		posDate = posDate.AddDays(-1)
	}

	return 0,req.NewUnprocessableEntityError("missing ATR for date: %v", origDate.String())
}

//=============================================================================

func calcTime(trades *[]db.Trade) []time.Time {
	var res []time.Time

	for _, t := range *trades {
		res = append(res, *t.ExitDate)
	}

	return res
}

//=============================================================================

func buildAtrMap(days []*platform.BarResult) map[types.Date]float64 {
	atrMap := make(map[types.Date]float64)

	for _, br := range days {
		date := types.ToDate(&br.Time)
		atr  := br.Atr
slog.Info("Timezone: "+br.Time.Location().String()+", Date: "+date.String())
		atrMap[date] = atr
	}

	return atrMap
}

//=============================================================================
