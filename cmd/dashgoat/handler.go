/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"net/http"
	"strings"
	"time"

	dg "github.com/iobear/dashgoat/common"
	"github.com/labstack/echo/v4"
)

// updateStatus - service update
func updateStatus(c echo.Context) error {

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	var result = map[string]string{}
	var postService dg.ServiceState

	if err := c.Bind(&postService); err != nil {
		return err
	}

	if !checkUpdatekey(postService.UpdateKey) {
		return c.JSON(http.StatusUnauthorized, "Check your updatekey!")
	}

	postService.UpdateKey = "valid"
	postService = dg.FilterUpdate(postService)
	postService = runDependOn(postService)

	// TODO
	// if err := validator.Validate(postService); err != nil {
	// 	return c.JSON(http.StatusBadRequest, ss.serviceStateList)
	// }

	strID := ""

	if postService.Host != "" && postService.Service != "" {
		strID = postService.Host + postService.Service
	} else {
		return c.JSON(http.StatusBadRequest, postService)
	}

	//	iSnewState(postService)

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
	if !isDashGoatReady() {
		return c.NoContent(http.StatusServiceUnavailable)
	}

	return c.JSON(http.StatusOK, readHostFacts())
}

func checkUpdatekey(key string) bool {

	return key == config.UpdateKey
}
