//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"encoding/json"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/allocation/builder"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================
//===
//=== Model
//===
//=============================================================================

type PortfolioBuilder struct {
	Name   db.BuilderName `json:"name"`
	Config map[string]any `json:"config"`
	Specs  map[string]any `json:"specs"`
}

//=============================================================================

type PortfolioBuilderSpec struct {
	Name   db.BuilderName `json:"name"`
	Config map[string]any `json:"config"`
}

//=============================================================================

type PortfolioBuilderResponse struct {
	Current   *PortfolioBuilder     `json:"current"`
	Available []*PortfolioBuilder   `json:"available"`
}

//=============================================================================

func (r *PortfolioBuilderResponse) Add(pb builder.PortfolioBuilder) {
	r.Available = append(r.Available, createBuilderInfo(pb))
}

//=============================================================================
//===
//=== Logic
//===
//=============================================================================

func GetPortfolioBuilder(tx *gorm.DB, c *auth.Context, id uint) (*PortfolioBuilderResponse,error) {
	p,err := getPortfolio(tx, c, id, "GetPortfolioBuilder")
	if err != nil {
		return nil,err
	}

	b,err := builder.Create(p.BuilderName, p.BuilderConfig)
	if err != nil {
		return nil,err
	}

	pbr := &PortfolioBuilderResponse{
		Current: createBuilderInfo(b),
	}

	pbr.Add(builder.NewNoneBuilder())
	pbr.Add(builder.NewSimpleBuilder())

	return pbr,nil
}

//=============================================================================

func SetPortfolioBuilder(tx *gorm.DB, c *auth.Context, id uint, spec *PortfolioBuilderSpec) error {
	p,err := getPortfolio(tx, c, id, "SetPortfolioBuilder")
	if err != nil {
		return err
	}

	b,err := builder.New(spec.Name, spec.Config)
	if err != nil {
		return req.NewBadRequestError("Invalid config for builder: %v", err.Error())
	}

	data,err := json.Marshal(b.Config())
	if err != nil {
		return req.NewServerErrorByError(err)
	}

	p.BuilderName   = spec.Name
	p.BuilderConfig = string(data)

	return db.UpdatePortfolio(tx, p)
}

//=============================================================================
//===
//=== Private
//===
//=============================================================================

func createBuilderInfo(b builder.PortfolioBuilder) *PortfolioBuilder {
	return &PortfolioBuilder{
		Name  : b.Name(),
		Config: b.Config(),
		Specs : b.Specs(),
	}
}

//=============================================================================
