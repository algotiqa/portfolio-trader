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

package quality

import (
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/portfolio-trader/pkg/platform"
	"github.com/algotiqa/types"
)

//=============================================================================
//===
//=== MarketRegime
//===
//=============================================================================

type MarketRegime interface {
	MapTrade(trade *db.Trade) (int, int)
}

//=============================================================================

func NewMarketRegime(results []*platform.BarResult, timeframe int, loc *time.Location) MarketRegime {
	if timeframe == 1440 {
		return newDailyMarketRegime(results, loc)
	}

	return newBarMarketRegime(results, loc)
}

//=============================================================================
//===
//=== DailyMarketRegime
//===
//=============================================================================

type DailyMarketRegime struct {
	dayMap map[types.Date]*platform.BarResult
	loc    *time.Location
}

//=============================================================================

func newDailyMarketRegime(results []*platform.BarResult, loc *time.Location) MarketRegime {
	res := &DailyMarketRegime{
		dayMap: map[types.Date]*platform.BarResult{},
		loc   : loc,
	}

	for _, br := range results {
		res.add(br)
	}

	return res
}

//=============================================================================

func (mr *DailyMarketRegime) add(br *platform.BarResult) {
	br.Time = br.Time.In(mr.loc)
	date := types.ToDate(&br.Time)
	mr.dayMap[date] = br
}

//=============================================================================

func (mr *DailyMarketRegime) MapTrade(trade *db.Trade) (int, int) {
	entryDate := trade.EntryDate.In(mr.loc)
	date      := types.ToDate(&entryDate)

	br, ok := mr.dayMap[date]


	if !ok || entryDate.After(br.Time) {
		//--- We probably need to go to the next day because there is a missing
		//--- day (maybe a holiday) but the trade started in the afternoon of that holiday.
		//--- If entryDate>endSessionDate, we need to go to the next session

		br2, ok2 := mr.dayMap[date.AddDays(1)]
		if !ok2 {
			//--- If the trade is not on the next day, the session may have been reduced
			//--- because of holidays. Let's check if the current day exists

			if !ok {
				return -1, -1
			}
		} else {
			br = br2
		}
	}

	return br.Direction, br.Volatility
}

//=============================================================================
//===
//=== BarMarketRegime
//===
//=============================================================================

type BarMarketRegime struct {
	dayMap map[types.Date]*BarDayMap
	loc    *time.Location
}

//=============================================================================

func newBarMarketRegime(results []*platform.BarResult, loc *time.Location) MarketRegime {
	res := &BarMarketRegime{
		dayMap: map[types.Date]*BarDayMap{},
		loc   : loc,
	}

	for _, br := range results {
		res.add(br)
	}

	return res
}

//=============================================================================

func (mr *BarMarketRegime) add(br *platform.BarResult) {
	br.Time = br.Time.In(mr.loc)
	date := types.ToDate(&br.Time)
	barMap, ok := mr.dayMap[date]

	if !ok {
		barMap = NewBarDayMap()
		mr.dayMap[date] = barMap
	}

	barMap.add(br)
}

//=============================================================================

func (mr *BarMarketRegime) MapTrade(trade *db.Trade) (int, int) {
	entryDate := trade.EntryDate.In(mr.loc)
	date      := types.ToDate(&entryDate)

	barMap, ok := mr.dayMap[date]

	if !ok {
		return -1, -1
	}

	return barMap.mapTrade(entryDate)
}

//=============================================================================
//===
//=== BarDayMap
//===
//=============================================================================

type BarDayMap struct {
	dayMap map[int]*platform.BarResult
}

//=============================================================================

func NewBarDayMap() *BarDayMap {
	return &BarDayMap{
		dayMap: map[int]*platform.BarResult{},
	}
}

//=============================================================================

func (bm *BarDayMap) add(br *platform.BarResult) {
	hh, mm, _ := br.Time.Clock()
	mins := hh * 60 + mm
	bm.dayMap[mins] = br
}

//=============================================================================

func (bm *BarDayMap) mapTrade(entryDate time.Time) (int, int) {
	hh, mm, _ := entryDate.Clock()
	mins := hh * 60 + mm
	br, ok := bm.dayMap[mins]

	if !ok {
		return -1, -1
	}

	return br.Direction, br.Volatility
}

//=============================================================================
