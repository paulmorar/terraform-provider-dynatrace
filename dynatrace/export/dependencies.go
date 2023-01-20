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
	"fmt"
	"regexp"
	"strings"
)

type Dependency interface {
	Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any)
	ResourceType() ResourceType
	DataSourceType() DataSourceType
}

var Dependencies = struct {
	ManagementZone   Dependency
	LegacyID         func(resourceType ResourceType) Dependency
	ID               func(resourceType ResourceType) Dependency
	DataSourceID     func(DataSourceType) Dependency
	RequestAttribute Dependency
}{
	ManagementZone:   &mgmzdep{ResourceTypes.ManagementZoneV2},
	LegacyID:         func(resourceType ResourceType) Dependency { return &legacyID{resourceType} },
	ID:               func(resourceType ResourceType) Dependency { return &iddep{resourceType} },
	DataSourceID:     func(dataSourceType DataSourceType) Dependency { return &dsid{dataSourceType} },
	RequestAttribute: &reqAttName{ResourceTypes.RequestAttribute},
}

type mgmzdep struct {
	resourceType ResourceType
}

func (me *mgmzdep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *mgmzdep) DataSourceType() DataSourceType {
	return ""
}

func (me *mgmzdep) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
	replacePattern := "${var.%s.%s.name}"
	if environment.Flags.Flat {
		replacePattern = "${%s.%s.name}"
	}
	resources := []any{}
	for _, resource := range environment.Module(me.resourceType).Resources {
		if resource.Status.IsOneOf(ResourceStati.Erronous, ResourceStati.Excluded) {
			continue
		}
		found := false

		m1 := regexp.MustCompile(fmt.Sprintf(`"managementZone": {([\S\s]*)"name":(.*)"%s"([\S\s]*)}`, resource.Name))
		replaced := m1.ReplaceAllString(s, fmt.Sprintf(`"managementZone": {$1"name":$2"%s"$3}`, fmt.Sprintf(replacePattern, me.resourceType, resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		m1 = regexp.MustCompile(fmt.Sprintf(`management_zones = \[(.*)\"%s\"(.*)\]`, resource.Name))
		replaced = m1.ReplaceAllString(s, fmt.Sprintf(`management_zones = [ $1"%s"$2 ]`, fmt.Sprintf(replacePattern, me.resourceType, resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}

		if found {
			resources = append(resources, resource)
		}
	}
	return s, resources
}

type legacyID struct {
	resourceType ResourceType
}

func (me *legacyID) ResourceType() ResourceType {
	return me.resourceType
}

func (me *legacyID) DataSourceType() DataSourceType {
	return ""
}

func (me *legacyID) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
	replacePattern := "${var.%s.%s.legacy_id}"
	if environment.Flags.Flat {
		replacePattern = "${%s.%s.legacy_id}"
	}
	resources := []any{}
	for _, resource := range environment.Module(me.resourceType).Resources {
		if len(resource.LegacyID) > 0 {
			found := false
			if strings.Contains(s, resource.LegacyID) {
				s = strings.ReplaceAll(s, resource.LegacyID, fmt.Sprintf(replacePattern, me.resourceType, resource.UniqueName))
				found = true
			}
			if found {
				resources = append(resources, resource)
			}
		}
	}
	return s, resources
}

type iddep struct {
	resourceType ResourceType
}

func (me *iddep) ResourceType() ResourceType {
	return me.resourceType
}

func (me *iddep) DataSourceType() DataSourceType {
	return ""
}

func (me *iddep) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
	replacePattern := "${var.%s.%s.id}"
	if environment.Flags.Flat {
		replacePattern = "${%s.%s.id}"
	}
	resources := []any{}
	for id, resource := range environment.Module(me.resourceType).Resources {
		if resource.Type == replacingIn {
			replacePattern = "${%s.%s.id}"
		}
		found := false
		if strings.Contains(s, id) {
			s = strings.ReplaceAll(s, id, fmt.Sprintf(replacePattern, me.resourceType, resource.UniqueName))
			found = true
		}
		if found {
			resources = append(resources, resource)
		}
	}
	return s, resources
}

type dsid struct {
	dataSourceType DataSourceType
}

func (me *dsid) ResourceType() ResourceType {
	return ""
}

func (me *dsid) DataSourceType() DataSourceType {
	return me.dataSourceType
}

func (me *dsid) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
	datasources := []any{}
	for id, datasource := range environment.DataSourceModule(me.dataSourceType).DataSources {
		found := false
		if strings.Contains(s, id) {
			s = strings.ReplaceAll(s, id, fmt.Sprintf("${data.%s.%s.id}", me.dataSourceType, datasource.UniqueName))
			found = true
		}
		if found {
			datasources = append(datasources, datasource)
		}
	}
	return s, datasources
}

type reqAttName struct {
	resourceType ResourceType
}

func (me *reqAttName) ResourceType() ResourceType {
	return me.resourceType
}

func (me *reqAttName) DataSourceType() DataSourceType {
	return ""
}

func (me *reqAttName) Replace(environment *Environment, s string, replacingIn ResourceType) (string, []any) {
	replacePattern := "${var.%s.%s.name}"
	if environment.Flags.Flat {
		replacePattern = "${%s.%s.name}"
	}
	resources := []any{}
	for _, resource := range environment.Module(me.resourceType).Resources {
		found := false
		m1 := regexp.MustCompile(fmt.Sprintf("request_attribute(.*)=(.*)\"%s\"", resource.Name))
		replaced := m1.ReplaceAllString(s, fmt.Sprintf("request_attribute$1=$2\"%s\"", fmt.Sprintf(replacePattern, me.resourceType, resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		m1 = regexp.MustCompile(fmt.Sprintf("{RequestAttribute:%s}", resource.Name))
		replaced = m1.ReplaceAllString(s, fmt.Sprintf("{RequestAttribute:%s}", fmt.Sprintf(replacePattern, me.resourceType, resource.UniqueName)))
		if replaced != s {
			s = replaced
			found = true
		}
		if found {
			resources = append(resources, resource)
		}
	}
	return s, resources
}
