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
	"encoding/json"
	"errors"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

type PortfolioBuilder interface {
	Name()   db.BuilderName
	Config() map[string]any
	Specs()  map[string]any
	Init(config map[string]any) error

	TunePosition(position float64, correlations []float64) float64
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

func Create(name db.BuilderName, config string) (PortfolioBuilder,error) {
	cfgMap := make(map[string]any)
	err := json.Unmarshal([]byte(config), &cfgMap)
	if err != nil {
		return nil,req.NewServerErrorByError(err)
	}

	b,err := New(name, cfgMap)
	if err != nil {
		return nil,req.NewServerErrorByError(err)
	}

	return b, nil
}

//=============================================================================
