//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

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

package trade

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type AnalysisResponse struct {
	TradingSystem *db.TradingSystem `json:"tradingSystem"`
	Trades        []*Entry          `json:"trades"`
}

//=============================================================================

func NewAnalysisResponse() *AnalysisResponse {
	return &AnalysisResponse{
	}
}

//=============================================================================

type Entry struct {
	TradeType    string     `json:"tradeType"`
	EntryDate    *time.Time `json:"entryDate"`
	EntryLabel   string     `json:"entryLabel"`
	ExitDate     *time.Time `json:"exitDate"`
	ExitLabel    string     `json:"exitLabel"`
	GrossReturn  float64    `json:"grossReturn"`
	MaxContracts int        `json:"maxContracts"`

	GrossEquity *EquityInfo  `json:"grossEquity"`
	NetEquity   *EquityInfo  `json:"netEquity"`
	Contracts   []int        `json:"contracts"`
}

//=============================================================================

type EquityInfo struct {
	Equity   []float64  `json:"equity"`
	Return   float64    `json:"return"`
	RunUp    float64    `json:"runUp"`
	Drawdown float64    `json:"drawdown"`
}

//=============================================================================
