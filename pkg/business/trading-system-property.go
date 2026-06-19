//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package business

import (
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

type TradingSystemTradingRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

type TradingSystemRunningRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

type TradingSystemActivationRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

type TradingSystemActiveRequest struct {
	Value bool `json:"value"`
}

//=============================================================================

const (
	ResponseStatusOk      = "ok"
	ResponseStatusSkipped = "skipped"
	ResponseStatusError   = "error"
)

//-----------------------------------------------------------------------------

type TradingSystemPropertyResponse struct {
	Status        string            `json:"status"`
	Message       string            `json:"message"`
	TradingSystem *db.TradingSystem `json:"tradingSystem"`
}

//=============================================================================

func SetTradingSystemTrading(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemTradingRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemTrading: Trading property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.Trading
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status: ResponseStatusSkipped,
		}, nil
	}

	if !oldValue && newValue {
		//--- Turning on
	} else {
		//--- Turning off
		if ts.Running {
			return &TradingSystemPropertyResponse{
				Status:  ResponseStatusError,
				Message: "Trading system must be stopped",
			}, nil
		}
	}

	ts.Trading = newValue
	updateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	c.Log.Info("SetTradingSystemTrading: Trading property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status:        ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================

func SetTradingSystemRunning(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemRunningRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemRunning: Running property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.Running
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status: ResponseStatusSkipped,
		}, nil
	}

	ts.Running = newValue
	updateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateLivePeriod(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateRewind(ts)
	if err != nil {
		return nil, err
	}

	c.Log.Info("SetTradingSystemRunning: Running property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status:        ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================

func SetTradingSystemActivation(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemActivationRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemActivation: Auto-activation property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.AutoActivation
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status: ResponseStatusSkipped,
		}, nil
	}

	ts.AutoActivation = newValue
	err = db.UpdateTradingSystem(tx, ts)

	c.Log.Info("SetTradingSystemActivation: Auto-activation property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status:        ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================

func SetTradingSystemActive(tx *gorm.DB, c *auth.Context, tsId uint, req *TradingSystemActiveRequest) (*TradingSystemPropertyResponse, error) {
	c.Log.Info("SetTradingSystemActive: Active property change request", "id", tsId, "value", req.Value)

	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	oldValue := ts.Active
	newValue := req.Value

	if oldValue == newValue {
		return &TradingSystemPropertyResponse{
			Status: ResponseStatusSkipped,
		}, nil
	}

	if ts.AutoActivation {
		return &TradingSystemPropertyResponse{
			Status:  ResponseStatusError,
			Message: "Trading system is in AUTOMATIC mode. Switch to MANUAL to change",
		}, nil
	}

	ts.Active = newValue
	updateStatus(ts)
	err = db.UpdateTradingSystem(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateLivePeriod(tx, ts)
	if err != nil {
		return nil, err
	}

	err = updateRewind(ts)

	c.Log.Info("SetTradingSystemActive: Active property changed", "id", tsId, "value", req.Value)

	return &TradingSystemPropertyResponse{
		Status:        ResponseStatusOk,
		TradingSystem: ts,
	}, err
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func updateStatus(ts *db.TradingSystem) {
	ts.SuggestedAction = db.TsActionNone

	if !ts.Running {
		ts.Status = db.TsStatusOff
	} else if ts.Active {
		ts.Status = db.TsStatusRunning
	} else {
		ts.Status = db.TsStatusPaused
	}
}

//============================================================================

func updateLivePeriod(tx *gorm.DB, ts *db.TradingSystem) error {
	lp := &db.LivePeriod{
		TradingSystemId: ts.Id,
		Period:          time.Now(),
		Active:          ts.Running && ts.Active,
	}

	return db.AddLivePeriod(tx, lp)
}

//============================================================================

func updateRewind(ts *db.TradingSystem) error {
	return nil
}

//=============================================================================
