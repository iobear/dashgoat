/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

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
		} else if postService.Probe <= ss.serviceStateList[strID].Probe { //Already reported
			return c.JSON(http.StatusAlreadyReported, "")
		}
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
	currentTime := int(time.Now().Unix())
	tmpServiceStateList := make(map[string]ServiceState)

	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	if len(ss.serviceStateList) == 0 {
		return c.JSON(http.StatusNoContent, "")
	}

	//Filter events with ttl defined
	for index, event := range ss.serviceStateList {

		if event.Ttl == 0 {
			tmpServiceStateList[index] = event
		} else if event.Ttl+int(event.Probe) < currentTime && event.Status != "ok" {
			event.Status = "ok"
			tmpServiceStateList[index] = event
		} else if event.Ttl+int(event.Probe) > currentTime {
			tmpServiceStateList[index] = event
		}

	}

	return c.JSON(http.StatusOK, tmpServiceStateList)

}

// getStatusList HPO/MSO
func getStatusListMSO(c echo.Context) error {
	currentTime := int(time.Now().Unix())

	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	var tmpServiceStateMSO ServiceStateMSO
	serviceStateMSOlist := make(map[string]ServiceStateMSO)

	for index, event := range ss.serviceStateList {
		tmpServiceStateMSO.Status = event.Status
		tmpServiceStateMSO.Message = "[" + event.Status + "] " + event.Service + " " + event.Host + "-" + event.Message

		//Filter events with ttl defined
		if event.Ttl == 0 {
			serviceStateMSOlist[index] = tmpServiceStateMSO
		} else if event.Ttl+int(event.Probe) > currentTime && event.Status != "ok" {
			tmpServiceStateMSO.Status = "ok"
			serviceStateMSOlist[index] = tmpServiceStateMSO
		} else if event.Ttl+int(event.Probe) > currentTime {
			serviceStateMSOlist[index] = tmpServiceStateMSO
		}
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

// deleteService - removes service from serviceStateList
func deleteService(c echo.Context) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	id := c.Param("id")
	id = strings.Replace(id, " ", "-", -1)

	_, mapContainsKey := ss.serviceStateList[id]

	if mapContainsKey {
		serviceStateToDelete := ss.serviceStateList[id]
		delete(ss.serviceStateList, id)
		go updateBuddy(serviceStateToDelete, id)
		return c.NoContent(http.StatusNoContent)
	}

	return c.NoContent(http.StatusNotFound)
}

func health(c echo.Context) error {
	if !dashGoatReady() {
		return c.NoContent(http.StatusServiceUnavailable)
	}

	return c.JSON(http.StatusOK, readHostFacts())
}

// validate and enrich input from POST
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

		if ss.Status == "ok" {
			ss.Severity = "info"

		} else {
			ss.Severity = "error"
		}
	}

	ss.Host = strings.Replace(ss.Host, " ", "", -1)
	ss.Service = strings.Replace(ss.Service, " ", "-", -1)

	return true
}
