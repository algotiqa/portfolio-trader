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
	"errors"

	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type PortfolioBuilder interface {
	Name()   db.BuilderName
	Config() map[string]any
	Specs()  map[string]any
	Init(config map[string]any) error

	//PositionInit(ts *TradingSnapshot)
	//PositionFor(ts *TradingSnapshot) int
}

//=============================================================================

func New(name db.BuilderName, config map[string]any) (PortfolioBuilder, error) {
	var b PortfolioBuilder
	switch name {
		case db.BuilderNone:
			b = NewNoneBuilder()

		case db.BuilderSimple:
			b = NewSimpleBuilder()

		default:
			return nil, errors.New("Unknown portfolio builder name: " + string(name))
	}

	return b, b.Init(config)
}

//=============================================================================

