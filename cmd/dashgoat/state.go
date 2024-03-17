/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	dg "github.com/iobear/dashgoat/common"
)

// iSnewState checks if state is changing
// Only call this method if you have ss.mutex lock
func iSnewState(checkss dg.ServiceState) (state string, new bool) {
	hostservice := checkss.Host + checkss.Service

	if _, ok := ss.serviceStateList[hostservice]; ok {

		if ss.serviceStateList[hostservice].Status == checkss.Status {
			return "", false
		}
		go reportStateChange(ss.serviceStateList[hostservice].Status, checkss)
		return checkss.Status, false
	}

	go reportStateChange("", checkss)
	return checkss.Status, true
}

// ReportStateChange
func reportStateChange(fromstate string, reportss dg.ServiceState) {
	logger.Info("statechange", "hostservice", reportss.Host+reportss.Service, "from", fromstate, "to", reportss.Status)

	if config.PagerdutyConfig.PdMode != "off" {
		if len(reportss.From) == 1 { // Check if I'm the first to know
			pdClient.CompilePdEvent(fromstate, reportss)
		}
	}
}
