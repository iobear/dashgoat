/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
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
		Ack           string   `json:"ack"`
		Ttl           int      `json:"ttl"`
		DependOn      string   `json:"dependon"`
		UpdateKey     string
	}
)

// ValidateUpdate
func filterUpdate(ss ServiceState) (ServiceState, error) {

	if ss.Host == "" || ss.Service == "" {
		return ss, fmt.Errorf("missing Host or Service value")
	}

	this_is_now := time.Now().Unix()
	if ss.Probe == 0 {
		ss.Probe = this_is_now
	}

	if ss.Change == 0 {
		ss.Change = this_is_now
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

	//Replace non alpha numeric characters with -
	ss.Host = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.' {
			return r
		}
		return '-' // Replace all other characters with '-'
	}, ss.Host)

	ss.Service = strings.Replace(ss.Service, " ", "-", -1)
	ss.Service = strings.ToLower(ss.Service)

	return ss, nil
}
