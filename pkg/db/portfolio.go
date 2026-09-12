//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
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

func GetPortfolioById(tx *gorm.DB, id uint) (*Portfolio, error) {
	var list []Portfolio
	res := tx.Find(&list, id)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	if len(list) == 1 {
		return &list[0], nil
	}

	return nil, nil
}

//=============================================================================

func UpdatePortfolio(tx *gorm.DB, p *Portfolio) error {
	return tx.Save(p).Error
}

//=============================================================================

func UpdateAccountInfo(tx *gorm.DB, accountId uint, values map[string]interface{}) error {
	return tx.Model(&Portfolio{}).
		Where("account_id", accountId).
		Updates(values).Error
}

//=============================================================================

func DeletePortfolio(tx *gorm.DB, id uint) error {
	return tx.Delete(&Portfolio{}, id).Error
}

//=============================================================================
