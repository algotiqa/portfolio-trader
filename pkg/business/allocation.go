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
	"log/slog"
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================
//===
//=== Model
//===
//=============================================================================

type AllocationSpec struct {
	PortfolioId uint `json:"portfolioId"`
}

//=============================================================================

type AllocationExt struct {
	db.Allocation
	Portfolio    *db.Portfolio               `json:"portfolio"`
	Filters      *[]*db.AllocationFilterFull `json:"filters"`
	Logs         *[]db.AllocationLog         `json:"logs"`
	Correlations *[]db.SystemCorrelationFull `json:"correlations"`
	CorrMatrix   *CorrelationMatrix          `json:"corrMatrix"`
}

//=============================================================================

type CorrelationMatrix struct {
	Names []string		`json:"names"`
	Cells [][]*float64	`json:"cells"`
}

//=============================================================================
//===
//=== Functions
//===
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

	//--- Get logs

	corr, err := db.GetSystemCorrelationsByAllocationId(tx, a.Id)
	if err != nil {
		c.Log.Error("GetAllocationById: Could not retrieve system correlations", "error", err.Error())
		return nil, err
	}

	//--- Put all together

	fils2 := setFilterComments(fils)

	ae := AllocationExt{
		Allocation  : *a,
		Portfolio   : p,
		Filters     : fils2,
		Logs        : logs,
		Correlations: corr,
		CorrMatrix  : buildCorrelationMatrix(fils2,corr),
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

var one = 1.0

func buildCorrelationMatrix(fils *[]*db.AllocationFilterFull, corr *[]db.SystemCorrelationFull) *CorrelationMatrix {
	matrix := &CorrelationMatrix{}

	//--- Get only active filters

	var idMap = make(map[uint]int)
	var list []*db.AllocationFilterFull

	for _, f := range *fils {
		if f.FilterPassed {
			idMap[f.TradingSystemId] = len(list)
			list         = append(list,         f)
			matrix.Names = append(matrix.Names, f.TsName)
		}
	}

	//--- Allocate full matrix, setting identity on diagonal

	size := len(list)

	for i,_ := range list {
		matrix.Cells = append(matrix.Cells, make([]*float64, size))
		matrix.Cells[i][i] = &one
	}

	//--- Fill upper diagonal with data

	for _,scf := range *corr {
		index1, ok1 := idMap[scf.TradingSystem1Id]
		index2, ok2 := idMap[scf.TradingSystem2Id]

		if !ok1 || !ok2 {
			slog.Error("buildCorrelationMatrix: Map lookup failure!", "allocationId", scf.AllocationId,
					   "id1", scf.TradingSystem1Id, "id2", scf.TradingSystem2Id)
		} else {
			if scf.Message == "" {
				if index1 > index2 {
					aux    := index1
					index1 = index2
					index2 = aux
				}

				matrix.Cells[index1][index2] = &scf.Correlation
			}
		}
	}

	return matrix
}

//=============================================================================
