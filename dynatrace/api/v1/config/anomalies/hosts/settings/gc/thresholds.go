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

package gc

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Thresholds Custom thresholds for high GC activity. If not set, automatic mode is used.
//
//	Meeting **any** of these conditions triggers an alert.
type Thresholds struct {
	GcSuspensionPercentage int32 `json:"gcSuspensionPercentage"` // GC suspension is higher than *X*% in 3 out of 5 samples.
	GcTimePercentage       int32 `json:"gcTimePercentage"`       // GC time is higher than *X*% in 3 out of 5 samples.
}

func (me *Thresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"suspension_percentage": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "GC suspension is higher than *X*% in 3 out of 5 samples",
		},
		"time_percentage": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "GC time is higher than *X*% in 3 out of 5 samples",
		},
	}
}

func (me *Thresholds) MarshalHCL() (map[string]any, error) {
	return map[string]any{
		"time_percentage":       int(me.GcSuspensionPercentage),
		"suspension_percentage": int(me.GcTimePercentage),
	}, nil
}

func (me *Thresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("time_percentage"); ok {
		me.GcSuspensionPercentage = int32(value.(int))
	}
	if value, ok := decoder.GetOk("suspension_percentage"); ok {
		me.GcTimePercentage = int32(value.(int))
	}
	return nil
}
