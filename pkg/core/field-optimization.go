//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"errors"
	"fmt"
	"math/rand"
)

//=============================================================================
//===
//=== FieldOptimization
//===
//=============================================================================

type FieldOptimization[T int|float64] struct {
	Enabled  bool  `json:"enabled"`
	CurValue T     `json:"curValue"`
	MinValue T     `json:"minValue"`
	MaxValue T     `json:"maxValue"`
	Step     T     `json:"step"`

	//--- Caching
	steps    *[]T
}

//=============================================================================

func (f *FieldOptimization[T]) StepsCount() uint {
	if !f.Enabled {
		return 1
	}

	return uint((f.MaxValue - f.MinValue) / f.Step) +1
}

//=============================================================================

func (f *FieldOptimization[T]) Steps() *[]T {
	if f.steps != nil {
		return f.steps
	}

	var list []T

	if !f.Enabled {
		list = append(list, f.CurValue)
	} else {
		for i := f.MinValue; i<= f.MaxValue; i += f.Step {
			list = append(list, i)
		}
	}

	f.steps = &list

	return &list
}

//=============================================================================

func (f *FieldOptimization[T]) ValidateWithSpec(spec *NumberParamSpec[T]) error {
	return f.Validate(spec.MinValue, spec.MaxValue)
}

//=============================================================================

func (f *FieldOptimization[T]) Validate(min, max T) error {
	if f.Enabled {
		if f.MinValue < min || f.MinValue > max {
			return fmt.Errorf("min value out of range [%v .. %v]", min, max)
		}

		if f.MaxValue < min || f.MaxValue > max {
			return fmt.Errorf("max value out of range [%v .. %v]", min, max)
		}

		if f.MinValue > f.MaxValue {
			return errors.New("min value greater than max value")
		}

		//--- Validate step

		if f.Step < 0 || f.Step > max {
			return fmt.Errorf("step value out of range (0 .. %v]", max)
		}

		if f.Step == 0 {
			return errors.New("step value cannot be zero")
		}
	} else {
		if f.CurValue < min || f.CurValue > max {
			return fmt.Errorf("current value out of range [%v .. %v]", min, max)
		}
	}

	return nil
}

//=============================================================================

func (f *FieldOptimization[T]) RandomValue() T {
	list := f.Steps()
	idx  := rand.Intn(len(*list))

	return (*list)[idx]
}

//=============================================================================
