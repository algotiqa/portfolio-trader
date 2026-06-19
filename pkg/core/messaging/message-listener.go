//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package messaging

import (
	"log/slog"

	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/portfolio-trader/pkg/core/messaging/inventory"
	"github.com/algotiqa/portfolio-trader/pkg/core/messaging/runtime"
	"github.com/algotiqa/portfolio-trader/pkg/core/messaging/system"
)

//=============================================================================

func InitMessageListener() {
	slog.Info("Starting message listeners...")

	go msg.ReceiveMessages(msg.QuInventoryToPortfolio, inventory.HandleMessage)
	go msg.ReceiveMessages(msg.QuRuntimeToPortfolio,   runtime.HandleMessage)
	go msg.ReceiveMessages(msg.QuSystemToPortfolio,    system.HandleMessage)
}

//=============================================================================
