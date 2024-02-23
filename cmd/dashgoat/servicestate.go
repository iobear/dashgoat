/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

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
