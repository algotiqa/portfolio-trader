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
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================
//===
//=== Specs
//===
//=============================================================================

//=============================================================================
//===
//=== Config
//===
//=============================================================================

type NoneConfig struct {
}

//=============================================================================
//===
//=== Builder
//===
//=============================================================================

type NoneBuilder struct {
	config *NoneConfig
}

//=============================================================================

func NewNoneBuilder() *NoneBuilder {
	return &NoneBuilder{
		config: &NoneConfig{
		},
	}
}

//=============================================================================

func (b *NoneBuilder) Name() db.BuilderName {
	return db.BuilderNone
}

//=============================================================================

func (b *NoneBuilder) Init(config map[string]any) error {
	return nil
}

//=============================================================================

func (b *NoneBuilder) Config() map[string]any {
	cfg := make(map[string]any)
	return cfg
}

//=============================================================================

func (b *NoneBuilder) Specs() map[string]any {
	specs := make(map[string]any)
	return specs
}

//=============================================================================
