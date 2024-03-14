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
		go reportStateChange(hostservice, ss.serviceStateList[hostservice].Status, checkss.Status)
		return checkss.Status, false
	}

	go reportStateChange(hostservice, "", checkss.Status)
	return checkss.Status, true
}

// reportStateChange, updates dependencies
func reportStateChange(hostservice string, from string, to string) {
	logger.Info("statechange", "hostservice", hostservice, "from", from, "to", to)

}
