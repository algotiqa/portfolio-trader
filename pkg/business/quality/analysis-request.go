//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package quality

import "github.com/algotiqa/portfolio-trader/pkg/core"

//=============================================================================

type AnalysisRequest struct {
	core.SelectedPeriod
	TimeframeType string `json:"timeframeType"`
	AtrLength     int    `json:"atrLength"     binding:"min=5,max=50"`
}

//=============================================================================
