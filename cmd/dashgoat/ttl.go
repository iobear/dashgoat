/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"time"
)

func readStatusList() map[string]ServiceState {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	statusListCopy := make(map[string]ServiceState, len(ss.serviceStateList))
	for key, serviceState := range ss.serviceStateList {
		statusListCopy[key] = serviceState
	}

	return statusListCopy
}

func updateServiceState(key string, serviceState ServiceState) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		ss.mutex.Unlock()
	}()

	ss.mutex.Lock()
	ss.serviceStateList[key] = serviceState
}

func ttlHousekeeping() {
	ticker := time.NewTicker(time.Second * 3) // adjust the interval as needed
	defer ticker.Stop()

	for range ticker.C {
		statusList := readStatusList()
		currentUnixTimestamp := time.Now().Unix()
		for key, serviceState := range statusList {
			if serviceState.Ttl <= 0 {
				continue
			}

			if serviceState.Probe+int64(serviceState.Ttl) <= currentUnixTimestamp {
				if config.TtlBehavior == "remove" {
					deleteServiceState(serviceState.Host + serviceState.Service)
				} else {
					serviceState = promoteStatus(serviceState, currentUnixTimestamp)
					updateServiceState(key, serviceState)
				}
			}

			if serviceState.Status == "ok" && (serviceState.Probe+int64(config.TtlOkDelete) <= currentUnixTimestamp) {
				deleteServiceState(key)
			}
		}
	}
}

func deleteServiceState(key string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		ss.mutex.Unlock()
	}()

	ss.mutex.Lock()
	delete(ss.serviceStateList, key)
}

func promoteStatus(serviceState ServiceState, currentUnixTimestamp int64) ServiceState {

	if config.TtlBehavior == "promoteonce" { // PromoteOnce
		serviceState.Ttl = 0
	}

	// Default PromoteOneStep
	statusHierarchy := []string{"critical", "error", "warning", "info", "ok"}
	for i, status := range statusHierarchy {
		if serviceState.Status == status && i < len(statusHierarchy)-1 {
			serviceState.Status = statusHierarchy[i+1]
			serviceState.Probe = currentUnixTimestamp
			break
		}
	}

	if config.TtlBehavior == "promotetook" { // PromoteToOk
		serviceState.Status = "ok"
	}

	return serviceState
}
