//=============================================================================
//===
//=== Copyright (C) 2025-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package platform

import (
	"log/slog"
	"strconv"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/req"
)

//=============================================================================

type EquityRequest struct {
	Username string            `json:"username"`
	Images   map[string][]byte `json:"images"`
}

//-----------------------------------------------------------------------------

func NewEquityRequest() *EquityRequest {
	return &EquityRequest{
		Images: map[string][]byte{},
	}
}

//=============================================================================
//===
//=== Public functions
//===
//=============================================================================

func SetEquityChart(id uint, er *EquityRequest) error {
	slog.Info("SetEquityChart: Sending equity chart to storage manager", "id", id, "username", er.Username)

	token, err := auth.Token()
	if err != nil {
		return err
	}

	client := req.GetDefaultClient()
	url := platform.Storage + "/v1/trading-systems/" + strconv.Itoa(int(id)) + "/equity-chart"

	err = req.DoPut(client, url, &er, "", token)
	if err != nil {
		slog.Error("SetEquityChart: Got an error when sending to storage-manager", "id", id, "error", err.Error())
		return req.NewServerError("Cannot communicate with storage-manager: %v", err.Error())
	}

	slog.Info("SetEquityChart: Equity chart saved", "id", id)
	return nil
}

//=============================================================================

func DeleteEquityChart(username string, id uint) error {
	slog.Info("DeleteEquityChart: Deleting equity chart from the storage manager", "id", id, "username", username)

	token, err := auth.Token()
	if err != nil {
		return err
	}

	client := req.GetDefaultClient()
	url := platform.Storage + "/v1/trading-systems/" + strconv.Itoa(int(id)) + "/equity-chart"
	er := EquityRequest{
		Username: username,
	}

	err = req.DoDelete(client, url, &er, "", token)
	if err != nil {
		slog.Error("DeleteEquityChart: Got an error when sending to storage-manager", "id", id, "error", err.Error())
		return req.NewServerError("Cannot communicate with storage-manager: %v", err.Error())
	}

	slog.Info("DeleteEquityChart: Equity chart deleted", "id", id, "username", username)
	return nil
}

//=============================================================================
