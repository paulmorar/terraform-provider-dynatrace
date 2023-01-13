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

package http

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings/validation"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/request"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Requests []*Request

func (me *Requests) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"request": {
			Type:        schema.TypeList,
			Description: "A HTTP request to be performed by the monitor.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Request).Schema()},
		},
	}
}

type Request struct {
	Description    *string              `json:"description,omitempty"`   // A short description of the event to appear in the web UI
	URL            string               `json:"url"`                     // The URL to check
	Method         string               `json:"method"`                  // The HTTP method of the request
	Authentication *Authentication      `json:"authentication"`          // Authentication options for this request
	RequestBody    *string              `json:"requestBody,omitempty"`   // The body of the HTTP request—you need to escape all JSON characters. \n\n Is set to null if the request method is GET, HEAD, or OPTIONS.
	Validation     *validation.Settings `json:"validation,omitempty"`    // Validation helps you verify that your HTTP monitor loads the expected content
	Configuration  *request.Config      `json:"configuration,omitempty"` // The setup of the monitor
	PreProcessing  *string              `json:"preProcessingScript,omitempty"`
	PostProcessing *string              `json:"postProcessingScript,omitempty"`
}

func (me *Request) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the event to appear in the web UI.",
			Optional:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL to check.",
			Required:    true,
		},
		"method": {
			Type:        schema.TypeString,
			Description: "The HTTP method of the request.",
			Required:    true,
		},
		"body": {
			Type:        schema.TypeString,
			Description: "The body of the HTTP request.",
			Optional:    true,
		},
		"pre_processing": {
			Type:        schema.TypeString,
			Description: "Javascript code to execute before sending the request.",
			Optional:    true,
		},
		"post_processing": {
			Type:        schema.TypeString,
			Description: "Javascript code to execute after sending the request.",
			Optional:    true,
		},
		"validation": {
			Type:        schema.TypeList,
			Description: "Validation helps you verify that your HTTP monitor loads the expected content",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(validation.Settings).Schema()},
		},
		"authentication": {
			Type:        schema.TypeList,
			Description: "Authentication options for this request",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Authentication).Schema()},
		},
		"configuration": {
			Type:        schema.TypeList,
			Description: "The setup of the monitor",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(request.Config).Schema()},
		},
	}
}

func (me *Request) MarshalHCL() (map[string]any, error) {
	result := map[string]any{}
	if me.Description != nil && len(*me.Description) > 0 {
		result["description"] = *me.Description
	}
	result["url"] = me.URL
	result["method"] = me.Method
	if me.RequestBody != nil && len(*me.RequestBody) > 0 {
		result["body"] = *me.RequestBody
	}
	if me.PreProcessing != nil && len(*me.PreProcessing) > 0 {
		result["pre_processing"] = *me.PreProcessing
	}
	if me.PostProcessing != nil && len(*me.PostProcessing) > 0 {
		result["post_processing"] = *me.PostProcessing
	}
	if me.Validation != nil {
		if marshalled, err := me.Validation.MarshalHCL(); err == nil {
			result["validation"] = []any{marshalled}
		} else {
			return nil, err
		}
	}
	if me.Authentication != nil {
		if marshalled, err := me.Authentication.MarshalHCL(); err == nil {
			result["authentication"] = []any{marshalled}
		} else {
			return nil, err
		}
	}
	if me.Configuration != nil {
		if marshalled, err := me.Configuration.MarshalHCL(); err == nil {
			result["configuration"] = []any{marshalled}
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (me *Request) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("description", &me.Description); err != nil {
		return err
	}
	if err := decoder.Decode("url", &me.URL); err != nil {
		return err
	}
	if err := decoder.Decode("method", &me.Method); err != nil {
		return err
	}
	if err := decoder.Decode("body", &me.RequestBody); err != nil {
		return err
	}
	if err := decoder.Decode("pre_processing", &me.PreProcessing); err != nil {
		return err
	}
	if err := decoder.Decode("post_processing", &me.PostProcessing); err != nil {
		return err
	}
	if result, ok := decoder.GetOk("validation.#"); ok && result.(int) == 1 {
		me.Validation = new(validation.Settings)
		if err := me.Validation.UnmarshalHCL(hcl.NewDecoder(decoder, "validation", 0)); err != nil {
			return err
		}
	}
	if result, ok := decoder.GetOk("authentication.#"); ok && result.(int) == 1 {
		me.Authentication = new(Authentication)
		if err := me.Authentication.UnmarshalHCL(hcl.NewDecoder(decoder, "authentication", 0)); err != nil {
			return err
		}
	}
	if result, ok := decoder.GetOk("configuration.#"); ok && result.(int) == 1 {
		me.Configuration = new(request.Config)
		if err := me.Configuration.UnmarshalHCL(hcl.NewDecoder(decoder, "configuration", 0)); err != nil {
			return err
		}
	}
	return nil
}
