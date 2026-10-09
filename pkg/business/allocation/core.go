//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package allocation

import (
	"encoding/json"
	"strconv"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter"
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/portfolio-trader/pkg/platform"
	"golang.org/x/exp/slices"
)

//=============================================================================
//===
//=== Public
//===
//=============================================================================

func (j *Job) BuildPortfolioAllocation() *Report {
	var survivors []*SystemPosition

	if len(j.tradingSystems) == 0 {
		j.log(LogLevelError, "No trading systems assigned to portfolio")
	} else {
		activeSystems := j.calcFilterOutcomes()
		survivors      = j.setupProcess(activeSystems)

		j.log(LogLevelInfo, "Healthy systems passing the filter: %v/%v", len(activeSystems), len(j.tradingSystems))
		j.log(LogLevelInfo, "Available capital: %v %v", j.portfolio.AvailableCapital(), j.portfolio.AccountCurrencyCode)

		for {
			if len(survivors) == 0 {
				j.log(LogLevelWarning, "There are no trading systems passing the checks")
				break
			}

			j.calcSystemCorrelation(survivors)
			newSurvivors := j.calcPosition(survivors)

			if len(newSurvivors) == len(survivors) {
				j.buildCorrelationMatrix(survivors)
				break
			}

			j.log(LogLevelInfo, "Survivors decreased (%v --> %v). Looping again with less systems", len(survivors), len(newSurvivors))
			survivors = newSurvivors
		}
	}

	return j.buildReport(survivors)
}

//=============================================================================
//===
//=== Private
//===
//=============================================================================

func (j *Job) calcFilterOutcomes() []*TradingSystemInfo {
	var list []*TradingSystemInfo

	for _, tsi := range j.tradingSystems {
		if tsi.healthy {
			tsi.filterPassed = true
			if tsi.filter != nil {
				tsi.filterPassed = filter.CalcActivation(tsi.system, tsi.filter, *tsi.trades)
			}

			//TODO: Add filter details
			j.addFilter(tsi, tsi.filterPassed, "")

			if tsi.filterPassed {
				list = append(list, tsi)
			}
		}
	}

	return list
}

//=============================================================================

func (j *Job) setupProcess(activeSystems []*TradingSystemInfo) []*SystemPosition {
	var survivors []*SystemPosition

	for _, tsi := range activeSystems {
		pos := tsi.position
		mod,err := buildModel(pos)
		if err != nil {
			j.addExcludedSystem(tsi, "Error building model: "+ err.Error())
			continue
		}

		risk,err := core.CalcRiskForPosition(pos.RiskPerUnit, pos.RiskValue, tsi.trades, tsi.system.CostPerOperation)
		if err != nil {
			j.addExcludedSystem(tsi, "Cannot detect risk for '%v'"+ tsi.system.Name +" ("+err.Error()+")")
			continue
		}

		atrVal := 0.0
		atrLen := mod.AtrLen()
		if atrLen > 0 {
			res,errp := platform.GetLatestDataProductAnalysis(tsi.system.DataProductId,j.spec.TradeYears*365, atrLen, 1440)
			if errp != nil {
				j.addExcludedSystem(tsi,"Cannot get ATR: "+errp.Error())
				continue
			}

			if len(res.BarResults) == 0 {
				j.addExcludedSystem(tsi,"No ATR values")
				continue
			}

			atrVal = res.BarResults[len(res.BarResults) -1].Atr

			j.log(LogLevelInfo, "For '%v' got an ATR value of '%v' ('%v' %v)", tsi.system.Name, atrVal,
								atrVal*tsi.system.PointValue, tsi.system.CurrencyCode)
		}

		sp := NewSystemPosition(tsi, mod, risk, atrVal)
		survivors = append(survivors, sp)
	}

	return survivors
}

//=============================================================================

func (j *Job) calcSystemCorrelation(list []*SystemPosition) {
	j.clearSystemCorrelations()

	if len(list) < 2 {
		return
	}

	for i := 0; i < len(list) -1; i++ {
		for k := i + 1; k < len(list); k++ {
			si := list[i].tsi
			sk := list[k].tsi
			corr,err := core.CalcCorrelation(si.returns, sk.returns)
			if err != nil {
				j.log(LogLevelWarning, "Cannot calc correlation between '%v' and '%v'. Assuming '0'. Error: %v",
					  si.system.Name, sk.system.Name, err.Error())
				corr = 0
			}

			j.addSystemCorrelation(si, sk, corr)
		}
	}
}

//=============================================================================

var one = 1.0

func (j *Job) buildCorrelationMatrix(activeSystems []*SystemPosition) {
	matrix := &CorrelationMatrix{}

	//--- Build lookup map

	var idMap = make(map[uint]int)

	for i, sp := range activeSystems {
		idMap[sp.tsi.system.Id] = i
		matrix.Names = append(matrix.Names, sp.tsi.system.Name)
	}

	//--- Allocate full matrix, setting identity on diagonal

	size := len(activeSystems)

	for i,_ := range activeSystems {
		matrix.Cells = append(matrix.Cells, make([]*float64, size))
		matrix.Cells[i][i] = &one
	}

	//--- Fill upper diagonal with data

	for _,scf := range j.correlations {
		index1, ok1 := idMap[scf.Ts1Id]
		index2, ok2 := idMap[scf.Ts2Id]

		if !ok1 || !ok2 {
			j.log(LogLevelError, "Map lookup failure!", "id1", scf.Ts1Id, "id2", scf.Ts2Id)
		} else {
			if index1 > index2 {
				aux    := index1
				index1 = index2
				index2 = aux
			}

			matrix.Cells[index1][index2] = &scf.Correlation
			matrix.Cells[index2][index1] = &scf.Correlation
		}
	}

	j.correlationMatrix = matrix
}

//=============================================================================

func (j *Job) calcPosition(activeSystems []*SystemPosition) []*SystemPosition {
	var survivors []*SystemPosition
	p := j.portfolio

	for _, sp := range activeSystems {
		tsi := sp.tsi

		snapshot := &model.TradingSnapshot{
			InitialCapital: p.AvailableCapital(),
			CurrentCapital: p.AvailableCapital(),
			RiskValue     : sp.Risk,
			AtrValue      : sp.AtrValue,
			PointValue    : sp.tsi.system.PointValue,
		}

		position := sp.mod.PositionFor(snapshot)
		sp.InitialPosition = position
		position = j.builder.TunePosition(position, j.getCorrelationsForSystem(tsi.system.Id))
		intPosition := int(position)

		//--- Check if we go above the max allowed position

		pos := sp.tsi.position

		if intPosition > pos.MaxUnits {
			j.log(LogLevelInfo, "Clamped position for '%v' : %v --> %v (max units exceeded)", tsi.system.Name, intPosition, pos.MaxUnits)
			intPosition = pos.MaxUnits
			position = float64(pos.MaxUnits)
		}

		//--- Check if we have enough margin to trade

		margin := tsi.system.MarginValue
		if pos.MarginOverride != nil {
			margin = *pos.MarginOverride
		}

		maxCapitalUnderMargin := snapshot.CurrentCapital * p.MaxMarginPerc/100
		if margin * float64(intPosition) >= maxCapitalUnderMargin {
			j.log(LogLevelInfo, "Clamped position for '%v' : %v --> %v (margin exceeded)", tsi.system.Name, intPosition, pos.MaxUnits)
			position = snapshot.CurrentCapital / margin
			intPosition = int(position)
		}

		if intPosition == 0 {
			sPos := strconv.FormatFloat(position, 'f', 3, 64)
			j.addExcludedSystem(tsi,"Position below 1: "+ sPos)
			continue
		}

		sp.FinalPosition = intPosition
		survivors = append(survivors, sp)
	}

	return survivors
}

//=============================================================================

func (j *Job) getCorrelationsForSystem(id uint) []float64 {
	var list []float64

	for _, sc := range j.correlations {
		if sc.Ts1Id == id || sc.Ts2Id == id {
			list = append(list, sc.Correlation)
		}
	}

	//--- Sort in descending mode
	slices.SortFunc[[]float64](list, func(a,b float64) int {
		if a < b {
			return 1
		}
		if a > b {
			return -1
		}
		return 0
	})

	return list
}

//=============================================================================

func (j *Job) buildReport(positions []*SystemPosition) *Report {
	//--- Sorts in descending order
	corr := j.correlations
	slices.SortFunc[[]*SystemCorrelation](corr, func(a,b *SystemCorrelation) int {
		if a.Correlation < b.Correlation {
			return 1
		}
		if a.Correlation > b.Correlation {
			return -1
		}
		return 0
	})

	//--- Remove some useless decimals

	for _, sp := range positions {
		sp.InitialPosition = core.Trunc2d(sp.InitialPosition)
	}

	return &Report{
		TradingSystems   : j.healthySystems,
		FilterOutcomes   : j.filterOutcomes,
		Correlations     : j.correlations,
		SystemPositions  : positions,
		ExcludedSystems  : j.excludedSystems,
		SystemActions    : j.buildSystemActions(),
		Logs             : j.logs,
		CorrelationMatrix: j.correlationMatrix,
	}
}

//=============================================================================

func (j *Job) buildSystemActions() []*SystemAction {
	var list []*SystemAction

	for _,fo := range j.filterOutcomes {
		if fo.Action == ActionTurnOn || fo.Action == ActionTurnOff {
			sa := &SystemAction{
				Id     : fo.TsId,
				Name   : fo.TsName,
				Action : fo.Action,
				Message: "Outcome from the filter",
			}
			list = append(list, sa)
		}
	}

	for _,es := range j.excludedSystems {
		sa := &SystemAction{
			Id     : es.Id,
			Name   : es.Name,
			Action : ActionTurnOff,
			Message: "System excluded: "+es.Name,
		}
		list = append(list, sa)
	}

	return list
}

//=============================================================================
//===
//=== General private functions
//===
//=============================================================================

func buildModel(pos *db.TradingPosition) (model.PositionModel,error) {
	var err error
	var mod model.PositionModel

	config := make(map[string]any)
	err = json.Unmarshal([]byte(pos.Config), &config)
	if err == nil {
		mod,err = model.New(pos.Model, config)
	}

	return mod, err
}

//=============================================================================

//func allocatePortfolio(job *AllocationJob) {
//	remaining, bestTs, bestIdx := getFilteredSystems(job.tradingSystems)
//	remaining = removeSystem(remaining, bestIdx)
//
//	currPor := NewAllocatedPortfolio()
//	currPor.Add(bestTs)
//	currSr := currPor.CalcSharpeRatio()
//
//	for ;len(remaining)>0; {
//		bestSr := -1.0
//		bestTs = nil
//
//		for i, candidate := range remaining {
//			newPor := currPor.Clone()
//			newPor.Add(candidate)
//			newSR := newPor.CalcSharpeRatio()
//
//			if newSR > bestSr {
//				bestSr  = newSR
//				bestTs  = candidate
//				bestIdx = i
//			}
//		}
//
//		//--- We exit if the bestTS doesn't provide a sharp increase in the SR
//
//		if bestSr * 0.9 < currSr {
//			break
//		}
//
//		//--- Ok, makes sense to add bestTS to the portfolio
//
//		currPor.Add(bestTs)
//		currSr = bestSr
//		remaining = removeSystem(remaining, bestIdx)
//	}
//
//	job.allocatedPortfolio = currPor
//}

//=============================================================================

//func getFilteredSystems(list []*TradingSystemInfo) ([]*TradingSystemInfo, *TradingSystemInfo, int) {
//	var ret []*TradingSystemInfo
//	var bestTs *TradingSystemInfo
//	var bestIdx int
//
//	for _, tsi := range list {
//		if tsi.filterPassed {
//			ret = append(ret, tsi)
//			if bestTs == nil || tsi.sharpeRatio > bestTs.sharpeRatio {
//				bestTs  = tsi
//				bestIdx = len(ret) -1
//			}
//		}
//	}
//
//	return ret, bestTs, bestIdx
//}

//=============================================================================

//func removeSystem(list []*TradingSystemInfo, index int) []*TradingSystemInfo {
//	return append(list[:index], list[index+1:]...)
//}

//=============================================================================
