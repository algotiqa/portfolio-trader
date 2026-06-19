//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package genetic

import (
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type Population struct {
	size       int
	candidates []*Candidate
}

//=============================================================================
//===
//=== Constructor
//===
//=============================================================================

func NewPopulation(size int, baseline *db.TradingFilter, fc *optimization.FilterConfig) *Population {
	p := &Population{
		size:       size,
		candidates: []*Candidate{},
	}

	for i := 0; i < size; i++ {
		p.Add(NewRandomCandidate(fc))
	}

	return p
}

//=============================================================================
//===
//=== Methods
//===
//=============================================================================

func (p *Population) Select(perc int) *Population {
	//TODO
	return nil
}

//=============================================================================

func (p *Population) Choose() *Candidate {
	//TODO
	return nil
}

//=============================================================================

func (p *Population) Add(c *Candidate) {
	p.candidates = append(p.candidates, c)
}

//=============================================================================

func (p *Population) AverageFitness(num int) float64 {
	//TODO
	return 0
}

//=============================================================================

func (p *Population) sort() {
	//TODO
}

//=============================================================================
