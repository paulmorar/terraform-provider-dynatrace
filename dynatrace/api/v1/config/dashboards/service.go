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

package dashboards

import (
	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services/cache"

	dashboards "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/dashboards/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/jsondashboards"
)

const SchemaID = "v1:config:dashboards"

func Service(credentials *api.Credentials) api.CRUDService[*dashboards.Dashboard] {
	return &service{service: cache.CRUD(jsondashboards.Service(credentials))}
}

type service struct {
	service api.CRUDService[*dashboards.JSONDashboard]
}

func (me *service) NoCache() bool {
	return true
}

func (me *service) List() (api.Stubs, error) {
	return me.service.List()
}

func (me *service) Get(id string, v *dashboards.Dashboard) error {
	var err error
	var data []byte
	jsondb := api.NewSettings(me.service.(api.RService[*dashboards.JSONDashboard]))
	if err = me.service.Get(id, jsondb); err != nil {
		return err
	}
	if data, err = api.ToJSON(jsondb); err != nil {
		return err
	}
	return api.FromJSON(data, v)
}

func (me *service) Validate(v *dashboards.Dashboard) error {
	var err error
	var data []byte
	jsondb := api.NewSettings(me.service.(api.RService[*dashboards.JSONDashboard]))
	if data, err = api.ToJSON(v); err != nil {
		return err
	}
	if err = api.FromJSON(data, jsondb); err != nil {
		return err
	}
	if validator, ok := me.service.(api.Validator[*dashboards.JSONDashboard]); ok {
		return validator.Validate(jsondb)
	}
	return nil
}

func (me *service) Create(v *dashboards.Dashboard) (*api.Stub, error) {
	var err error
	var data []byte
	jsondb := api.NewSettings(me.service.(api.RService[*dashboards.JSONDashboard]))
	if data, err = api.ToJSON(v); err != nil {
		return nil, err
	}
	if err = api.FromJSON(data, jsondb); err != nil {
		return nil, err
	}
	return me.service.Create(jsondb)
}

func (me *service) Update(id string, v *dashboards.Dashboard) error {
	var err error
	var data []byte
	jsondb := api.NewSettings(me.service.(api.RService[*dashboards.JSONDashboard]))
	if data, err = api.ToJSON(v); err != nil {
		return err
	}
	if err = api.FromJSON(data, jsondb); err != nil {
		return err
	}
	return me.service.Update(id, jsondb)
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) SchemaID() string {
	return SchemaID
}
