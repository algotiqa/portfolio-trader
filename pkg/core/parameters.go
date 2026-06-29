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
	"fmt"
)

//=============================================================================
//===
//=== ParamSpec
//===
//=============================================================================

type ParamSpecType string

//-----------------------------------------------------------------------------

const (
	PSTString ParamSpecType = "string"
	PSTList   ParamSpecType = "list"
	PSTInt    ParamSpecType = "int"
	PSTReal   ParamSpecType = "real"
)

//-----------------------------------------------------------------------------

type ParamSpec struct {
	Name     string        `json:"name"`
	Type     ParamSpecType `json:"type"`
	Required bool          `json:"required"`
}

//=============================================================================
//===
//=== NumberParamSpec
//===
//=============================================================================

type NumberParamSpec[T int|float64] struct {
	ParamSpec
	MinValue  T `json:"minValue"`
	MaxValue  T `json:"maxValue"`
	DefValue *T `json:"defValue,omitempty"`
}

//=============================================================================

func NewNumberParamSpec[T int|float64](name string, required bool, minValue, maxValue T, defValue *T) *NumberParamSpec[T] {
	return &NumberParamSpec[T]{
		ParamSpec: ParamSpec{
			Name     : name,
			Type     : PSTInt,
			Required : required,
		},
		MinValue: minValue,
		MaxValue: maxValue,
		DefValue: defValue,
	}
}

//=============================================================================

func (s *NumberParamSpec[T]) Validate(value *T) (*T,error) {
	if value == nil {
		if s.Required {
			if s.DefValue == nil {
				return nil,errors.New("missing required parameter: "+ s.Name)
			}

			return s.DefValue, nil
		}

		return value,nil
	}

	if *value < s.MinValue {
		return nil,fmt.Errorf("parameter '%v' is below its minimum (%v): %v", s.Name, s.MinValue, *value)
	}

	if *value > s.MaxValue {
		return nil,fmt.Errorf("parameter '%v' is above its maximum (%v): %v", s.Name, s.MaxValue, *value)
	}

	return value,nil
}

//=============================================================================
//===
//=== ListParamSpec
//===
//=============================================================================

type ListParamSpec[T ~string] struct {
	ParamSpec
	Domain   []T `json:"domain"`
	DefValue T   `json:"defValue,omitempty"`
}

//=============================================================================

func NewListParamSpec[T ~string](name string, required bool, domain []T, defValue T) *ListParamSpec[T] {
	return &ListParamSpec[T]{
		ParamSpec: ParamSpec{
			Name     : name,
			Type     : PSTList,
			Required : required,
		},
		Domain  : domain,
		DefValue: defValue,
	}
}

//=============================================================================

func (s *ListParamSpec[T]) Validate(value T) (T,error) {
	if value == "" {
		if s.Required {
			if s.DefValue == "" {
				return "",errors.New("missing required parameter: "+ s.Name)
			}

			return s.DefValue, nil
		}

		return value,nil
	}

	found := false
	for _, d := range s.Domain {
		if value == d {
			found = true
			break
		}
	}
	if !found {
		return "",fmt.Errorf("value not allowed for parameter '%v': %v", s.Name, value)
	}

	return value,nil
}

//=============================================================================
