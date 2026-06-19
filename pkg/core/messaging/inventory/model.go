//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package inventory

import "github.com/algotiqa/types"

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
	InstanceCode         string `json:"instanceCode"`
	SupportsData         bool   `json:"supportsData"`
	SupportsBroker       bool   `json:"supportsBroker"`
	SupportsMultipleData bool   `json:"supportsMultipleData"`
	SupportsInventory    bool   `json:"supportsInventory"`
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
	Id     uint   `json:"id"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
}

//=============================================================================

type TradingSession struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Config string `json:"config"`
}

//=============================================================================

type TradingSystem struct {
	Id               uint       `json:"id"`
	Username         string     `json:"username"`
	Name             string     `json:"name"`
	Timeframe        int        `json:"timeframe"`
	DataProductId    uint       `json:"dataProductId"`
	BrokerProductId  uint       `json:"brokerProductId"`
	TradingSessionId uint       `json:"tradingSessionId"`
	AgentProfileId   *uint      `json:"agentProfileId"`
	StrategyType     string     `json:"strategyType"`
	Overnight        bool       `json:"overnight"`
	Tags             string     `json:"tags"`
	ExternalRef      string     `json:"externalRef"`
	Finalized        bool       `json:"finalized"`
	InSampleFrom     types.Date `json:"inSampleFrom"`
	InSampleTo       types.Date `json:"inSampleTo"`
	EngineCode       string     `json:"engineCode"`
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
