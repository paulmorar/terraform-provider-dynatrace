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

package browser

import (
	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors"
	browser "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/browser/settings"
)

const SchemaID = "v1:synthetic:monitors:browser"

func Service(credentials *api.Credentials) api.CRUDService[*browser.SyntheticMonitor] {
	return api.NewCRUDService(credentials, SchemaID, &api.ServiceOptions[*browser.SyntheticMonitor]{
		Get:            api.Path("/api/v1/synthetic/monitors/%s"),
		List:           api.Path("/api/v1/synthetic/monitors?type=BROWSER"),
		CreateURL:      func(v *browser.SyntheticMonitor) string { return "/api/v1/synthetic/monitors" },
		Stubs:          &monitors.Monitors{},
		HasNoValidator: true,
		CreateConfirm:  20,
	})
}
