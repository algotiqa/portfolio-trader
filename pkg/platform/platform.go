//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package platform

import (
	"github.com/algotiqa/core"
	"github.com/algotiqa/portfolio-trader/pkg/app"
)

//=============================================================================

var platform *core.Platform

//=============================================================================
//===
//=== Init function
//===
//=============================================================================

func InitPlatform(cfg *app.Config) {
	platform = &cfg.Platform
}

//=============================================================================
