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

package keyrequests

import (
	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services/cache"
	v2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services/v2"
	services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/services"
	keyrequests "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/keyrequests/settings"
)

const SchemaID = "builtin:settings.subscriptions.service"

func Service(credentials *api.Credentials) api.CRUDService[*keyrequests.KeyRequest] {
	topologyService := cache.Read(services.Service(credentials))
	return v2.Service(credentials, SchemaID, &v2.ServiceOptions[*keyrequests.KeyRequest]{
		Name: func(id string, v *keyrequests.KeyRequest) (string, error) {
			service := api.NewSettings(topologyService)
			if err := topologyService.Get(v.ServiceID, service); err != nil {
				return "", err
			}
			return "Key Requests for " + service.DisplayName, nil
		},
	})
}
