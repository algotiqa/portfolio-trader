//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/algotiqa/types"
)

//=============================================================================
//===
//=== Entities
//===
//=============================================================================

type Portfolio struct {
	Id       uint   `json:"id" gorm:"primaryKey"`
	ParentId uint   `json:"parentId"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

//=============================================================================

type TsStatus int8

const (
	TsStatusOff     TsStatus = 0
	TsStatusPaused  TsStatus = 1
	TsStatusRunning TsStatus = 2
	TsStatusIdle    TsStatus = 3
	TsStatusBroken  TsStatus = 4
)

//-----------------------------------------------------------------------------

type TsSuggAction int8

const (
	TsActionNone    TsSuggAction = 0
	TsActionTurnOff TsSuggAction = 1
	TsActionTurnOn  TsSuggAction = 2
	TsActionCheck   TsSuggAction = 3
)

//-----------------------------------------------------------------------------

type TradingSystem struct {
	Id               uint         `json:"id" gorm:"primaryKey"`
	Username         string       `json:"username"`
	Name             string       `json:"name"`
	Timeframe        int          `json:"timeframe"`
	DataProductId    uint         `json:"dataProductId"`
	DataSymbol       string       `json:"dataSymbol"`
	BrokerProductId  uint         `json:"brokerProductId"`
	BrokerSymbol     string       `json:"brokerSymbol"`
	PointValue       float64      `json:"pointValue"`
	CostPerOperation float64      `json:"costPerOperation"`
	MarginValue      float64      `json:"marginValue"`
	Increment        float64      `json:"increment"`
	MarketType       string       `json:"marketType"`
	CurrencyId       uint         `json:"currencyId"`
	CurrencyCode     string       `json:"currencyCode"`
	CurrencySymbol   string       `json:"currencySymbol"`
	TradingSessionId uint         `json:"tradingSessionId"`
	SessionName      string       `json:"sessionName"`
	SessionConfig    string       `json:"sessionConfig"`
	AgentProfileId   *uint        `json:"agentProfileId"`
	ExternalRef      string       `json:"externalRef"`
	StrategyType     string       `json:"strategyType"`
	Overnight        bool         `json:"overnight"`
	Tags             string       `json:"tags"`
	Finalized        bool         `json:"finalized"`
	Trading          bool         `json:"trading"`
	Running          bool         `json:"running"`
	AutoActivation   bool         `json:"autoActivation"`
	Active           bool         `json:"active"`
	Status           TsStatus     `json:"status"`
	SuggestedAction  TsSuggAction `json:"suggestedAction"`
	FirstTrade       *time.Time   `json:"firstTrade"`
	LastTrade        *time.Time   `json:"lastTrade"`
	LastNetProfit    float64      `json:"lastNetProfit"`
	LastNetAvgTrade  float64      `json:"lastNetAvgTrade"`
	LastNumTrades    int          `json:"lastNumTrades"`
	PortfolioId      *uint        `json:"portfolioId"`
	Timezone         string       `json:"timezone"`
	InSampleFrom     types.Date   `json:"inSampleFrom"`
	InSampleTo       types.Date   `json:"inSampleTo"`
	EngineCode       string       `json:"engineCode"`
}

//=============================================================================

type TradingFilter struct {
	TradingSystemId  uint `json:"omit" gorm:"primaryKey"`
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

type RpuType string

//-----------------------------------------------------------------------------

const (
	RpuStopLoss   RpuType = "stopLoss"
	RpuMaxLoss    RpuType = "maxLoss"
	RpuAvgLoss    RpuType = "avgLoss"
	RpuFixedValue RpuType = "fixedValue"
)

//-----------------------------------------------------------------------------

var RpuDomain = []RpuType{
	RpuStopLoss, RpuMaxLoss, RpuAvgLoss, RpuFixedValue,
}

//-----------------------------------------------------------------------------

type ModelName string

//-----------------------------------------------------------------------------

const (
	ModelFixedUnit         ModelName = "FU"
	ModelPercentRisk       ModelName = "PR"
	ModelPercentVolatility ModelName = "PV"
	ModelMarketMoney       ModelName = "MM"
)

//-----------------------------------------------------------------------------

type TradingPosition struct {
	TradingSystemId   uint        `json:"omit" gorm:"primaryKey"`
	InitialCapital    float64     `json:"initialCapital"`
	RuinPercentage    float64     `json:"ruinPercentage"`
	MarginOverride    *float64    `json:"marginOverride"`
	MaxUnits          int         `json:"maxUnits"`
	RiskPerUnit       RpuType     `json:"riskPerUnit"`
	RiskValue         *float64    `json:"riskValue"`
	Model             ModelName   `json:"model"`
	Config            string      `json:"config"`
}

//=============================================================================

const (
	TradeTypeLong  = "LO"
	TradeTypeShort = "SH"
	TradeTypeAll   = "**"
)

//-----------------------------------------------------------------------------

type Trade struct {
	Id                 int64      `json:"id" gorm:"primaryKey"`
	TradingSystemId    uint       `json:"tradingSystemId"`
	TradeType          string     `json:"tradeType"`
	EntryDate          *time.Time `json:"entryDate"`
	EntryPrice         float64    `json:"entryPrice"`
	EntryLabel         string     `json:"entryLabel"`
	ExitDate           *time.Time `json:"exitDate"`
	ExitPrice          float64    `json:"exitPrice"`
	ExitLabel          string     `json:"exitLabel"`
	GrossReturn        float64    `json:"grossReturn"`
	MaxContracts       int        `json:"maxContracts"`
	EntryDateAtBroker  *time.Time `json:"entryDateAtBroker"`
	EntryPriceAtBroker float64    `json:"entryPriceAtBroker"`
	ExitDateAtBroker   *time.Time `json:"exitDateAtBroker"`
	ExitPriceAtBroker  float64    `json:"exitPriceAtBroker"`
}

//-----------------------------------------------------------------------------

func (t Trade) String() string {
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v|%v|%v",
		t.TradeType,
		t.EntryDate.UTC(), t.EntryPrice, t.EntryLabel,
		t.ExitDate.UTC(), t.ExitPrice, t.ExitLabel,
		t.GrossReturn, t.MaxContracts)
}

//=============================================================================

type EquityBar struct {
	TradeId     int64      `json:"tradeId"`
	Date        time.Time  `json:"date"`
	GrossReturn float64    `json:"grossReturn"`
	Contracts   int        `json:"contracts"`
}

//=============================================================================

type LivePeriod struct {
	Id              uint      `json:"id" gorm:"primaryKey"`
	TradingSystemId uint      `json:"tradingSystemId"`
	Period          time.Time `json:"period"`
	Active          bool      `json:"active"`
}

//=============================================================================
//===
//=== Table names
//===
//=============================================================================

func (TradingSystem)   TableName() string { return "trading_system"   }
func (TradingFilter)   TableName() string { return "trading_filter"   }
func (TradingPosition) TableName() string { return "trading_position" }
func (Trade)           TableName() string { return "trade"            }
func (Portfolio)       TableName() string { return "portfolio"        }
func (EquityBar)       TableName() string { return "equity_bar"       }
func (LivePeriod)      TableName() string { return "live_period"      }

//=============================================================================
//===
//=== ParamMap type
//===
//=============================================================================

type ParamMap map[string]any

//=============================================================================

func (pm *ParamMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, pm)

	return err
}

//=============================================================================

func (pm ParamMap) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

//=============================================================================
