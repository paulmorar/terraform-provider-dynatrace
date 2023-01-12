/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package hcl_test

import (
	"fmt"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
)

type Shape interface {
	UnmarshalHCL(hcl.Decoder) error
}

type BaseShape struct {
	Type string
}

func (me *BaseShape) UnmarshalHCL(hcl.Decoder) error {
	return nil
}

type Rectangle struct {
	BaseShape
	Width  int
	Height int
}

func (me *Rectangle) UnmarshalHCL(decoder hcl.Decoder) error {
	return nil
}

type Square struct {
	BaseShape
	Length int
}

func (me *Square) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.Decode("length", &me.Length)
}

type ShapeWrapper struct {
	Shape Shape
}

func (me *ShapeWrapper) MarshalHCL(decoder hcl.Decoder) (map[string]any, error) {
	properties := hcl.Properties{}
	switch cmp := me.Shape.(type) {
	case *Square:
		if err := properties.Encode("square", cmp); err != nil {
			return nil, err
		}
		return properties, nil
	case *Rectangle:
		if err := properties.Encode("rectangle", cmp); err != nil {
			return nil, err
		}
		return properties, nil
	case *BaseShape:
		if err := properties.Encode("generic", cmp); err != nil {
			return nil, err
		}
		return properties, nil
	default:
		return nil, fmt.Errorf("cannot HCL marshal (x) objects of type %T", cmp)
	}
}

func (me *ShapeWrapper) UnmarshalHCL(decoder hcl.Decoder) error {
	var err error
	var shape any
	if shape, err = decoder.DecodeAny(map[string]any{
		"square":    new(Square),
		"rectangle": new(Rectangle),
		"generic":   new(BaseShape),
	}); err != nil {
		return err
	}
	me.Shape = shape.(Shape)
	return nil
}

func TestDecodeInheritance(t *testing.T) {
	decoder := hcl.NewDecoder(&testDecoder{
		Values: map[string]any{
			"square.#":        1,
			"square.0.length": 3,
		},
	})
	wrapper := &ShapeWrapper{}
	if err := wrapper.UnmarshalHCL(decoder); err != nil {
		t.Error(err)
	}
	if square, ok := wrapper.Shape.(*Square); ok {
		if square.Length != 3 {
			t.Errorf("Square.Length: expected: %d, actual: %d", 3, square.Length)
		}
	} else {
		t.Error("Square expected")
	}
}
