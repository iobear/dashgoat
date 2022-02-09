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

	//AppHealth holds health data
	AppHealth struct {
		DashAPI  string `json:"dashapi"`
		DashName string `json:"dashname"`
	}
)

var appHealthResult *AppHealth

//updateStatus - service update
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
		}
	}

	ss.serviceStateList[strID] = postService

	result["id"] = strID

	return c.JSON(http.StatusOK, result)
}

//getStatus - get status of service with service id
func getStatus(c echo.Context) error {

	id := c.Param("id")
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	return c.JSON(http.StatusOK, ss.serviceStateList[id])

}

//getStatusList - return all data
func getStatusList(c echo.Context) error {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	return c.JSON(http.StatusOK, ss.serviceStateList)

}

//getStatusList HPO/MSO
func getStatusListMSO(c echo.Context) error {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	var tmpServiceStateMSO ServiceStateMSO
	serviceStateMSOlist := make(map[string]ServiceStateMSO)

	for index, event := range ss.serviceStateList {
		tmpServiceStateMSO.Status = event.Status
		tmpServiceStateMSO.Message = event.Status + " - " + event.Service + " " + event.Host + "-" + event.Message

		serviceStateMSOlist[index] = tmpServiceStateMSO
	}

	return c.JSON(http.StatusOK, serviceStateMSOlist)

}

//serviceFilter - list services with item value of..
func serviceFilter(c echo.Context) error {
	//placeholder func
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	return c.JSON(http.StatusOK, ss.serviceStateList)
}

//getUniq - list unique values of service items
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

//deleteService - removes service from serviceStateList
func deleteService(c echo.Context) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	id := c.Param("id")
	id = strings.Replace(id, " ", "-", -1)

	delete(ss.serviceStateList, id)

	return c.NoContent(http.StatusNoContent)
}

func health(c echo.Context) error {
	appHealthResult = &AppHealth{}

	appHealthResult.DashAPI = "1.1.10"
	appHealthResult.DashName = dashName

	return c.JSON(http.StatusOK, appHealthResult)
}
