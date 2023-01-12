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

type record struct {
	Value string
}

func (me *record) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("value", &me.Value); err != nil {
		return err
	}
	return nil
}

type records []*record

func (me *records) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("records", me); err != nil {
		return err
	}
	return nil
}

type testDecoder struct {
	Values map[string]any
}

func (me *testDecoder) Decode(key string, v any) error {
	return fmt.Errorf("Decode(%v, %T)", key, v)
}

func (me *testDecoder) DecodeAll(map[string]any) error {
	return fmt.Errorf("DecodeAll(%v)", "...")
}

func (me *testDecoder) DecodeSlice(key string, v any) error {
	return fmt.Errorf("DecodeSlice(%v, %T)", key, v)
}

func (me *testDecoder) Get(key string) any {
	return nil
}

func (me *testDecoder) GetChange(key string) (any, any) {
	return nil, nil
}

func (me *testDecoder) GetOk(key string) (any, bool) {
	if value, found := me.Values[key]; found {
		// fmt.Printf("GetOk(%v) => %v\n", key, value)
		return value, true
	}
	// fmt.Printf("GetOk(%v) not found\n", key)
	return nil, false
}

func (me *testDecoder) GetStringSet(key string) []string {
	return nil
}

func (me *testDecoder) DecodeAny(map[string]any) (any, error) {
	return nil, nil
}

func TestDecodeSlice(t *testing.T) {
	recs := records{}
	decoder := hcl.NewDecoder(&testDecoder{
		Values: map[string]any{
			"rectangle":       2,
			"records.0.value": "value0",
			"records.1.value": "value1",
		},
	})

	if err := recs.UnmarshalHCL(decoder); err != nil {
		t.Error(err)
	}
	for idx, rec := range recs {
		fmt.Printf("%d: %v\n", idx, rec.Value)
	}

}
