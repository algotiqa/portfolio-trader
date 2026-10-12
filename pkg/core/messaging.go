//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"log/slog"

	"github.com/algotiqa/core/dbms"
	"github.com/algotiqa/core/msg"
	"gorm.io/gorm"
)

//=============================================================================

func SendEvent(user string, code string, params map[string]any) error {
	return dbms.RunInTransaction(func(tx *gorm.DB) error {
		err := msg.SendEventByCode(user, code, params, tx)
		if err != nil {
			slog.Error("SendEvent: Cannot send event", "code",code, "error", err.Error())
		}
		return err
	})
}

//=============================================================================
