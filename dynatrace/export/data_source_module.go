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

package export

import (
	"log"
	"path"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type DataSourceModule struct {
	Environment *Environment
	Type        DataSourceType
	DataSources map[string]*DataSource
	namer       UniqueNamer
	Status      ModuleStatus
	Error       error
	Descriptor  *DataSourceDescriptor
	Service     settings.RService[settings.Settings]
}

func (me *DataSourceModule) DataSource(id string) *DataSource {
	if stored, found := me.DataSources[id]; found {
		return stored
	}
	res := &DataSource{ID: id, Type: me.Type, Module: me, Status: DataSourceStati.Discovered}
	me.DataSources[id] = res
	return res
}

func (me *DataSourceModule) GetFolder(base string) string {
	return path.Join(base, me.Type.Trim())
}

func (me *DataSourceModule) Discover(credentials *settings.Credentials) error {
	if me.Status.IsOneOf(ModuleStati.Downloaded, ModuleStati.Discovered, ModuleStati.Erronous) {
		return nil
	}

	if me.Descriptor == nil {
		descriptor := AllDataSources[me.Type]
		me.Descriptor = &descriptor
	}

	if me.Service == nil {
		me.Service = me.Descriptor.Service(credentials)
	}

	var err error

	var stubs []*settings.Stub
	log.Println("Discovering \"" + me.Type + "\" ...")
	if stubs, err = me.Service.List(); err != nil {
		if strings.Contains(err.Error(), "Token is missing required scope") {
			me.Status = ModuleStati.Erronous
			me.Error = err
			return nil
		}
		if strings.Contains(err.Error(), "No schema with topic identifier") {
			me.Status = ModuleStati.Erronous
			me.Error = err
			return nil
		}
		return err
	}
	for _, stub := range stubs {
		me.DataSource(stub.ID).SetName(stub.Name)
	}
	me.Status = ModuleStati.Discovered
	// log.Println("   ", fmt.Sprintf("%d items found", len(stubs)))
	return nil
}
