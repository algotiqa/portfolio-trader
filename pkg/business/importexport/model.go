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
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Trading system
//===
//=============================================================================

type TradingSystem struct {
	Id               uint           `json:"id"`
	Trading          bool           `json:"trading"`
	Running          bool           `json:"running"`
	AutoActivation   bool           `json:"autoActivation"`
	Active           bool           `json:"active"`
	Status           db.TsStatus    `json:"status"`
	FirstTrade       *time.Time     `json:"firstTrade"`
	LastTrade        *time.Time     `json:"lastTrade"`
	LastNetProfit    float64        `json:"lastNetProfit"`
	LastNetAvgTrade  float64        `json:"lastNetAvgTrade"`
	LastNumTrades    int            `json:"lastNumTrades"`
	TradingFilter    *TradingFilter `json:"tradingFilter"`
	Trades           []*Trade       `json:"trades"`
	LivePeriods      []*LivePeriod  `json:"livePeriods"`
}

//=============================================================================

type TradingFilter struct {
	EquAvgEnabled    bool `json:"equAvgEnabled"`
	EquAvgLen        int  `json:"equAvgLen"`
	PosProEnabled    bool `json:"posProEnabled"`
	PosProLen        int  `json:"posProLen"`
	WinPerEnabled    bool `json:"winPerEnabled"`
	WinPerLen        int  `json:"winPerLen"`
	WinPerValue      int  `json:"winPerValue"`
	OldNewEnabled    bool `json:"oldNewEnabled"`
	OldNewOldLen     int  `json:"oldNewOldLen"`
	OldNewOldPerc    int  `json:"oldNewOldPerc"`
	OldNewNewLen     int  `json:"oldNewNewLen"`
	TrendlineEnabled bool `json:"trendlineEnabled"`
	TrendlineLen     int  `json:"trendlineLen"`
	TrendlineValue   int  `json:"trendlineValue"`
	DrawdownEnabled  bool `json:"drawdownEnabled"`
	DrawdownMin      int  `json:"drawdownMin"`
	DrawdownMax      int  `json:"drawdownMax"`
}

//=============================================================================

type Trade struct {
	TradeType          string       `json:"tradeType"`
	EntryDate          *time.Time   `json:"entryDate"`
	EntryPrice         float64      `json:"entryPrice"`
	EntryLabel         string       `json:"entryLabel"`
	ExitDate           *time.Time   `json:"exitDate"`
	ExitPrice          float64      `json:"exitPrice"`
	ExitLabel          string       `json:"exitLabel"`
	GrossReturn        float64      `json:"grossReturn"`
	MaxContracts       int          `json:"maxContracts"`
	EntryDateAtBroker  *time.Time   `json:"entryDateAtBroker"`
	EntryPriceAtBroker float64      `json:"entryPriceAtBroker"`
	ExitDateAtBroker   *time.Time   `json:"exitDateAtBroker"`
	ExitPriceAtBroker  float64      `json:"exitPriceAtBroker"`
	Equity             []*EquityBar `json:"equity"`
}

//=============================================================================

type EquityBar struct {
	Date        time.Time `json:"date"`
	GrossReturn float64    `json:"grossReturn"`
	Contracts   int        `json:"contracts"`
}

//=============================================================================

type LivePeriod struct {
	Period  time.Time `json:"period"`
	Active  bool      `json:"active"`
}

//=============================================================================
//===
//=== ExportedData
//===
//=============================================================================

type ExportedData struct {
	TradingSystems []*EncodedSystem `json:"tradingSystems"`
}

//=============================================================================

type EncodedSystem struct {
	Id       uint   `json:"id"`
	JsonData []byte `json:"jsonData"`
}

//=============================================================================
//===
//=== Constructors
//===
//=============================================================================

func NewTradingSystem(ts *db.TradingSystem) *TradingSystem {
	return &TradingSystem{
		Id             : ts.Id,
		Trading        : ts.Trading,
		Running        : ts.Running,
		AutoActivation : ts.AutoActivation,
		Active         : ts.Active,
		Status         : ts.Status,
		FirstTrade     : ts.FirstTrade,
		LastTrade      : ts.LastTrade,
		LastNetProfit  : ts.LastNetProfit,
		LastNetAvgTrade: ts.LastNetAvgTrade,
		LastNumTrades  : ts.LastNumTrades,
	}
}

//=============================================================================

func NewTradingFilter(f *db.TradingFilter) *TradingFilter {
	return &TradingFilter{
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
	}
}

//=============================================================================

func NewTrade(f *db.Trade) *Trade {
	return &Trade{
		TradeType:          f.TradeType,
		EntryDate:          f.EntryDate,
		EntryPrice:         f.EntryPrice,
		EntryLabel:         f.EntryLabel,
		ExitDate:           f.ExitDate,
		ExitPrice:          f.ExitPrice,
		ExitLabel:          f.ExitLabel,
		GrossProfit:        f.GrossProfit,
		Contracts:          f.Contracts,
		EntryDateAtBroker:  f.EntryDateAtBroker,
		EntryPriceAtBroker: f.EntryPriceAtBroker,
		ExitDateAtBroker:   f.ExitDateAtBroker,
		ExitPriceAtBroker:  f.ExitPriceAtBroker,
	}
}

//=============================================================================

func NewDailyReturn(f *db.DailyReturn) *DailyReturn {
	return &DailyReturn{
		Day        : f.Day,
		GrossProfit: f.GrossProfit,
		Trades     : f.Trades,
	}
}

//=============================================================================

func NewLivePeriod(f *db.LivePeriod) *LivePeriod {
	return &LivePeriod{
		Period:  f.Period,
		Active:  f.Active,
	}
}

//=============================================================================
