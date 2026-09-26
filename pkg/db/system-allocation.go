//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package db

import (
	"github.com/algotiqa/core/req"
	"gorm.io/gorm"
)

//=============================================================================

func GetSystemCorrelationsByAllocationId(tx *gorm.DB, allId uint) (*[]SystemCorrelationFull, error) {
	var list []SystemCorrelationFull

	filter := map[string]any{}
	filter["allocation_id"] = allId

	res := tx.Model(&SystemCorrelationFull{}).Select("system_correlation.*, " +
		"ts1.name AS trading_system1_name, ts2.name AS trading_system2_name").
		Joins("LEFT JOIN trading_system ts1 ON system_correlation.trading_system1_id = ts1.id").
		Joins("LEFT JOIN trading_system ts2 ON system_correlation.trading_system2_id = ts2.id").
		Where(filter).Order("correlation DESC").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddSystemCorrelation(tx *gorm.DB, sc *SystemCorrelation) error {
	err := tx.Create(sc).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================
