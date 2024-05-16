/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// heartBeat update
func heartBeat(c echo.Context) error {

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	var post_service_state ServiceState

	urnkey := c.Param("urnkey")
	if urnkey == "" {
		return c.JSON(http.StatusUnauthorized, "Missing urnkey")
	}

	if !checkUrnKey(urnkey) {
		return c.JSON(http.StatusUnauthorized, "Check your urnkey")
	}
	post_service_state.UpdateKey = "valid"

	host := c.Param("host")
	if host == "" {
		return c.JSON(http.StatusBadRequest, "Missing host")
	}

	post_service_state.Host = host
	post_service_state.Service = "heartbeat"

	nextupdatesec := c.Param("nextupdatesec")
	if nextupdatesec == "" {
		return c.JSON(http.StatusBadRequest, "Missing nextupdatesec")
	}

	sec_number, err := strconv.ParseUint(nextupdatesec, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "nextupdatesec not number")
	}
	post_service_state.NextUpdateSec = int(sec_number)

	host_service := post_service_state.Host + post_service_state.Service

	tags := c.Param("tags")
	if tags != "" {
		post_service_state.Tags = parseTags(tags)
	}
	post_service_state.Status = "ok"

	post_service_state.From = append(post_service_state.From, "heartbeat")

	post_service_state, err = filterUpdate(post_service_state)
	if err != nil {
		logger.Error("filterUpdate", "msg", err)
		return err
	}

	this_is_now := time.Now().Unix()
	change := iSnewState(post_service_state) // Informs abount state change
	if change {
		post_service_state.Status = "ok"
		post_service_state.Change = this_is_now
	} else {
		post_service_state.Change = ss.serviceStateList[host_service].Change
	}

	post_service_state = runDependOn(post_service_state)

	ss.serviceStateList[host_service] = post_service_state //commit change to service state

	go updateBuddy(post_service_state, "")

	return c.JSON(http.StatusOK, host_service)
}

// updateStatus - service update
func updateStatus(c echo.Context) error {

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	var result = map[string]string{}
	var post_service_state ServiceState

	if err := c.Bind(&post_service_state); err != nil {
		return err
	}

	if !checkUpdatekey(post_service_state.UpdateKey) {
		return c.JSON(http.StatusUnauthorized, "Check your updatekey")
	}

	post_service_state.UpdateKey = "valid"
	post_service_state, err := filterUpdate(post_service_state)
	if err != nil {
		return c.JSON(http.StatusBadRequest, post_service_state)
	}
	post_service_state = runDependOn(post_service_state)

	host_service := ""

	if post_service_state.Host != "" && post_service_state.Service != "" {
		host_service = post_service_state.Host + post_service_state.Service
	} else {
		return c.JSON(http.StatusBadRequest, post_service_state)
	}

	if len(post_service_state.From) == 0 { //From can't be empty
		post_service_state.From = append(post_service_state.From, "127.0.0.1")
	}

	post_service_state, err = filterUpdate(post_service_state)
	if err != nil {
		logger.Error("filterUpdate", "msg", err)
		return err
	}

	change := iSnewState(post_service_state) // Informs abount state change
	if change {
		post_service_state.Change = time.Now().Unix()
	} else {
		post_service_state.Change = ss.serviceStateList[host_service].Change
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
		return c.NoContent(http.StatusNoContent)
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

func checkUrnKey(key string) bool {
	return key == config.UrnKey
}

func parseTags(tags string) []string {

	slice_result := strings.Split(tags, ",")

	return slice_result
}
