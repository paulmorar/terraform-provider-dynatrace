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

package common

import (
	"sort"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TagFilters []*TagFilter

func (me TagFilters) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A Tag Filter",
			Elem:        &schema.Resource{Schema: new(TagFilter).Schema()},
		},
	}
}

func (me TagFilters) MarshalHCL(properties hcl.Properties) error {
	tagFilters := TagFilters{}
	tagFilters = append(tagFilters, me...)
	sort.Slice(tagFilters, func(i int, j int) bool {
		a := tagFilters[i]
		b := tagFilters[j]
		return strings.Compare(a.Key, b.Key) > 0
	})
	filters := []any{}
	for _, filter := range tagFilters {
		marshalled := hcl.Properties{}
		if err := filter.MarshalHCL(marshalled); err == nil {
			filters = append(filters, marshalled)
		} else {
			return err
		}
	}
	if len(filters) > 0 {
		properties["filter"] = filters
	}
	return nil
}

func (me *TagFilters) UnmarshalHCL(decoder hcl.Decoder) error {
	nme := TagFilters{}
	if result, ok := decoder.GetOk("filter.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(TagFilter)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "filter", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	*me = nme
	return nil
}
