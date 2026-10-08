//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package builder

import (
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Specs
//===
//=============================================================================

var DefCorrThresh = 0.7

var SpecCorrThresh = core.NewNumberParamSpec[float64]( "corrThresh", true, 0.01, 0.99, &DefCorrThresh)

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type SimpleConfig struct {
	corrThresh float64
}

//=============================================================================
//===
//=== Builder
//===
//=============================================================================

type SimpleBuilder struct {
	config *SimpleConfig
}

//=============================================================================

func NewSimpleBuilder() *SimpleBuilder {
	return &SimpleBuilder{
		config: &SimpleConfig{
			corrThresh: DefCorrThresh,
		},
	}
}

//=============================================================================

func NewSimpleBuilderWithParams(corrThresh float64) *SimpleBuilder {
	return &SimpleBuilder{
		config: &SimpleConfig{
			corrThresh: corrThresh,
		},
	}
}

//=============================================================================

func (b *SimpleBuilder) Name() db.BuilderName {
	return db.BuilderSimple
}

//=============================================================================

func (b *SimpleBuilder) Init(config map[string]any) error {
	corrThresh,err := core.MapNumber[float64](config, SpecCorrThresh)
	if err != nil {
		return err
	}

	b.config.corrThresh = *corrThresh
	return nil
}

//=============================================================================

func (b *SimpleBuilder) Config() map[string]any {
	cfg := make(map[string]any)
	cfg[SpecCorrThresh.Name] = b.config.corrThresh

	return cfg
}

//=============================================================================

func (b *SimpleBuilder) Specs() map[string]any {
	specs := make(map[string]any)

	specs[SpecCorrThresh.Name] = SpecCorrThresh

	return specs
}

//=============================================================================
