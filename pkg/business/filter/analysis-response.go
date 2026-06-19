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

	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== FilterAnalysisResponse
//===
//=============================================================================

type AnalysisResponse struct {
	TradingSystem TradingSystem     `json:"tradingSystem"`
	Filter        *db.TradingFilter `json:"filter"`
	Summary       Summary           `json:"summary"`
	Equities      Equities          `json:"equities"`
	Activations   *Activations      `json:"activations"`
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
	NetProfit          []float64   `json:"netProfit"`
	UnfilteredEquity   []float64   `json:"unfilteredEquity"`
	FilteredEquity     []float64   `json:"filteredEquity"`
	UnfilteredDrawdown []float64   `json:"unfilteredDrawdown"`
	FilteredDrawdown   []float64   `json:"filteredDrawdown"`
	FilterActivation   []int8      `json:"filterActivation"`
	Average            *core.Serie `json:"average"`
}

//=============================================================================

type Activation struct {
	Time   []time.Time `json:"time"`
	Values []int8      `json:"values"`
}

//-----------------------------------------------------------------------------

func (p *Activation) AddPoint(time time.Time, value int8) {
	p.Time = append(p.Time, time)
	p.Values = append(p.Values, value)
}

//-----------------------------------------------------------------------------

func (p *Activation) IsLastActive() bool {
	return p.Values[len(p.Values)-1] != 0
}

//=============================================================================

type Activations struct {
	EquityVsAverage   *Activation `json:"equityVsAverage"`
	PositiveProfit    *Activation `json:"positiveProfit"`
	WinningPercentage *Activation `json:"winningPercentage"`
	OldVsNew          *Activation `json:"oldVsNew"`
	Trendline         *Activation `json:"trendline"`
	Drawdown          *Activation `json:"drawdown"`
}

//-----------------------------------------------------------------------------

func (a *Activations) IsLastActive() bool {
	equVsAvg := true
	posProf := true
	winPerc := true
	oldNew := true
	trendline := true
	drawdown := true

	if a.EquityVsAverage != nil {
		equVsAvg = a.EquityVsAverage.IsLastActive()
	}

	if a.PositiveProfit != nil {
		posProf = a.PositiveProfit.IsLastActive()
	}

	if a.WinningPercentage != nil {
		winPerc = a.WinningPercentage.IsLastActive()
	}

	if a.OldVsNew != nil {
		oldNew = a.OldVsNew.IsLastActive()
	}

	if a.Trendline != nil {
		trendline = a.Trendline.IsLastActive()
	}

	if a.Drawdown != nil {
		drawdown = a.Drawdown.IsLastActive()
	}

	return equVsAvg && posProf && winPerc && oldNew && trendline && drawdown
}

//=============================================================================
