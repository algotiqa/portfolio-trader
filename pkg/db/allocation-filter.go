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

func GetAllocationFiltersByAllocationId(tx *gorm.DB, allId uint) (*[]AllocationFilterFull, error) {
	var list []AllocationFilterFull

	filter := map[string]any{}
	filter["allocation_id"] = allId

	res := tx.Model(&AllocationFilterFull{}).Select("allocation_filter.*, " +
		"name as ts_name, " +
		"data_symbol as ts_data_symbol, broker_symbol as ts_broker_symbol, " +
		"market_type as ts_market_type, strategy_type as ts_strategy_type, running as ts_running").
		Joins("LEFT JOIN trading_system ON allocation_filter.trading_system_id = trading_system.id").
		Where(filter).Order("ts_market_type, ts_broker_symbol, ts_name").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddAllocationFilter(tx *gorm.DB, af *AllocationFilter) error {
	err := tx.Create(af).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================
