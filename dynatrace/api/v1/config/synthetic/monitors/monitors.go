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
	"strings"

	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
)

type Monitors struct {
	Monitors []*MonitorCollectionElement `json:"monitors"` // The list of synthetic monitors
}

func (me *Monitors) ToStubs() api.Stubs {
	stubs := api.Stubs{}
	if len(me.Monitors) > 0 {
		for _, monitor := range me.Monitors {
			if !strings.Contains(monitor.Name, "synchronizing credentials with") {
				stubs = append(stubs, &api.Stub{ID: monitor.EntityID, Name: monitor.Name})
			}
		}
	}
	return stubs
}
