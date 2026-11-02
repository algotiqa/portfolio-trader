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
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type Evaluator struct {
	activCombiner   *Combiner
	deactCombiner   *Combiner
	combinedFinal   []int
	sampledFinal    []int
	finalActivation *criteria.Activation
}

//=============================================================================

func NewEvaluator(f *TradingFilter) (*Evaluator,error) {
	activComb, err := NewCombiner(&f.Activation, criteria.RunModeActivation)
	if err != nil {
		return nil, err
	}

	deactComb, err := NewCombiner(&f.Deactivation, criteria.RunModeDeactivation)
	if err != nil {
		return nil, err
	}

	return &Evaluator{
		activCombiner: activComb,
		deactCombiner: deactComb,
	}, nil
}

//=============================================================================

func (e *Evaluator) Eval(s *criteria.Snapshot, p *db.Portfolio) {
	e.activCombiner.Eval(s)
	e.deactCombiner.Eval(s)
	e.combineOutcomes()

	sampler := NewPeriodicSampler(p)
	e.sampledFinal, e.finalActivation = sampler.Sample(s.Time, e.combinedFinal)
}

//=============================================================================

func (e *Evaluator) Outcomes() *Outcomes {
	return &Outcomes{
		Activations     : e.activCombiner.activations,
		Deactivations   : e.deactCombiner.activations,
		CombinedActiv   : e.activCombiner.outcome,
		CombinedDeact   : e.deactCombiner.outcome,
		CombinedFinal   : e.combinedFinal,
		SampledFinal    : e.sampledFinal,
		FinalActivation : e.finalActivation,
	}
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func (e *Evaluator) combineOutcomes() {
	activ := e.activCombiner.outcome
	deact := e.deactCombiner.outcome

	e.combinedFinal = make([]int, len(activ))

	prev := 1

	for i,_ := range activ {
		activValue := activ[i]
		deactValue := deact[i]

		if deactValue == 1 {
			e.combinedFinal[i] = 0
			prev         = 0
		} else if activValue == 1 {
			e.combinedFinal[i] = 1
			prev         = 1
		} else {
			e.combinedFinal[i] = prev
		}
	}
}

//=============================================================================
