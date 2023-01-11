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

package v2

import (
	"encoding/json"
	"fmt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"log"

	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"

	"net/url"
)

func Service[T api.Settings](credentials *api.Credentials, schemaID string, options ...*ServiceOptions[T]) api.CRUDService[T] {
	var opts *ServiceOptions[T]
	if len(options) > 0 {
		opts = options[0]
	}
	return &service[T]{
		schemaID: schemaID,
		client:   rest.DefaultClient(credentials.URL, credentials.Token),
		options:  opts,
	}
}

type SettingsObjectUpdate struct {
	SchemaVersion string `json:"schemaVersion,omitempty"`
	Value         any    `json:"value"`
}

type SettingsObjectCreate struct {
	SchemaVersion string `json:"schemaVersion,omitempty"`
	SchemaID      string `json:"schemaId"`
	Scope         string `json:"scope"`
	Value         any    `json:"value"`
}

type SettingsObjectCreateResponse struct {
	ObjectID string `json:"objectId"`
}

type service[T api.Settings] struct {
	schemaID string
	client   rest.Client
	options  *ServiceOptions[T]
}

func (me *service[T]) LegacyID() func(id string) string {
	if me.options != nil && me.options.LegacyID != nil {
		return me.options.LegacyID
	}
	return nil
}

func (me *service[T]) Get(id string, v T) error {
	var err error
	var settingsObject SettingsObject

	req := me.client.Get(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id))).Expect(200)
	if err = req.Finish(&settingsObject); err != nil {
		return err
	}

	if err = json.Unmarshal(settingsObject.Value, v); err != nil {
		return err
	}
	if scopeAware, ok := any(v).(ScopeAware); ok {
		scopeAware.SetScope(settingsObject.Scope)
	}
	if me.options != nil && me.options.LegacyID != nil {
		api.SetLegacyID(id, me.options.LegacyID, v)
	}

	return nil
}

func (me *service[T]) List() (api.Stubs, error) {
	var err error

	stubs := api.Stubs{}
	nextPage := true

	var nextPageKey *string
	for nextPage {
		var sol SettingsObjectList
		var urlStr string
		if nextPageKey != nil {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?nextPageKey=%s", url.QueryEscape(*nextPageKey))
		} else {
			urlStr = fmt.Sprintf("/api/v2/settings/objects?schemaIds=%s&fields=%s&pageSize=500", url.QueryEscape(me.SchemaID()), url.QueryEscape("objectId,value,scope,schemaVersion"))
		}
		req := me.client.Get(urlStr, 200)
		if err = req.Finish(&sol); err != nil {
			return nil, err
		}

		if len(sol.Items) > 0 {
			for _, item := range sol.Items {
				newItem := api.NewSettings[T](me)
				if err = json.Unmarshal(item.Value, &newItem); err != nil {
					return nil, err
				}
				if me.options != nil && me.options.LegacyID != nil {
					api.SetLegacyID(item.ObjectID, me.options.LegacyID, newItem)
				}
				if scopeAware, ok := any(newItem).(ScopeAware); ok {
					scopeAware.SetScope(item.Scope)
				}
				var itemName string
				if me.options != nil && me.options.Name != nil {
					if itemName, err = me.options.Name(item.ObjectID, newItem); err != nil {
						log.Println(err.Error())
						itemName = api.Name(newItem)
					}
				} else {
					itemName = api.Name(newItem)
				}
				stub := &api.Stub{ID: item.ObjectID, Name: itemName, Value: newItem, LegacyID: api.GetLegacyID(newItem)}
				if len(itemName) > 0 {
					stubs = append(stubs, stub)
				}
			}
		}
		nextPageKey = sol.NextPageKey
		nextPage = (nextPageKey != nil)
	}

	return stubs, nil
}

func (me *service[T]) Validate(v T) error {
	return nil // Settings 2.0 doesn't offer validation
}

func (me *service[T]) Create(v T) (*api.Stub, error) {
	soc := SettingsObjectCreate{
		SchemaID: me.schemaID,
		Scope:    "environment",
		Value:    v,
	}
	if scopeAware, ok := any(v).(ScopeAware); ok {
		soc.Scope = scopeAware.GetScope()
	}

	req := me.client.Post("/api/v2/settings/objects", []SettingsObjectCreate{soc}).Expect(200)
	objectID := []SettingsObjectCreateResponse{}

	if err := req.Finish(&objectID); err != nil {
		return nil, err
	}
	itemName := api.Name(v)
	stub := &api.Stub{ID: objectID[0].ObjectID, Name: itemName}
	return stub, nil
}

func (me *service[T]) Update(id string, v T) error {
	sou := SettingsObjectUpdate{Value: v}
	return me.client.Put(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), &sou, 200).Finish()
}

func (me *service[T]) Delete(id string) error {
	return me.client.Delete(fmt.Sprintf("/api/v2/settings/objects/%s", url.PathEscape(id)), 204).Finish()
}

func (me *service[T]) Name() string {
	return me.SchemaID()
}

func (me *service[T]) SchemaID() string {
	return me.schemaID
}
