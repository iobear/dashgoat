/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import "strings"

// iSnewState checks if state is changing
// Only call this method if you have ss.mutex lock
func iSnewState(checkss ServiceState) (change string, new_service bool) {
	hostservice := strings.ToLower(checkss.Host) + strings.ToLower(checkss.Service)

	if _, ok := ss.serviceStateList[hostservice]; ok {

		current_status := ss.serviceStateList[hostservice].Status

		// no change
		if current_status == checkss.Status {
			return "", false
		}

		// change
		go reportStateChange(current_status, checkss)
		return checkss.Status, false
	}

	// change, new service
	go reportStateChange("", checkss)
	return checkss.Status, true
}

// ReportStateChange
func reportStateChange(fromstate string, reportss ServiceState) {
	logger.Info("statechange", "hostservice", reportss.Host+reportss.Service, "from", fromstate, "to", reportss.Status)

	if config.PagerdutyConfig.PdMode != "off" {
		//if len(reportss.From) == 1 { // Check if I'm the first to know
		pdClient.CompilePdEvent(fromstate, reportss)
		//}
	}
}
