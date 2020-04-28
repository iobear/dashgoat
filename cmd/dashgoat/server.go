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
		mutex         sync.RWMutex
	}

	//AppHealth holds health data
	AppHealth struct {
		APIVersion string `json:"apiversion"`
	}
)

var (
	serviceStateList = map[string]*ServiceState{}
	serviceList      = []string{}
	tagList          = []string{}
	appHealthResult  *AppHealth
)

//newPost - init new post
func newPost() *ServiceState {
	return &ServiceState{}
}

//updateStatus - service update
func (ss *ServiceState) updateStatus(c echo.Context) error {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	var result = map[string]string{}

	if err := c.Bind(ss); err != nil {
		c.Logger().Error(ss)
		return err
	}

	if ss.validateUpdate() == false {
		c.Logger().Error(ss)
		return c.JSON(http.StatusUnauthorized, "Check your updatekey!")
	}

	if err := validator.Validate(ss); err != nil {
		c.Logger().Error(ss)
		return c.JSON(http.StatusBadRequest, ss)
	}

	strID := ""

	if ss.Host != "" && ss.Service != "" {
		strID = ss.Host + ss.Service
	} else {
		strID = ss.Host + ss.Service
		c.Logger().Error(strID)
		c.Logger().Error(ss)
		return c.JSON(http.StatusBadRequest, ss)
	}

	if serviceStateList[strID] != nil {

		if ss.Status != serviceStateList[strID].Status {
			ss.Change = time.Now().Unix()
			c.Logger().Info("change")

		} else {
			ss.Change = serviceStateList[strID].Change
			c.Logger().Info("no change")

		}
	}

	c.Logger().Info(serviceStateList)

	serviceStateList[strID] = ss
	c.Logger().Info(serviceStateList)

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
	appHealthResult = &AppHealth{}

	appHealthResult.APIVersion = "1.0.11"

	return c.JSON(http.StatusOK, appHealthResult)
}
