//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package position

//=============================================================================
//===
//=== AnalysisResponse
//===
//=============================================================================

type AnalysisResponse struct {
	TradingSystem  *TradingSystem  `json:"tradingSystem"`
	Params         *Parameters     `json:"params"`
	Baseline       *AnalysisResult `json:"baseline"`
	Current        *AnalysisResult `json:"current"`
	Selected       *AnalysisResult `json:"selected"`
	ParamSpecs     map[string]any  `json:"paramSpecs"`
}

//=============================================================================

type TradingSystem struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
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
	Return            float64   `json:"return"`
	MaxDrawdown       float64   `json:"maxDrawdown"`
	ReturnDrawdRatio  float64   `json:"returnDrawdRatio"`
}

//=============================================================================
