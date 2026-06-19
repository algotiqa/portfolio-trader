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
//=== Note: Bars are NOT ordered by date (possibly too many records and useless to do)

func FindEquityBarsByTradesId(tx *gorm.DB, tradeIds []int64) (*[]EquityBar, error) {
	var list []EquityBar

	res := tx.Find(&list, "trade_id in ?", tradeIds)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddEquityBar(tx *gorm.DB, eb *EquityBar) error {
	err := tx.Create(eb).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================

func DeleteAllEquityBarsByTradingSystemId(tx *gorm.DB, id uint) error {
	query := "DELETE from equity_bar WHERE trade_id in ( SELECT id FROM trade WHERE trading_system_id = ? )"
	res   := tx.Exec(query, id)
	return req.NewServerErrorByError(res.Error)
}

//=============================================================================
