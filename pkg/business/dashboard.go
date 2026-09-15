//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"sort"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================
//===
//=== Dashboard structure
//===
//=============================================================================

type DashboardSummary struct {
	TradingSystemCount int                       `json:"tradingSystemCount"`
	TotalNetProfit     float64                   `json:"totalNetProfit"`
	TotalTrades        int                       `json:"totalTrades"`
	AllSystems         []*DashboardItem[int]     `json:"allSystems"`
	ByMarket           []*DashboardItem[int]     `json:"byMarket"`
	ByStatus           []*DashboardItem[int]     `json:"byStatus"`
	ByCurrency         []*DashboardItem[float64] `json:"byCurrency"`
	TopSystemsByProfit []*DashboardTopSystem     `json:"topSystemsByProfit"`
	TopSystemsByTrades []*DashboardTopSystem     `json:"topSystemsByTrades"`
}

//=============================================================================

type DashboardItem[T int|float64] struct {
	Name  string `json:"name"`
	Value T      `json:"value"`
}

//=============================================================================

type DashboardTopSystem struct {
	Id    uint    `json:"id"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

//=============================================================================
//===
//=== Dashboard implementation
//===
//=============================================================================

func GetDashboardSummary(tx *gorm.DB, c *auth.Context) (*DashboardSummary, error) {
	filter := map[string]any{}

	//--- Multi-tenant scoping: non-admin users only see their own rows

	if !c.Session.IsAdmin() {
		filter["username"] = c.Session.Username
	}

	//--- Load all matching trading systems (no paging)

	list, err := db.GetTradingSystems(tx, filter, 0, -1)
	if err != nil {
		return nil, req.NewServerErrorByError(err)
	}

	res := &DashboardSummary{
	}
	allSystems := map[string]int{}
	byMarket   := map[string]int{}
	byStatus   := map[db.TsStatus]int{}
	byCurrency := map[string]float64{}

	var activeSystems []*db.TradingSystem

	for _, ts := range *list {
		res.TradingSystemCount++

		if ! ts.Finalized {
			allSystems["DV"]++
		} else {
			if ts.Trading {
				allSystems["TR"]++
				activeSystems = append(activeSystems, &ts)

				res.TotalNetProfit += ts.LastNetProfit
				res.TotalTrades    += ts.LastNumTrades

				byMarket  [ts.MarketType  ]++
				byStatus  [ts.Status      ]++
				byCurrency[ts.CurrencyCode] += ts.LastNetProfit
			} else {
				allSystems["AR"]++
			}
		}
	}

	res.AllSystems = toAllSystems(allSystems)
	res.ByStatus   = toStatus(byStatus)
	res.ByMarket   = toItem(byMarket)
	res.ByCurrency = toItem(byCurrency)

	res.TopSystemsByProfit = buildTopSystemsByProfit(activeSystems)
	res.TopSystemsByTrades = buildTopSystemsByTrades(activeSystems)

	return res, nil
}

//=============================================================================
//---
//--- Helpers
//---
//=============================================================================

func toAllSystems(m map[string]int) []*DashboardItem[int] {
	var res []*DashboardItem[int]
	res = append(res, &DashboardItem[int]{Name: "DV", Value: m["DV"]})
	res = append(res, &DashboardItem[int]{Name: "TR", Value: m["TR"]})
	res = append(res, &DashboardItem[int]{Name: "AR", Value: m["AR"]})

	return res
}

//=============================================================================

func toStatus(m map[db.TsStatus]int) []*DashboardItem[int] {
	var res []*DashboardItem[int]
	res = append(res, &DashboardItem[int]{Name: "0", Value: m[db.TsStatusOff]})
	res = append(res, &DashboardItem[int]{Name: "1", Value: m[db.TsStatusWaiting]})
	res = append(res, &DashboardItem[int]{Name: "2", Value: m[db.TsStatusActive]})
	res = append(res, &DashboardItem[int]{Name: "3", Value: m[db.TsStatusIdle]})
	res = append(res, &DashboardItem[int]{Name: "4", Value: m[db.TsStatusBroken]})

	return res
}

//=============================================================================

func toItem[T int|float64](m map[string]T) []*DashboardItem[T] {
	var res []*DashboardItem[T]

	for k, v := range m {
		res = append(res, &DashboardItem[T]{Name: k, Value: v})
	}

	if res != nil {
		sort.Slice(res, func(i int, j int) bool {
			return res[i].Value > res[j].Value
		})
	}

	return res
}

//=============================================================================

func buildTopSystemsByProfit(list []*db.TradingSystem) []*DashboardTopSystem {
	var res []*DashboardTopSystem

	for _, ts := range list {
		res = append(res, &DashboardTopSystem{
			Id   : ts.Id,
			Name : ts.Name,
			Value: ts.LastNetProfit,
		})
	}

	if res != nil {
		sort.Slice(res, func(i int, j int) bool {
			return res[i].Value > res[j].Value
		})
	}

	if len(res) > 10 {
		res = res[:10]
	}

	return res
}

//=============================================================================

func buildTopSystemsByTrades(list []*db.TradingSystem) []*DashboardTopSystem {
	var res []*DashboardTopSystem

	for _, ts := range list {
		res = append(res, &DashboardTopSystem{
			Id   : ts.Id,
			Name : ts.Name,
			Value: float64(ts.LastNumTrades),
		})
	}

	if res != nil {
		sort.Slice(res, func(i int, j int) bool {
			return res[i].Value > res[j].Value
		})
	}

	if len(res) > 10 {
		res = res[:10]
	}

	return res
}

//=============================================================================
