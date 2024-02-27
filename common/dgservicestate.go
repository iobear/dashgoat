/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package common

import (
	"strings"
	"time"
)

type (
	//ServiceState struct for validating post input
	ServiceState struct {
		Service       string   `json:"service" validate:"min=1,max=100"`
		Host          string   `json:"host" validate:"min=1,max=100"`
		Status        string   `json:"status" validate:"min=1,max=10,regexp=^[a-z]*$"`
		Message       string   `json:"message" validate:"max=255"`
		Severity      string   `json:"severity" validate:"max=10"`
		NextUpdateSec int      `json:"nextupdatesec" validate:"max=605000"`
		Tags          []string `json:"tags" validate:"max=20"`
		Probe         int64    `json:"probe"`
		Change        int64    `json:"change"`
		From          []string `json:"from"`
		Ttl           int      `json:"ttl"`
		DependOn      string   `json:"dependon"`
		UpdateKey     string
	}
)

// ValidateUpdate
func FilterUpdate(ss ServiceState) ServiceState {

	if ss.Probe == 0 {
		ss.Probe = time.Now().Unix()
	}

	msglength := len(ss.Message)
	if msglength > 254 {
		ss.Message = string(ss.Message[0:254])
	}

	severitylen := len(ss.Severity)
	if severitylen > 10 {
		ss.Severity = string(ss.Severity[0:10])
	}
	ss.Severity = strings.ToLower(ss.Severity)

	statuslen := len(ss.Status)
	if statuslen > 10 {
		ss.Status = string(ss.Status[0:10])
	}
	ss.Status = strings.ToLower(ss.Status)

	if ss.Severity == "" {

		if ss.Status == "ok" || ss.Status == "info" {
			ss.Severity = "info"
		} else {
			ss.Severity = "error"
		}
	}

	ss.Host = strings.Replace(ss.Host, " ", "", -1)
	ss.Host = strings.ToLower(ss.Host)
	ss.Service = strings.Replace(ss.Service, " ", "-", -1)
	ss.Service = strings.ToLower(ss.Service)

	return ss
}
