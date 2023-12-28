/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
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

// updateStatus - service update
func updateStatus(c echo.Context) error {

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	var result = map[string]string{}
	var postService ServiceState

	if err := c.Bind(&postService); err != nil {
		return err
	}

	if !postService.validateUpdate() {
		return c.JSON(http.StatusUnauthorized, "Check your updatekey!")
	}

	// if err := validator.Validate(postService); err != nil {
	// 	return c.JSON(http.StatusBadRequest, ss.serviceStateList)
	// }

	strID := ""

	if postService.Host != "" && postService.Service != "" {
		strID = postService.Host + postService.Service
	} else {
		return c.JSON(http.StatusBadRequest, postService)
	}

	if _, ok := ss.serviceStateList[strID]; ok {

		if postService.Change == 0 {
			if postService.Status != ss.serviceStateList[strID].Status {
				postService.Change = time.Now().Unix()
			} else {
				postService.Change = ss.serviceStateList[strID].Change
			}
		} else if postService.Probe <= ss.serviceStateList[strID].Probe { // Already reported
			return c.JSON(http.StatusAlreadyReported, "")
		}
	}

	if _, exists := ss.serviceStateList[strID]; !exists {
		postService.Change = time.Now().Unix()
	}

	ss.serviceStateList[strID] = postService

	go updateBuddy(postService, "")

	result["id"] = strID

	return c.JSON(http.StatusOK, result)
}

// getStatus - get status of service with service id
func getStatus(c echo.Context) error {

	id := c.Param("id")
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	return c.JSON(http.StatusOK, ss.serviceStateList[id])

}

// getStatusList - return all data
func getStatusList(c echo.Context) error {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	if len(ss.serviceStateList) == 0 {
		return c.JSON(http.StatusNoContent, "")
	}

	return c.JSON(http.StatusOK, ss.serviceStateList)

}

// getStatusList HPO/MSO
func getStatusListMSO(c echo.Context) error {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	var tmpServiceStateMSO ServiceStateMSO
	serviceStateMSOlist := make(map[string]ServiceStateMSO)

	for index, event := range ss.serviceStateList {
		tmpServiceStateMSO.Status = event.Status
		tmpServiceStateMSO.Message = "[" + event.Status + "] " + event.Service + " " + event.Host + "-" + event.Message

		serviceStateMSOlist[index] = tmpServiceStateMSO
	}

	return c.JSON(http.StatusOK, serviceStateMSOlist)

}

// serviceFilter - list services with item value of..
func serviceFilter(c echo.Context) error {
	//placeholder func
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	return c.JSON(http.StatusOK, ss.serviceStateList)
}

// getUniq - list unique values of service items
func getUniq(c echo.Context) error {

	item := c.Param("serviceitem")
	var result = []string{}

	if item == "id" {
		result = listServiceIDs()

	} else {
		result = uniqList(item)

	}

	return c.JSON(http.StatusOK, result)
}

// deleteServiceHandler - removes service from serviceStateList
func deleteServiceHandler(c echo.Context) error {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	id := c.Param("id")
	id = strings.Replace(id, " ", "-", -1)

	_, mapContainsKey := ss.serviceStateList[id]

	if mapContainsKey {
		go deleteServiceState(id)
		return c.NoContent(http.StatusNoContent)
	}

	return c.NoContent(http.StatusNotFound)
}

// health of dashGoat app
func health(c echo.Context) error {
	if !dashGoatReady() {
		return c.NoContent(http.StatusServiceUnavailable)
	}

	return c.JSON(http.StatusOK, readHostFacts())
}

// validateUpdate, validate and enrich input from POST
func (ss *ServiceState) validateUpdate() bool {

	if ss.UpdateKey == config.UpdateKey {
		ss.UpdateKey = "valid"
	} else {
		return false
	}

	if ss.Probe == 0 {
		ss.Probe = time.Now().Unix()
	}

	msglength := len(ss.Message)
	if msglength > 254 {
		ss.Message = string(ss.Message[0:254])
	}

	severitylen := len(ss.Severity)
	if severitylen > 10 {
		ss.Severity = string(ss.Severity[0:10])
	}
	ss.Severity = strings.ToLower(ss.Severity)

	statuslen := len(ss.Status)
	if statuslen > 10 {
		ss.Status = string(ss.Status[0:10])
	}
	ss.Status = strings.ToLower(ss.Status)

	if ss.Severity == "" {

		if ss.Status == "ok" || ss.Status == "info" {
			ss.Severity = "info"
		} else {
			ss.Severity = "error"
		}
	}

	ss.Host = strings.Replace(ss.Host, " ", "", -1)
	ss.Host = strings.ToLower(ss.Host)
	ss.Service = strings.Replace(ss.Service, " ", "-", -1)
	ss.Service = strings.ToLower(ss.Service)

	if ss.Status != "ok" && ss.DependOn != "" {
		msg := isDependOnError(ss.DependOn)
		if msg == "down" {
			ss.Severity = "info"
			ss.Status = "info"
			ss.Message = "( " + ss.DependOn + " down ) " + ss.Message
		} else if msg != "" {
			ss.Severity = "info"
			ss.Status = "info"
			ss.Message = "( " + ss.DependOn + " ) " + msg + ss.Message
		}
	}

	return true
}
