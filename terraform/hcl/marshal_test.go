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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
)

type StringEnum string

func (me StringEnum) Ref() *StringEnum {
	return &me
}

func TestMarshaller(t *testing.T) {
	{
		properties := hcl.Properties{}
		obj := struct{ EnumRef *StringEnum }{EnumRef: StringEnum("asdf").Ref()}
		if err := properties.Marshal(hcl.VoidDecoder(), "asdf", obj.EnumRef); err != nil {
			t.Error(err)
		}
		if len(properties) != 1 {
			t.Fail()
		}
		if stored, found := properties["asdf"]; !found {
			t.Fail()
		} else {
			switch tStored := stored.(type) {
			case string:
				if tStored != "asdf" {
					t.Fail()
				}
			default:
				t.Fail()
			}
		}

	}

	{
		properties := hcl.Properties{}
		obj := struct{ EnumRef *StringEnum }{EnumRef: nil}
		if err := properties.Marshal(hcl.VoidDecoder(), "asdf", obj.EnumRef); err != nil {
			t.Error(err)
		}
		if len(properties) != 0 {
			t.Fail()
		}

	}

	{
		properties := hcl.Properties{}
		obj := struct{ OptString *string }{OptString: opt.NewString("asdf")}
		if err := properties.Marshal(hcl.VoidDecoder(), "asdf", obj.OptString); err != nil {
			t.Error(err)
		}
		if len(properties) != 1 {
			t.Fail()
		}
		if stored, found := properties["asdf"]; !found {
			t.Fail()
		} else {
			switch tStored := stored.(type) {
			case string:
				if tStored != "asdf" {
					t.Fail()
				}
			default:
				t.Fail()
			}
		}

	}
	{
		properties := hcl.Properties{}
		obj := struct{ OptString *string }{OptString: nil}
		if err := properties.Marshal(hcl.VoidDecoder(), "asdf", obj.OptString); err != nil {
			t.Error(err)
		}
		if len(properties) != 0 {
			t.Fail()
		}
	}

	{
		properties := hcl.Properties{}
		obj := struct{ EnumRef StringEnum }{EnumRef: StringEnum("asdf")}
		if err := properties.Marshal(hcl.VoidDecoder(), "asdf", obj.EnumRef); err != nil {
			t.Error(err)
		}
		if len(properties) != 1 {
			t.Fail()
		}
		if stored, found := properties["asdf"]; !found {
			t.Fail()
		} else {
			switch tStored := stored.(type) {
			case string:
				if tStored != "asdf" {
					t.Fail()
				}
			default:
				t.Fail()
			}
		}

	}
}
