//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func getTradingSystemAndCheckAccess(tx *gorm.DB, c *auth.Context, id uint) (*db.TradingSystem, error) {
	ts, err := db.GetTradingSystemById(tx, id)
	if err != nil {
		c.Log.Error("getTradingSystem: Cannot get the trading system", "id", id, "error", err)
		return nil, err
	}

	if ts == nil {
		return nil, req.NewNotFoundError("Trading system was not found: %v", id)
	}

	if !c.Session.IsAdmin() {
		if ts.Username != c.Session.Username {
			return nil, req.NewForbiddenError("Trading system not owned by user: %v", id)
		}
	}

	return ts, nil
}

//=============================================================================
