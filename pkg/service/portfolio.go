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

func getPortfolios(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetPortfolios(tx, c, filter, offset, limit)

			if err != nil {
				return err
			}

			return c.ReturnList(list, offset, limit, len(*list))
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getPortfolioTree(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetPortfolioTree(tx, c, filter, offset, limit)

			if err != nil {
				return err
			}

			return c.ReturnObject(list)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getPortfolioMonitoring(c *auth.Context) {
	params := business.PortfolioMonitoringParams{}
	err := c.BindParamsFromBody(&params)

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			result, err := business.GetPortfolioMonitoring(tx, &params)

			if err != nil {
				return err
			}

			return c.ReturnObject(result)
		})
	}

	c.ReturnError(err)
}

//=============================================================================
