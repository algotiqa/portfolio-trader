//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package db

import (
	"github.com/algotiqa/core/req"
	"gorm.io/gorm"
)

//=============================================================================
//=== Note: Bars are NOT ordered by date (possibly too many records and useless to do)

func FindEquityBarsByTradingSystemId(tx *gorm.DB, tsId uint) (*[]EquityBar, error) {
	var list []EquityBar

	query := "trade_id in ( SELECT id FROM trade WHERE trading_system_id = ? )"
	res := tx.Find(&list, query, tsId)

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
