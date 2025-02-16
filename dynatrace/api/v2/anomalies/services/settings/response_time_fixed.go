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

package services

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No documentation available
type ResponseTimeFixed struct {
	ResponseTimeSlowest    *ResponseTimeFixedSlowest `json:"responseTimeSlowest"`    // Slowest 10%. Alert if the average response time of the slowest 10% of requests degrades beyond this threshold:
	OverAlertingProtection *OverAlertingProtection   `json:"overAlertingProtection"` // Avoid over-alerting
	Sensitivity            Sensitivity               `json:"sensitivity"`            // Sensitivity
	ResponseTimeAll        *ResponseTimeFixedAll     `json:"responseTimeAll"`        // All requests. Alert if the average response time of all requests degrades beyond this threshold:
}

func (me *ResponseTimeFixed) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"response_time_slowest": {
			Type:        schema.TypeList,
			Description: "Slowest 10%. Alert if the average response time of the slowest 10% of requests degrades beyond this threshold:",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ResponseTimeFixedSlowest).Schema()},
			Required:    true,
		},
		"over_alerting_protection": {
			Type:        schema.TypeList,
			Description: "Avoid over-alerting",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(OverAlertingProtection).Schema()},
			Required:    true,
		},
		"sensitivity": {
			Type:        schema.TypeString,
			Description: "Sensitivity",
			Required:    true,
		},
		"response_time_all": {
			Type:        schema.TypeList,
			Description: "All requests. Alert if the average response time of all requests degrades beyond this threshold:",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(ResponseTimeFixedAll).Schema()},
			Required:    true,
		},
	}
}

func (me *ResponseTimeFixed) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"response_time_slowest":    me.ResponseTimeSlowest,
		"over_alerting_protection": me.OverAlertingProtection,
		"sensitivity":              me.Sensitivity,
		"response_time_all":        me.ResponseTimeAll,
	})
}

func (me *ResponseTimeFixed) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"response_time_slowest":    &me.ResponseTimeSlowest,
		"over_alerting_protection": &me.OverAlertingProtection,
		"sensitivity":              &me.Sensitivity,
		"response_time_all":        &me.ResponseTimeAll,
	})
}
