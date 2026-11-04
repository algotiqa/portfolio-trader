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
	"testing"

	"github.com/algotiqa/portfolio-trader/pkg/db"
	"golang.org/x/exp/slices"

	"github.com/algotiqa/portfolio-trader/pkg/business/filter/criteria"
)

//=============================================================================

func hysteresisFilter() *TradingFilter {
	return &TradingFilter{
		Activation: Config{
			CombinationMode: CombinationModeOne,
			NetReturn:       criteria.NetReturnConfig{Enabled: true, Length: 1, Value: 0},
		},
		Deactivation: Config{
			CombinationMode: CombinationModeOne,
			MaxDrawdown:     criteria.MaxDrawdownConfig{Enabled: true, Value: -5},
		},
	}
}

//=============================================================================

func TestEvaluatorHysteresis(t *testing.T) {
	eval, err := NewEvaluator(hysteresisFilter())
	if err != nil {
		t.Fatalf("cannot create evaluator: %v", err)
	}

	s := testSnapshot([]float64{100, 110, 90, 95}, nil, []float64{-6, -1, -6, -1})

	eval.Eval(s, &db.Portfolio{
		PeriodType: db.PeriodTypeWeekly,
	})

	expected := []int{0, 1, 0, 1}
	combined := eval.Outcomes().CombinedFinal

	if !slices.Equal(combined, expected) {
		t.Errorf("Bad combined outcome. Expected %v but got %v", expected, combined)
	}
}

//=============================================================================

func TestEvaluatorStartsActive(t *testing.T) {
	f := &TradingFilter{
		Activation:   Config{CombinationMode: CombinationModeOne},
		Deactivation: Config{CombinationMode: CombinationModeOne},
	}

	eval, err := NewEvaluator(f)
	if err != nil {
		t.Fatalf("cannot create evaluator: %v", err)
	}

	s := testSnapshot([]float64{1, 2}, []float64{0, 0}, []float64{0, 0})

	eval.Eval(s, &db.Portfolio{
		PeriodType: db.PeriodTypeWeekly,
	})

	if !slices.Equal(eval.Outcomes().CombinedFinal, []int{1, 1}) {
		t.Errorf("The evaluator must start active, got %v", eval.Outcomes().CombinedFinal)
	}
}

//=============================================================================

func TestEvaluatorOutcomes(t *testing.T) {
	eval, err := NewEvaluator(hysteresisFilter())
	if err != nil {
		t.Fatalf("cannot create evaluator: %v", err)
	}

	s := testSnapshot([]float64{100, 110, 90, 95}, nil, []float64{-6, -1, -6, -1})

	eval.Eval(s, &db.Portfolio{
		PeriodType: db.PeriodTypeWeekly,
	})

	out := eval.Outcomes()

	if len(out.Activations) != 1 || out.Activations[0].Name != criteria.NetReturn {
		t.Errorf("Bad activations: %v", out.Activations)
	}

	if len(out.Deactivations) != 1 || out.Deactivations[0].Name != criteria.MaxDrawdown {
		t.Errorf("Bad deactivations: %v", out.Deactivations)
	}
}

//=============================================================================
