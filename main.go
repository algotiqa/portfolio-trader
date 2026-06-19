//=============================================================================
//===
//=== Copyright (C) 2023-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package main

import (
	"log/slog"

	"github.com/algotiqa/core/auth"
	"github.com/algotiqa/core/boot"
	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/core/msg"
	"github.com/algotiqa/core/req"
	"github.com/algotiqa/portfolio-trader/pkg/app"
	"github.com/algotiqa/portfolio-trader/pkg/core/messaging"
	"github.com/algotiqa/portfolio-trader/pkg/core/process"
	"github.com/algotiqa/portfolio-trader/pkg/platform"
	"github.com/algotiqa/portfolio-trader/pkg/service"
)

//=============================================================================

const component = "portfolio-trader"

//=============================================================================

func main() {
	cfg := &app.Config{}
	boot.ReadConfig(component, cfg)
	logger := boot.InitLogger(component, &cfg.Application)
	engine := boot.InitEngine(logger, &cfg.Application)
	initClients()
	auth.InitAuthentication(&cfg.Authentication)
	dbms.InitDatabase(&cfg.Database)
	msg.InitMessaging(&cfg.Messaging)
	service.Init(engine, cfg, logger)
	process.Init(cfg)
	messaging.InitMessageListener()
	platform.InitPlatform(cfg)
	boot.RunHttpServer(engine, &cfg.Application)
}

//=============================================================================

func initClients() {
	slog.Info("Initializing clients...")
	req.AddDefaultClient("ca.crt", "server.crt", "server.key")
}

//=============================================================================
