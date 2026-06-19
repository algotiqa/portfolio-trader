//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package simulation

import "github.com/algotiqa/portfolio-trader/pkg/core"

//=============================================================================

type Request struct {
	core.SelectedPeriod
	Runs   int `json:"runs"   binding:"max=50000"`
	Width  int `json:"width"  binding:"max=4000"`
	Height int `json:"height" binding:"max=3000"`
}

//=============================================================================
