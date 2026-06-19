//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

func BuildGrossProfits(trades *[]db.Trade, tradeType string) (*[]time.Time, *[]float64) {
	timeSlice := []time.Time{}
	equSlice := []float64{}

	for _, tr := range *trades {
		if tradeType == db.TradeTypeAll || tr.TradeType == tradeType {
			etime := tr.ExitDate
			gross := tr.GrossReturn

			timeSlice = append(timeSlice, *etime)
			equSlice = append(equSlice, gross)
		}
	}

	return &timeSlice, &equSlice
}

//=============================================================================

func BuildNetProfits(grossProfits *[]float64, costPerOper float64) *[]float64 {
	netSlice := []float64{}

	for _, gross := range *grossProfits {
		net := gross - 2*costPerOper
		netSlice = append(netSlice, net)
	}

	return &netSlice
}

//=============================================================================

func BuildEquity(profits *[]float64) *[]float64 {
	equity := []float64{}
	value := 0.0

	for _, profit := range *profits {
		value += profit

		equity = append(equity, value)
	}

	return &equity
}

//=============================================================================

func BuildDrawDown(equity *[]float64) (*[]float64, float64) {
	maxProfit    := 0.0
	currDrawDown := 0.0
	maxDrawDown  := 0.0
	drawDown     := []float64{}

	for _, currProfit := range *equity {
		if currProfit >= maxProfit {
			maxProfit = currProfit
			currDrawDown = 0
		} else {
			currDrawDown = currProfit - maxProfit
		}

		drawDown = append(drawDown, currDrawDown)

		if currDrawDown < maxDrawDown {
			maxDrawDown = currDrawDown
		}
	}

	return &drawDown, maxDrawDown
}

//=============================================================================

func CalcWinningPercentage(profits []float64, filter []int8) float64 {
	tot := 0
	pos := 0

	for i, profit := range profits {
		if profit != 0 {
			if filter == nil || filter[i] == 1 {
				tot++
				if profit > 0 {
					pos++
				}
			}
		}
	}

	if tot == 0 {
		return 0
	}

	return float64(pos*10000/tot) / 100
}

//=============================================================================

func CalcAverageTrade(profits []float64, filter []int8) float64 {
	sum := 0.0
	num := 0.0

	for i, profit := range profits {
		if profit != 0 {
			if filter == nil || filter[i] == 1 {
				sum += profit
				num++
			}
		}
	}

	return Trunc2d(sum/num)
}

//=============================================================================

func CalcMin(data []float64) float64 {
	minv := data[0]
	for _, value := range data {
		if value < minv {
			minv = value
		}
	}

	return minv
}

//=============================================================================

func CalcRunUpAndDrawdown(data []float64) (float64,float64) {
	ru := 0.0
	dd := 0.0

	for _, value := range data {
		if value > ru {
			ru = value
		}

		if value < dd {
			dd = value
		}
	}

	return ru, dd
}

//=============================================================================
