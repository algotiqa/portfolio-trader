//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package filter

import (
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
)

//=============================================================================

type Combiner struct {
	combinationMode CombinationMode
	criterias       []criteria.FilteringCriteria
	activations     []*criteria.Activation
	outcome         []int
}

//=============================================================================

func NewCombiner(cfg *Config, runMode criteria.RunMode) (*Combiner,error) {
	netRet,err := criteria.NewNetReturnCriteria(&cfg.NetReturn, runMode)
	if err != nil {
		return nil, err
	}

	avgTrade,err := criteria.NewAverageTradeCriteria(&cfg.AverageTrade, runMode)
	if err != nil {
		return nil, err
	}

	equVsAvg,err := criteria.NewEquityVsAverageCriteria(&cfg.EquityVsAverage, runMode)
	if err != nil {
		return nil, err
	}

	winPerc,err := criteria.NewWinningPercentageCriteria(&cfg.WinningPercentage, runMode)
	if err != nil {
		return nil, err
	}

	maxDD,err := criteria.NewMaxDrawdownCriteria(&cfg.MaxDrawdown, runMode)
	if err != nil {
		return nil, err
	}

	trend,err := criteria.NewTrendLineCriteria(&cfg.TrendLine, runMode)
	if err != nil {
		return nil, err
	}

	profFac,err := criteria.NewProfitFactorCriteria(&cfg.ProfitFactor, runMode)
	if err != nil {
		return nil, err
	}

	sharpe,err := criteria.NewSharpeRatioCriteria(&cfg.SharpeRatio, runMode)
	if err != nil {
		return nil, err
	}

	sqn100,err := criteria.NewSqn100Criteria(&cfg.Sqn100, runMode)
	if err != nil {
		return nil, err
	}

	retDD,err := criteria.NewReturnMaxDDRatioCriteria(&cfg.ReturnMaxDDRatio, runMode)
	if err != nil {
		return nil, err
	}

	if cfg.CombinationMode != CombinationModeOne && cfg.CombinationMode != CombinationModeMost && cfg.CombinationMode != CombinationModeAll {
		return nil, req.NewBadRequestError("unknown filter combination mode: %v", cfg.CombinationMode)
	}

	c :=&Combiner{
		combinationMode: cfg.CombinationMode,
	}

	c.criterias = append(c.criterias, netRet, avgTrade, equVsAvg, winPerc, maxDD, trend, profFac, sharpe, sqn100, retDD)

	return c, nil
}

//=============================================================================

func (c *Combiner) Eval(s *criteria.Snapshot) {
	for _,cr := range c.criterias {
		cr.Eval(s)
		a := cr.Activation()
		if a != nil {
			c.activations = append(c.activations, a)
		}
	}

	c.combineActivations(s)
}

//=============================================================================

func (c *Combiner) Activations() []*criteria.Activation {
	return c.activations
}

//=============================================================================

func (c *Combiner) Outcome() []int {
	return c.outcome
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func (c *Combiner) combineActivations(s *criteria.Snapshot) {
	c.outcome = make([]int, len(s.Time))

	for i, _ := range s.Time {
		c.outcome[i] = c.calcActivationAt(i)
	}
}

//=============================================================================

func (c *Combiner) calcActivationAt(i int) int {
	outcomes := 0
	size     := len(c.activations)

	if size == 0 {
		return 0
	}

	for _, a := range c.activations {
		outcomes += a.Active(i)
	}

	if c.combinationMode == CombinationModeOne && outcomes > 0 {
		return 1
	}

	if c.combinationMode == CombinationModeAll && outcomes == size {
		return 1
	}

	if c.combinationMode == CombinationModeMost && outcomes >= size/2 +1 {
		return 1
	}

	return 0
}

//=============================================================================
