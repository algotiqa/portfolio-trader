//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package service

import (
	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/portfolio-trader/pkg/business"
	"gorm.io/gorm"
)

//=============================================================================

func getDashboardSummary(c *auth.Context) {
	err := dbms.RunInTransaction(func(tx *gorm.DB) error {
		summary, err := business.GetDashboardSummary(tx, c)
		if err != nil {
			return err
		}

		return c.ReturnObject(summary)
	})

	c.ReturnError(err)
}

//=============================================================================
