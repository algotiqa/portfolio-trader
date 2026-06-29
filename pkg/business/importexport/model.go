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
	Id               uint                `json:"id"`
	Trading          bool                `json:"trading"`
	Running          bool                `json:"running"`
	AutoActivation   bool                `json:"autoActivation"`
	Active           bool                `json:"active"`
	Status           db.TsStatus         `json:"status"`
	FirstTrade       *time.Time          `json:"firstTrade"`
	LastTrade        *time.Time          `json:"lastTrade"`
	LastNetProfit    float64             `json:"lastNetProfit"`
	LastNetAvgTrade  float64             `json:"lastNetAvgTrade"`
	LastNumTrades    int                 `json:"lastNumTrades"`
	TradingFilter    *db.TradingFilter   `json:"tradingFilter"`
	TradingPosition  *db.TradingPosition `json:"tradingPosition"`
	Trades           []*Trade            `json:"trades"`
	LivePeriods      []*LivePeriod       `json:"livePeriods"`
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
