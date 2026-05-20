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
	"gorm.io/gorm"
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

func ImportTradingSystem(tx *gorm.DB, ts *db.TradingSystem, data []byte) error {
	its := TradingSystem{}
	err := json.Unmarshal(data, &its)
	if err == nil {
		err = updateTradingSystem(tx, ts, &its)
		if err == nil {
			err = setTradingFilter(tx, ts.Id, its.TradingFilter)
			if err == nil {
				err = addTrades(tx, ts.Id, its.Trades)
				if err == nil {
					err = addDailyReturns(tx, ts.Id, its.DailyReturns)
					if err == nil {
						err = addLivePeriods(tx, ts.Id, its.LivePeriods)
					}
				}
			}
		}
	}

	return err
}

//=============================================================================

func updateTradingSystem(tx *gorm.DB, ts *db.TradingSystem, its *TradingSystem) error {
	ts.Trading         = its.Trading
	ts.AutoActivation  = its.AutoActivation
	ts.Active          = its.Active
	ts.FirstTrade      = its.FirstTrade
	ts.LastTrade       = its.LastTrade
	ts.LastNetProfit   = its.LastNetProfit
	ts.LastNetAvgTrade = its.LastNetAvgTrade
	ts.LastNumTrades   = its.LastNumTrades

	//--- These cannot be set (an imported trading system must be off)
	ts.Running         = false
	ts.Status          = db.TsStatusOff

	return db.UpdateTradingSystem(tx, ts)
}

//=============================================================================

func setTradingFilter(tx *gorm.DB, id uint, f *TradingFilter) error {
	return db.SetTradingFilter(tx, &db.TradingFilter{
		TradingSystemId : id,
		EquAvgEnabled   : f.EquAvgEnabled,
		EquAvgLen       : f.EquAvgLen,
		PosProEnabled   : f.PosProEnabled,
		PosProLen       : f.PosProLen,
		WinPerEnabled   : f.WinPerEnabled,
		WinPerLen       : f.WinPerLen,
		WinPerValue     : f.WinPerValue,
		OldNewEnabled   : f.OldNewEnabled,
		OldNewOldLen    : f.OldNewOldLen,
		OldNewOldPerc   : f.OldNewOldPerc,
		OldNewNewLen    : f.OldNewNewLen,
		TrendlineEnabled: f.TrendlineEnabled,
		TrendlineLen    : f.TrendlineLen,
		TrendlineValue  : f.TrendlineValue,
		DrawdownEnabled : f.DrawdownEnabled,
		DrawdownMin     : f.DrawdownMin,
		DrawdownMax     : f.DrawdownMax,
	})
}

//=============================================================================

func addTrades(tx *gorm.DB, id uint, list []*Trade) error {
	for _, t := range list {
		err := db.AddTrade(tx, convertTrade(id,t))
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func convertTrade(id uint, t *Trade) *db.Trade {
	return &db.Trade{
		TradingSystemId   : id,
		TradeType         : t.TradeType,
		EntryDate         : t.EntryDate,
		EntryPrice        : t.EntryPrice,
		EntryLabel        : t.EntryLabel,
		ExitDate          : t.ExitDate,
		ExitPrice         : t.ExitPrice,
		ExitLabel         : t.ExitLabel,
		GrossProfit       : t.GrossProfit,
		Contracts         : t.Contracts,
		EntryDateAtBroker : t.EntryDateAtBroker,
		EntryPriceAtBroker: t.EntryPriceAtBroker,
		ExitDateAtBroker  : t.ExitDateAtBroker,
		ExitPriceAtBroker : t.ExitPriceAtBroker,
	}
}

//=============================================================================

func addDailyReturns(tx *gorm.DB, id uint, list []*DailyReturn) error {
	for _, dr := range list {
		err := db.AddDailyReturn(tx, convertDailyReturn(id,dr))
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func convertDailyReturn(id uint, d *DailyReturn) *db.DailyReturn {
	return &db.DailyReturn{
		TradingSystemId: id,
		Day            : d.Day,
		GrossProfit    : d.GrossProfit,
		Trades         : d.Trades,
	}
}

//=============================================================================

func addLivePeriods(tx *gorm.DB, id uint, list []*LivePeriod) error {
	for _, lp := range list {
		err := db.AddLivePeriod(tx, convertLivePeriod(id,lp))
		if err != nil {
			return err
		}
	}

	return nil
}

//=============================================================================

func convertLivePeriod(id uint, l *LivePeriod) *db.LivePeriod {
	return &db.LivePeriod{
		TradingSystemId: id,
		Period         : l.Period,
		Active         : l.Active,
	}
}

//=============================================================================
