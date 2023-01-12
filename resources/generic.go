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

package resources

import (
	"context"
	"os"

	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NewGeneric(resourceType export.ResourceType) *Generic {
	descriptor := export.AllResources[resourceType]
	return &Generic{Type: resourceType, Descriptor: descriptor}
}

type Generic struct {
	Type       export.ResourceType
	Descriptor export.ResourceDescriptor
}

func (me *Generic) Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        me.Descriptor.NewSettings().Schema(),
		CreateContext: logging.Enable(me.Create),
		UpdateContext: logging.Enable(me.Update),
		ReadContext:   logging.Enable(me.Read),
		DeleteContext: logging.Enable(me.Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func (me *Generic) createCredentials(m any) *api.Credentials {
	conf := m.(*config.ProviderConfiguration)
	return &api.Credentials{
		Token: conf.APIToken,
		URL:   conf.EnvironmentURL,
	}
}

func (me *Generic) Settings() api.Settings {
	return me.Descriptor.NewSettings()
}

func (me *Generic) Service(m any) api.CRUDService[api.Settings] {
	return me.Descriptor.Service(me.createCredentials(m))
}

func (me *Generic) Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	settings := me.Settings()
	if err := settings.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	stub, err := me.Service(m).Create(settings)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(stub.ID)
	return me.Read(ctx, d, m)
}

func (me *Generic) Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	settings := me.Settings()
	if err := settings.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	if err := me.Service(m).Update(d.Id(), settings); err != nil {
		return diag.FromErr(err)
	}
	return me.Read(ctx, d, m)
}

func (me *Generic) Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var err error
	var restLogFile *os.File
	restLogFileName := os.Getenv("DT_REST_DEBUG_LOG")
	if len(restLogFileName) > 0 {
		if restLogFile, err = os.Create(restLogFileName); err != nil {
			return diag.FromErr(err)
		}
		rest.SetLogWriter(restLogFile)
	}
	settings := me.Settings()
	// if os.Getenv("CACHE_OFFLINE_MODE") != "true" {
	// 	if _, ok := settings.(*vault.Credentials); ok {
	// 		return diag.Diagnostics{}
	// 	}
	// 	if _, ok := settings.(*notifications.Notification); ok {
	// 		return diag.Diagnostics{}
	// 	}
	// }
	service := me.Service(m)
	if err := service.Get(d.Id(), settings); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if restError.Code == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
		}
		return diag.FromErr(err)
	}
	if preparer, ok := settings.(MarshalPreparer); ok {
		preparer.PrepareMarshalHCL(hcl.DecoderFrom(d))
	}
	if os.Getenv("DT_TERRAFORM_IMPORT") == "true" {
		if demoSettings, ok := settings.(api.DemoSettings); ok {
			demoSettings.FillDemoValues()
		}
	}
	marshalled, err := settings.MarshalHCL(hcl.DecoderFrom(d))
	if err != nil {
		return diag.FromErr(err)
	}
	for k, v := range marshalled {
		d.Set(k, v)
	}
	return diag.Diagnostics{}
}

func (me *Generic) Delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	if err := me.Service(m).Delete(d.Id()); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if restError.Code == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
		}
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

type MarshalPreparer interface {
	PrepareMarshalHCL(hcl.Decoder) error
}
