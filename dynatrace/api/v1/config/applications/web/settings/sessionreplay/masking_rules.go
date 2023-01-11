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

package sessionreplay

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MaskingRules The masking rules defining how data is hidden
type MaskingRules []*MaskingRule

func (me *MaskingRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Description: "The masking rule defining how data is hidden",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(MaskingRule).Schema()},
		},
	}
}

func (me MaskingRules) MarshalHCL() (map[string]interface{}, error) {
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
		result["rule"] = entries
	}
	return result, nil
}

func (me *MaskingRules) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("rule", me); err != nil {
		return err
	}
	return nil
}

// MaskingRule The masking rule defining how data is hidden
type MaskingRule struct {
	Type                  MaskingRuleType `json:"maskingRuleType"`       // The type of the masking rule
	Selector              string          `json:"selector"`              // The selector for the element or the attribute to be masked. \n\nSpecify a CSS expression for an element or a [regular expression](https://dt-url.net/k9e0iaq) for an attribute
	UserInteractionHidden bool            `json:"userInteractionHidden"` // Interactions with the element are (`true`) or are not (`false) masked
}

func (me *MaskingRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the masking rule",
			Required:    true,
		},
		"selector": {
			Type:        schema.TypeString,
			Description: "The selector for the element or the attribute to be masked. \n\nSpecify a CSS expression for an element or a [regular expression](https://dt-url.net/k9e0iaq) for an attribute",
			Required:    true,
		},
		"user_interaction_hidden": {
			Type:        schema.TypeBool,
			Description: "Interactions with the element are (`true`) or are not (`false) masked",
			Optional:    true,
		},
	}
}

func (me *MaskingRule) MarshalHCL() (map[string]interface{}, error) {
	return hcl.Properties{}.EncodeAll(map[string]interface{}{
		"type":                    me.Type,
		"selector":                me.Selector,
		"user_interaction_hidden": me.UserInteractionHidden,
	})
}

func (me *MaskingRule) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]interface{}{
		"type":                    &me.Type,
		"selector":                &me.Selector,
		"user_interaction_hidden": &me.UserInteractionHidden,
	})
}

type MaskingRuleType string

var MaskingRuleTypes = struct {
	Attribute MaskingRuleType
	Element   MaskingRuleType
}{
	"ATTRIBUTE",
	"ELEMENT",
}
