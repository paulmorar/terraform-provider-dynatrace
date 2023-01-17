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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
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
	stngs := me.Descriptor.NewSettings()
	sch := stngs.Schema()
	// implicitUpate := false
	// stnt := reflect.ValueOf(stngs).Elem().Type()
	// for idx := 0; idx < stnt.NumField(); idx++ {
	// 	field := stnt.Field(idx)
	// 	if field.Type == implicitUpdateType {
	// 		implicitUpate = true
	// 		break
	// 	}
	// }
	// if implicitUpate {
	// 	sch["replaced_value"] = &schema.Schema{
	// 		Type:        schema.TypeString,
	// 		Description: "for internal use only",
	// 		Optional:    true,
	// 		Computed:    true,
	// 	}
	// }

	return &schema.Resource{
		Schema:        sch,
		CreateContext: logging.Enable(me.Create),
		UpdateContext: logging.Enable(me.Update),
		ReadContext:   logging.Enable(me.Read),
		DeleteContext: logging.Enable(me.Delete),
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func (me *Generic) createCredentials(m any) *settings.Credentials {
	conf := m.(*config.ProviderConfiguration)
	return &settings.Credentials{
		Token: conf.APIToken,
		URL:   conf.EnvironmentURL,
	}
}

func (me *Generic) Settings() settings.Settings {
	return me.Descriptor.NewSettings()
}

func (me *Generic) Service(m any) settings.CRUDService[settings.Settings] {
	return me.Descriptor.Service(me.createCredentials(m))
}

func (me *Generic) Create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sttngs := me.Settings()
	if err := sttngs.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	stub, err := me.Service(m).Create(sttngs)
	if err != nil {
		// if restError, ok := err.(rest.Error); ok {
		// 	if len(restError.ConstraintViolations) == 1 {
		// 		if restError.ConstraintViolations[0].Message == "Management zone with this name already exists. Please provide a different one." {
		// 			stubs, e2 := me.Service(m).List()
		// 			if e2 != nil {
		// 				return diag.FromErr(err)
		// 			}
		// 			foundID := ""
		// 			for _, stub := range stubs {
		// 				if settings.Name(sttngs) == stub.Name {
		// 					foundID = stub.ID
		// 					break
		// 				}
		// 			}
		// 			if foundID == "" {
		// 				return diag.FromErr(err)
		// 			}
		// 			d.SetId(foundID)
		// 			replaceSettings := me.Settings()
		// 			if e2 = me.Service(m).Get(foundID, replaceSettings); e2 != nil {
		// 				return diag.FromErr(err)
		// 			}
		// 			if err := me.Service(m).Update(foundID, sttngs); err != nil {
		// 				return diag.FromErr(err)
		// 			}
		// 			data, e2 := json.Marshal(replaceSettings)
		// 			if e2 != nil {
		// 				return diag.FromErr(err)
		// 			}
		// 			buf := new(bytes.Buffer)
		// 			// writer, e2 := gzip.NewWriterLevel(buf, gzip.BestCompression)
		// 			writer, e2 := zlib.NewWriterLevel(buf, gzip.BestCompression)
		// 			if e2 != nil {
		// 				return diag.FromErr(err)
		// 			}
		// 			writer.Write(data)
		// 			writer.Flush()

		// 			if e2 != nil {
		// 				return diag.FromErr(err)
		// 			}
		// 			d.Set("replaced_value", base64.StdEncoding.EncodeToString(buf.Bytes()))
		// 			return me.Read(ctx, d, m)
		// 		}
		// 	}
		// }
		return diag.FromErr(err)
	}
	d.SetId(stub.ID)
	return me.Read(ctx, d, m)
}

func (me *Generic) Update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sttngs := me.Settings()
	if err := sttngs.UnmarshalHCL(hcl.DecoderFrom(d)); err != nil {
		return diag.FromErr(err)
	}
	if err := me.Service(m).Update(d.Id(), sttngs); err != nil {
		return diag.FromErr(err)
	}
	return me.Read(ctx, d, m)
}

func (me *Generic) Read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	var err error
	sttngs := me.Settings()
	// if os.Getenv("CACHE_OFFLINE_MODE") != "true" {
	// 	if _, ok := settings.(*vault.Credentials); ok {
	// 		return diag.Diagnostics{}
	// 	}
	// 	if _, ok := settings.(*notifications.Notification); ok {
	// 		return diag.Diagnostics{}
	// 	}
	// }
	service := me.Service(m)
	if err := service.Get(d.Id(), sttngs); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if restError.Code == 404 {
				d.SetId("")
				return diag.Diagnostics{}
			}
		}
		return diag.FromErr(err)
	}
	if preparer, ok := sttngs.(MarshalPreparer); ok {
		preparer.PrepareMarshalHCL(hcl.DecoderFrom(d))
	}
	if os.Getenv("DT_TERRAFORM_IMPORT") == "true" {
		if demoSettings, ok := sttngs.(settings.DemoSettings); ok {
			demoSettings.FillDemoValues()
		}
	}
	marshalled := hcl.Properties{}
	err = sttngs.MarshalHCL(marshalled)
	attributes := Attributes{}
	attributes.collect("", map[string]any(marshalled))
	stateAttributes := NewAttributes(d.State().Attributes)
	for key, value := range attributes {
		if value == "${state.secret_value}" {
			matches := stateAttributes.MatchingKeys(key)
			siblings := attributes.Siblings(key)
			for _, m := range matches {
				sibs := stateAttributes.Siblings(m)
				if sibs.Contains(siblings...) {
					store(marshalled, key, stateAttributes[m])
					break
				}
			}
		}
	}
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
