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

func getAllocations(c *auth.Context) {
	filter := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		var pId int
		pId,err = c.GetParamAsInt("portfolioId", 0)
		if err == nil {
			if pId > 0 {
				filter["portfolio_id"] = pId
			}

			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				list, err := business.GetAllocations(tx, c, filter, offset, limit)

				if err != nil {
					return err
				}

				return c.ReturnList(list, offset, limit, len(*list))
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func getAllocationById(c *auth.Context) {
	id, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ap, errx := business.GetAllocationById(tx, c, id)

			if errx != nil {
				return errx
			}

			return c.ReturnObject(&ap)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func addAllocation(c *auth.Context) {
	spec:= business.AllocationSpec{}
	err := c.BindParamsFromBody(&spec)

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			a,errx := business.AddAllocation(tx, c, &spec)

			if errx != nil {
				return errx
			}

			return c.ReturnObject(a)
		})
	}

	c.ReturnError(err)
}

//=============================================================================
