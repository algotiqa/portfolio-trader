//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
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
		TradeType         : f.TradeType,
		EntryDate         : f.EntryDate,
		EntryPrice        : f.EntryPrice,
		EntryLabel        : f.EntryLabel,
		ExitDate          : f.ExitDate,
		ExitPrice         : f.ExitPrice,
		ExitLabel         : f.ExitLabel,
		GrossReturn       : f.GrossReturn,
		MaxContracts      : f.MaxContracts,
		EntryDateAtBroker : f.EntryDateAtBroker,
		EntryPriceAtBroker: f.EntryPriceAtBroker,
		ExitDateAtBroker  : f.ExitDateAtBroker,
		ExitPriceAtBroker : f.ExitPriceAtBroker,
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
