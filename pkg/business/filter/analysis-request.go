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
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
)

//=============================================================================
//===
//=== AnalysisRequest
//===
//=============================================================================

type AnalysisRequest struct {
	Period  core.SelectedPeriod   `json:"period"`
	Filter  *model.TradingFilter  `json:"filter,omitempty"`
}

//=============================================================================
