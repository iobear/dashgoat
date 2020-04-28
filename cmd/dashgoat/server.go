package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/validator.v2"
)

type (
	//ServiceState struct for validating post input
	ServiceState struct {
		Service       string   `json:"service" validate:"min=1,max=40"`
		Host          string   `json:"host" validate:"min=1,max=40"`
		Status        string   `json:"status" validate:"min=1,max=10,regexp=^[a-z]*$"`
		Message       string   `json:"message" validate:"max=255"`
		Severity      string   `json:"severity" validate:"max=10"`
		NextUpdateSec int      `json:"nextupdatesec" validate:"max=605000"`
		Tags          []string `json:"tags" validate:"max=20"`
		Seen          int64    `json:"seen"`
		Change        int64    `json:"change"`
		UpdateKey     string
	}

	// Services contai
	Services struct {
		mutex            sync.RWMutex
		serviceStateList map[string]ServiceState
	}

	//AppHealth holds health data
	AppHealth struct {
		APIVersion string `json:"apiversion"`
	}
)

var (
	serviceList     = []string{}
	tagList         = []string{}
	appHealthResult *AppHealth
)

//updateStatus - service update
func updateStatus(c echo.Context) error {

	ss.mutex.Lock()

	defer ss.mutex.Unlock()

	var result = map[string]string{}

	var postService ServiceState

	if err := c.Bind(&postService); err != nil {
		c.Logger().Error(ss.serviceStateList)
		return err
	}

	if postService.validateUpdate() == false {
		c.Logger().Error(ss.serviceStateList)
		return c.JSON(http.StatusUnauthorized, "Check your updatekey!")
	}

	if err := validator.Validate(postService); err != nil {
		c.Logger().Error(ss.serviceStateList)
		return c.JSON(http.StatusBadRequest, ss.serviceStateList)
	}

	strID := ""

	if postService.Host != "" && postService.Service != "" {
		strID = postService.Host + postService.Service
	} else {
		strID = postService.Host + postService.Service
		c.Logger().Error(strID)
		c.Logger().Error(postService)
		return c.JSON(http.StatusBadRequest, postService)
	}

	if _, ok := ss.serviceStateList[strID]; ok {

		if postService.Status != ss.serviceStateList[strID].Status {
			postService.Change = time.Now().Unix()
			c.Logger().Info("change")

		} else {
			postService.Change = ss.serviceStateList[strID].Change
			c.Logger().Info("no change")

		}
	}

	c.Logger().Info(ss.serviceStateList)

	ss.serviceStateList[strID] = postService
	c.Logger().Info(ss.serviceStateList)

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
	delete(ss.serviceStateList, id)

	return c.NoContent(http.StatusNoContent)
}

func health(c echo.Context) error {
	appHealthResult = &AppHealth{}

	appHealthResult.APIVersion = "1.0.11"

	return c.JSON(http.StatusOK, appHealthResult)
}
