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

package event

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Locators a list of locators identifying the desired element
type Locators []*Locator

func (me *Locators) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"locator": {
			Type:        schema.TypeList,
			Description: "A locator dentifyies the desired element",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Locator).Schema()},
		},
	}
}

func (me Locators) MarshalHCL() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	entries := []interface{}{}
	for _, entry := range me {
		if marshalled, err := entry.MarshalHCL(); err == nil {
			entries = append(entries, marshalled)
		} else {
			return nil, err
		}
		result["locator"] = entries
	}
	return result, nil
}

func (me *Locators) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("locator", me); err != nil {
		return err
	}
	return nil
}

type Locator struct {
	Type  LocatorType `json:"type"`  // Defines where to look for an element. `css` (CSS Selector) or `dom` (Javascript code)
	Value string      `json:"value"` // The name of the element to be found
}

func (me *Locator) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Defines where to look for an element. `css` (CSS Selector) or `dom` (Javascript code)",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The name of the element to be found",
			Required:    true,
		},
	}
}

func (me *Locator) MarshalHCL() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	result["type"] = string(me.Type)
	result["value"] = me.Value
	return result, nil
}

func (me *Locator) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	if err := decoder.Decode("value", &me.Value); err != nil {
		return err
	}
	return nil
}

// LocatorType defines where to look for an element. `css` (CSS Selector) or `dom` (Javascript code)
type LocatorType string

// LocatorTypes offers the known enum values
var LocatorTypes = struct {
	ContentMatch LocatorType
	ElementMatch LocatorType
}{
	"css",
	"dom",
}
