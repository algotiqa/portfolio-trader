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
	"encoding/json"
	"errors"
	"time"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/portfolio-trader/pkg/business/position"
	"github.com/algotiqa/portfolio-trader/pkg/business/position/model"
	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"gorm.io/gorm"
)

//=============================================================================

func SetTradingPosition(tx *gorm.DB, c *auth.Context, tsId uint, p *position.TradingPosition) error {
	_, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	err = p.Params.Validate()
	if err != nil {
		return err
	}

	_,tp,err := convertPosition(&p.Model, &p.Params)
	if err != nil {
		return err
	}

	tp.TradingSystemId = tsId

	return db.SetTradingPosition(tx, tp)
}

//=============================================================================

func RunPositionAnalysis(tx *gorm.DB, c *auth.Context, tsId uint, par *position.AnalysisRequest) (*position.AnalysisResponse, error) {
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return nil, err
	}

	pos, err := db.GetTradingPositionByTsId(tx, tsId)
	if err != nil {
		return nil, err
	}

	curConfig := make(map[string]any)
	err = json.Unmarshal([]byte(pos.Config), &curConfig)
	if err != nil {
		return nil, errors.New("cannot parse position config: " + err.Error())
	}

	curModel,err := model.New(pos.Model, curConfig)
	if err != nil {
		return nil, errors.New("cannot build position model: " + err.Error())
	}

	//--- Get selected position model (if provided)

	var selModel model.PositionModel

	if par.Model != nil && par.Model.Name != "" && par.Params != nil {
		err = par.Params.Validate()
		if err != nil {
			return nil, err
		}

		selModel,pos,err = convertPosition(par.Model, par.Params)
		if err != nil {
			return nil, err
		}
	}

	//--- Other needed data

	fromTime, toTime, err := core.CalcSelectedPeriod(&par.Period, time.UTC)
	if err != nil {
		return nil, err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return nil, err
	}

	return position.RunAnalysis(ts, curModel, selModel, trades, pos)
}

//=============================================================================

func StartPositionOptimization(tx *gorm.DB, c *auth.Context, tsId uint, oreq *position.OptimizationRequest) error {
	ts, err := getTradingSystemAndCheckAccess(tx, c, tsId)
	if err != nil {
		return err
	}

	err = oreq.Validate()
	if err != nil {
		return err
	}

	//--- Other needed data

	fromTime, toTime, err := core.CalcSelectedPeriod(&oreq.RunConfig.Period, time.UTC)
	if err != nil {
		return err
	}

	trades, err := db.FindTradesByTsIdFromTime(tx, ts.Id, fromTime, toTime)
	if err != nil {
		return err
	}

	err = position.StartOptimization(ts, trades, oreq)
	if err == nil {
		c.Log.Info("StartPositionOptimization: Starting optimization", "tsId", ts.Id, "tsName", ts.Name)
	}

	return err
}

//=============================================================================

func StopPositionOptimization(c *auth.Context, tsId uint) error {
	c.Log.Info("StopPositionOptimization: Stopping optimization", "tsId", tsId)
	err := position.StopOptimization(tsId)

	return err
}

//=============================================================================

func GetPositionOptimizationInfo(c *auth.Context, tsId uint) (*position.OptimizationResponse, error) {
	info := position.GetOptimizationInfo(tsId)
	return position.NewOptimizationResponse(info), nil
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func convertPosition(m *position.Model, p *position.Parameters) (model.PositionModel,*db.TradingPosition,error) {
	mod,err := model.New(m.Name, m.Config)
	if err != nil {
		return nil, nil, err
	}

	data,err := json.Marshal(mod.Config())
	if err != nil {
		return nil, nil, err
	}

	tp := &db.TradingPosition{
		InitialCapital : *p.InitialCapital,
		MaxTolDrawdPerc: *p.MaxTolDrawdPerc,
		MarginOverride : p.MarginOverride,
		MaxUnits       : *p.MaxUnits,
		RiskPerUnit    : p.RiskPerUnit,
		RiskValue      : p.RiskValue,
		Model          : mod.Name(),
		Config         : string(data),
	}

	return mod,tp,nil
}

//=============================================================================
