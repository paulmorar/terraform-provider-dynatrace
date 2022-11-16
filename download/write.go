package download

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/version"
)

func ProcessWrite(dlConfig DownloadConfig, resourceDataMap ResourceData, dataSourceDataMap DataSourceData, replacedIDs ReplacedIDs) error {
	var err error

	if err := os.RemoveAll(dlConfig.TargetFolder); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	os.MkdirAll(dlConfig.TargetFolder, os.ModePerm)

	var mainFile *os.File
	mainFileName := dlConfig.TargetFolder + "/" + "main.tf"
	if mainFile, err = os.Create(mainFileName); err != nil {
		return err
	}

	for resName, resources := range resourceDataMap {
		if len(resources) == 0 {
			continue
		}
		if resName == "dynatrace_dashboard_sharing" {
			continue
		}
		resFolder := strings.TrimPrefix(resName, "dynatrace_")
		resNameCnt := NewNameCounter()
		resNameCnt.Replace(func(s string, cnt int) string {
			return fmt.Sprintf("%s_%d", s, cnt)
		})
		os.MkdirAll(dlConfig.TargetFolder+"/"+resFolder, os.ModePerm)
		if err = resourceDataMap.WriteResourceSeparate(dlConfig, resName, resFolder, resources, resNameCnt); err != nil {
			return err
		}
		if ResourceInfoMap[resName].HardcodedIds != nil && dlConfig.References {
			dataSourceDataMap.WriteDataSource(dlConfig, resName, resFolder, replacedIDs)
		}

		if err := writeNestedProviderFile(dlConfig.TargetFolder, resFolder); err != nil {
			return err
		}

		if err := writeMainFile(mainFile, resourceDataMap, resName, resFolder, replacedIDs, dlConfig.References); err != nil {
			return err
		}
	}
	mainFile.Close()

	if err := writeProviderFile(dlConfig.TargetFolder); err != nil {
		return err
	}

	return nil
}

var forbiddenFileNameChars = []string{"<", ">", ":", "\"", "/", "|", "?", "*", "	", "\r", "\n", "\f", "\v"}

func escf(s string) string {
	for _, ch := range forbiddenFileNameChars {
		s = strings.ReplaceAll(s, ch, "_")
	}
	return s
}

func escape(s string) string {
	result := ""
	for _, c := range s {
		if unicode.IsLetter(c) {
			result = result + string(c)
		} else if unicode.IsDigit(c) {
			result = result + string(c)
		} else if c == '-' {
			result = result + string(c)
		} else if c == '_' {
			result = result + string(c)
		} else {
			result = result + "_"
		}
	}
	return result
}

func writeProviderFile(targetFolder string) error {
	var err error
	var providerFile *os.File
	fileName := targetFolder + "/" + "providers.tf"
	if providerFile, err = os.Create(fileName); err != nil {
		return err
	}
	content := `terraform {
  required_providers {
    dynatrace = {
    version = "${version}"
    source = "dynatrace-oss/dynatrace"
    }
  }
}
	
# provider "dynatrace" {
#   alias        = "default"
#   dt_env_url   = "https://########.live.dynatrace.com/"
#   dt_api_token = "dt0c01.#########################################################################################"
# }
	`
	content = strings.Replace(content, "${version}", version.Current, 1)
	if _, err := providerFile.WriteString(content); err != nil {
		providerFile.Close()
		return err
	}

	return nil
}

func writeNestedProviderFile(targetFolder string, resFolder string) error {
	var err error
	var providerFile *os.File
	fileName := targetFolder + "/" + resFolder + "/" + "providers.tf"
	if providerFile, err = os.Create(fileName); err != nil {
		return err
	}
	content := `terraform {
  required_providers {
    dynatrace = {
    version = "${version}"
    source = "dynatrace-oss/dynatrace"
    }
  }
}
	`
	content = strings.Replace(content, "${version}", version.Current, 1)
	if _, err := providerFile.WriteString(content); err != nil {
		providerFile.Close()
		return err
	}

	return nil
}

func writeMainFile(file *os.File, resourceDataMap ResourceData, resName string, resFolder string, replacedIDs ReplacedIDs, dependsOn bool) error {
	var content string
	modules := map[string]bool{}
	if dependsOn && ResourceInfoMap[resName].HardcodedIds != nil {
		for _, hcName := range ResourceInfoMap[resName].HardcodedIds {
			for _, ids := range replacedIDs[resName] {
				if hcName == ids.RefDS && ids.RefRes != "" {
					if _, ok := modules[ids.RefRes]; !ok {
						modules[ids.RefRes] = true
					}
				}
			}
		}
		for module := range modules {
			if _, ok := resourceDataMap[module]; !ok {
				delete(modules, module)
			}
		}

		if len(modules) > 0 {
			content = `module "${resource_folder}" {
  source = "./${resource_folder}"
  depends_on = [${modules}]
  providers = {
    dynatrace = dynatrace.default
  }
}

`
			var modulesStr string
			for str := range modules {
				modulesStr = modulesStr + "module." + strings.TrimPrefix(str, "dynatrace_") + ", "

			}
			content = strings.Replace(content, "${modules}", strings.TrimSuffix(modulesStr, ", "), 1)
		}
	}
	if !dependsOn || len(modules) == 0 {
		content = `module "${resource_folder}" {
  source = "./${resource_folder}"
  providers = {
    dynatrace = dynatrace.default
  }
}

`
	}
	content = strings.Replace(content, "${resource_folder}", resFolder, 2)

	if _, err := file.WriteString(content); err != nil {
		file.Close()
		return err
	}

	return nil
}
