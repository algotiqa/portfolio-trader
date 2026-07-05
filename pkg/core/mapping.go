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
	"errors"
	"reflect"
)

//=============================================================================

func MapNumber[T int|float64](data map[string]any, spec *NumberParamSpec[T]) (*T,error) {
	key := spec.Name

	v, ok := data[key]
	if !ok {
		if spec.Required {
			return nil,errors.New("missing required parameter: " + key)
		}

		return spec.DefValue,nil
	}

	v = fixForInt[T](v)

	val, ok := v.(T)
	if !ok {
		return nil,errors.New("invalid type for parameter: " + key)
	}

	if val < spec.MinValue || val > spec.MaxValue {
		return nil,errors.New("value out of range for parameter: " + key)
	}

	return &val,nil
}

//=============================================================================

func fixForInt[T int|float64](v any) any {
	switch v.(type) {
		case float64:
			tType := reflect.TypeOf((*T)(nil)).Elem()
			if tType.Kind() == reflect.Int {
				val,_ := v.(float64)
				return int(val)
			}
	}

	return v
}

//=============================================================================

func MapList[T ~string](data map[string]any, spec *ListParamSpec[T]) (T,error) {
	key := spec.Name

	v, ok := data[key]
	if !ok {
		if spec.Required {
			return "",errors.New("missing required parameter: " + key)
		}

		return spec.DefValue,nil
	}

	val, ok := v.(string)
	if !ok {
		return "",errors.New("invalid type for parameter: " + key)
	}

	tval := T(val)

	if spec.Domain != nil {
		found := false
		for _, d := range spec.Domain {
			if tval == d {
				found = true
				break
			}
		}
		if !found {
			return "",errors.New("value not in domain for parameter: " + key)
		}
	}

	return tval,nil
}

//=============================================================================
