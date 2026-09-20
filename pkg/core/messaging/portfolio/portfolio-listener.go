//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package portfolio

import (
	"encoding/json"
	"log/slog"

	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/portfolio-trader/pkg/db"
)

//=============================================================================

func HandleMessage(m *msg.Message) bool {
	slog.Info("handleInternalMessage: New internal message received", "source", m.Source, "type", m.Type)

	if m.Type == msg.TypeNewJob {
		if m.Source == msg.SourceAllocationJob {
			return startAllocationJob(m)
		}
	}

	slog.Error("handleInternalMessage: Dropping message with unknown source/type!", "source", m.Source, "type", m.Type)
	return true
}

//=============================================================================

func startAllocationJob(m *msg.Message) bool {
	all := db.Allocation{}
	err := json.Unmarshal(m.Entity, &all)
	if err != nil {
		slog.Error("startAllocationJob: Dropping badly formatted message for allocation job!", "entity", string(m.Entity))
		return true
	}

	return runAllocationJob(&all)
}

//=============================================================================
//===
//=== New allocation job
//===
//=============================================================================

func runAllocationJob(a *db.Allocation) bool {
	slog.Info("runAllocationJob: Starting new allocation", "id", a.Id, "portfolioId", a.PortfolioId)

	err := calcAllocation(a)

	if err != nil {
		slog.Error("runAllocationJob: Raised error while creating allocation", "error", err.Error())
	} else {
		slog.Info("runAllocationJob: Operation complete", "id", a.Id)
	}

	return err == nil
}

//=============================================================================
