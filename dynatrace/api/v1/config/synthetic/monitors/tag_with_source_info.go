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

package monitors

import (
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TagsWithSourceInfo []*TagWithSourceInfo

func (me *TagsWithSourceInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tag": {
			Type:        schema.TypeList,
			Description: "Tag with source of a Dynatrace entity.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(TagWithSourceInfo).Schema(),
			},
		},
	}
}

func (me TagsWithSourceInfo) MarshalHCL(properties hcl.Properties) error {
	entries := []any{}
	sorted := TagsWithSourceInfo{}
	sorted = append(sorted, me...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	for _, tag := range sorted {
		marshalled := hcl.Properties{}
		if err := tag.MarshalHCL(marshalled); err == nil {
			entries = append(entries, marshalled)
		} else {
			return err
		}
	}
	properties["tag"] = entries
	return nil
}

func (me *TagsWithSourceInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("tag.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(TagWithSourceInfo)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "tag", idx)); err != nil {
				return err
			}
			*me = append(*me, entry)
		}
	}
	return nil
}

// Tag with source of a Dynatrace entity
type TagWithSourceInfo struct {
	Source  *TagSource `json:"source,omitempty"` // The source of the tag, such as USER, RULE_BASED or AUTO
	Context TagContext `json:"context"`          // The origin of the tag, such as AWS or Cloud Foundry. \n\n Custom tags use the `CONTEXTLESS` value
	Key     string     `json:"key"`              // The key of the tag. \n\n Custom tags have the tag value here
	Value   *string    `json:"value"`            // The value of the tag. \n\n Not applicable to custom tags
}

func (me *TagWithSourceInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"source": {
			Type:        schema.TypeString,
			Description: "The source of the tag. Supported values are `USER`, `RULE_BASED` and `AUTO`.",
			Optional:    true,
		},
		"context": {
			Type:        schema.TypeString,
			Description: "The origin of the tag. Supported values are `AWS`, `AWS_GENERIC`, `AZURE`, `CLOUD_FOUNDRY`, `CONTEXTLESS`, `ENVIRONMENT`, `GOOGLE_CLOUD` and `KUBERNETES`.\n\nCustom tags use the `CONTEXTLESS` value.",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The key of the tag.\n\nCustom tags have the tag value here.",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: " The value of the tag.\n\nNot applicable to custom tags.",
			Optional:    true,
		},
	}
}

func (me TagWithSourceInfo) MarshalHCL(properties hcl.Properties) error {
	if me.Source != nil {
		properties["source"] = string(*me.Source)
	}
	properties["context"] = string(me.Context)
	properties["key"] = me.Key
	if me.Value != nil {
		properties["value"] = *me.Value
	}
	return nil
}

func (me *TagWithSourceInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("source", &me.Source); err != nil {
		return err
	}
	if err := decoder.Decode("context", &me.Context); err != nil {
		return err
	}
	if err := decoder.Decode("key", &me.Key); err != nil {
		return err
	}
	if err := decoder.Decode("value", &me.Value); err != nil {
		return err
	}
	return nil
}

// TagSource The source of the tag, such as USER, RULE_BASED or AUTO
type TagSource string

// TagSources offers the known enum values
var TagSources = struct {
	Auto      TagSource
	RuleBased TagSource
	User      TagSource
}{
	"AUTO",
	"RULE_BASED",
	"USER",
}

// TagContext The origin of the tag, such as AWS or Cloud Foundry. \n\n Custom tags use the `CONTEXTLESS` value
type TagContext string

// TagContexts offers the known enum values
var TagContexts = struct {
	AWS          TagContext
	AWSGeneric   TagContext
	Azure        TagContext
	CloudFoundry TagContext
	ContextLess  TagContext
	Environment  TagContext
	GoogleCloud  TagContext
	Kubernetes   TagContext
}{
	"AWS",
	"AWS_GENERIC",
	"AZURE",
	"CLOUD_FOUNDRY",
	"CONTEXTLESS",
	"ENVIRONMENT",
	"GOOGLE_CLOUD",
	"KUBERNETES",
}
