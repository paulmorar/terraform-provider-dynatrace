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

package hosts

import (
	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
	hosts "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/hosts/settings"
)

const SchemaID = "v1:environment:hosts"
const BasePath = "/api/v1/entity/infrastructure/hosts"

func Service(credentials *api.Credentials) api.RService[*hosts.Host] {
	return api.NewCRUDService(
		credentials,
		SchemaID,
		api.DefaultServiceOptions[*hosts.Host](BasePath).WithStubs(new(hosts.Hosts)),
	)
}
