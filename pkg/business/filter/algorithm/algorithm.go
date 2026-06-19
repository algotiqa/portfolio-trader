//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package algorithm

import (
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/genetic"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/optimization"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/algorithm/simple"
)

//=============================================================================

const Simple  = "simple"
const Genetic = "genetic"

//=============================================================================

func New(name string) optimization.Algorithm {
	switch name {
	case Simple:
		return simple.New()

	case Genetic:
		return genetic.New()

	default:
		panic("Unknown optimization algorithm : " + name)
	}
}

//=============================================================================
