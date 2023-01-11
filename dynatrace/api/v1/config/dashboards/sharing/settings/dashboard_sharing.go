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

package sharing

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DashboardSharing represents sharing configuration of the dashboard
type DashboardSharing struct {
	DashboardID  string           `json:"id"`           // The Dynatrace entity ID of the dashboard
	Permissions  SharePermissions `json:"permissions"`  // Access permissions of the dashboard
	PublicAccess *AnonymousAccess `json:"publicAccess"` // Configuration of the [anonymous access](https://dt-url.net/ov03sf1) to the dashboard
	Preset       bool             `json:"preset"`       // If `true` the dashboard will be marked as preset
	Enabled      bool             `json:"enabled"`      // The dashboard is shared (`true`) or private (`false`)

	// not part of payload - used by export
	Name string `json:"-"`
}

func (me *DashboardSharing) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dashboard_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Dynatrace entity ID of the dashboard",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "The dashboard is shared (`true`) or private (`false`)",
		},
		"preset": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If `true` the dashboard will be marked as preset",
		},
		"permissions": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(SharePermissions).Schema()},
			Description: "Access permissions of the dashboard",
		},
		"public": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(AnonymousAccess).Schema()},
			Description: "Configuration of the [anonymous access](https://dt-url.net/ov03sf1) to the dashboard",
		},
	}
}

// MarshalHCL has no documentation
func (me *DashboardSharing) MarshalHCL() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	m["dashboard_id"] = me.DashboardID
	m["enabled"] = me.Enabled
	m["preset"] = me.Preset
	if len(me.Permissions) > 0 {
		if marshalled, err := me.Permissions.MarshalHCL(); err != nil {
			return nil, err
		} else {
			m["permissions"] = []interface{}{marshalled}
		}
	}
	if me.PublicAccess != nil && !me.PublicAccess.IsEmpty() {
		if marshalled, err := me.PublicAccess.MarshalHCL(); err != nil {
			return nil, err
		} else {
			m["public"] = []interface{}{marshalled}
		}
	}
	return m, nil
}

// UnmarshalHCL has no documentation
func (me *DashboardSharing) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("dashboard_id"); ok {
		me.DashboardID = value.(string)
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	} else {
		me.Enabled = false
	}
	if value, ok := decoder.GetOk("preset"); ok {
		me.Preset = value.(bool)
	} else {
		me.Preset = false
	}
	if value, ok := decoder.GetOk("permissions.#"); ok {
		count := value.(int)
		if count != 0 {
			if value, ok := decoder.GetOk("permissions.0.permission.#"); ok {
				count := value.(int)
				if count != 0 {
					me.Permissions = SharePermissions{}
					if value, ok := decoder.GetOk("permissions.0.permission"); ok {
						permissionSet := value.(*schema.Set)
						for _, permissionRes := range permissionSet.List() {
							hash := permissionSet.F(permissionRes)
							permission := new(SharePermission)
							if err := permission.UnmarshalHCL(hcl.NewDecoder(decoder, fmt.Sprintf("permissions.0.permission.%d", hash))); err != nil {
								return err
							} else {
								me.Permissions = append(me.Permissions, permission)
							}
						}
					}
				}
			}
		}
	} else {
		me.Permissions = nil
	}
	if len(me.Permissions) == 0 {
		me.Permissions = nil
	}
	if me.Permissions == nil {
		me.Permissions = SharePermissions{}
	}
	me.PublicAccess = &AnonymousAccess{
		ManagementZoneIDs: []string{},
		URLs:              map[string]string{},
	}
	if value, ok := decoder.GetOk("public.#"); ok {
		count := value.(int)
		if count != 0 {
			anonAccess := &AnonymousAccess{}
			anonAccess.UnmarshalHCL(hcl.NewDecoder(decoder, "public.0"))
			me.PublicAccess = anonAccess
		}
	}
	return nil
}
