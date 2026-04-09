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

package platform

import (
	"fmt"
	"net/url"
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
)

//=============================================================================
//===
//=== Response
//===
//=============================================================================

const (
	DirectionStrongBear = -2
	DirectionBear       = -1
	DirectionNeutral    = 0
	DirectionBull       = 1
	DirectionStrongBull = 2
)

const (
	VolatilityQuiet        = 0
	VolatilityNormal       = 1
	VolatilityVolatile     = 2
	VolatilityVeryVolatile = 3
)

//=============================================================================

type DataProductAnalysisResponse struct {
	Id         uint          `json:"id"`
	Symbol     string        `json:"symbol"`
	From       types.Date    `json:"from"`
	To         types.Date    `json:"to"`
	Bars       int           `json:"bars"`
	Timeframe  int           `json:"timeframe"`
	AtrLength  int           `json:"atrLength"`
	BarResults []*BarResult  `json:"barResults"`
}

//=============================================================================

type BarResult struct {
	Time            time.Time  `json:"time"`
	Close           float64    `json:"close"`
	BarChangePerc   float64    `json:"barChangePerc"`
	TrueRange       float64    `json:"trueRange"`
	Sqn100          float64    `json:"sqn100"`
	Atr             float64    `json:"atr"`
	AtrPerc         float64    `json:"atrPerc"`
	AtrMeanPerc     float64    `json:"atrMeanPerc"`
	AtrStdDevPerc   float64    `json:"atrStdDevPerc"`
	Direction       int        `json:"direction"`
	Volatility      int        `json:"volatility"`
}

//=============================================================================
//===
//=== Public functions
//===
//=============================================================================

func AnalyzeDataProduct(c *auth.Context, ts *db.TradingSystem, from,to *time.Time, atrLen int, timeframe int) (*DataProductAnalysisResponse, error) {
	id := ts.DataProductId
	c.Log.Info("AnalyzeDataProduct: Asking data product analysis to data collector", "id", id, "from", from, "to", to)

	token  := c.Token
	client := req.GetClient("bf")
	srvUrl := fmt.Sprintf("%s/v1/data-products/%d/analysis?timeframe=%d&sessionId=%d&atrLen=%d",
						platform.Data, id, timeframe, ts.TradingSessionId, atrLen)
	if from != nil {
		srvUrl = srvUrl + "&from="+ url.QueryEscape(from.Format(time.DateTime))
	}
	if to != nil {
		srvUrl = srvUrl + "&to="+ url.QueryEscape(to.Format(time.DateTime))
	}

	var res DataProductAnalysisResponse
	err := req.DoGet(client, srvUrl, &res, token)
	if err != nil {
		c.Log.Error("AnalyzeDataProduct: Got an error when accessing the data-collector", "id", id, "error", err.Error())
		return nil, err
	}

	c.Log.Info("AnalyzeDataProduct: Analysis received", "id", id)
	return &res, nil
}

//=============================================================================

