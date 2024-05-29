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
			time_now := time.Now().Unix()
			for _, value := range check_slice {
				if time_now > value.TimeoutEpoch {
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
func updateEventLostProbe(host_service string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	tmp_struct := ss.serviceStateList[host_service]
	tmp_struct.Message = "Lost probe heartbeat"
	tmp_struct.Severity = config.ProbeTimeoutStatus
	tmp_struct.Status = config.ProbeTimeoutStatus

	change := iSnewState(tmp_struct) // Informs abount state change
	if change {
		tmp_struct.Change = time.Now().Unix()
	} else {
		tmp_struct.Change = ss.serviceStateList[host_service].Change
	}

	ss.serviceStateList[host_service] = tmp_struct
}

func findProbeInterval() int {
	interval_max := 60
	interval_min := 2
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
