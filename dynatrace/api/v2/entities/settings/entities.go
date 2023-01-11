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

package entities

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Entities []*Entity // A list of monitored entities.

func (me Entities) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity": {
			Type:        schema.TypeList,
			Description: "A list of monitored entities.",
			Elem:        &schema.Resource{Schema: new(Entity).Schema()},
			Optional:    true,
		},
	}
}

func (me *Entities) MarshalHCL() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	entries := []interface{}{}
	for _, entry := range *me {
		if marshalled, err := entry.MarshalHCL(); err == nil {
			entries = append(entries, marshalled)
		} else {
			return nil, err
		}
	}
	result["entity"] = entries
	return result, nil
}

func (me *Entities) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeSlice("entity", me); err != nil {
		return err
	}
	return nil
}

type Entity struct {
	EntityId    *string `json:"entityId,omitempty"`    // The ID of the entity.
	Type        *string `json:"type,omitempty"`        // The type of the entity.
	DisplayName *string `json:"displayName,omitempty"` // The name of the entity, displayed in the UI.
	Tags        Tags    `json:"tags,omitempty"`        // A set of tags assigned to the entity.
}

func (me *Entity) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_id": {
			Type:        schema.TypeString,
			Description: "The ID of the entity.",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the entity.",
			Optional:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "The name of the entity, displayed in the UI.",
			Optional:    true,
		},
		"tags": {
			Type:        schema.TypeList,
			Description: "A set of tags assigned to the entity.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(Tags).Schema(),
			},
		},
	}
}

func (me *Entity) MarshalHCL() (map[string]interface{}, error) {
	properties := hcl.Properties{}

	return properties.EncodeAll(map[string]interface{}{
		"entity_id":    me.EntityId,
		"type":         me.Type,
		"display_name": me.DisplayName,
		"tags":         me.Tags,
	})
}

func (me *Entity) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]interface{}{
		"entity_id":    &me.EntityId,
		"type":         &me.Type,
		"display_name": &me.DisplayName,
		"tags":         &me.Tags,
	})
}
