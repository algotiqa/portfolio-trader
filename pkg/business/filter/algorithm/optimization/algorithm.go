//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package optimization

import "github.com/algotiqa/portfolio-trader/pkg/db"

//=============================================================================

type Context interface {
	FilterConfig() *FilterConfig
	AlgorithmConfig() *AlgorithmConfig
	IsStopping() bool
	Baseline() db.TradingFilter

	RunAnalysis(filter *db.TradingFilter) float64
	LogInfo(message string)
}

//=============================================================================

type Algorithm interface {
	Init(ctx Context)
	Optimize()
	StepsCount() uint
}

//=============================================================================
//===
//=== Algorithm configurations
//===
//=============================================================================

type SimpleConfig struct {
}

//=============================================================================

type GeneticConfig struct {
	PopulationSize int  `json:"populationSize"`
	MinSteps       uint `json:"minSteps"`
}

//=============================================================================

type AlgorithmConfig struct {
	Simple  SimpleConfig  `json:"simple"`
	Genetic GeneticConfig `json:"genetic"`
}

//=============================================================================
