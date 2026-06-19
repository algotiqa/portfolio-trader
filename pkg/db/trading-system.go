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
	"time"

	"github.com/algotiqa/core/req"
	"gorm.io/gorm"
)

//=============================================================================

func GetTradingSystems(tx *gorm.DB, filter map[string]any, offset int, limit int) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Where(filter).Offset(offset).Limit(limit).Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemById(tx *gorm.DB, id uint) (*TradingSystem, error) {
	var list []TradingSystem
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

func GetTradingSystemsByUser(tx *gorm.DB, name string) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, "username = ?", name)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetTradingSystemsById(tx *gorm.DB, username string, ids []uint) (*[]TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, "username = ? and id in ?", username, ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func GetUsersWithTradingSystems(tx *gorm.DB) ([]string, error) {
	var list []string
	res := tx.Table("trading_system").Distinct("username").Scan(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return list, nil
}

//=============================================================================

func GetTradingSystemsByIdsAsMap(tx *gorm.DB, ids []uint) (map[uint]*TradingSystem, error) {
	var list []TradingSystem
	res := tx.Find(&list, "id in ?", ids)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	tsMap := map[uint]*TradingSystem{}

	for _, ts := range list {
		tsAux := ts
		tsMap[ts.Id] = &tsAux
	}

	return tsMap, nil
}

//=============================================================================

func GetTradingSystemsInIdle(tx *gorm.DB, days int) (*[]TradingSystem, error) {
	var list []TradingSystem
	date := time.Now().Add(-time.Hour * 24 * time.Duration(days))

	res := tx.
		Where("running    = ?", true).
		Where("active     = ?", true).
		Where("last_trade < ?", date).
		Find(&list)

	if res.Error != nil {
		return nil, req.NewServerErrorByError(res.Error)
	}

	return &list, nil
}

//=============================================================================

func UpdateTradingSystem(tx *gorm.DB, ts *TradingSystem) error {
	return tx.Save(ts).Error
}

//=============================================================================

func UpdateDataProductInfo(tx *gorm.DB, dataProductId uint, values map[string]interface{}) error {
	return tx.Model(&TradingSystem{}).
		Where("data_product_id", dataProductId).
		Updates(values).Error
}

//=============================================================================

func UpdateBrokerProductInfo(tx *gorm.DB, brokerProductId uint, values map[string]interface{}) error {
	return tx.Model(&TradingSystem{}).
		Where("broker_product_id", brokerProductId).
		Updates(values).Error
}

//=============================================================================

func DeleteTradingSystem(tx *gorm.DB, id uint) error {
	return tx.Delete(&TradingSystem{}, id).Error
}

//=============================================================================
