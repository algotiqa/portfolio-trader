//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package statusupdater

import (
	"log/slog"
	"time"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/portfolio-trader/pkg/app"
	"github.com/algotiqa/portfolio-trader/pkg/consts"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func Init(cfg *app.Config) *time.Ticker {
	ticker := time.NewTicker(4 * time.Hour)

	go func() {
		//--- Wait 2 secs to allow the system to boot properly
		time.Sleep(2 * time.Second)
		run(cfg)

		for range ticker.C {
			run(cfg)
		}
	}()

	return ticker
}

//=============================================================================

func run(cfg *app.Config) {
	slog.Info("StatusUpdater: Starting")
	start := time.Now()

	list, err := GetTradingSystemsInIdle()
	if err != nil {
		slog.Error("StatusUpdater:Cannot get list of trading systems. Update aborted", "error", err)
	} else {
		slog.Info("StatusUpdater: Processing trading systems", "count", len(*list))

		for _, ts := range *list {
			err = updateTradingSystem(&ts)
			if err != nil {
				slog.Error("StatusUpdater:Cannot update trading system", "id", ts.Id, "error", err)
			}
		}
	}

	duration := time.Now().Sub(start).Seconds()
	slog.Info("StatusUpdater: Ended", "seconds", duration)
}

//=============================================================================

func GetTradingSystemsInIdle() (*[]db.TradingSystem, error) {
	var list *[]db.TradingSystem
	var err error

	err = dbms.RunInTransaction(func(tx *gorm.DB) error {
		list, err = db.GetTradingSystemsInIdle(tx, consts.IdleDays)
		return err
	})

	return list, err
}

//=============================================================================

func updateTradingSystem(ts *db.TradingSystem) error {
	if ts.Status == db.TsStatusRunning {
		ts.Status = db.TsStatusIdle
	} else if ts.Status == db.TsStatusIdle {
		brokenDate := time.Now().Add(-time.Hour * 24 * time.Duration(consts.BrokenDays))

		if ts.LastTrade.Before(brokenDate) {
			ts.Status = db.TsStatusBroken
			ts.SuggestedAction = db.TsActionCheck
		}
	} else {
		return nil
	}

	return dbms.RunInTransaction(func(tx *gorm.DB) error {
		return db.UpdateTradingSystem(tx, ts)
	})
}

//=============================================================================
