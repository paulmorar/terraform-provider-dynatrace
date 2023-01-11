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

package managementzones

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DimensionConditions []*DimensionCondition

func (me *DimensionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "Dimension conditions",
			Elem:        &schema.Resource{Schema: new(DimensionCondition).Schema()},
		},
	}
}

func (me DimensionConditions) MarshalHCL() (map[string]interface{}, error) {
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
		result["condition"] = entries
	}
	return result, nil
}

func (me *DimensionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("condition"); ok {

		entrySet := value.(*schema.Set)

		for _, entryMap := range entrySet.List() {
			hash := entrySet.F(entryMap)
			entry := new(DimensionCondition)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "condition", hash)); err != nil {
				return err
			}
			*me = append(*me, entry)
		}
	}
	return nil
}

// No documentation available
type DimensionCondition struct {
	ConditionType DimensionConditionType `json:"conditionType"` // Type
	Key           *string                `json:"key,omitempty"` // Key
	RuleMatcher   DimensionOperator      `json:"ruleMatcher"`   // Operator
	Value         string                 `json:"value"`         // Value
}

func (me *DimensionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition_type": {
			Type:        schema.TypeString,
			Description: "Type",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Key",
			Optional:    true,
		},
		"rule_matcher": {
			Type:        schema.TypeString,
			Description: "Operator",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Value",
			Required:    true,
		},
	}
}

func (me *DimensionCondition) MarshalHCL() (map[string]interface{}, error) {
	properties := hcl.Properties{}

	return properties.EncodeAll(map[string]interface{}{
		"condition_type": me.ConditionType,
		"key":            me.Key,
		"rule_matcher":   me.RuleMatcher,
		"value":          me.Value,
	})
}

func (me *DimensionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]interface{}{
		"condition_type": &me.ConditionType,
		"key":            &me.Key,
		"rule_matcher":   &me.RuleMatcher,
		"value":          &me.Value,
	})
}
