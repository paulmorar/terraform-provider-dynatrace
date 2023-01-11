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

package service

import (
	"fmt"
	api "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/services"
	settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/metrics/calculated/service/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"net/url"
	"strings"
	"time"
)

const SchemaID = "v1:config:calculated-metrics:service"

func Service(credentials *api.Credentials) api.CRUDService[*settings.CalculatedServiceMetric] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(id string, v *settings.CalculatedServiceMetric) error {
	return me.client.Get(fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), 200).Finish(v)
}

func (me *service) List() (api.Stubs, error) {
	var err error

	req := me.client.Get("/api/config/v1/calculatedMetrics/service", 200)
	var stubList api.StubList
	if err = req.Finish(&stubList); err != nil {
		return nil, err
	}
	return stubList.Values, nil
}

func (me *service) Validate(v *settings.CalculatedServiceMetric) error {
	var err error

	client := me.client

	retry := true
	maxAttempts := 64
	attempts := 0

	for retry {
		attempts = attempts + 1
		req := client.Post("/api/config/v1/calculatedMetrics/service/validator", v, 204)
		if err = req.Finish(); err != nil {
			if !strings.Contains(err.Error(), "Metric definition must specify a known request attribute") {
				return err
			}
			// log.Println(".... request attribute is not fully known yet to cluster - retrying")
			if attempts < maxAttempts {
				time.Sleep(500 * time.Millisecond)
			} else {
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}

func (me *service) Create(v *settings.CalculatedServiceMetric) (*api.Stub, error) {
	var err error

	client := me.client

	retry := true
	maxAttempts := 64
	attempts := 0
	var stub api.Stub

	for retry {
		attempts = attempts + 1
		req := client.Post("/api/config/v1/calculatedMetrics/service", v, 201)
		if err = req.Finish(&stub); err != nil {
			if !strings.Contains(err.Error(), "Metric definition must specify a known request attribute") {
				return nil, err
			}
			// log.Println(".... request attribute is not fully known yet to cluster - retrying")
			if attempts < maxAttempts {
				time.Sleep(500 * time.Millisecond)
			} else {
				return nil, err
			}
		} else {
			return &stub, nil
		}
	}
	return &stub, nil
}

func (me *service) Update(id string, v *settings.CalculatedServiceMetric) error {
	if err := me.client.Put(fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), v, 204).Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	for {
		if err := me.client.Delete(fmt.Sprintf("/api/config/v1/calculatedMetrics/service/%s", url.PathEscape(id)), 204).Finish(); err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("Service metric with %s not found", id)) {
				return nil
			}
			if !strings.Contains(err.Error(), "Could not delete configuration") {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}
