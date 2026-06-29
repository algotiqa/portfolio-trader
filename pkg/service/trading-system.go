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
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter"
	"github.com/algotiqa/portfolio-trader/pkg/business/performance"
	"github.com/algotiqa/portfolio-trader/pkg/business/position"
	"github.com/algotiqa/portfolio-trader/pkg/business/quality"
	"github.com/algotiqa/portfolio-trader/pkg/business/simulation"
	"github.com/algotiqa/portfolio-trader/pkg/business/trade"
	"gorm.io/gorm"
)

//=============================================================================

func getTradingSystems(c *auth.Context) {
	filt := map[string]any{}
	offset, limit, err := c.GetPagingParams()

	if err == nil {
		details, err := c.GetParamAsBool("details", false)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				list, err := business.GetTradingSystems(tx, c, filt, offset, limit, details)

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

func getTradingSystem(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			ts, err2 := business.GetTradingSystem(tx, c, tsId)

			if err2 != nil {
				return err2
			}

			return c.ReturnObject(&ts)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func getTrades(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = dbms.RunInTransaction(func(tx *gorm.DB) error {
			list, err := business.GetTrades(tx, c, tsId)

			if err != nil {
				return err
			}

			return c.ReturnObject(&list)
		})
	}

	c.ReturnError(err)
}

//=============================================================================

func setTradingFilters(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		filters := filter.TradingFilter{}
		err = c.BindParamsFromBody(&filters)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				err = business.SetTradingFilters(tx, c, tsId, &filters)

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

func runFilterAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		rq := filter.AnalysisRequest{}
		err = c.BindParamsFromBody(&rq)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				rep, errx := business.RunFilterAnalysis(tx, c, tsId, &rq)

				if errx != nil {
					return errx
				}

				return c.ReturnObject(rep)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func runPerformanceAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		rq := performance.AnalysisRequest{}
		err = c.BindParamsFromBody(&rq)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				res, errx := business.RunPerformanceAnalysis(tx, c, tsId, &rq)

				if errx != nil {
					return errx
				}

				return c.ReturnObject(res)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func runQualityAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		rq := quality.AnalysisRequest{}
		err = c.BindParamsFromBody(&rq)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				res, errx := business.RunQualityAnalysis(tx, c, tsId, &rq)

				if errx != nil {
					return errx
				}

				return c.ReturnObject(res)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func runTradeAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		rq := trade.AnalysisRequest{}
		err = c.BindParamsFromBody(&rq)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				res, errx := business.RunTradeAnalysis(tx, c, tsId, &rq)

				if errx != nil {
					return errx
				}

				return c.ReturnObject(res)
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func exportTradingSystems(c *auth.Context) {
	ids,err := c.GetIdsFromUrl()
	if err == nil {
		if len(ids) == 0 {
			err = req.NewBadRequestError("Parameter 'id' is missing or empty")
		} else {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				res, terr := business.ExportTradingSystems(tx, c, ids)
				if terr == nil {
					return c.ReturnObject(res)
				}
				return terr
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================
//===
//=== Filter optimization
//===
//=============================================================================

func startFilterOptimization(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		rq := filter.OptimizationRequest{}
		err = c.BindParamsFromBody(&rq)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				err := business.StartFilterOptimization(tx, c, tsId, &rq)

				if err != nil {
					return err
				}

				return c.ReturnObject(NewStatusOkResponse())
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func stopFilterOptimization(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		err = business.StopFilterOptimization(c, tsId)

		if err != nil {
			c.ReturnError(err)
		} else {
			_ = c.ReturnObject(NewStatusOkResponse())
		}

		return
	}

	c.ReturnError(err)
}

//=============================================================================

func getFilterOptimizationInfo(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		res, err := business.GetFilterOptimizationInfo(c, tsId)

		if err != nil {
			c.ReturnError(err)
		} else {
			_ = c.ReturnObject(res)
		}

		return
	}

	c.ReturnError(err)
}

//=============================================================================
//===
//=== Simulation
//===
//=============================================================================

func startSimulation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		rq := simulation.Request{}
		err = c.BindParamsFromBody(&rq)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				err2 := business.StartSimulation(tx, c, tsId, &rq)

				if err2 != nil {
					return err2
				}

				return c.ReturnObject(NewStatusOkResponse())
			})
		}
	}

	c.ReturnError(err)
}

//=============================================================================

func stopSimulation(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		_ = business.StopSimulation(c, tsId)
		_ = c.ReturnObject(NewStatusOkResponse())
	}

	c.ReturnError(err)
}

//=============================================================================

func getSimulationResult(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		res := business.GetSimulationResult(c, tsId)
		_ = c.ReturnObject(res)
		return
	}

	c.ReturnError(err)
}

//=============================================================================
//===
//=== Position sizing
//===
//=============================================================================

func setTradingPosition(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		pos := position.TradingPosition{}
		err = c.BindParamsFromBody(&pos)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				err = business.SetTradingPosition(tx, c, tsId, &pos)

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

func runPositionAnalysis(c *auth.Context) {
	tsId, err := c.GetIdFromUrl()

	if err == nil {
		rq := position.AnalysisRequest{}
		err = c.BindParamsFromBody(&rq)

		if err == nil {
			err = dbms.RunInTransaction(func(tx *gorm.DB) error {
				rep, err := business.RunPositionAnalysis(tx, c, tsId, &rq)

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
