//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package performance

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/core/stats"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
)

//=============================================================================

func GetPerformanceAnalysis(ts *db.TradingSystem, trades *[]db.Trade, livePeriods *[]db.LivePeriod) *AnalysisResponse {
	res := AnalysisResponse{}
	res.TradingSystem = ts
	res.Trades        = trades

	allEq,   allMaxGrossDD,   allMaxNetDD   := calcEquities(ts, trades, db.TradeTypeAll)
	longEq,  longMaxGrossDD,  longMaxNetDD  := calcEquities(ts, trades, db.TradeTypeLong)
	shortEq, shortMaxGrossDD, shortMaxNetDD := calcEquities(ts, trades, db.TradeTypeShort)

	res.AllEquities   = allEq
	res.LongEquities  = longEq
	res.ShortEquities = shortEq

	res.Gross.MaxDrawdown.Total = allMaxGrossDD
	res.Gross.MaxDrawdown.Long  = longMaxGrossDD
	res.Gross.MaxDrawdown.Short = shortMaxGrossDD
	res.Net  .MaxDrawdown.Total = allMaxNetDD
	res.Net  .MaxDrawdown.Long  = longMaxNetDD
	res.Net  .MaxDrawdown.Short = shortMaxNetDD

	res.Gross.Return.Total = calcProfit(allEq  .GrossEquity)
	res.Gross.Return.Long  = calcProfit(longEq .GrossEquity)
	res.Gross.Return.Short = calcProfit(shortEq.GrossEquity)
	res.Net  .Return.Total = calcProfit(allEq  .NetEquity)
	res.Net  .Return.Long  = calcProfit(longEq .NetEquity)
	res.Net  .Return.Short = calcProfit(shortEq.NetEquity)

	res.Gross.AverageTrade.Total = calcAvgTrade(res.Gross.Return.Total, allEq  .Trades)
	res.Gross.AverageTrade.Long  = calcAvgTrade(res.Gross.Return.Long,  longEq .Trades)
	res.Gross.AverageTrade.Short = calcAvgTrade(res.Gross.Return.Short, shortEq.Trades)
	res.Net  .AverageTrade.Total = calcAvgTrade(res.Net  .Return.Total, allEq  .Trades)
	res.Net  .AverageTrade.Long  = calcAvgTrade(res.Net  .Return.Long,  longEq .Trades)
	res.Net  .AverageTrade.Short = calcAvgTrade(res.Net  .Return.Short, shortEq.Trades)

	calcAggregates(&res)
	updateGeneralInfo(&res)
	calcDistributions(&res)
	calcRolling(&res)
	calcLivePeriods(&res, livePeriods)

	return &res
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcEquities(ts *db.TradingSystem, trades *[]db.Trade, tradeType string) (*Equities, float64, float64) {
	timeSlice, grossProfits := core.BuildGrossProfits(trades, tradeType)
	netProfits := core.BuildNetProfits(grossProfits, ts.CostPerOperation)

	grossEquity := core.BuildEquity(grossProfits)
	netEquity   := core.BuildEquity(netProfits)

	grossDD, maxGrossDD := core.BuildDrawDown(grossEquity)
	netDD,   maxNetDD   := core.BuildDrawDown(netEquity)

	return &Equities{
		Time:          timeSlice,
		GrossEquity:   grossEquity,
		NetEquity:     netEquity,
		GrossDrawdown: grossDD,
		NetDrawdown:   netDD,
		Trades:        len(*timeSlice),
	}, maxGrossDD, maxNetDD
}

//=============================================================================

func calcProfit(equity *[]float64) float64 {
	if equity == nil || len(*equity) == 0 {
		return 0
	}

	return (*equity)[len(*equity)-1]
}

//=============================================================================

func calcAvgTrade(value float64, count int) float64 {
	if count == 0 {
		return 0
	}

	return core.Trunc2d(value / float64(count))
}

//=============================================================================
//=== Timezone shifting
//=============================================================================

func calcAggregates(res *AnalysisResponse) {
	calcYearAggregates(res)
}

//=============================================================================

func calcYearAggregates(res *AnalysisResponse) {
	cost := float64(res.TradingSystem.CostPerOperation)
	list := []*AnnualAggregate{}

	var currYear *AnnualAggregate

	for _, tr := range *res.Trades {
		if currYear == nil {
			//--- Beginning of a new year

			currYear = NewAggregate(&tr, cost)
			list = append(list, currYear)
		} else {
			if currYear.Year == tr.ExitDate.Year() {
				//--- Continue on the current year
				currYear.addTrade(&tr, cost)
			} else {
				//--- Continue on the new year

				currYear.consolidate()
				currYear = NewAggregate(&tr, cost)
				list = append(list, currYear)
			}
		}
	}

	if currYear != nil {
		currYear.consolidate()
	}

	res.Aggregates.Annual = &list
}

//=============================================================================
//=== General information
//=============================================================================

func updateGeneralInfo(res *AnalysisResponse) {
	calcFromToDates(res)
}

//=============================================================================

func calcFromToDates(res *AnalysisResponse) {
	numTrades := len(*res.Trades)
	if numTrades > 0 {
		firstTrade := (*res.Trades)[0].ExitDate
		lastTrade := (*res.Trades)[numTrades-1].ExitDate

		res.General.FromDate = NewDate(firstTrade)
		res.General.ToDate   = NewDate(lastTrade)
	}
}

//=============================================================================

func NewDate(t *time.Time) types.Date {
	return types.Date(t.Year()*10000 + int(t.Month())*100 + t.Day())
}

//=============================================================================
//=== Metrics
//=============================================================================

func calcDistributions(res *AnalysisResponse) {
	dist := &res.Distributions

	//--- All (gross + net)

	_, allGross := core.BuildGrossProfits(res.Trades, db.TradeTypeAll)
	allNet := core.BuildNetProfits(allGross, res.TradingSystem.CostPerOperation)

	dist.TradesAllGross = calcDistribution(*allGross)
	dist.TradesAllNet   = calcDistribution(*allNet)

	//--- Long (gross + net)

	_, longGross := core.BuildGrossProfits(res.Trades, db.TradeTypeLong)
	longNet := core.BuildNetProfits(longGross, res.TradingSystem.CostPerOperation)

	dist.TradesLongGross = calcDistribution(*longGross)
	dist.TradesLongNet   = calcDistribution(*longNet)

	//--- Short (gross + net)

	_, shortGross := core.BuildGrossProfits(res.Trades, db.TradeTypeShort)
	shortNet := core.BuildNetProfits(shortGross, res.TradingSystem.CostPerOperation)

	dist.TradesShortGross = calcDistribution(*shortGross)
	dist.TradesShortNet   = calcDistribution(*shortNet)
}

//=============================================================================

func calcDistribution(data []float64) *Distribution {
	if data == nil || len(data) == 0 {
		return nil
	}

	mean     := stats.Mean(data)
	median   := stats.Median(data)
	stdDev   := stats.StdDev(data, mean)
	sharpeR  := stats.SharpeRatio(mean, stdDev)
	skewness := stats.Skewness(mean, median, stdDev)
	percen   := stats.NewPercentile(data)

	perc01 := percen.Get(1) - mean
	perc30 := percen.Get(30) - mean
	perc70 := percen.Get(70) - mean
	perc99 := percen.Get(99) - mean

	lowerPercRatio := perc01 / perc30
	upperPercRatio := perc99 / perc70

	return &Distribution{
		Mean:        core.Trunc2d(mean),
		Median:      core.Trunc2d(median),
		StandardDev: core.Trunc2d(stdDev),
		SharpeRatio: core.Trunc2d(sharpeR),
		LowerTail:   core.Trunc2d(lowerPercRatio / 4.43),
		UpperTail:   core.Trunc2d(upperPercRatio / 4.43),
		Skewness:    core.Trunc2d(skewness),
		Histogram:   stats.NewHistogram(data),
	}
}

//=============================================================================

func calcRolling(res *AnalysisResponse) {
	costPerOper := res.TradingSystem.CostPerOperation

	for _, tr := range *res.Trades {
		year := tr.EntryDate.Year()
		dow := int(tr.EntryDate.Weekday())
		mon := int(tr.EntryDate.Month()) - 1

		dowRI := &res.Rolling.Daily[dow]
		monRI := &res.Rolling.Monthly[mon]

		updateRollingInfo(&tr, dowRI, costPerOper)
		updateRollingInfo(&tr, monRI, costPerOper)

		res.Rolling.DayYoY   = updateYoY(res.Rolling.DayYoY,   year, &tr, dow, costPerOper, 7)
		res.Rolling.MonthYoY = updateYoY(res.Rolling.MonthYoY, year, &tr, mon, costPerOper, 12)
	}
}

//=============================================================================

func updateRollingInfo(tr *db.Trade, ri *RollingInfo, costPerOper float64) {
	ri.Trades.Total++
	ri.GrossReturns.Total += tr.GrossReturn
	ri.NetReturns  .Total += tr.GrossReturn - 2*costPerOper

	if tr.TradeType == db.TradeTypeLong {
		ri.Trades.Long++
		ri.GrossReturns.Long += tr.GrossReturn
		ri.NetReturns  .Long += tr.GrossReturn - 2*costPerOper
	} else {
		ri.Trades.Short++
		ri.GrossReturns.Short += tr.GrossReturn
		ri.NetReturns  .Short += tr.GrossReturn - 2*costPerOper
	}
}

//=============================================================================

func updateYoY(list []*YoYRolling, year int, tr *db.Trade, slot int, costPerOper float64, slots int) []*YoYRolling {
	var yoy *YoYRolling

	if list == nil || list[len(list)-1].Year != year {
		yoy = &YoYRolling{
			Year: year,
		}

		list = append(list, yoy)

		for i := 0; i < slots; i++ {
			yoy.Data = append(yoy.Data, &RollingInfo{})
		}
	} else {
		yoy = list[len(list)-1]
	}

	updateRollingInfo(tr, yoy.Data[slot], costPerOper)

	return list
}

//=============================================================================

func calcLivePeriods(res *AnalysisResponse, livePeriods *[]db.LivePeriod) {
	list := []*LivePeriod{}

	if livePeriods != nil && len(*livePeriods) > 0 {
		var currPer *LivePeriod

		for _, lp := range *livePeriods {
			//--- Skip initial inactive states
			if currPer == nil && !lp.Active {
				continue
			}

			//--- First active state
			if currPer == nil && lp.Active {
				currPer = &LivePeriod{
					From: &lp.Period,
				}
			}

			//--- We have a starting active period. Skip other active periods
			if currPer != nil && lp.Active {
				continue
			}

			//--- Close the active period
			if currPer != nil && !lp.Active {
				currPer.To = &lp.Period
				list = append(list, currPer)
				currPer = nil
			}
		}

		if currPer != nil {
			list = append(list, currPer)
		}
	}

	res.LivePeriods = list
}

//=============================================================================
