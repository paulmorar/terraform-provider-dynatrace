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

package metricevents

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DimensionFilters []*DimensionFilter // Dimension filter definitions

func (me *DimensionFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Dimension filter definitions",
			Elem:        &schema.Resource{Schema: new(DimensionFilter).Schema()},
		},
	}
}

func (me DimensionFilters) MarshalHCL() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if len(me) > 0 {
		entries := []interface{}{}
		for _, entry := range me {
			if marshalled, err := entry.MarshalHCL(); err == nil {
				entries = append(entries, marshalled)
			} else {
				return nil, err
			}
		}
		result["filter"] = entries
	}
	return result, nil
}

func (me *DimensionFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("filter"); ok {

		entrySet := value.(*schema.Set)

		for _, entryMap := range entrySet.List() {
			hash := entrySet.F(entryMap)
			entry := new(DimensionFilter)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "filter", hash)); err != nil {
				return err
			}
			*me = append(*me, entry)
		}
	}
	return nil
}

type DimensionFilter struct {
	DimensionKey   string `json:"dimensionKey"`   // The key of the dimension filter
	DimensionValue string `json:"dimensionValue"` // The value of the dimension filter
}

func (me *DimensionFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dimension_key": {
			Type:        schema.TypeString,
			Description: "The key of the dimension filter",
			Required:    true,
		},
		"dimension_value": {
			Type:        schema.TypeString,
			Description: "The value of the dimension filter",
			Required:    true,
		},
	}
}

func (me *DimensionFilter) MarshalHCL() (map[string]interface{}, error) {
	properties := hcl.Properties{}

	return properties.EncodeAll(map[string]interface{}{
		"dimension_key":   me.DimensionKey,
		"dimension_value": me.DimensionValue,
	})
}

func (me *DimensionFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]interface{}{
		"dimension_key":   &me.DimensionKey,
		"dimension_value": &me.DimensionValue,
	})
}
