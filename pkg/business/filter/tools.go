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
	"encoding/json"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/business/filter/model"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

func ConvertToDbFilter(f *model.TradingFilter) (*db.TradingFilter,error) {
	config,err := json.Marshal(f)
	if err != nil {
		return nil,req.NewServerErrorByError(err)
	}

	return &db.TradingFilter{
		Config: string(config),
	},nil
}

//=============================================================================

func ConvertFromDbFilter(dbf *db.TradingFilter) (*model.TradingFilter,error) {
	var f model.TradingFilter
	err := json.Unmarshal([]byte(dbf.Config), &f)
	if err != nil {
		return nil,req.NewServerErrorByError(err)
	}

	return &f,nil
}

//=============================================================================
