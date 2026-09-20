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

func GetAllocations(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]AllocationFull, error) {
	var list []AllocationFull

	res := tx.Model(&AllocationFull{}).Select("allocation.*, " +
		"portfolio.name as portfolio_name, " +
		"account_id, account_code, account_name, account_currency_code").
		Joins("LEFT JOIN portfolio ON allocation.portfolio_id = portfolio.id").
		Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetAllocationById(tx *gorm.DB, id uint) (*Allocation, error) {
	var list []Allocation
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

func AddAllocation(tx *gorm.DB, a *Allocation) error {
	err := tx.Create(a).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================

func UpdateAllocation(tx *gorm.DB, a *Allocation) error {
	return tx.Save(a).Error
}

//=============================================================================
