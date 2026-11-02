//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== FilterAnalysisResponse
//===
//=============================================================================

type AnalysisResponse struct {
	TradingSystem    TradingSystem `json:"tradingSystem"`
	Filter          *TradingFilter `json:"filter"`
	Summary          Summary       `json:"summary"`
	Equities         Equities      `json:"equities"`
	Outcomes        *Outcomes      `json:"outcomes"`
	Portfolio       *db.Portfolio  `json:"portfolio"`
}

//=============================================================================

type TradingSystem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

//=============================================================================

type Summary struct {
	UnfProfit       float64 `json:"unfProfit"`
	FilProfit       float64 `json:"filProfit"`
	UnfMaxDrawdown  float64 `json:"unfMaxDrawdown"`
	FilMaxDrawdown  float64 `json:"filMaxDrawdown"`
	UnfWinningPerc  float64 `json:"unfWinningPerc"`
	FilWinningPerc  float64 `json:"filWinningPerc"`
	UnfAverageTrade float64 `json:"unfAverageTrade"`
	FilAverageTrade float64 `json:"filAverageTrade"`
}

//=============================================================================

type Equities struct {
	Time               []time.Time `json:"time"`
	UnfilteredEquity   []float64   `json:"unfilteredEquity"`
	FilteredEquity     []float64   `json:"filteredEquity"`
	UnfilteredDrawdown []float64   `json:"unfilteredDrawdown"`
	FilteredDrawdown   []float64   `json:"filteredDrawdown"`
}

//=============================================================================

type Outcomes struct {
	Activations     []*criteria.Activation  `json:"activations"`
	Deactivations   []*criteria.Activation  `json:"deactivations"`
	CombinedActiv   []int                   `json:"combinedActiv"`
	CombinedDeact   []int                   `json:"combinedDeact"`
	CombinedFinal   []int                   `json:"combinedFinal"`
	SampledFinal    []int                   `json:"sampledFinal"`
	FinalActivation *criteria.Activation    `json:"finalActivation"`
}

//=============================================================================
