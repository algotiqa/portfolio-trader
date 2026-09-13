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

func getAssignableTradingSystems(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			list, errt := business.GetAssignableTradingSystems(tx, c, id)

			if errt != nil {
				return errt
			}

			return c.ReturnList(list, 0, 5000, len(*list))
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getAssignedTradingSystems(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			list, errt := business.GetAssignedTradingSystems(tx, c, id)

			if errt != nil {
				return errt
			}

			return c.ReturnList(list, 0, 5000, len(*list))
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func assignTradingSystemsToPortfolio(c *auth.Context) {
	var list []uint
	err := c.BindParamsFromBody(&list)

	if err == nil {
		var id uint
		id, err = c.GetIdFromUrl()
		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				errt := business.AssignTradingSystemsToPortfolio(tx, c, id, list)

				if errt != nil {
					return errt
				}

				return c.ReturnObject("")
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func unassignTradingSystemsFromPortfolio(c *auth.Context) {
	var list []uint
	err := c.BindParamsFromBody(&list)

	if err == nil {
		var id uint
		id, err = c.GetIdFromUrl()
		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				errt := business.UnassignTradingSystemsFromPortfolio(tx, c, id, list)

				if errt != nil {
					return errt
				}

				return c.ReturnObject("")
			})
		}
	}

	c.ReturnError(err)
}

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
