//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package quality

import (
	"log/slog"
	"math"
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/portfolio-trader/pkg/platform"
	"golang.org/x/exp/stats"
)

//=============================================================================

const (
	DirectionAll  = 3
	VolatilityAll = 4
)

//=============================================================================

func GetQualityAnalysis(ts *db.TradingSystem, trades *[]db.Trade, man *platform.DataProductAnalysisResponse,
						timeframe int) (*AnalysisResponse, error) {
	res := NewAnalysisResponse()
	res.TradingSystem = ts

	if man.BarResults == nil || len(man.BarResults) == 0 {
		return &AnalysisResponse{
			TradingSystem: ts,
		}, nil
	}

	prodLoc,err := time.LoadLocation(ts.Timezone)
	if err != nil {
		slog.Error("GetQualityAnalysis: Failed to load timezone '%s'", ts.Timezone)
		return nil, err
	}

	marketRegime := NewMarketRegime(man.BarResults, timeframe, prodLoc)

	//--- Calc all standard metrics

	for dir := platform.DirectionStrongBear; dir <= platform.DirectionStrongBull; dir++ {
		for vol := platform.VolatilityQuiet; vol <= platform.VolatilityVeryVolatile; vol++ {
			calcQualityCell(res, trades, dir, vol, ts.CostPerOperation, marketRegime)
		}
	}

	//--- Calc summary by direction

	for dir := platform.DirectionStrongBear; dir <= platform.DirectionStrongBull; dir++ {
		calcQualityCell(res, trades, dir, VolatilityAll, ts.CostPerOperation, marketRegime)
	}

	//--- Calc summary by volatility

	for vol := platform.VolatilityQuiet; vol <= platform.VolatilityVeryVolatile; vol++ {
		calcQualityCell(res, trades, DirectionAll, vol, ts.CostPerOperation, marketRegime)
	}

	//--- Calc overall

	calcQualityCell(res, trades, DirectionAll, VolatilityAll, ts.CostPerOperation, marketRegime)

	return res, nil
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func calcQualityCell(res *AnalysisResponse, trades *[]db.Trade, dir int, vol int, costPerOperation float64,
					 marketRegime MarketRegime) {
	res.QualityAllGross  [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeAll,   dir, vol, 0, marketRegime)
	res.QualityLongGross [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeLong,  dir, vol, 0, marketRegime)
	res.QualityShortGross[dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeShort, dir, vol, 0, marketRegime)

	res.QualityAllNet  [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeAll,   dir, vol, costPerOperation, marketRegime)
	res.QualityLongNet [dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeLong,  dir, vol, costPerOperation, marketRegime)
	res.QualityShortNet[dir+2][vol] = calcQualityMetrics(trades, db.TradeTypeShort, dir, vol, costPerOperation, marketRegime)
}

//=============================================================================

func calcQualityMetrics(trades *[]db.Trade, tradeType string, direction int, volatility int, costPerOper float64,
						marketRegime MarketRegime) *Metrics {

	returns   := core.GetReturns(trades, tradeType, costPerOper)
	risk, err := core.CalcRisk(returns, costPerOper)
	if err != nil {
		return &Metrics{
			Trades: 0,
		}
	}

	//--- Step 1: Collect relevant trades

	var list []float64

	for _, t := range *trades {
		if t.TradeType == tradeType || tradeType == db.TradeTypeAll {
			tradeDir, tradeVol := marketRegime.MapTrade(&t)
			if tradeVol == -1 {
				slog.Warn("calcQualityMetrics: Cannot find market bar for trade", "tradeId", t.Id)
				continue
			}

			if direction == DirectionAll || direction == tradeDir {
				if volatility == VolatilityAll || volatility == tradeVol {
					ret := t.GrossReturn - 2*costPerOper
					list = append(list, ret/risk)
				}
			}
		}
	}

	//--- Step 2: Calc metrics

	cell := &Metrics{
		Trades: len(list),
	}

	if len(list) > 0 {
		cell.TradesPerc = core.Trunc2d(100 * float64(len(list)) / float64(len(*trades)))
		calcMetrics(list, cell)
	}

	return cell
}

//=============================================================================

func calcMetrics(list []float64, cell *Metrics) {
	mean, stdd := stats.MeanAndStdDev(list)
	listLen := float64(len(list))
	capLen := math.Min(listLen, 100)

	if stdd > 0.0 {
		cell.Sqn    = core.Trunc2d(mean / stdd * math.Sqrt(listLen))
		cell.Sqn100 = core.Trunc2d(mean / stdd * math.Sqrt(capLen))
	}

	equity := core.BuildEquity(&list)
	_, maxDD := core.BuildDrawDown(equity)

	cell.MaxDrawdown = core.Trunc2d(maxDD)
}

//=============================================================================
