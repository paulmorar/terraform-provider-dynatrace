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

package dashboards

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DashboardFilter represents filters, applied to a dashboard
type DashboardFilter struct {
	Timeframe      *string                    `json:"timeframe,omitempty"` // the default timeframe of the dashboard
	ManagementZone *EntityRef                 `json:"managementZone,omitempty"`
	Unknowns       map[string]json.RawMessage `json:"-"`
}

func (me *DashboardFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"timeframe": {
			Type:        schema.TypeString,
			Description: "the default timeframe of the dashboard",
			Optional:    true,
		},
		"management_zone": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "the management zone this dashboard applies to",
			Elem: &schema.Resource{
				Schema: new(EntityRef).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *DashboardFilter) MarshalHCL() (map[string]interface{}, error) {
	result := map[string]interface{}{}

	if len(me.Unknowns) > 0 {
		data, err := json.Marshal(me.Unknowns)
		if err != nil {
			return nil, err
		}
		result["unknowns"] = string(data)
	}
	if me.Timeframe != nil {
		result["timeframe"] = opt.String(me.Timeframe)
	}
	if me.ManagementZone != nil {
		if marshalled, err := me.ManagementZone.MarshalHCL(); err == nil {
			result["management_zone"] = []interface{}{marshalled}
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (me *DashboardFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "timeframe")
		delete(me.Unknowns, "management_zone")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}

	if value, ok := decoder.GetOk("timeframe"); ok {
		me.Timeframe = opt.NewString(value.(string))
	}

	if _, ok := decoder.GetOk("management_zone.#"); ok {
		me.ManagementZone = new(EntityRef)
		if err := me.ManagementZone.UnmarshalHCL(hcl.NewDecoder(decoder, "management_zone", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *DashboardFilter) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.Timeframe != nil {
		rawMessage, err := json.Marshal(me.Timeframe)
		if err != nil {
			return nil, err
		}
		m["timeframe"] = rawMessage
	}
	if me.ManagementZone != nil {
		rawMessage, err := json.Marshal(me.ManagementZone)
		if err != nil {
			return nil, err
		}
		m["managementZone"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *DashboardFilter) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("timeframe", &me.Timeframe); err != nil {
		return err
	}
	if err := m.Unmarshal("managementZone", &me.ManagementZone); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
