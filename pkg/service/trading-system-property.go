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

func setTradingSystemTrading(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemTradingRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemTrading(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingSystemRunning(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemRunningRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemRunning(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingSystemActivation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemActivationRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemActivation(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingSystemActive(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		req := business.TradingSystemActiveRequest{}
		err = c.BindParamsFromBody(&req)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.SetTradingSystemActive(tx, c, tsId, &req)

				if err != nil {
					return err
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================
