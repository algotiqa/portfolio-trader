//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

type AllocationSpec struct {
	PortfolioId uint `json:"portfolioId"`
}

//=============================================================================

type AllocationExt struct {
	db.Allocation
	Portfolio   *db.Portfolio               `json:"portfolio"`
	Filters     *[]*db.AllocationFilterFull `json:"filters"`
	Logs        *[]db.AllocationLog         `json:"logs"`
}

//=============================================================================

func GetAllocations(tx *gorm.DB, c *auth.Context, filter map[string]any, offset int, limit int) (*[]db.AllocationFull, error) {
	if !c.Session.IsAdmin() {
		filter["portfolio.username"] = c.Session.Username
	}

	return db.GetAllocations(tx, filter, offset, limit)
}

//=============================================================================

func GetAllocationById(tx *gorm.DB, c *auth.Context, id uint) (*AllocationExt, error) {
	c.Log.Info("GetAllocationById: Getting allocation", "id", id)

	a, err := getAllocation(tx, c, id)
	if err != nil {
		return nil, err
	}

	//--- Get portfolio

	p,errp := getPortfolio(tx, c, a.PortfolioId, "GetAllocationById")
	if errp != nil {
		return nil, errp
	}

	//--- Get filters

	fils, err := db.GetAllocationFiltersByAllocationId(tx, a.Id)
	if err != nil {
		c.Log.Error("GetAllocationById: Could not retrieve allocation filters", "error", err.Error())
		return nil, err
	}

	//--- Get logs

	logs, err := db.GetAllocationLogsByAllocationId(tx, a.Id)
	if err != nil {
		c.Log.Error("GetAllocationById: Could not retrieve allocation logs", "error", err.Error())
		return nil, err
	}

	//--- Put all together

	ae := AllocationExt{
		Allocation: *a,
		Portfolio : p,
		Filters   : setFilterComments(fils),
		Logs      : logs,
	}

	return &ae, nil
}

//=============================================================================

func AddAllocation(tx *gorm.DB, c *auth.Context, as *AllocationSpec) (*db.Allocation,error) {
	p, err := getPortfolio(tx, c, as.PortfolioId, "AddAllocation")
	if err != nil {
		return nil,err
	}

	a := db.Allocation{
		PortfolioId       : as.PortfolioId,
		RunDate           : time.Now(),
		RunType           : db.RunTypeManual,
		AccountPerc       : p.AccountPerc,
		MaxMarginPerc     : p.MaxMarginPerc,
		CorrelationPeriod : p.CorrelationPeriod,
		AccountCapital    : p.AccountCurrentCapital,
		Status            : db.AllocStatusWaiting,
	}

	err = db.AddAllocation(tx, &a)
	if err != nil {
		return nil,err
	}

	err = sendAllocationMessage(tx, c, &a)

	return &a, err
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func getAllocation(tx *gorm.DB, c *auth.Context, id uint) (*db.Allocation, error) {
	a, err := db.GetAllocationById(tx, id)

	if err != nil {
		c.Log.Error("getAllocation: Could not retrieve allocation", "error", err.Error())
		return nil, err
	}

	if a == nil {
		c.Log.Error("getAllocation: Allocation was not found", "id", id)
		return nil, req.NewNotFoundError("Allocation was not found: %v", id)
	}

	//--- We don't need to check for ownership as this is done by getPortfolio function

	return a, nil
}

//=============================================================================

func sendAllocationMessage(tx *gorm.DB, c *auth.Context, job *db.Allocation) error {
	err := msg.SendMessage(msg.ExPortfolio, msg.SourceAllocationJob, msg.TypeNewJob, job, tx)

	if err != nil {
		c.Log.Error("sendAllocationMessage: Could not publish the allocation message", "error", err.Error())
		return err
	}

	return nil
}

//=============================================================================

func setFilterComments(list *[]db.AllocationFilterFull) *[]*db.AllocationFilterFull {
	var res []*db.AllocationFilterFull

	for _, f := range *list {
		if f.TsRunning && !f.FilterPassed {
			f.Comment = "Turn OFF"
		} else if !f.TsRunning && f.FilterPassed {
			f.Comment = "Turn ON"
		}
		res = append(res, &f)
	}

	return &res
}

//=============================================================================
