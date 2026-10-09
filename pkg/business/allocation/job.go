//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package allocation

import (
	"fmt"

	"github.com/algotiqa/portfolio-trader/pkg/business/allocation/builder"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Public
//===
//=============================================================================

func NewJob(spec *JobSpec) *Job {
	return &Job{
		spec: spec,
	}
}

//=============================================================================

func (j *Job) SetData(p *db.Portfolio, b builder.PortfolioBuilder) {
	j.portfolio = p
	j.builder   = b
}

//=============================================================================

func (j *Job) AddTradingSystem(tsi *TradingSystemInfo) {
	tsi.healthy = true

	ts := &TradingSystem{
		Id          : tsi.system.Id,
		Name        : tsi.system.Name,
		DataSymbol  : tsi.system.DataSymbol,
		BrokerSymbol: tsi.system.BrokerSymbol,
		MarketType  : tsi.system.MarketType,
		StrategyType: tsi.system.StrategyType,
		HealthLevel : HealthLevelOk,
	}

	if len(*tsi.returns) < 2 {
		ts.HealthLevel = HealthLevelError
		ts.Message     = format("Missing daily returns to calculate correlations")
		tsi.healthy    = false
		j.addExcludedSystem(tsi, ts.Message)

	} else {
		if len(*tsi.returns) < MinDailyReturns {
			ts.HealthLevel = HealthLevelWarn
			ts.Message     = format("Insufficient daily returns to calculate good correlations: %v < %v", len(*tsi.returns), MinDailyReturns)
		}
	}

	j.tradingSystems = append(j.tradingSystems, tsi)
	j.healthySystems = append(j.healthySystems, ts)
}

//=============================================================================
//===
//=== Auxiliary
//===
//=============================================================================

func (j *Job) log(level LogLevel, message string, params ...any) {
	log := &Log{
		Level  : level,
		Message: format(message, params...),
	}

	j.logs = append(j.logs, log)
}

//=============================================================================

func (j *Job) addFilter(tsi *TradingSystemInfo, filterPassed bool, comment string) {
	af := &FilterOutcome{
		TsId          : tsi.system.Id,
		TsName        : tsi.system.Name,
		TsDataSymbol  : tsi.system.DataSymbol,
		TsBrokerSymbol: tsi.system.BrokerSymbol,
		TsMarketType  : tsi.system.MarketType,
		TsStrategyType: tsi.system.StrategyType,
		TsRunning     : tsi.system.Running,
		FilterPassed  : filterPassed,
		Action        : ActionNone,
		Comment       : comment,
	}

	if af.TsRunning && !af.FilterPassed {
		af.Action = ActionTurnOff
	} else if !af.TsRunning && af.FilterPassed {
		af.Action = ActionTurnOn
	}

	j.filterOutcomes = append(j.filterOutcomes, af)
}

//=============================================================================

func (j *Job) addExcludedSystem(tsi *TradingSystemInfo, message string) {
	es := NewExcludedSystem(tsi, message)
	j.excludedSystems = append(j.excludedSystems, es)
	j.log(LogLevelInfo, "Excluding '%v'. Reason: %v", tsi.system.Name, message)
}

//=============================================================================

func (j *Job) clearSystemCorrelations() {
	j.correlations = nil
}

//=============================================================================

func (j *Job) addSystemCorrelation(tsi1, tsi2 *TradingSystemInfo, correlation float64) {
	sc := &SystemCorrelation{
		Ts1Id      : tsi1.system.Id,
		Ts1Name    : tsi1.system.Name,
		Ts2Id      : tsi2.system.Id,
		Ts2Name    : tsi2.system.Name,
		Correlation: core.Trunc2d(correlation),
	}

	j.correlations = append(j.correlations, sc)
}

//=============================================================================

func format(message string, params ...any) string {
	if len(params) > 0 && params[0] != nil {
		return fmt.Sprintf(message, params...)
	}

	return message
}

//=============================================================================
