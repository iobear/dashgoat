/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"strings"

	dg "github.com/iobear/dashgoat/common"
)

// Looks for service or tag that a second service depends on
// beware - Only call this method if you have ss.mutex lock
func isDependOnError(search_host_key string) string {

	if config.DisableDependOn {
		return ""
	}

	count_ok := 0
	count_error := 0

	search := strings.ToLower(strings.TrimSpace(search_host_key))

	for statekey := range ss.serviceStateList {
		if ss.serviceStateList[statekey].Host == search || contains(ss.serviceStateList[statekey].Tags, search) {
			if ss.serviceStateList[statekey].Status == "error" || ss.serviceStateList[statekey].Status == "critical" {
				count_error++
			} else {
				count_ok++
			}
		}
	}

	if count_error == 0 {
		return ""
	}
	if count_error > 0 && count_ok == 0 {
		return "down"
	}
	if count_error > 0 && count_ok > 0 {
		return fmt.Sprintf("partly down %d/%d ", count_error, count_ok+count_error)
	}

	return ""
}

func runDependOn(ss dg.ServiceState) dg.ServiceState {

	if ss.Status != "ok" && ss.DependOn != "" {
		msg := isDependOnError(ss.DependOn)
		if msg == "down" {
			ss.Severity = "warning"
			ss.Status = "warning"
			ss.Message = "( " + ss.DependOn + " down ) " + ss.Message
		} else if msg != "" {
			ss.Severity = "info"
			ss.Status = "info"
			ss.Message = "( " + ss.DependOn + " ) " + msg + ss.Message
		}
	}

	return ss
}
