//=============================================================================
/*
Copyright © 2026 Andrea Carboni andrea.carboni71@gmail.com

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
