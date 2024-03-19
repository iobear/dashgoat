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
	var post_service_state dg.ServiceState

	if err := c.Bind(&post_service_state); err != nil {
		return err
	}

	if !checkUpdatekey(post_service_state.UpdateKey) {
		return c.JSON(http.StatusUnauthorized, "Check your updatekey!")
	}

	post_service_state.UpdateKey = "valid"
	post_service_state = dg.FilterUpdate(post_service_state)
	post_service_state = runDependOn(post_service_state)

	// TODO
	// if err := validator.Validate(postService); err != nil {
	// 	return c.JSON(http.StatusBadRequest, ss.serviceStateList)
	// }

	host_service := ""

	if post_service_state.Host != "" && post_service_state.Service != "" {
		host_service = post_service_state.Host + post_service_state.Service
	} else {
		return c.JSON(http.StatusBadRequest, post_service_state)
	}

	if len(post_service_state.From) == 0 { //From can't be empty
		post_service_state.From = append(post_service_state.From, "127.0.0.1")
	}

	iSnewState(post_service_state) // Informs abount state change

	if _, ok := ss.serviceStateList[host_service]; ok {

		if post_service_state.Change == 0 {
			if post_service_state.Status != ss.serviceStateList[host_service].Status {
				post_service_state.Change = time.Now().Unix()
			} else {
				post_service_state.Change = ss.serviceStateList[host_service].Change
			}
		} else if post_service_state.Probe <= ss.serviceStateList[host_service].Probe { // Already reported
			return c.JSON(http.StatusAlreadyReported, "")
		}
	}

	if _, exists := ss.serviceStateList[host_service]; !exists {
		post_service_state.Change = time.Now().Unix()
	}

	ss.serviceStateList[host_service] = post_service_state

	go updateBuddy(post_service_state, "")

	result["id"] = host_service

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

// getStatusList HPO/MSO - simplified status list
func getStatusListMSO(c echo.Context) error {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	var tmp_service_state_MSO ServiceStateMSO
	service_state_MSO_list := make(map[string]ServiceStateMSO)

	for index, event := range ss.serviceStateList {
		tmp_service_state_MSO.Status = event.Status
		tmp_service_state_MSO.Message = "[" + event.Status + "] " + event.Service + " " + event.Host + "-" + event.Message

		service_state_MSO_list[index] = tmp_service_state_MSO
	}

	return c.JSON(http.StatusOK, service_state_MSO_list)

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
