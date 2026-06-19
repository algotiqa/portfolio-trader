//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import "github.com/algotiqa/portfolio-trader/pkg/core"

//=============================================================================
//===
//=== AnalysisRequest
//===
//=============================================================================

type AnalysisRequest struct {
	core.SelectedPeriod
	Filter    *TradingFilter  `json:"filter,omitempty"`
}

//=============================================================================

type TradingFilter struct {
	EquAvgEnabled    bool   `json:"equAvgEnabled"`
	EquAvgLen        int    `json:"equAvgLen"`
	PosProEnabled    bool   `json:"posProEnabled"`
	PosProLen        int    `json:"posProLen"`
	WinPerEnabled    bool   `json:"winPerEnabled"`
	WinPerLen        int    `json:"winPerLen"`
	WinPerValue      int    `json:"winPerValue"`
	OldNewEnabled    bool   `json:"oldNewEnabled"`
	OldNewOldLen     int    `json:"oldNewOldLen"`
	OldNewOldPerc    int    `json:"oldNewOldPerc"`
	OldNewNewLen     int    `json:"oldNewNewLen"`
	TrendlineEnabled bool   `json:"trendlineEnabled"`
	TrendlineLen     int    `json:"trendlineLen"`
	TrendlineValue   int    `json:"trendlineValue"`
	DrawdownEnabled  bool   `json:"drawdownEnabled"`
	DrawdownMin      int    `json:"drawdownMin"`
	DrawdownMax      int    `json:"drawdownMax"`
}

//=============================================================================
