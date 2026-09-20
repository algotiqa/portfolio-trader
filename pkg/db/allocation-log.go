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

func GetAllocationLogsByAllocationId(tx *gorm.DB, allId uint) (*[]AllocationLog, error) {
	var list []AllocationLog

	filter := map[string]any{}
	filter["allocation_id"] = allId

	res := tx.Where(filter).Order("id").Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func AddAllocationLog(tx *gorm.DB, al *AllocationLog) error {
	err := tx.Create(al).Error
	return req.NewServerErrorByError(err)
}

//=============================================================================
