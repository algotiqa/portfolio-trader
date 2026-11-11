//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package criteria

import "github.com/algotiqa/core/req"

//=============================================================================
//===
//=== Constant ranges
//===
//=============================================================================

const MinTradesLength     = 1
const MaxTradesLength     = 1000

const MinWinningPercentage= 0
const MaxWinningPercentage= 100

const MinDrawdownValue    = -100000
const MaxDrawdownValue    = 0

const MinTrendlineValue   = -10
const MaxTrendlineValue   = 10

const MinProfitFactor     = 0
const MaxProfitFactor     = 10

const MinReturnMaxDDValue = 0
const MaxReturnMaxDDValue = 10

//=============================================================================
//===
//=== Criteria config
//===
//=============================================================================

type Config struct {
	Enabled bool    `json:"enabled"`
	Length  int     `json:"length"`
	Value   float64 `json:"value"`
}

//=============================================================================

func (c *Config) validate() error {
	if c.Length < MinTradesLength || c.Length > MaxTradesLength {
		return req.NewBadRequestError("length must be between %v and %v", MinTradesLength, MaxTradesLength)
	}

	return nil
}

//=============================================================================
