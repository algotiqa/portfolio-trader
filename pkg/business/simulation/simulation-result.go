//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
