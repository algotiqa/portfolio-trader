//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package inventory

import (
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
)

//=============================================================================
//=== Entities
//=============================================================================

type DataProduct struct {
	Id           uint   `json:"id"`
	ConnectionId uint   `json:"connectionId"`
	ExchangeId   uint   `json:"exchangeId"`
	Username     string `json:"username"`
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	MarketType   string `json:"marketType"`
	ProductType  string `json:"productType"`
}

//=============================================================================

type BrokerProduct struct {
	Id               uint    `json:"id"`
	ConnectionId     uint    `json:"connectionId"`
	ExchangeId       uint    `json:"exchangeId"`
	Username         string  `json:"username"`
	Symbol           string  `json:"symbol"`
	Name             string  `json:"name"`
	PointValue       float64 `json:"pointValue"`
	CostPerOperation float64 `json:"costPerOperation"`
	MarginValue      float64 `json:"marginValue"`
	Increment        float64 `json:"increment"`
	MarketType       string  `json:"marketType"`
	ProductType      string  `json:"productType"`
}

//=============================================================================

type Connection struct {
	Id                   uint   `json:"id"`
	Username             string `json:"username"`
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	SystemCode           string `json:"systemCode"`
	SystemName           string `json:"systemName"`
	SystemConfig         string `json:"systemConfig"`
	SupportsData         bool   `json:"supportsData"`
	SupportsBroker       bool   `json:"supportsBroker"`
	SupportsMultipleData bool   `json:"supportsMultipleData"`
	SupportsInventory    bool   `json:"supportsInventory"`
	SupportsAccount      bool   `json:"supportsAccount"`
}

//=============================================================================

type Exchange struct {
	Id         uint   `json:"id"`
	CurrencyId uint   `json:"currencyId"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Timezone   string `json:"timezone"`
	Url        string `json:"url"`
}

//=============================================================================

type Currency struct {
	Id        uint    `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	LastValue float64 `json:"lastValue"`
}

//=============================================================================

type TradingSession struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Config string `json:"config"`
}

//=============================================================================

type TradingSystem struct {
	Id               uint          `json:"id"`
	Username         string        `json:"username"`
	DataProductId    uint          `json:"dataProductId"`
	BrokerProductId  uint          `json:"brokerProductId"`
	TradingSessionId uint          `json:"tradingSessionId"`
	PortfolioId      *uint         `json:"portfolioId"`
	AgentProfileId   *uint         `json:"agentProfileId"`
	Name             string        `json:"name"`
	Timeframe        int           `json:"timeframe"`
	StrategyType     string        `json:"strategyType"`
	Overnight        bool          `json:"overnight"`
	Tags             string        `json:"tags"`
	ExternalRef      string        `json:"externalRef"`
	Finalized        bool          `json:"finalized"`
	InSampleFrom     types.Date    `json:"inSampleFrom"`
	InSampleTo       types.Date    `json:"inSampleTo"`
	EngineCode       db.EngineCode `json:"engineCode"`
}

//=============================================================================

type Account struct {
	Id              uint    `json:"id"`
	Username        string  `json:"username"`
	ConnectionId    uint    `json:"connectionId"`
	CurrencyId      uint    `json:"currencyId"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	CurrentCapital  float64 `json:"currentCapital"`
	SupportsAccount bool    `json:"supportsAccount"`
	StatusMessage   string  `json:"statusMessage"`
}

//=============================================================================

type Portfolio struct {
	Id             uint              `json:"id"`
	Username       string            `json:"username"`
	AccountId      uint              `json:"accountId"`
	Name           string            `json:"name"`
	Management     db.ManagementType `json:"management"`
	AccountPerc    float64           `json:"accountPerc"`
	MaxMarginPerc  float64           `json:"maxMarginPerc"`
}

//=============================================================================
//=== Messages
//=============================================================================

type DataProductMessage struct {
	DataProduct DataProduct `json:"dataProduct"`
	Connection  Connection  `json:"connection"`
	Exchange    Exchange    `json:"exchange"`
}

//=============================================================================

type BrokerProductMessage struct {
	BrokerProduct BrokerProduct `json:"brokerProduct"`
	Connection    Connection    `json:"connection"`
	Exchange      Exchange      `json:"exchange"`
}

//=============================================================================

type TradingSystemMessage struct {
	TradingSystem  TradingSystem  `json:"tradingSystem"`
	DataProduct    DataProduct    `json:"dataProduct"`
	BrokerProduct  BrokerProduct  `json:"brokerProduct"`
	Currency       Currency       `json:"currency"`
	TradingSession TradingSession `json:"tradingSession"`
	Exchange       Exchange       `json:"exchange"`
	PortfolioPack  []byte         `json:"portfolioPack"`
	StoragePack    []byte         `json:"storagePack"`
}

//=============================================================================

type AccountMessage struct {
	Account    Account    `json:"account"`
	Connection Connection `json:"connection"`
	Currency   Currency   `json:"currency"`
}

//=============================================================================

type PortfolioMessage struct {
	Portfolio  Portfolio `json:"portfolio"`
	Account    Account   `json:"account"`
	Currency   Currency  `json:"currency"`
}

//=============================================================================
