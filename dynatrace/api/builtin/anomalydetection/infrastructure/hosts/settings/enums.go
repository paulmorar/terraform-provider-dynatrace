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

package hosts

type ConnectionLostDetectionSensitivity string

var ConnectionLostDetectionSensitivities = struct {
	AlertOnGracefulShutdown     ConnectionLostDetectionSensitivity
	DontAlertOnGracefulShutdown ConnectionLostDetectionSensitivity
}{
	ConnectionLostDetectionSensitivity("ALERT_ON_GRACEFUL_SHUTDOWN"),
	ConnectionLostDetectionSensitivity("DONT_ALERT_ON_GRACEFUL_SHUTDOWN"),
}

type DetectionMode string

var DetectionModes = struct {
	Auto   DetectionMode
	Custom DetectionMode
}{
	DetectionMode("Auto"),
	DetectionMode("Custom"),
}
