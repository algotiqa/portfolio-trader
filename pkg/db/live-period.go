//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
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

func FindLivePeriodsByTradingSystemId(tx *gorm.DB, tsId uint) (*[]LivePeriod, error) {
	var list []LivePeriod

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Order("period").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func FindLivePeriodsByTradingSystemsId(tx *gorm.DB, ids []uint) (*[]LivePeriod, error) {
	var list []LivePeriod
	res := tx.Find(&list, "trading_system_id in ?", ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddLivePeriod(tx *gorm.DB, lp *LivePeriod) error {
	err := tx.Create(lp).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================

func DeleteAllLivePeriodsByTradingSystemId(tx *gorm.DB, id uint) error {
	err := tx.Delete(&LivePeriod{}, "trading_system_id", id).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================
