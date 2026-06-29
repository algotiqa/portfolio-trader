//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import "errors"

//=============================================================================

func MapReal(data map[string]any, key string, target *float64, min, max float64, required bool, defValue float64) error {
	v, ok := data[key]
	if !ok {
		if required {
			return errors.New("missing required parameter: " + key)
		}
		*target = defValue
		return nil
	}

	val, ok := v.(float64)
	if !ok {
		return errors.New("invalid type for parameter: " + key)
	}

	if val < min || val > max {
		return errors.New("value out of range for parameter: " + key)
	}

	*target = val
	return nil
}

//=============================================================================

func MapInt(data map[string]any, key string, target *int, min, max int, required bool, defValue int) error {
	v, ok := data[key]
	if !ok {
		if required {
			return errors.New("missing required parameter: " + key)
		}
		*target = defValue
		return nil
	}

	val, ok := v.(int)
	if !ok {
		return errors.New("invalid type for parameter: " + key)
	}

	if val < min || val > max {
		return errors.New("value out of range for parameter: " + key)
	}

	*target = val
	return nil
}

//=============================================================================

func MapString[T ~string](data map[string]any, key string, target *T, domain []T, required bool, defValue T) error {
	v, ok := data[key]
	if !ok {
		if required {
			return errors.New("missing required parameter: " + key)
		}
		*target = defValue
		return nil
	}

	val, ok := v.(string)
	if !ok {
		return errors.New("invalid type for parameter: " + key)
	}

	tval := T(val)

	if domain != nil {
		found := false
		for _, d := range domain {
			if tval == d {
				found = true
				break
			}
		}
		if !found {
			return errors.New("value not in domain for parameter: " + key)
		}
	}

	*target = tval
	return nil
}

//=============================================================================
