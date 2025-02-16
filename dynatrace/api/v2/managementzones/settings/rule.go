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
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Rules []*Rule

func (me *Rules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "A management zone rule",
			Elem:        &schema.Resource{Schema: new(Rule).Schema()},
		},
	}
}

func (me Rules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *Rules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

// No documentation available
type Rule struct {
	Enabled        bool                         `json:"enabled"`                  // Enabled
	Type           RuleType                     `json:"type"`                     // Rule type
	AttributeRule  *ManagementZoneAttributeRule `json:"attributeRule,omitempty"`  // No documentation available
	DimensionRule  *DimensionRule               `json:"dimensionRule,omitempty"`  // No documentation available
	EntitySelector string                       `json:"entitySelector,omitempty"` // Entity selector. The documentation of the entity selector can be found [here](https://dt-url.net/apientityselector).
}

func (me *Rule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enabled",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Rule type",
			Required:    true,
		},
		"attribute_rule": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ManagementZoneAttributeRule).Schema()},
			Optional:    true,
		},
		"dimension_rule": {
			Type:        schema.TypeList,
			Description: "No documentation available",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(DimensionRule).Schema()},
			Optional:    true,
		},
		"entity_selector": {
			Type:             schema.TypeString,
			Description:      "Entity selector. The documentation of the entity selector can be found [here](https://dt-url.net/apientityselector).",
			Optional:         true,
			DiffSuppressFunc: hcl.SuppressEOT,
		},
	}
}

func (me *Rule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":         me.Enabled,
		"type":            me.Type,
		"attribute_rule":  me.AttributeRule,
		"dimension_rule":  me.DimensionRule,
		"entity_selector": me.EntitySelector,
	})
}

func (me *Rule) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("type"); !ok {
		return errors.New("invalid")
	}
	return decoder.DecodeAll(map[string]any{
		"enabled":         &me.Enabled,
		"type":            &me.Type,
		"attribute_rule":  &me.AttributeRule,
		"dimension_rule":  &me.DimensionRule,
		"entity_selector": &me.EntitySelector,
	})
}
