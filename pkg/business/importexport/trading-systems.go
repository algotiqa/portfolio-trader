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

package importexport

import (
	"encoding/json"

	"github.com/algotiqa/portfolio-trader/pkg/db"
	"golang.org/x/exp/maps"
)

//=============================================================================

func BuildTradingSystems(systems *[]db.TradingSystem, filters *[]db.TradingFilter, trades *[]db.Trade,
						 dailys *[]db.DailyReturn, periods *[]db.LivePeriod) []*TradingSystem {

	tsMap := map[uint]*TradingSystem{}

	for _,ts := range *systems {
		tsMap[ts.Id] = NewTradingSystem(&ts)
	}

	for _, f := range *filters {
		ts,ok := tsMap[f.TradingSystemId]
		if ok {
			ts.TradingFilter = NewTradingFilter(&f)
		}
	}

	for _, tr := range *trades {
		ts,ok := tsMap[tr.TradingSystemId]
		if ok {
			ts.Trades = append(ts.Trades, NewTrade(&tr))
		}
	}

	for _, dr := range *dailys {
		ts,ok := tsMap[dr.TradingSystemId]
		if ok {
			ts.DailyReturns = append(ts.DailyReturns, NewDailyReturn(&dr))
		}
	}

	for _, lp := range *periods {
		ts,ok := tsMap[lp.TradingSystemId]
		if ok {
			ts.LivePeriods = append(ts.LivePeriods, NewLivePeriod(&lp))
		}
	}

	return maps.Values(tsMap)
}

//=============================================================================

func EncodeTradingSystems(list []*TradingSystem) (*ExportedData, error) {
	ed := &ExportedData{}

	for _, ts := range list {
		data, err := json.MarshalIndent(ts, "", "\t")
		if err != nil {
			return nil, err
		}

		es := &EncodedSystem{
			Id: ts.Id,
			JsonData: data,
		}

		ed.TradingSystems = append(ed.TradingSystems, es)
	}

	return ed, nil
}

//=============================================================================
