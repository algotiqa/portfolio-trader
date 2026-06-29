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

func GetTradingPositionByTsId(tx *gorm.DB, tsId uint) (*TradingPosition, error) {
	var list []TradingPosition

	filter := map[string]any{}
	filter["trading_system_id"] = tsId

	res := tx.Where(filter).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(list) == 0 {
		return nil, req.NewServerError("trading position not found for tsId=%v", tsId)
	}

	return &list[0], nil
}

//=============================================================================

func GetTradingPositionsByTsIds(tx *gorm.DB, ids []uint) (*[]TradingPosition, error) {
	var list []TradingPosition
	res := tx.Find(&list, "trading_system_id in ?", ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func SetTradingPosition(tx *gorm.DB, tp *TradingPosition) error {
	return tx.Save(tp).Error
}

//=============================================================================

func DeleteTradingPosition(tx *gorm.DB, id uint) error {
	return tx.Delete(&TradingPosition{}, id).Error
}

//=============================================================================
