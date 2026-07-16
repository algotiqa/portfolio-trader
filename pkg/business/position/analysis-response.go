//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package position

import "time"

//=============================================================================
//===
//=== AnalysisResponse
//===
//=============================================================================

type AnalysisResponse struct {
	TradingSystem    *TradingSystem  `json:"tradingSystem"`
	Params           *Parameters     `json:"params"`
	Baseline         *AnalysisResult `json:"baseline"`
	Current          *AnalysisResult `json:"current"`
	Selected         *AnalysisResult `json:"selected"`
	Time             []time.Time     `json:"time"`
	ParamSpecs       map[string]any  `json:"paramSpecs"`
	ModelSpecs       map[string]any  `json:"modelSpecs"`
	UsedMargin       float64         `json:"usedMargin"`
	GrossRisk        float64         `json:"grossRisk"`
	NetRisk          float64         `json:"netRisk"`
	NoLosses         bool            `json:"noLosses"`
	CostPerOperation float64         `json:"costPerOperation"`
}

//=============================================================================

type TradingSystem struct {
	Id     uint    `json:"id"`
	Name   string  `json:"name"`
}

//=============================================================================

type AnalysisResult struct {
	Model  *Model            `json:"model"`
	Gross  *ModelPerformance `json:"gross"`
	Net    *ModelPerformance `json:"net"`
}

//=============================================================================

type ModelPerformance struct {
	Equity            []float64 `json:"equity"`
	DrawdownPerc      []float64 `json:"drawdownPerc"`
	Positions         []int     `json:"positions"`
	Return            float64   `json:"return"`
	MaxDrawdown       float64   `json:"maxDrawdown"`
	MaxDrawdownPerc   float64   `json:"maxDrawdownPerc"`
	ReturnDrawdRatio  float64   `json:"returnDrawdRatio"`
	ReturnOnAccount   float64   `json:"returnOnAccount"`
	Ruined            bool      `json:"ruined"`
}

//=============================================================================
