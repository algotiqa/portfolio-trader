//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"time"

	"github.com/algotiqa/core/req"
	"github.com/algotiqa/types"
)

//=============================================================================

type SelectedPeriod struct {
	DaysBack *int       `json:"daysBack"`
	FromDate types.Date `json:"fromDate"`
	ToDate   types.Date `json:"toDate"`
}

//=============================================================================

func CalcSelectedPeriod(period *SelectedPeriod, loc *time.Location) (*time.Time, *time.Time, error) {
	daysBack := period.DaysBack
	fromDate := period.FromDate
	toDate   := period.ToDate

	if daysBack == nil {
		var from *time.Time
		var to   *time.Time

		if !fromDate.IsNil() {
			if !fromDate.IsValid() {
				return nil, nil, req.NewBadRequestError("Invalid fromDate parameter: %d", fromDate)
			}

			tt := fromDate.ToDateTime(false, loc)
			from = &tt
		}

		if !toDate.IsNil() {
			if !toDate.IsValid() {
				return nil, nil, req.NewBadRequestError("Invalid toDate parameter: %d", toDate)
			}

			tt := toDate.ToDateTime(true, loc)
			to = &tt
		}

		return from, to, nil
	}

	//--- All

	if *daysBack == 0 {
		return nil, nil, nil
	}

	//--- Specific last days

	if *daysBack > 0 && *daysBack <= 20000 {
		fromTime := time.Now().UTC()
		back     := time.Hour * time.Duration(24 * *daysBack)
		fromTime = fromTime.Add(-back)

		return &fromTime, nil, nil
	}

	//--- Custom range

	return nil, nil, req.NewBadRequestError("Invalid daysBack parameter: %d", daysBack)
}

//=============================================================================
