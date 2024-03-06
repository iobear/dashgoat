/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"strconv"
	"time"
)

type (
	TimerEpoch struct {
		HostService  string
		TimeoutEpoch int64
	}
)

// lostProbeTimer - handles lost probes
func lostProbeTimer() {

	for {
		interval := findProbeInterval()
		count, check_slice := listProbeTimeout()

		if count > 0 {
			timeNow := time.Now().Unix()
			for _, value := range check_slice {
				if timeNow > value.TimeoutEpoch {
					updateEventLostProbe(value.HostService)
				}
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}

}

// listProbeTimeout - returns list of services with a timout defined
func listProbeTimeout() (int, []TimerEpoch) {
	var result []TimerEpoch

	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	count := len(ss.serviceStateList)

	if count == 0 {
		return count, result
	}

	var tmpTimeoutEpoch TimerEpoch
	for index, event := range ss.serviceStateList {
		if event.NextUpdateSec > 0 {
			tmpTimeoutEpoch.HostService = index
			tmpTimeoutEpoch.TimeoutEpoch = event.Probe + int64(event.NextUpdateSec)
			result = append(result, tmpTimeoutEpoch)
		}
	}

	return count, result
}

// updateEventLostProbe - sets Message , Severity and Status accordingly
func updateEventLostProbe(hostService string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	tmpStruct := ss.serviceStateList[hostService]
	tmpStruct.Message = "Lost probe heartbeat"
	tmpStruct.Severity = "error"
	tmpStruct.Status = "error"

	ss.serviceStateList[hostService] = tmpStruct
}

func findProbeInterval() int {
	interval_max := 60
	interval_min := 4
	result := interval_min

	result_slice := uniqList("nextupdatesec")

	if len(result_slice) < 1 {
		return result
	}

	for _, value := range result_slice {
		i, _ := strconv.ParseInt(value, 10, 64)
		if i != 0 {
			if int(i) < interval_max {
				result = int(i)
			}
		}
	}

	if result < interval_min {
		result = interval_min
	}

	if result > interval_max {
		result = interval_max
	}

	return result
}
