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
	"log/slog"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/auth/roles"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/app"
	"github.com/gin-gonic/gin"
)

//=============================================================================

func Init(router *gin.Engine, cfg *app.Config, logger *slog.Logger) {

	ctrl := auth.NewOidcController(cfg.Authentication.Authority, req.GetDefaultClient(), logger, cfg)

	router.GET   ("/api/portfolio/v1/trading-systems",                           ctrl.Secure(getTradingSystems,          roles.Admin_User_Service))
	router.GET   ("/api/portfolio/v1/trading-systems/:id",                       ctrl.Secure(getTradingSystem,           roles.Admin_User_Service))
	router.GET   ("/api/portfolio/v1/trading-systems/:id/trades",                ctrl.Secure(getTrades,                  roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/filters",               ctrl.Secure(setTradingFilters,          roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/position",              ctrl.Secure(setTradingPosition,         roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/trading",               ctrl.Secure(setTradingSystemTrading,    roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/running",               ctrl.Secure(setTradingSystemRunning,    roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/activation",            ctrl.Secure(setTradingSystemActivation, roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/active",                ctrl.Secure(setTradingSystemActive,     roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/performance-analysis",  ctrl.Secure(runPerformanceAnalysis,     roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/quality-analysis",      ctrl.Secure(runQualityAnalysis,         roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/trade-analysis",        ctrl.Secure(runTradeAnalysis,           roles.Admin_User_Service))

	router.POST  ("/api/portfolio/v1/trading-systems/:id/filter-analysis",       ctrl.Secure(runFilterAnalysis,          roles.Admin_User_Service))
	router.GET   ("/api/portfolio/v1/trading-systems/:id/filter-optimization",   ctrl.Secure(getFilterOptimizationInfo,  roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/filter-optimization",   ctrl.Secure(startFilterOptimization,    roles.Admin_User_Service))
	router.DELETE("/api/portfolio/v1/trading-systems/:id/filter-optimization",   ctrl.Secure(stopFilterOptimization,     roles.Admin_User_Service))

	router.GET   ("/api/portfolio/v1/trading-systems/:id/simulation",            ctrl.Secure(getSimulationResult,        roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/trading-systems/:id/simulation",            ctrl.Secure(startSimulation,            roles.Admin_User_Service))
	router.DELETE("/api/portfolio/v1/trading-systems/:id/simulation",            ctrl.Secure(stopSimulation,             roles.Admin_User_Service))

	router.POST  ("/api/portfolio/v1/trading-systems/:id/position-analysis",     ctrl.Secure(runPositionAnalysis,        roles.Admin_User_Service))

	router.GET   ("/api/portfolio/v1/trading-systems/export",                    ctrl.Secure(exportTradingSystems,       roles.Admin_User_Service))

	router.GET   ("/api/inventory/v1/portfolios",                                ctrl.Secure(getPortfolios,              roles.Admin_User_Service))
	router.GET   ("/api/inventory/v1/portfolio/tree",                            ctrl.Secure(getPortfolioTree,           roles.Admin_User_Service))
	router.POST  ("/api/portfolio/v1/portfolio/monitoring",                      ctrl.Secure(getPortfolioMonitoring,     roles.Admin_User_Service))
}

//=============================================================================

func NewStatusOkResponse() any {
	return struct {
		status string
	}{status: "ok"}
}

//=============================================================================
