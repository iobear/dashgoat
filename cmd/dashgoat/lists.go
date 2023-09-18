/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import "fmt"

// listServiceIDs - list all reported services
func listServiceIDs() []string {

	var resultList = []string{}
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	for key := range ss.serviceStateList {
		resultList = append(resultList, key)
	}

	return resultList
}

// uniqList - search struct item, return a list of unique values
func uniqList(item string) []string {

	var resultList = []string{}
	var resultStr = ""
	var listOfServices = listServiceIDs()
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	for _, id := range listOfServices {

		if item == "service" {
			resultStr = ss.serviceStateList[id].Service

		} else if item == "host" {
			resultStr = ss.serviceStateList[id].Host

		} else if item == "status" {
			resultStr = ss.serviceStateList[id].Status

		} else if item == "message" {
			resultStr = ss.serviceStateList[id].Message

		} else if item == "severity" {
			resultStr = ss.serviceStateList[id].Severity

		} else if item == "nextupdatesec" {
			resultInt := ss.serviceStateList[id].NextUpdateSec
			resultStr = fmt.Sprintf("%d", resultInt)

		} else if item == "probe" {
			int64Unix := ss.serviceStateList[id].Probe
			resultStr = fmt.Sprintf("%d", int64Unix)

		} else if item == "change" {
			int64Unix := ss.serviceStateList[id].Change
			resultStr = fmt.Sprintf("%d", int64Unix)
		}

		if item == "tags" {
			for _, value := range ss.serviceStateList[id].Tags {
				if indexOf(resultList, value) == -1 {
					resultList = append(resultList, value)
				}
			}

		} else {
			if indexOf(resultList, resultStr) == -1 {
				resultList = append(resultList, resultStr)
			}
		}
	}

	return resultList
}
