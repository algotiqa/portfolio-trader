//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package quality

import (
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type AnalysisResponse struct {
	TradingSystem     *db.TradingSystem `json:"tradingSystem"`
	QualityAllGross   *[6][5]*Metrics   `json:"qualityAllGross"`
	QualityLongGross  *[6][5]*Metrics   `json:"qualityLongGross"`
	QualityShortGross *[6][5]*Metrics   `json:"qualityShortGross"`
	QualityAllNet     *[6][5]*Metrics   `json:"qualityAllNet"`
	QualityLongNet    *[6][5]*Metrics   `json:"qualityLongNet"`
	QualityShortNet   *[6][5]*Metrics   `json:"qualityShortNet"`
}

//=============================================================================

func NewAnalysisResponse() *AnalysisResponse {
	return &AnalysisResponse{
		QualityAllGross:   &[6][5]*Metrics{},
		QualityLongGross:  &[6][5]*Metrics{},
		QualityShortGross: &[6][5]*Metrics{},
		QualityAllNet:     &[6][5]*Metrics{},
		QualityLongNet:    &[6][5]*Metrics{},
		QualityShortNet:   &[6][5]*Metrics{},
	}
}

//=============================================================================

type Metrics struct {
	Sqn         float64 `json:"sqn"`
	Sqn100      float64 `json:"sqn100"`
	Trades      int     `json:"trades"`
	TradesPerc  float64 `json:"tradesPerc"`
	MaxDrawdown float64 `json:"maxDrawdown"`
}

//=============================================================================
