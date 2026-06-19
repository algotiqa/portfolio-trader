//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import "time"

//=============================================================================
//===
//=== PortfolioMonitoringResponse
//===
//=============================================================================

type PortfolioMonitoringParams struct {
	TsIds  []uint `form:"tsIds"  binding:"required,min=1,dive"`
	Period    int `form:"period" binding:"required,min=1,max=5000"`
}

//=============================================================================

type BaseMonitoring struct {
	Time          *[]time.Time `json:"time"`
	GrossProfit   *[]float64   `json:"grossProfit"`
	NetProfit     *[]float64   `json:"netProfit"`
	GrossDrawdown *[]float64   `json:"grossDrawdown"`
	NetDrawdown   *[]float64   `json:"netDrawdown"`
}

//=============================================================================

type TradingSystemMonitoring struct {
	BaseMonitoring
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//=============================================================================

func NewTradingSystemMonitoring(size int) *TradingSystemMonitoring {
	tsa := &TradingSystemMonitoring{}

	return tsa
}

//=============================================================================

type PortfolioMonitoringResponse struct {
	BaseMonitoring
	TradingSystems []*TradingSystemMonitoring `json:"tradingSystems"`
}

//=============================================================================
