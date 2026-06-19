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

type geneticAlgorithm struct {
	ctx    optimization.Context
	config *optimization.GeneticConfig
}

//=============================================================================

func New() optimization.Algorithm {
	return &geneticAlgorithm{}
}

//=============================================================================
//===
//=== Genetic algorithm implementation
//===
//=============================================================================

func (ga *geneticAlgorithm) Init(ctx optimization.Context) {
	ga.ctx = ctx
	ga.config = &ctx.AlgorithmConfig().Genetic

	//--- TODO Validate config and return an error in case of errors

	ga.config = &optimization.GeneticConfig{
		PopulationSize: 1000,
		MinSteps:       1000,
	}
}

//=============================================================================

func (ga *geneticAlgorithm) StepsCount() uint {
	return ga.config.MinSteps
}

//=============================================================================

func (ga *geneticAlgorithm) Optimize() {

}

//=============================================================================
