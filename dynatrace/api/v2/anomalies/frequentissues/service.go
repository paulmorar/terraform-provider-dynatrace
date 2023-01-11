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

package frequentissues

import (
	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
	v2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services/v2"
	frequentissues "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/anomalies/frequentissues/settings"
)

const SchemaID = "builtin:anomaly-detection.frequent-issues"

func Service(credentials *api.Credentials) api.CRUDService[*frequentissues.FrequentIssues] {
	return v2.Service[*frequentissues.FrequentIssues](credentials, SchemaID)
}
