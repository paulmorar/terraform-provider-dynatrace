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
	"reflect"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services/cache"

	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/aws/iam"
	ds_services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/services"
)

type DataSourceDescriptor struct {
	Service   func(credentials *api.Credentials) api.RService[api.Settings]
	protoType api.Settings
}

func (me DataSourceDescriptor) NewSettings() api.Settings {
	return reflect.New(reflect.TypeOf(me.protoType).Elem()).Interface().(api.Settings)
}

func NewDataSourceDescriptor[T api.Settings](fn func(credentials *api.Credentials) api.RService[T]) DataSourceDescriptor {
	return DataSourceDescriptor{
		Service: func(credentials *api.Credentials) api.RService[api.Settings] {
			return &api.GenericRService[T]{Service: cache.Read(fn(credentials))}
		},
		protoType: newSettingsRead(fn),
	}
}

func newSettingsRead[T api.Settings](sfn func(credentials *api.Credentials) api.RService[T]) T {
	var proto T
	return reflect.New(reflect.TypeOf(proto).Elem()).Interface().(T)
}

var AllDataSources = map[DataSourceType]DataSourceDescriptor{
	DataSourceTypes.Service:          NewDataSourceDescriptor(ds_services.Service),
	DataSourceTypes.AWSIAMExternalID: NewDataSourceDescriptor(iam.Service),
}
