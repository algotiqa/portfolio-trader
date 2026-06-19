//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package db

import (
	"time"

	"github.com/algotiqa/core/req"
	"gorm.io/gorm"
)

//=============================================================================

func FindTradesByTradingSystemId(tx *gorm.DB, tsId uint) (*[]Trade, error) {
	var list []Trade

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Order("entry_date,exit_date").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func FindTradesByTradingSystemsId(tx *gorm.DB, ids []uint) (*[]Trade, error) {
	var list []Trade
	res := tx.Find(&list, "trading_system_id in ?", ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

// We MUST order by entry_date because we may have cases like:
// entry_date,       exit_date
// 2025-1-1 03:00    2025-1-2 06:00
// 2025-1-2 06:00    2025-1-2 06:00    <-- fake trade
// 2025-1-2 06:00    2025-1-4 02:00
// Ordering by exit_date, the second record could come first

func FindTradesByTsIdFromTime(tx *gorm.DB, tsId uint, fromTime *time.Time, toTime *time.Time) (*[]Trade, error) {
	to := time.Now().UTC()
	from := to.Add(-50 * 365 * 24 * time.Hour)

	if fromTime != nil {
		from = *fromTime
	}

	if toTime != nil {
		to = *toTime
	}

	var list []Trade

	//--- WHERE condition must be exit_date otherwise we loose trades started in the past and ended after fromTime
	query := "trading_system_id = ? and exit_date >= ? and exit_date<= ?"
	res := tx.Order("exit_date,entry_date").Find(&list, query, tsId, from, to)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func FindTradesFromTime(tx *gorm.DB, tsIds []uint, fromTime time.Time) (*[]Trade, error) {
	var list []Trade

	res := tx.Find(&list, "trading_system_id in ? and entry_date >= ?", tsIds, fromTime)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddTrade(tx *gorm.DB, tr *Trade) error {
	err := tx.Create(tr).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================

func DeleteAllTradesByTradingSystemId(tx *gorm.DB, id uint) error {
	err := tx.Delete(&Trade{}, "trading_system_id", id).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================
