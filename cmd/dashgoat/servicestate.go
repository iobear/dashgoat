/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import "sync"

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

	//ServiceState for MSO output
	ServiceStateMSO struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	// Services map
	Services struct {
		mutex            sync.RWMutex
		serviceStateList map[string]ServiceState
	}
)

// deleteService - ss unlock before use - removes service from serviceStateList
func deleteServiceState(id string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	serviceStateToDelete := ss.serviceStateList[id]

	delete(ss.serviceStateList, id)
	go updateBuddy(serviceStateToDelete, id) //Tell buddy to delete

	if !config.DisableMetrics {
		deleteServiceMetric(id) //Delete metric
	}

}
