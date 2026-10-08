//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package allocation

import (
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

const MinDailyReturns = 15

//=============================================================================
//===
//=== JobSpec
//===
//=============================================================================

type JobSpec struct {
	AccountCapital float64
	AccountPerc    float64
	MaxMarginPerc  float64
	TradeYears     int
}

//=============================================================================
//===
//=== Job
//===
//=============================================================================

type Job struct {
	spec               *JobSpec
	portfolio          *db.Portfolio
	tradingSystems     []*TradingSystemInfo
	healthySystems     []*TradingSystem
	filterOutcomes     []*FilterOutcome
	logs               []*Log
	correlations       []*SystemCorrelation
	excludedSystems    []*ExcludedSystem
	correlationMatrix  *CorrelationMatrix
}

//=============================================================================

type TradingSystemInfo struct {
	system        *db.TradingSystem
	filter        *db.TradingFilter
	position      *db.TradingPosition
	trades        *[]db.Trade
	returns       *[]db.DailyReturn
	healthy       bool
	filterPassed  bool
}

//=============================================================================

func NewTradingSystemInfo(ts *db.TradingSystem, tsf *db.TradingFilter, tsp *db.TradingPosition,
						  trades *[]db.Trade, returns *[]db.DailyReturn) *TradingSystemInfo {
	return &TradingSystemInfo{
		system      : ts,
		filter      : tsf,
		position    : tsp,
		trades      : trades,
		returns     : returns,
	}
}

//=============================================================================

type Action string

//-----------------------------------------------------------------------------

const (
	ActionNone    Action = "none"
	ActionTurnOn  Action = "turnOn"
	ActionTurnOff Action = "turnOff"
)

//-----------------------------------------------------------------------------

type FilterOutcome struct {
	TsId             uint      `json:"tsId"`
	TsName           string    `json:"tsName"`
	TsDataSymbol     string    `json:"tsDataSymbol"`
	TsBrokerSymbol   string    `json:"tsBrokerSymbol"`
	TsMarketType     string    `json:"tsMarketType"`
	TsStrategyType   string    `json:"tsStrategyType"`
	TsRunning        bool      `json:"tsRunning"`
	FilterPassed     bool      `json:"filterPassed"`
	Action           Action    `json:"action"`
	Comment          string    `json:"comment"`
}

//=============================================================================

type LogLevel string

const (
	LogLevelInfo    LogLevel = "I"
	LogLevelWarning LogLevel = "W"
	LogLevelError   LogLevel = "E"
)

//-----------------------------------------------------------------------------

type Log struct {
	Level    LogLevel `json:"level"`
	Message  string   `json:"message"`
}

//=============================================================================

type SystemCorrelation struct {
	Ts1Id        uint     `json:"ts1Id"`
	Ts2Id        uint     `json:"ts2Id"`
	Ts1Name      string   `json:"ts1Name"`
	Ts2Name      string   `json:"ts2Name"`
	Correlation  float64  `json:"correlation"`
}

//=============================================================================

type CorrelationMatrix struct {
	Names []string		`json:"names"`
	Cells [][]*float64	`json:"cells"`
}

//=============================================================================

type SystemPosition struct {
	tsi      *TradingSystemInfo
	mod      model.PositionModel
	risk     float64
	atrValue float64
	scaling  float64
	position float64
}

//=============================================================================

func NewSystemPosition(tsi *TradingSystemInfo, mod model.PositionModel, risk, atrValue float64) *SystemPosition {
	return &SystemPosition{
		tsi     : tsi,
		mod     : mod,
		risk    : risk,
		atrValue: atrValue,
		scaling : 1,
		position: 0,
	}
}

//=============================================================================
//===
//=== Report
//===
//=============================================================================

type Report struct {
	TradingSystems     []*TradingSystem      `json:"tradingSystems"`
	FilterOutcomes     []*FilterOutcome      `json:"filterOutcomes"`
	Logs               []*Log                `json:"logs"`
	Correlations       []*SystemCorrelation  `json:"correlations"`
	CorrelationMatrix  *CorrelationMatrix    `json:"correlationMatrix"`
}

//=============================================================================

type HealthLevel string

const (
	HealthLevelOk    HealthLevel = "O"
	HealthLevelWarn  HealthLevel = "W"
	HealthLevelError HealthLevel = "E"
)

//-----------------------------------------------------------------------------

type TradingSystem struct {
	Id             uint        `json:"id"`
	Name           string      `json:"name"`
	DataSymbol     string      `json:"dataSymbol"`
	BrokerSymbol   string      `json:"brokerSymbol"`
	MarketType     string      `json:"marketType"`
	StrategyType   string      `json:"strategyType"`
	HealthLevel    HealthLevel `json:"healthLevel"`
	Message        string      `json:"message"`
}

//=============================================================================

type ExcludedSystem struct {
	Id             uint        `json:"id"`
	Name           string      `json:"name"`
	DataSymbol     string      `json:"dataSymbol"`
	BrokerSymbol   string      `json:"brokerSymbol"`
	Message        string      `json:"message"`
}

//=============================================================================

func NewExcludedSystem(tsi *TradingSystemInfo, message string) *ExcludedSystem {
	return &ExcludedSystem{
		Id          : tsi.system.Id,
		Name        : tsi.system.Name,
		DataSymbol  : tsi.system.DataSymbol,
		BrokerSymbol: tsi.system.BrokerSymbol,
		Message     : message,
	}
}

//=============================================================================
