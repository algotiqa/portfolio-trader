//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package criteria

//=============================================================================
//===
//=== Base criteria
//===
//=============================================================================

type BaseCriteria struct {
	runMode     RunMode
	activation *Activation
}

//=============================================================================

func (c *BaseCriteria) Activation() *Activation {
	return c.activation
}

//=============================================================================
//===
//=== Private functions
//===
//=============================================================================

func (c *BaseCriteria) evaluate(value float64, skipValue *float64) {
	for i, m := range c.activation.Metrics {
		if skipValue != nil && *skipValue == m {
			continue
		}

		switch c.runMode {
			case RunModeActivation:
				if m >= value {
					c.activation.Flags[i] = 1
				}

			case RunModeDeactivation:
				if m < value {
					c.activation.Flags[i] = 1
				}
		}
	}
}

//=============================================================================
