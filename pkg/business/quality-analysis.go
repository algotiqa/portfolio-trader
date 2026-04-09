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

package business

import (
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/quality"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/portfolio-trader/pkg/platform"
	"gorm.io/gorm"
)

//=============================================================================

const (
	TimeframeSystem = "ts"
	TimeframeDaily  = "daily"
)
//=============================================================================

func RunQualityAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, req *quality.AnalysisRequest) (*quality.AnalysisResponse, error) {
	//--- Get trading system

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	fromDate,toDate,err := core.CalcSelectedPeriod(&req.SelectedPeriod, time.UTC)
	if err != nil {
		return nil, err
	}

	timeframe, err := parseTimeframeType(req.TimeframeType, ts)
	if err != nil {
		return nil, err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	if fromDate != nil {
		//--- We are going back 100 solar days in the past but aggregation on data collector is done on
		//--- 100 trading days. So, 100 trading days are (roughly) 150 solar days (we take some buffer)
		tmp := fromDate.Add(-time.Hour * time.Duration(24 * 150))
		fromDate = &tmp
	}

	man, err := platform.AnalyzeDataProduct(c, ts, fromDate, toDate, req.AtrLength, timeframe)
	if err != nil {
		return nil, err
	}

	return quality.GetQualityAnalysis(ts, trades, man, timeframe)
}

//=============================================================================

func parseTimeframeType(tfType string, ts *db.TradingSystem) (int, error) {
	if tfType == "" || tfType == TimeframeDaily {
		return 1440, nil
	}

	if tfType == TimeframeSystem {
		return ts.Timeframe, nil
	}

	return 0, req.NewBadRequestError("invalid timeframe type: " + tfType +" (must be one of 'ts','daily' )")
}

//=============================================================================
