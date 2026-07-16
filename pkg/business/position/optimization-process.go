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
	"math/rand"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/core/stats"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== OptimizationProcess
//===
//=============================================================================

type OptimizationProcess struct {
	ts         *db.TradingSystem
	trades     *[]db.Trade
	optReq     *OptimizationRequest
	info       *OptimizationInfo
	stopping    bool
}

//=============================================================================

type RunSpec struct {
	InitialCapital    float64
	MaxTolDrawdPerc   float64
	MaxUnits          int
	UsedMargin        float64
	RiskValue         float64
	TradesPerYear     int
	TradesToGenerate  int
	CostPerOper       float64
	RunsPerSimul      int
}

//=============================================================================

type RunStepResult struct {
	FinalCapital     float64
	Return           float64
	ReturnOnAccount  float64
	MaxDrawdownPerc  float64
	ReturnDrawdRatio float64
	Failed           bool
}

//=============================================================================

type ExecutionResult struct {
	InitialCapital        float64        `json:"initialCapital"`
	MedianFinalCapital    float64        `json:"medianFinalCapital"`
	MedianReturn          float64        `json:"medianReturn"`
	MedianReturnOnAcc     float64        `json:"medianReturnOnAcc"`
	MedianMaxDrawdownPerc float64        `json:"medianMaxDrawdownPerc"`
	MedianRetDrawdRatio   float64        `json:"medianRetDrawdRatio"`
	ProbOfSuccess         float64        `json:"probOfSuccess"`
	ProbOfFailure         float64        `json:"probOfFailure"`
	Model                 string         `json:"model"`
	Config                map[string]any `json:"config"`
}

//=============================================================================

func (er *ExecutionResult) FitnessValue() float64 {
	return er.ProbOfSuccess - er.ProbOfFailure
}

//=============================================================================
//===
//=== Public methods
//===
//=============================================================================

func (op *OptimizationProcess) Start() error {
	netRisk,err := calcRiskValue(op.optReq.Params.RiskPerUnit, op.optReq.Params.RiskValue, op.trades, op.ts.CostPerOperation)
	if err != nil {
		return req.NewUnprocessableEntityError(err.Error())
	}

	margin := op.ts.MarginValue

	if op.optReq.Params.MarginOverride != nil {
		margin = *op.optReq.Params.MarginOverride
	}

	tradesPerY, totTrades := calcTradesToGenerate(op.trades, op.optReq.RunConfig)

	spec := &RunSpec{
		InitialCapital   : *op.optReq.Params.InitialCapital,
		MaxTolDrawdPerc  : *op.optReq.Params.MaxTolDrawdPerc,
		MaxUnits         : *op.optReq.Params.MaxUnits,
		UsedMargin       : margin,
		RiskValue        : netRisk,
		TradesPerYear    : tradesPerY,
		TradesToGenerate : totTrades,
		CostPerOper      : op.ts.CostPerOperation,
		RunsPerSimul     : *op.optReq.RunConfig.RunsPerSimul,
	}

	totalSteps := op.calcTotalSteps()
	op.info = NewOptimizationInfo(spec, op.optReq.ModelConfig, op.optReq.Targets, totalSteps, MaxResultSize)

	go op.optimize(spec)
	return nil
}

//=============================================================================

func (op *OptimizationProcess) Stop() {
	slog.Info("Stop: Stopping position optimization process", "tsId", op.ts.Id)
	op.stopping = true
}

//=============================================================================

func (op *OptimizationProcess) GetInfo() *OptimizationInfo {
	return op.info
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func calcTradesToGenerate(trades *[]db.Trade, config *RunConfig) (int, int) {
	if config.TradesPerYear != nil {
		return *config.TradesPerYear, *config.TradesPerYear * *config.ProjectionYears
	}

	firstDate := (*trades)[0].ExitDate
	lastDate  := (*trades)[len(*trades) -1].ExitDate

	days := lastDate.Sub(*firstDate).Hours() / 24

	tradesPerYear := int(float64(len(*trades)) * 365 / days)
	if tradesPerYear < 1 {
		tradesPerYear = 1
	}

	return tradesPerYear, tradesPerYear* *config.ProjectionYears
}

//=============================================================================

func (op *OptimizationProcess) calcTotalSteps() uint {
	count := uint(0)

	mc := op.optReq.ModelConfig

	if mc.EnableFixedUnit {
		count += uint(len(*(mc.Units.Steps())))
	}

	if mc.EnablePercentRisk {
		count += uint(len(*(mc.RiskPerTrade.Steps())))
	}

	if mc.EnablePercentVol {
		count += uint(len(*(mc.AverageLength.Steps()))) * uint(len(*(mc.MaxVolatility.Steps())))
	}

	if mc.EnableMarketMoney {
		count += uint(len(*(mc.RiskPerTradeOnCap.Steps()))) * uint(len(*(mc.RiskPerTradeOnEarn.Steps()))) * uint(len(*(mc.PercentageOnCapital.Steps())))
	}

	return count
}

//=============================================================================
//--- GoRoutine

func (op *OptimizationProcess) optimize(spec *RunSpec) {
	slog.Info("optimize: Started position optimization", "tsId", op.ts.Id, "tsName", op.ts.Name)

	if !op.stopping {
		op.runFixedUnit(spec)
		if !op.stopping {
			op.runPercentRisk(spec)
			if !op.stopping {
				op.runPercentVol(spec)
				if !op.stopping {
					op.runMarketMoney(spec)
				}
			}
		}
	}

	op.info.finish(op.stopping)
	slog.Info("optimize: Position optimization complete.")
}

//=============================================================================

func (op *OptimizationProcess) runFixedUnit(spec *RunSpec) {
	mc := op.optReq.ModelConfig

	if mc.EnableFixedUnit {
		for _, units := range *mc.Units.Steps() {
			mod := model.NewFixedUnitModelWithParams(units)
			er  := op.executeSimulation(mod, spec)

			op.info.addResult(er, op.optReq.Targets)

			if op.stopping {
				break
			}
		}
	}
}

//=============================================================================

func (op *OptimizationProcess) runPercentRisk(spec *RunSpec) {
	mc := op.optReq.ModelConfig

	if mc.EnablePercentRisk {
		for _, riskPerTrade := range *mc.RiskPerTrade.Steps() {
			mod := model.NewPercentRiskModelWithParams(riskPerTrade)
			er  := op.executeSimulation(mod, spec)

			op.info.addResult(er, op.optReq.Targets)

			if op.stopping {
				break
			}
		}
	}
}

//=============================================================================

func (op *OptimizationProcess) runPercentVol(spec *RunSpec) {
	mc := op.optReq.ModelConfig

	if mc.EnablePercentVol {
		for _, avgLen := range *mc.AverageLength.Steps() {
			//--- Load ATR data here...

			for _, maxVol := range *mc.MaxVolatility.Steps() {
				mod := model.NewPercentVolatilityModelWithParams(avgLen, maxVol)
				er  := op.executeSimulation(mod, spec)

				op.info.addResult(er, op.optReq.Targets)

				if op.stopping {
					break
				}
			}
		}
	}
}

//=============================================================================

func (op *OptimizationProcess) runMarketMoney(spec *RunSpec) {
	mc := op.optReq.ModelConfig

	if mc.EnableMarketMoney {
		for _, riskOnCap := range *mc.RiskPerTradeOnCap.Steps() {
			for _, riskOnEarn := range *mc.RiskPerTradeOnEarn.Steps() {
				for _, percOnCap := range *mc.PercentageOnCapital.Steps() {
					mod := model.NewMarketMoneyModelWithParams(riskOnCap, riskOnEarn, model.McTypePercOnCapital, percOnCap)
					er  := op.executeSimulation(mod, spec)

					op.info.addResult(er, op.optReq.Targets)

					if op.stopping {
						break
					}
				}
			}
		}
	}
}

//=============================================================================

func (op *OptimizationProcess) executeSimulation(mod model.PositionModel, spec *RunSpec) *ExecutionResult {
	failCount := 0.0
	succCount := 0.0

	var finalCapitals  []float64
	var returns        []float64
	var returnOnAccs   []float64
	var maxDrawdPercs  []float64
	var retDrawdRatios []float64

	for i:=0; i<spec.RunsPerSimul; i++ {
		res := op.runStep(mod, spec)

		finalCapitals = append(finalCapitals , res.FinalCapital)
		returns       = append(returns       , res.Return)
		returnOnAccs  = append(returnOnAccs  , res.ReturnOnAccount)
		maxDrawdPercs = append(maxDrawdPercs , res.MaxDrawdownPerc)
		retDrawdRatios= append(retDrawdRatios, res.ReturnDrawdRatio)

		if res.Failed {
			failCount++
		}

		if res.ReturnOnAccount >= *op.optReq.Targets.DesReturnOnAccount {
			succCount++
		}
	}

	propOfFail := failCount / float64(spec.RunsPerSimul)
	probOfSucc := succCount / float64(spec.RunsPerSimul)

	return &ExecutionResult{
		InitialCapital       : spec.InitialCapital,
		MedianFinalCapital   : stats.Median(finalCapitals),
		MedianReturn         : stats.Median(returns),
		MedianReturnOnAcc    : core.Trunc2d(stats.Median(returnOnAccs)),
		MedianMaxDrawdownPerc: core.Trunc2d(stats.Median(maxDrawdPercs)),
		MedianRetDrawdRatio  : core.Trunc2d(stats.Median(retDrawdRatios)),
		ProbOfFailure        : core.Trunc2d(propOfFail * 100),
		ProbOfSuccess        : core.Trunc2d(probOfSucc * 100),
		Model                : string(mod.Name()),
		Config               : mod.Config(),
	}
}

//=============================================================================

func (op *OptimizationProcess) runStep(mod model.PositionModel, spec *RunSpec) *RunStepResult {
	trades := *op.trades

	snapshot := &model.TradingSnapshot{
		InitialCapital: spec.InitialCapital,
		CurrentCapital: spec.InitialCapital,
		RiskValue     : spec.RiskValue,
		AtrValue      : 0,
	}

	mod.PositionInit(snapshot)

	maxEquity := 0.0
	atrMap    := make(map[int]float64) //TODO
	failed    := false
	size      := len(trades)
	margin    := spec.UsedMargin

	drawdownPerc    := 0.0
	maxDrawdownPerc := 0.0

	//--- Main loop -----------------------------------

	for i:=0; i<spec.TradesToGenerate; i++ {
		trade := trades[rand.Intn(size)]

		snapshot.AtrValue = calcAtr(&trade, atrMap)
		position := mod.PositionFor(snapshot)

		//--- Check if we go above the max allowed position
		if position > spec.MaxUnits {
			position = spec.MaxUnits
		}

		//--- Check if we have enough margin to trade
		if margin * float64(position) >= snapshot.CurrentCapital {
			position = int(snapshot.CurrentCapital / margin)
		}

		currReturn := (trade.GrossReturn - 2 * spec.CostPerOper) * float64(position)
		snapshot.CurrentCapital += currReturn

		if snapshot.CurrentCapital > maxEquity {
			maxEquity = snapshot.CurrentCapital
		} else {
			drawdownPerc = 0
			if maxEquity > 0 {
				drawdownPerc = 1 - snapshot.CurrentCapital / maxEquity

				if drawdownPerc > maxDrawdownPerc {
					maxDrawdownPerc = drawdownPerc
				}
			}
		}

		if snapshot.CurrentCapital < margin || maxDrawdownPerc * 100 > spec.MaxTolDrawdPerc {
			failed = true
			break
		}
	}

	//--- Return metrics --------------------------------------

	finalReturn := snapshot.CurrentCapital - snapshot.InitialCapital
	retOnAcc    := finalReturn / snapshot.InitialCapital
	ratio       := 0.0

	if maxDrawdownPerc != 0 {
		ratio = retOnAcc / maxDrawdownPerc
	}

	return &RunStepResult{
		FinalCapital    : snapshot.CurrentCapital,
		Return          : finalReturn,
		MaxDrawdownPerc : maxDrawdownPerc *100,
		ReturnOnAccount : retOnAcc * 100,
		ReturnDrawdRatio: ratio,
		Failed          : failed,
	}
}

//=============================================================================
