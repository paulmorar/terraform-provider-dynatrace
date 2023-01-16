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
	"io"
	"os"
	"path"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hclgen"
)

type DataSourceStatus string

func (me DataSourceStatus) IsOneOf(stati ...DataSourceStatus) bool {
	if len(stati) == 0 {
		return false
	}
	for _, status := range stati {
		if me == status {
			return true
		}
	}
	return false
}

var DataSourceStati = struct {
	Downloaded DataSourceStatus
	Erronous   DataSourceStatus
	Discovered DataSourceStatus
}{
	"Downloaded",
	"Erronous",
	"Discovered",
}

type DataSource struct {
	ID         string
	Name       string
	UniqueName string
	Type       DataSourceType
	Module     *DataSourceModule
	Status     DataSourceStatus
	Error      error
}

func (me *DataSource) SetName(name string) *DataSource {
	me.Name = name
	terraformName := toTerraformName(name)
	me.UniqueName = me.Module.namer.Name(terraformName)
	me.Status = DataSourceStati.Discovered
	return me
}

func (me *DataSource) GetFile(base string) string {
	filename := fmt.Sprintf("%s.%s.tf", strings.TrimSpace(me.UniqueName), me.Type.Trim())
	filename = fileSystemName(filename)
	return path.Join(me.Module.GetFolder(base), filename)
}

func (me *DataSource) Download(credentials *settings.Credentials, outputFolder string) error {
	if me.Status.IsOneOf(DataSourceStati.Erronous, DataSourceStati.Downloaded) {
		return nil
	}

	var err error

	if me.Module.Status == ModuleStati.Erronous {
		me.Status = DataSourceStati.Erronous
	}

	if me.Module.Status == ModuleStati.Untouched {
		if err = me.Module.Discover(credentials); err != nil {
			return err
		}
	}

	var service = me.Module.Service

	settings := me.Module.Descriptor.NewSettings()
	if err = service.Get(me.ID, settings); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if strings.HasPrefix(restError.Message, "Editing or deleting a non user specific dashboard preset is not allowed.") {
				me.Status = DataSourceStati.Erronous
				me.Error = err
				return nil
			}
			if restError.Code == 404 {
				me.Status = DataSourceStati.Erronous
				me.Error = err
				return nil
			}
			if strings.HasPrefix(restError.Message, "Token is missing required scope") {
				me.Status = DataSourceStati.Erronous
				me.Error = err
				return nil
			}
		}
		return err
	}
	os.MkdirAll(me.Module.GetFolder(outputFolder), os.ModePerm)
	var outputFile *os.File
	outputFileName := me.GetFile(outputFolder)
	if outputFile, err = os.Create(outputFileName); err != nil {
		return err
	}
	defer outputFile.Close()

	finalComments := []string{"ID " + me.ID}

	if err = hclgen.ExportDataSource(settings, outputFile, string(me.Type), me.UniqueName, finalComments...); err != nil {
		return err
	}
	me.Status = DataSourceStati.Downloaded
	return nil
}

func (me *DataSource) DownloadTo(credentials *settings.Credentials, w io.Writer) error {
	if me.Status.IsOneOf(DataSourceStati.Erronous) {
		return nil
	}

	var err error

	if me.Module.Status == ModuleStati.Erronous {
		me.Status = DataSourceStati.Erronous
	}

	if me.Module.Status == ModuleStati.Untouched {
		if err = me.Module.Discover(credentials); err != nil {
			return err
		}
	}

	var service = me.Module.Service

	settings := me.Module.Descriptor.NewSettings()
	if err = service.Get(me.ID, settings); err != nil {
		if restError, ok := err.(rest.Error); ok {
			if strings.HasPrefix(restError.Message, "Editing or deleting a non user specific dashboard preset is not allowed.") {
				me.Status = DataSourceStati.Erronous
				me.Error = err
				return nil
			}
			if restError.Code == 404 {
				me.Status = DataSourceStati.Erronous
				me.Error = err
				return nil
			}
			if strings.HasPrefix(restError.Message, "Token is missing required scope") {
				me.Status = DataSourceStati.Erronous
				me.Error = err
				return nil
			}
		}
		return err
	}

	finalComments := []string{"ID " + me.ID}

	if err = hclgen.ExportDataSource(settings, w, string(me.Type), me.UniqueName, finalComments...); err != nil {
		return err
	}
	me.Status = DataSourceStati.Downloaded
	return nil
}
