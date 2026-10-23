//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"math"
	"slices"
	"time"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"golang.org/x/exp/stats"
)

//=============================================================================

type Metrics struct {
	NetReturn     float64 `json:"netReturn"`
	Trades        int     `json:"trades"`
	AverageTrade  float64 `json:"averageTrade"`
	ProfitFactor  float64 `json:"profitFactor"`
	WinningPerc   float64 `json:"winningPerc"`
	MaxDrawdown   float64 `json:"maxDrawdown"`
	SharpeRatio   float64 `json:"sharpeRatio"`
	Sqn100        float64 `json:"sqn100"`
	ReturnDDRatio float64 `json:"returnDDRatio"`
}

//=============================================================================

func Trunc2d(value float64) float64 {
	return float64(int(value*100)) / 100
}

//=============================================================================

func Trunc2dSlice(value []float64) {
	for i, v := range value {
		value[i] = Trunc2d(v)
	}
}

//=============================================================================

func GetLocation(timezone string, ts *db.TradingSystem) (*time.Location, error) {
	if timezone == "exchange" {
		timezone = ts.Timezone
	}

	return time.LoadLocation(timezone)
}

//=============================================================================

func CalcRisk(returns []float64, costPerOper float64) (float64, error) {
	if returns == nil || len(returns) == 0 {
		return 0, req.NewUnprocessableEntityError("no losses found")
	}

	//--- Step 1: calculate distribution of losses

	counts := map[float64]int{}

	for _, ret := range returns {
		if ret < NetReturn(0, costPerOper) {
			count, ok := counts[ret]
			if !ok {
				counts[ret] = 1
			} else {
				counts[ret] = count + 1
			}
		}
	}

	if len(counts) == 0 {
		return 0, req.NewUnprocessableEntityError("no losses found")
	}

	//--- Step 2: Order all losses by their count

	var losses []*riskInfo

	for loss, count := range counts {
		losses = append(losses, &riskInfo{ loss, count })
	}

	slices.SortFunc(losses, func(a, b *riskInfo) int {
		return b.count - a.count
	})

	//--- Step 3: find the stop loss (the loss with most hits)

	if len(losses) == 1 {
		return math.Abs(losses[0].loss), nil
	}

	if losses[0].count > losses[1].count * 3 {
		return math.Abs(losses[0].loss), nil
	}

	// Step 4: Stop loss not identifies (maybe varies depending on ATR). Calc average

	sum := 0.0

	for _, ri := range losses {
		sum += ri.loss
	}

	loss := math.Abs(sum / float64(len(losses)))
	return Trunc2d(loss), nil
}

//-----------------------------------------------------------------------------

type riskInfo struct {
	loss  float64
	count int
}

//=============================================================================

func GetReturns(trades *[]db.Trade, tradeType string, costPerOper float64, fromDate *time.Time) []float64 {
	var list []float64

	for _, t := range *trades {
		if tradeType == db.TradeTypeAll || tradeType == t.TradeType {
			if fromDate == nil || t.ExitDate.After(*fromDate) {
				returns := NetReturn(t.GrossReturn, costPerOper)
				list = append(list, returns)
			}
		}
	}

	return list
}

//=============================================================================

func CalcRMultiple(returns []float64, risk float64) []float64 {
	var list []float64

	for _, r := range returns {
		list = append(list, r/risk)
	}

	return list
}

//=============================================================================

func CalcSQN(list []float64) (float64,float64){
	if len(list) == 0 {
		return 0.0, 0.0
	}

	mean, stdd := stats.MeanAndStdDev(list)
	listLen := float64(len(list))
	capLen := math.Min(listLen, 100)

	sqn    := 0.0
	sqn100 := 0.0

	if stdd > 0.0 {
		sqn    = Trunc2d(mean / stdd * math.Sqrt(listLen))
		sqn100 = Trunc2d(mean / stdd * math.Sqrt(capLen))
	}

	return sqn, sqn100
}

//=============================================================================

func CalcSharpeRatio(list []float64) float64 {
	sr := 0.0

	if len(list) == 0 {
		return sr
	}

	mean, stdd := stats.MeanAndStdDev(list)

	if stdd > 0.0 {
		sr = Trunc2d(mean / stdd)
	}

	return sr
}

//=============================================================================

func CalcRiskForPosition(riskPerUnit db.RpuType, riskValue *float64, trades *[]db.Trade, costPerOper float64) (float64,error) {
	if riskPerUnit == db.RpuFixedValue {
		return NetRisk(*riskValue, costPerOper),nil
	}

	returns := GetReturns(trades, db.TradeTypeAll, costPerOper, nil)

	if riskPerUnit == db.RpuStopLoss {
		return CalcRisk(returns, costPerOper)
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

	if len(losses) == 0 {
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

func NormalizeContracts(list *[]db.Trade) {
	for index := range *list {
		t := &(*list)[index]

		if t.MaxContracts > 0 {
			t.GrossReturn /= float64(t.MaxContracts)
		}
	}
}

//=============================================================================

func NetReturn(grossReturn float64, costPerOper float64) float64 {
	return grossReturn - 2*costPerOper
}

//=============================================================================

func NetRisk(grossRisk float64, costPerOper float64) float64 {
	return grossRisk + 2*costPerOper
}

//=============================================================================

func CalcMetrics(trades *[]db.Trade, costPerOper float64, startDate *time.Time) *Metrics {
	netTotal  := 0.0
	netProf   := 0.0
	netLoss   := 0.0
	countTotal:= 0
	countProf := 0
	countLoss := 0

	var returns []float64

	for _, t := range *trades {
		if startDate == nil || t.ExitDate.After(*startDate) {
			netRet := NetReturn(t.GrossReturn, costPerOper)
			netTotal += netRet
			countTotal++
			returns = append(returns, netRet)

			if netRet > 0 {
				netProf += netRet
				countProf++
			}
			if netRet < 0 {
				netLoss += netRet
				countLoss++
			}
		}
	}

	//--- Calc average trade and winning %

	avgTrade    := 0.0
	winningPerc := 0.0

	if countTotal > 0 {
		avgTrade    = netTotal / float64(countTotal)
		winningPerc = float64(countProf) / float64(countTotal) * 100
	}

	//--- Calc profit factor

	profFactor := 10.0
	if netLoss != 0 {
		profFactor = netProf / math.Abs(netLoss)
	}

	//--- Calc sharpe and sqn

	sharpeRatio := CalcSharpeRatio(returns)
	_, sqn100   := CalcSQN(returns)

	//--- Calc max drawdown

	equity   := BuildEquity(&returns)
	_, maxDD := BuildDrawDown(equity)

	//--- Calc return/maxDrawdown ratio

	retDDRatio := 10.0
	if maxDD != 0 {
		retDDRatio = netTotal / math.Abs(maxDD)
	}

	return &Metrics{
		NetReturn    : Trunc2d(netTotal),
		Trades       : countTotal,
		AverageTrade : Trunc2d(avgTrade),
		ProfitFactor : Trunc2d(profFactor),
		WinningPerc  : Trunc2d(winningPerc),
		MaxDrawdown  : Trunc2d(maxDD),
		SharpeRatio  : sharpeRatio,
		Sqn100       : sqn100,
		ReturnDDRatio: Trunc2d(retDDRatio),
	}
}

//=============================================================================
