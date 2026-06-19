//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package genetic

import "github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"

//=============================================================================

type Candidate struct {
	parts []Part
}

//=============================================================================

func NewRandomCandidate(fc *optimization.FilterConfig) *Candidate {
	c := &Candidate{
		parts: []Part{},
	}

	if fc.EnablePosProfit {
		c.parts = append(c.parts, NewPosProfitPart(fc))
	}

	if fc.EnableEquAvg {
		c.parts = append(c.parts, NewEquityVsAvgPart(fc))
	}

	if fc.EnableOldNew {
		c.parts = append(c.parts, NewOldVsNewPart(fc))
	}

	if fc.EnableWinPerc {
		c.parts = append(c.parts, NewWinningPercPart(fc))
	}

	if fc.EnableTrendline {
		c.parts = append(c.parts, NewTrendlinePart(fc))
	}

	return c
}

//=============================================================================

func (c *Candidate) CrossOver(c2 *Candidate) *Candidate {
	//TODO
	return nil
}

//=============================================================================

func (c *Candidate) Mutate() {

}

//=============================================================================
