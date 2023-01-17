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

package request

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HeadersSection struct {
	Headers      Headers  `json:"addHeaders"`
	Restrictions []string `json:"toRequests,omitempty"`
}

func (me *HeadersSection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"header": {
			Type:        schema.TypeList,
			Description: "contains an HTTP header of the request",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Header).Schema()},
		},
		"restrictions": {
			Type:        schema.TypeSet,
			Description: "Restrict applying headers to a set of URLs",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *HeadersSection) MarshalHCL(properties hcl.Properties) error {
	if len(me.Headers) > 0 {
		entries := []any{}
		for _, header := range me.Headers {
			marshalled := hcl.Properties{}

			if err := header.MarshalHCL(marshalled); err == nil {
				entries = append(entries, marshalled)
			} else {
				return err
			}
		}
		properties["header"] = entries
	}
	if len(me.Restrictions) > 0 {
		properties["restrictions"] = me.Restrictions
	}
	return nil
}

func (me *HeadersSection) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("header.#"); ok {
		me.Headers = Headers{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Header)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "header", idx)); err != nil {
				return err
			}
			me.Headers = append(me.Headers, entry)
		}
	}
	if err := decoder.Decode("restrictions", &me.Restrictions); err != nil {
		return err
	}
	return nil
}

// Headers is a list of request headers
type Headers []*Header

func (me *Headers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"header": {
			Type:        schema.TypeList,
			Description: "contains an HTTP header of the request",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Header).Schema()},
		},
	}
}

func (me Headers) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		entries := []any{}
		for _, header := range me {
			marshalled := hcl.Properties{}

			if err := header.MarshalHCL(marshalled); err == nil {
				entries = append(entries, marshalled)
			} else {
				return err
			}
		}
		properties["header"] = entries
	}
	return nil
}

func (me *Headers) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("header.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Header)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "header", idx)); err != nil {
				return err
			}
			*me = append(*me, entry)
		}
	}
	return nil
}

// Header contains an HTTP header of the request
type Header struct {
	Name  string `json:"name"`  // The key of the header
	Value string `json:"value"` // The value of the header
}

func (me *Header) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The key of the header",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the header",
			Required:    true,
		},
	}
}

func (me *Header) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{"name": me.Name, "value": me.Value})
}

func (me *Header) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("name", &me.Name); err != nil {
		return err
	}
	if err := decoder.Decode("value", &me.Value); err != nil {
		return err
	}
	return nil
}
