package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/validator.v2"
)

type (
	serviceState struct {
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
		mutex         sync.RWMutex
	}

	appHealth struct {
		APIVersion string `json:"apiversion"`
	}
)

var (
	serviceStateList = map[string]*serviceState{}
	fromPost         *serviceState
	serviceList      = []string{}
	tagList          = []string{}
	appHealthResult  *appHealth
)

//updateStatus - service update
func updateStatus(c echo.Context) error {
	var result = map[string]string{}
	fromPost = &serviceState{}

	fromPost.mutex.Lock()
	defer fromPost.mutex.Unlock()

	if err := c.Bind(fromPost); err != nil {
		return err
	}

	if validateUpdate() == false {
		return c.JSON(http.StatusUnauthorized, "Check your updatekey!")
	}

	if err := validator.Validate(fromPost); err != nil {
		return c.JSON(http.StatusBadRequest, fromPost)
	}

	strID := ""

	if fromPost.Host != "" && fromPost.Service != "" {
		strID = fromPost.Host + fromPost.Service
	} else {
		return c.JSON(http.StatusBadRequest, fromPost)
	}

	if serviceStateList[strID] != nil {

		if fromPost.Status != serviceStateList[strID].Status {
			fromPost.Change = time.Now().Unix()

		} else {
			fromPost.Change = serviceStateList[strID].Change
		}

	}

	serviceStateList[strID] = fromPost
	result["id"] = strID

	return c.JSON(http.StatusOK, result)
}

//getStatus - get status of service with service id
func getStatus(c echo.Context) error {

	id := c.Param("id")
	return c.JSON(http.StatusOK, serviceStateList[id])

}

//getStatusList - return all data
func getStatusList(c echo.Context) error {

	return c.JSON(http.StatusOK, serviceStateList)

}

//serviceFilter - list services with item value of..
func serviceFilter(c echo.Context) error {
	//placeholder func
	return c.JSON(http.StatusOK, serviceStateList)
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

	id := c.Param("id")
	delete(serviceStateList, id)

	return c.NoContent(http.StatusNoContent)
}

func health(c echo.Context) error {
	appHealthResult = &appHealth{}

	if err := c.Bind(appHealthResult); err != nil {
		return err
	}

	appHealthResult.APIVersion = "1.0.5"

	return c.JSON(http.StatusOK, appHealthResult)
}
