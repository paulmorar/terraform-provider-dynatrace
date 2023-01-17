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

package notifications

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WebHookConfig Configuration of the custom WebHook notification.
type WebHookConfig struct {
	BaseNotificationConfig
	AcceptAnyCertificate     bool          `json:"acceptAnyCertificate"`               // Accept any, including self-signed and invalid, SSL certificate (`true`) or only trusted (`false`) certificates.
	Headers                  []*HTTPHeader `json:"headers,omitempty"`                  // A list of the additional HTTP headers.
	Payload                  string        `json:"payload"`                            // The content of the notification message.  You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemDetailsJSON}`: All problem event details, including root cause, as a JSON object.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas.
	URL                      string        `json:"url"`                                // The URL of the WebHook endpoint.
	NotifyEventMergesEnabled *bool         `json:"notifyEventMergesEnabled,omitempty"` // Call webhook if new events merge into existing problems
}

func (me *WebHookConfig) GetType() Type {
	return Types.Webhook
}

func (me *WebHookConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration",
			Required:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "The configuration is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"notify_event_merges": {
			Type:        schema.TypeBool,
			Description: "Call webhook if new events merge into existing problems",
			Optional:    true,
		},
		"alerting_profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
		"accept_any_certificate": {
			Type:        schema.TypeBool,
			Description: "Accept any, including self-signed and invalid, SSL certificate (`true`) or only trusted (`false`) certificates",
			Required:    true,
		},
		"header": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A list of the additional HTTP headers",
			Elem:        &schema.Resource{Schema: new(HTTPHeader).Schema()},
		},
		"payload": {
			Type:        schema.TypeString,
			Description: "The content of the notification message.  You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemDetailsJSON}`: All problem event details, including root cause, as a JSON object.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the WebHook endpoint",
			Required:    true,
		},
	}
}

func (me *WebHookConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	properties["name"] = me.Name
	properties["active"] = me.Active
	properties["alerting_profile"] = me.AlertingProfile
	properties["accept_any_certificate"] = me.AcceptAnyCertificate
	if me.NotifyEventMergesEnabled != nil {
		properties["notify_event_merges"] = *me.NotifyEventMergesEnabled
	}
	if len(me.Headers) > 0 {
		if err := properties.EncodeSlice("header", me.Headers); err != nil {
			return err
		}
	}
	properties["payload"] = me.Payload
	properties["url"] = me.URL

	return nil
}

func (me *WebHookConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "active")
		delete(me.Unknowns, "alertingProfile")
		delete(me.Unknowns, "acceptAnyCertificate")
		delete(me.Unknowns, "notifyEventMergesEnabled")
		delete(me.Unknowns, "header")
		delete(me.Unknowns, "payload")
		delete(me.Unknowns, "url")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("active"); ok {
		me.Active = value.(bool)
	}
	if value, ok := decoder.GetOk("alerting_profile"); ok {
		me.AlertingProfile = value.(string)
	}
	if value, ok := decoder.GetOk("notify_event_merges"); ok {
		me.NotifyEventMergesEnabled = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("accept_any_certificate"); ok {
		me.AcceptAnyCertificate = value.(bool)
	}
	if err := decoder.DecodeSlice("header", &me.Headers); err != nil {
		return err
	}
	if value, ok := decoder.GetOk("payload"); ok {
		me.Payload = value.(string)
	}
	if value, ok := decoder.GetOk("url"); ok {
		me.URL = value.(string)
	}
	return nil
}

func (me *WebHookConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":                       me.ID,
		"name":                     me.Name,
		"type":                     me.GetType(),
		"alertingProfile":          me.AlertingProfile,
		"notifyEventMergesEnabled": me.NotifyEventMergesEnabled,
		"active":                   me.Active,
		"acceptAnyCertificate":     me.AcceptAnyCertificate,
		"headers":                  me.Headers,
		"payload":                  me.Payload,
		"url":                      me.URL,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *WebHookConfig) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"id":                       &me.ID,
		"name":                     &me.Name,
		"type":                     &me.Type,
		"active":                   &me.Active,
		"alertingProfile":          &me.AlertingProfile,
		"notifyEventMergesEnabled": &me.NotifyEventMergesEnabled,
		"acceptAnyCertificate":     &me.AcceptAnyCertificate,
		"headers":                  &me.Headers,
		"payload":                  &me.Payload,
		"url":                      &me.URL,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
