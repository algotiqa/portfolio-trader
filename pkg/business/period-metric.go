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
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

type PeriodMetricResponse struct {
	TradingSystems []*TradingSystemMetric `json:"tradingSystems"`
}

//=============================================================================

type TradingSystemMetric struct {
	Name      string        `json:"name"`
	Running   bool          `json:"running"`
	Metrics3m *core.Metrics `json:"metrics3m"`
	Metrics4m *core.Metrics `json:"metrics4m"`
	Metrics6m *core.Metrics `json:"metrics6m"`
	Metrics9m *core.Metrics `json:"metrics9m"`
	Metrics1y *core.Metrics `json:"metrics1y"`
	Metrics5y *core.Metrics `json:"metrics5y"`
}

//=============================================================================

func ComputeMetrics(tx *gorm.DB, c *auth.Context, portfolioId string) (*PeriodMetricResponse, error) {
	filter := map[string]any{}
	filter["username" ] = c.Session.Username
	filter["finalized"] = true
	filter["trading"  ] = true

	if portfolioId != "" {
		filter["portfolio_id"] = portfolioId
	}

	res := &PeriodMetricResponse{}
	list,err := db.GetTradingSystems(tx, filter, 0, 5000)

	if err == nil {
		for _, ts := range *list {
			var trades *[]db.Trade
			trades,err = db.FindTradesByTradingSystemId(tx, ts.Id)
			if err != nil {
				break
			}

			tsm := computeMetrics(&ts, trades)
			res.TradingSystems = append(res.TradingSystems, tsm)
		}
	}

	return res, err
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func computeMetrics(ts *db.TradingSystem, trades *[]db.Trade) *TradingSystemMetric {
	core.NormalizeContracts(trades)

	tsm := &TradingSystemMetric{
		Name     : ts.Name,
		Running  : ts.Running,
		Metrics3m: computeMetricsOnPeriod(ts, trades, 3*30),
		Metrics4m: computeMetricsOnPeriod(ts, trades, 4*30),
		Metrics6m: computeMetricsOnPeriod(ts, trades, 6*30),
		Metrics9m: computeMetricsOnPeriod(ts, trades, 9*30),
		Metrics1y: computeMetricsOnPeriod(ts, trades, 1*365),
		Metrics5y: computeMetricsOnPeriod(ts, trades, 5*365),
	}

	return tsm
}

//=============================================================================

func computeMetricsOnPeriod(ts *db.TradingSystem, trades *[]db.Trade, daysBack int) *core.Metrics {
	back := time.Duration(daysBack) * 24 * time.Hour
	date := time.Now().Add(-back)

	return core.CalcMetrics(trades, ts.CostPerOperation, &date)
}

//=============================================================================
