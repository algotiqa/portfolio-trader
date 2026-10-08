//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
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

func getPortfolioBuilder(c *auth.Context) {
	pId, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			pb, err2 := business.GetPortfolioBuilder(tx, c, pId)

			if err2 != nil {
				return err2
			}

			return c.ReturnObject(pb)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func setPortfolioBuilder(c *auth.Context) {
	pId, err := c.GetIdFromUrl()

	if err == nil {
		spec := business.PortfolioBuilderSpec{}
		err = c.BindParamsFromBody(&spec)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				err = business.SetPortfolioBuilder(tx, c, pId, &spec)

				if err != nil {
					return err
				}

				return c.ReturnObject("")
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================
