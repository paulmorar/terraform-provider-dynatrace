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
	"context"
	"time"

	"github.com/dtcookie/dynatrace/api/config/synthetic/monitors"
	"github.com/dtcookie/hcl"
	"github.com/dtcookie/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/hcl2sdk"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/logging"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource produces terraform resource definition for Management Zones
func HTTPResource() *schema.Resource {
	return &schema.Resource{
		Schema:        hcl2sdk.Convert(new(monitors.HTTPSyntheticMonitorUpdate).Schema()),
		CreateContext: logging.Enable(HTTPCreate),
		UpdateContext: logging.Enable(HTTPUpdate),
		ReadContext:   logging.Enable(HTTPRead),
		DeleteContext: logging.Enable(Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

// HTTPCreate expects the configuration within the given ResourceData and sends it to the Dynatrace Server in order to create that resource
func HTTPCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := new(monitors.HTTPSyntheticMonitorUpdate)
	if err := config.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	config.ID = nil
	id, err := NewService(m).CreateHTTP(config)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*id)
	return httpReRead(ctx, d, m)
}

// HTTPUpdate expects the configuration within the given ResourceData and send them to the Dynatrace Server in order to update that resource
func HTTPUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := new(monitors.HTTPSyntheticMonitorUpdate)
	if err := config.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	config.ID = opt.NewString(d.Id())
	if err := NewService(m).UpdateHTTP(config); err != nil {
		return diag.FromErr(err)
	}
	return Read(ctx, d, m)
}

func httpReRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	srvc := NewService(m)
	attempts := 0
	var config *monitors.HTTPSyntheticMonitorUpdate
	var err error
	for attempts < 30 {
		if config, err = srvc.GetHTTP(d.Id()); err == nil {
			attempts = 30
		} else {
			attempts++
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return diag.FromErr(err)
	}
	marshalled, err := config.MarshalHCL()
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}

// HTTPRead queries the Dynatrace Server for the configuration
func HTTPRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config, err := NewService(m).GetHTTP(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	marshalled, err := config.MarshalHCL()
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}
