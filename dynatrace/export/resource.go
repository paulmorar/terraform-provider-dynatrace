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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hclgen"
)

type Resource struct {
	ID                   string
	LegacyID             string
	Name                 string
	UniqueName           string
	Type                 ResourceType
	Module               *Module
	Status               ResourceStatus
	Error                error
	ResourceReferences   []*Resource
	DataSourceReferences []*DataSource
}

func (me *Resource) SetName(name string) *Resource {
	me.Name = name
	terraformName := toTerraformName(name)
	me.UniqueName = me.Module.namer.Name(terraformName)
	me.Status = ResourceStati.Discovered
	return me
}

func (me *Resource) GetResourceReferences() []*Resource {
	resources := map[string]*Resource{}
	if len(me.ResourceReferences) == 0 {
		return []*Resource{}
	}
	for _, resource := range me.ResourceReferences {
		if !resource.Status.IsOneOf(ResourceStati.PostProcessed, ResourceStati.Downloaded) {
			continue
		}
		key := fmt.Sprintf("%s.%s", resource.ID, resource.Type)
		resources[key] = resource
		for _, resource := range resource.GetResourceReferences() {
			if !resource.Status.IsOneOf(ResourceStati.PostProcessed, ResourceStati.Downloaded) {
				continue
			}
			key := fmt.Sprintf("%s.%s", resource.ID, resource.Type)
			resources[key] = resource
		}
	}
	result := []*Resource{}
	for _, resource := range resources {
		result = append(result, resource)
	}
	return result
}

func (me *Resource) RefersTo(other *Resource) bool {
	if other == nil {
		return false
	}
	if len(me.ResourceReferences) == 0 {
		return false
	}
	for _, reference := range me.ResourceReferences {
		if reference.Type == other.Type && reference.ID == other.ID {
			return true
		}
	}
	return false
}

func (me *Resource) CreateFile() (*os.File, error) {
	return os.Create(me.GetFile())
}

func (me *Resource) ReadFile() ([]byte, error) {
	return os.ReadFile(me.GetFile())
}

func (me *Resource) GetFileName() string {
	return fileSystemName(fmt.Sprintf("%s.%s.tf", strings.TrimSpace(me.UniqueName), me.Type.Trim()))
}

func (me *Resource) GetFile() string {
	return path.Join(me.Module.GetFolder(), me.GetFileName())
}

func (me *Resource) GetAttentionFile() string {
	return path.Join(me.Module.GetAttentionFolder(false), me.GetFileName())
}

func (me *Resource) GetFlawedFile() string {
	return path.Join(me.Module.GetFlawedFolder(false), me.GetFileName())
}

func (me *Resource) CreateFlawedFile() (*os.File, error) {
	flawedFile := me.GetFlawedFile()
	absdir, _ := filepath.Abs(path.Dir(flawedFile))
	os.MkdirAll(absdir, os.ModePerm)
	return os.Create(flawedFile)
}

func (me *Resource) Download() error {
	if me.Status.IsOneOf(ResourceStati.Erronous, ResourceStati.Excluded, ResourceStati.Downloaded, ResourceStati.PostProcessed) {
		return nil
	}

	var err error

	if me.Module.Status == ModuleStati.Erronous {
		me.Status = ResourceStati.Erronous
	}

	if me.Module.Status == ModuleStati.Untouched {
		if err = me.Module.Discover(); err != nil {
			return err
		}
	}

	if me.Module.Descriptor.except != nil {
		if me.Module.Descriptor.except(me.ID, me.Name) {
			me.Status = ResourceStati.Excluded
			return nil
		}
	}

	var service = me.Module.Service

	settngs := me.Module.Descriptor.NewSettings()
	if err = service.Get(me.ID, settngs); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if strings.HasPrefix(restError.Message, "Editing or deleting a non user specific dashboard preset is not allowed.") {
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
			if restError.Code == 404 {
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
			if strings.HasPrefix(restError.Message, "Token is missing required scope") {
				me.Status = ResourceStati.Erronous
				me.Error = err
				return nil
			}
		}
		return err
	}
	legacyID := settings.GetLegacyID(settngs)
	if legacyID != nil {
		me.LegacyID = *legacyID
	}
	comments := settings.FillDemoValues(settngs)
	if len(comments) > 0 {
		for _, comment := range comments {
			if strings.HasPrefix(comment, "FLAWED SETTINGS") {
				me.Status = ResourceStati.Erronous
			}
		}
	}

	me.Module.MkdirAll()
	var outputFile *os.File
	if me.Status != ResourceStati.Erronous {
		if outputFile, err = me.CreateFile(); err != nil {
			return err
		}
		defer outputFile.Close()
	} else {
		if outputFile, err = me.CreateFlawedFile(); err != nil {
			return err
		}
		defer outputFile.Close()
	}

	finalComments := []string{}
	if me.Module.Environment.Flags.PersistIDs {
		finalComments = []string{"ID " + me.ID}
		if legacyID := settings.ClearLegacyID(settngs); legacyID != nil {
			finalComments = append(finalComments, "LEGACY_ID "+*legacyID)
		}
	}
	if len(comments) > 0 {
		for _, comment := range comments {
			if len(finalComments) > 0 {
				finalComments = append(finalComments, "")
			}
			finalComments = append(finalComments, "ATTENTION "+comment)
		}
	}

	if err = hclgen.ExportResource(settngs, outputFile, string(me.Type), me.UniqueName, finalComments...); err != nil {
		return err
	}
	if me.Status != ResourceStati.Erronous && len(comments) > 0 {
		orig, _ := filepath.Abs(me.GetFile())
		att, _ := filepath.Abs(me.GetAttentionFile())
		absdir, _ := filepath.Abs(path.Dir(me.GetAttentionFile()))
		os.MkdirAll(absdir, os.ModePerm)
		os.Link(orig, att)
	}
	if me.Status != ResourceStati.Erronous {
		me.Status = ResourceStati.Downloaded
	}
	return nil
}

func (me *Resource) PostProcess() error {
	if me.Status.IsOneOf(ResourceStati.Erronous, ResourceStati.Excluded, ResourceStati.PostProcessed) {
		return nil
	}
	var err error
	if me.Status == ResourceStati.Discovered {
		if err = me.Download(); err != nil {
			return err
		}
	}
	me.Status = ResourceStati.PostProcessed
	if !me.Module.Environment.Flags.FollowReferences {
		return nil
	}
	descriptor := me.Module.Descriptor
	if len(descriptor.Dependencies) == 0 {
		return nil
	}
	for _, dependency := range descriptor.Dependencies {
		resourceType := dependency.ResourceType()
		if len(resourceType) == 0 {
			// dataSourceType := dependency.DataSourceType()
			// module := me.Module.Environment.DataSourceModule(dataSourceType)
			// if module.Status == ModuleStati.Erronous {
			// 	continue
			// }
			// if err = module.Discover(credentials); err != nil {
			// 	return err
			// }
		} else {
			module := me.Module.Environment.Module(resourceType)
			if module.Status == ModuleStati.Erronous {
				continue
			}
			if !module.Status.IsOneOf(ModuleStati.Downloaded, ModuleStati.Discovered, ModuleStati.Erronous) {
				if err = module.Discover(); err != nil {
					return err
				}
			}
		}
		var err error
		var data []byte
		if data, err = me.ReadFile(); err != nil {
			return err
		}
		var foundItemsInFileContents []any
		var modifiedFileContents string

		fileContents := string(data)
		idx := strings.Index(fileContents, "\" {")
		fileHeader := fileContents[:idx]
		fileBody := fileContents[idx:]

		if modifiedFileContents, foundItemsInFileContents = dependency.Replace(me.Module.Environment, fileBody, me.Type); len(foundItemsInFileContents) > 0 {
			var outputFile *os.File
			if outputFile, err = me.CreateFile(); err != nil {
				return err
			}
			modifiedFileContents = fileHeader + modifiedFileContents
			defer func() {
				outputFile.Close()
				format(outputFile.Name(), false)
			}()
			if _, err = outputFile.Write([]byte(modifiedFileContents)); err != nil {
				return err
			}
			for _, item := range foundItemsInFileContents {
				switch typedItem := item.(type) {
				case *Resource:
					if err = typedItem.Download(); err != nil {
						return err
					}
					me.ResourceReferences = append(me.ResourceReferences, typedItem)
				case *DataSource:
					// me.DataSourceReferences = append(me.DataSourceReferences, typedItem)
				}
			}
		} else {
			format(me.GetFile(), false)
		}
	}
	return nil
}
