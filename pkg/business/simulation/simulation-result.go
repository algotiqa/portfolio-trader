//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package simulation

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
)

//=============================================================================

type Result struct {
	FirstTradeDate types.Date `json:"firstTradeDate"`
	LastTradeDate  types.Date `json:"lastTradeDate"`
	Runs           int        `json:"runs"`

	CostPerOper    float64    `json:"costPerOper"`
	CurrencyCode   string     `json:"currencyCode"`
	Status         string     `json:"status"`
	StartTime      time.Time  `json:"startTime"`
	EndTime        time.Time  `json:"endTime"`
	Step           int        `json:"step"`

	GrossAll   *Details `json:"grossAll"`
	GrossLong  *Details `json:"grossLong"`
	GrossShort *Details `json:"grossShort"`
	NetAll     *Details `json:"netAll"`
	NetLong    *Details `json:"netLong"`
	NetShort   *Details `json:"netShort"`
}

//=============================================================================

func NewResult(first, last types.Date, runs int, ts *db.TradingSystem) *Result {
	return &Result{
		FirstTradeDate: first,
		LastTradeDate : last,
		Runs          : runs,
		Status        : SimStatusRunning,
		StartTime     : time.Now(),
		CostPerOper   : ts.CostPerOperation,
		CurrencyCode  : ts.CurrencyCode,
	}
}

//=============================================================================

type Details struct {
	DetectedRisk         float64       `json:"detectedRisk"`
	NumberOfTrades       int           `json:"numberOfTrades"`
	EquitiesImage        string        `json:"equitiesImage"`
	MaxDrawdownDistr     *Distribution `json:"maxDrawdownDistr"`
	MaxDrawdownProb      *Distribution `json:"maxDrawdownProb"`
	EquityReturn         float64       `json:"equityReturn"`
	EquityMaxDD          float64       `json:"equityMaxDD"`
	EquityReturnDDRatio  float64       `json:"equityReturnDDRatio"`
	EquityAverageTrade   float64       `json:"equityAverageTrade"`
	MedianReturn         float64       `json:"medianReturn"`
	MedianMaxDD          float64       `json:"medianMaxDD"`
	MedianReturnDDRatio  float64       `json:"medianReturnDDRatio"`
	MedianAverageTrade   float64       `json:"medianAverageTrade"`
}

//=============================================================================

type Distribution struct {
	XAxis []string  `json:"xAxis"`
	YAxis []float64 `json:"yAxis"`
}

//=============================================================================
