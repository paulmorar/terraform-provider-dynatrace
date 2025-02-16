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

package vault

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"

	vault "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/vault/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the credential. Possible values are `CERTIFICATE`, `PUBLIC_CERTIFICATE`, `TOKEN`, `USERNAME_PASSWORD` and `UNKNOWN`. If not specified all credential types will match",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the credential as shown within the Dynatrace WebUI. If not specified all names will match",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The scope of the credential. Possible values are `ALL`, `EXTENSION` and `SYNTHETIC`. If not specified all scopes will match.",
			},
		},
	}
}

func DataSourceRead(d *schema.ResourceData, m any) (err error) {
	name := ""
	typ := ""
	scope := ""
	if value, ok := d.GetOk("name"); ok {
		name = value.(string)
	}
	if value, ok := d.GetOk("type"); ok {
		typ = value.(string)
	}
	if value, ok := d.GetOk("scope"); ok {
		scope = value.(string)
	}
	if name == "" && typ == "" && scope == "" {
		return fmt.Errorf("at least one of `name`, `type` or `scope` needs to be specified as a non empty string")
	}

	service := export.Service(config.Credentials(m), export.ResourceTypes.Credentials)
	var stubs settings.Stubs
	if stubs, err = service.List(); err != nil {
		return err
	}
	if len(stubs) == 0 {
		d.SetId("")
	}
	for _, stub := range stubs {
		if name != "" && stub.Name != name {
			continue
		}
		var credentials vault.Credentials
		if err = service.Get(stub.ID, &credentials); err != nil {
			return err
		}
		if scope != "" && string(credentials.Scope) != scope {
			continue
		}
		if typ != "" && string(credentials.Type) != typ {
			continue
		}
		d.Set("scope", string(credentials.Scope))
		d.Set("type", string(credentials.Type))
		d.Set("name", stub.Name)
		d.SetId(stub.ID)
		return nil
	}

	d.SetId("")
	return nil
}
