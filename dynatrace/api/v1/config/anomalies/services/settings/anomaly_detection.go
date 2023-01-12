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

package services

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/detection"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/failurerate"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/load"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/common/responsetime"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AnomalyDetection Dynatrace automatically detects service-related performance anomalies such as response time degradations and failure rate increases. Use these settings to configure detection sensitivity, set alert thresholds, or disable alerting for certain services.
type AnomalyDetection struct {
	LoadSpike               *load.SpikeDetection    `json:"loadSpike,omitempty"`     // The configuration of load spikes detection.
	ResponseTimeDegradation *responsetime.Detection `json:"responseTimeDegradation"` // Configuration of response time degradation detection.
	FailureRateIncrease     *failurerate.Detection  `json:"failureRateIncrease"`     // Configuration of failure rate increase detection.
	LoadDrop                *load.DropDetection     `json:"loadDrop,omitempty"`      // The configuration of load drops detection.
}

func (me *AnomalyDetection) Name() string {
	return "service_anomalies"
}

func (me *AnomalyDetection) getFailureRateIncrease() *failurerate.Detection {
	if me.FailureRateIncrease == nil {
		return &failurerate.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	return me.FailureRateIncrease
}

func (me *AnomalyDetection) getResponseTimeDegradation() *responsetime.Detection {
	if me.ResponseTimeDegradation == nil {
		return &responsetime.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	return me.ResponseTimeDegradation
}

func (me *AnomalyDetection) getLoadSpike() *load.SpikeDetection {
	if me.LoadSpike == nil {
		return &load.SpikeDetection{Enabled: false}
	}
	return me.LoadSpike
}

func (me *AnomalyDetection) getLoadDrop() *load.DropDetection {
	if me.LoadDrop == nil {
		return &load.DropDetection{Enabled: false}
	}
	return me.LoadDrop
}

func (me *AnomalyDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"load": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration of load spikes detection. Detecting load spikes will be disabled if this block is omitted.",
			Elem:        &schema.Resource{Schema: new(load.Detection).Schema()},
		},
		"response_times": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of response time degradation detection. Detecting response time degradation will be disabled if this block is omitted.",
			Elem:        &schema.Resource{Schema: new(responsetime.Detection).Schema()},
		},
		"failure_rates": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of failure rate increase detection. Detecting failure rate increases will be disabled if this block is omitted.",
			Elem:        &schema.Resource{Schema: new(failurerate.Detection).Schema()},
		},
		"load_drops": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The configuration of load drops detection. Detecting load drops will be disabled if this block is omitted.",
			Elem:        &schema.Resource{Schema: new(load.DropDetection).Schema()},
		},
	}
}

func (me *AnomalyDetection) MarshalHCL(decoder hcl.Decoder) (map[string]any, error) {
	result := map[string]any{}

	loadDetection := &load.Detection{
		Drops:  me.LoadDrop,
		Spikes: me.LoadSpike,
	}
	if !loadDetection.IsEmpty() {
		if marshalled, err := loadDetection.MarshalHCL(decoder); err == nil {
			result["load"] = []any{marshalled}
		} else {
			return nil, err
		}
	}
	if me.ResponseTimeDegradation != nil && string(me.ResponseTimeDegradation.DetectionMode) != string(detection.Modes.DontDetect) {
		if marshalled, err := me.ResponseTimeDegradation.MarshalHCL(decoder); err == nil {
			result["response_times"] = []any{marshalled}
		} else {
			return nil, err
		}
	}
	if me.FailureRateIncrease != nil && string(me.FailureRateIncrease.DetectionMode) != string(detection.Modes.DontDetect) {
		if marshalled, err := me.FailureRateIncrease.MarshalHCL(decoder); err == nil {
			result["failure_rates"] = []any{marshalled}
		} else {
			return nil, err
		}
	}
	if me.LoadDrop != nil && me.LoadDrop.Enabled {
		if marshalled, err := me.LoadDrop.MarshalHCL(decoder); err == nil {
			result["load_drops"] = []any{marshalled}
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (me *AnomalyDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("load.#"); ok {
		loadDetection := new(load.Detection)
		if err := loadDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "load", 0)); err != nil {
			return err
		}
		me.LoadSpike = loadDetection.Spikes
		me.LoadDrop = loadDetection.Drops
	} else {
		me.LoadSpike = &load.SpikeDetection{Enabled: false}
	}
	if _, ok := decoder.GetOk("response_times.#"); ok {
		me.ResponseTimeDegradation = new(responsetime.Detection)
		if err := me.ResponseTimeDegradation.UnmarshalHCL(hcl.NewDecoder(decoder, "response_times", 0)); err != nil {
			return err
		}
	} else {
		me.ResponseTimeDegradation = &responsetime.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	if _, ok := decoder.GetOk("failure_rates.#"); ok {
		me.FailureRateIncrease = new(failurerate.Detection)
		if err := me.FailureRateIncrease.UnmarshalHCL(hcl.NewDecoder(decoder, "failure_rates", 0)); err != nil {
			return err
		}
	} else {
		me.FailureRateIncrease = &failurerate.Detection{DetectionMode: detection.Modes.DontDetect}
	}
	if _, ok := decoder.GetOk("load_drops.#"); ok {
		me.LoadDrop = new(load.DropDetection)
		if err := me.LoadDrop.UnmarshalHCL(hcl.NewDecoder(decoder, "load_drops", 0)); err != nil {
			return err
		}
	} else {
		me.LoadDrop = &load.DropDetection{Enabled: false}
	}
	return nil
}

func (me *AnomalyDetection) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"loadSpike":               me.getLoadSpike(),
		"responseTimeDegradation": me.getResponseTimeDegradation(),
		"failureRateIncrease":     me.getFailureRateIncrease(),
		"loadDrop":                me.getLoadDrop(),
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *AnomalyDetection) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"loadSpike":               &me.LoadSpike,
		"responseTimeDegradation": &me.ResponseTimeDegradation,
		"failureRateIncrease":     &me.FailureRateIncrease,
		"loadDrop":                &me.LoadDrop,
	}); err != nil {
		return err
	}
	return nil
}
