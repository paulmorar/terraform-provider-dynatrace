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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/aws/iam"
	ds_services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/services"
)

type DataSourceDescriptor struct {
	Service   func(credentials *settings.Credentials) settings.RService[settings.Settings]
	protoType settings.Settings
}

func (me DataSourceDescriptor) NewSettings() settings.Settings {
	return reflect.New(reflect.TypeOf(me.protoType).Elem()).Interface().(settings.Settings)
}

func NewDataSourceDescriptor[T settings.Settings](fn func(credentials *settings.Credentials) settings.RService[T]) DataSourceDescriptor {
	return DataSourceDescriptor{
		Service: func(credentials *settings.Credentials) settings.RService[settings.Settings] {
			return &settings.GenericRService[T]{Service: cache.Read(fn(credentials))}
		},
		protoType: newSettingsRead(fn),
	}
}

func newSettingsRead[T settings.Settings](sfn func(credentials *settings.Credentials) settings.RService[T]) T {
	var proto T
	return reflect.New(reflect.TypeOf(proto).Elem()).Interface().(T)
}

var AllDataSources = map[DataSourceType]DataSourceDescriptor{
	DataSourceTypes.Service:          NewDataSourceDescriptor(ds_services.Service),
	DataSourceTypes.AWSIAMExternalID: NewDataSourceDescriptor(iam.Service),
}
