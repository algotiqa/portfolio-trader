//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import "time"

//=============================================================================
//===
//=== Plot
//===
//=============================================================================

type Serie struct {
	Time   []time.Time `json:"time"`
	Values []float64   `json:"values"`
}

//-----------------------------------------------------------------------------

func (p *Serie) AddPoint(t time.Time, value float64) {
	p.Time   = append(p.Time,   t)
	p.Values = append(p.Values, value)
}

//=============================================================================
