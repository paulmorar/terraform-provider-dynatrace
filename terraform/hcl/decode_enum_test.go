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
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
)

type Enum string

type EnumContainer struct {
	Enum    Enum
	OptEnum *Enum
}

func (me *EnumContainer) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("enum", &me.Enum); err != nil {
		return err
	}
	if err := decoder.Decode("opt_enum", &me.OptEnum); err != nil {
		return err
	}
	return nil
}

func TestDecodeEnum(t *testing.T) {
	decoder := hcl.NewDecoder(&testDecoder{
		Values: map[string]any{
			"enum":     "Test",
			"opt_enum": "OptTest",
		},
	})
	ec := &EnumContainer{}
	if err := ec.UnmarshalHCL(decoder); err != nil {
		t.Error(err)
	}
	if string(ec.Enum) != "Test" {
		t.Errorf("expected: %v, actual: %v", "Test", ec.Enum)
	}
	if string(*ec.OptEnum) != "OptTest" {
		t.Errorf("expected: %v, actual: %v", "OptTest", ec.OptEnum)
	}
}
