//=============================================================================
//===
//=== Copyright (C) 2027-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package db

import (
	"time"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/types"
	"gorm.io/gorm"
)

//=============================================================================

func FindDailyReturnsByTradingSystemId(tx *gorm.DB, tsId uint) (*[]DailyReturn, error) {
	var list []DailyReturn

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Order("date").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func FindDailyReturnsByTradingSystemsId(tx *gorm.DB, ids []uint) (*[]DailyReturn, error) {
	var list []DailyReturn
	res := tx.Find(&list, "trading_system_id in ?", ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func FindDailyReturnsByTsIdFromTime(tx *gorm.DB, tsId uint, fromDate *types.Date, toDate *types.Date) (*[]DailyReturn, error) {
	toT   := time.Now().UTC()
	fromT := toT.Add(-50 * 365 * 24 * time.Hour) //--- 50 years back

	from := types.ToDate(&fromT)
	to   := types.ToDate(&toT)

	if fromDate != nil {
		from = *fromDate
	}

	if toDate != nil {
		to = *toDate
	}

	var list []DailyReturn

	query := "trading_system_id = ? and date >= ? and date<= ?"
	res := tx.Order("date").Find(&list, query, tsId, from, to)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddDailyReturn(tx *gorm.DB, dr *DailyReturn) error {
	err := tx.Create(dr).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================

func DeleteAllDailyReturnsByTradingSystemId(tx *gorm.DB, id uint) error {
	err := tx.Delete(&DailyReturn{}, "trading_system_id", id).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================
