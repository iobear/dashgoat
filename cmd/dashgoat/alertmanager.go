/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	// HookMessage is a JSON body from the alertmanager send HTTP POST requests.
	// Learn more to see: #https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
	HookMessage struct {
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
		Status            string            `json:"status"`
		Receiver          string            `json:"receiver"`
		GroupLabels       map[string]string `json:"groupLabels"`
		CommonLabels      map[string]string `json:"commonLabels"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
		Alerts            []Alert           `json:"alerts"`
	}

	// Alert is a single alert.
	Alert struct {
		Status       string            `json:"status"`
		Labels       map[string]string `json:"labels"`
		Annotations  map[string]string `json:"annotations"`
		StartsAt     string            `json:"startsAt,omitempty"`
		EndsAt       string            `json:"EndsAt,omitempty"`
		GeneratorURL string            `json:"generatorURL"`
	}
)

func fromAlertmanager(c echo.Context) error {

	dec := json.NewDecoder(c.Request().Body)
	defer c.Request().Body.Close()

	var message HookMessage
	if err := dec.Decode(&message); err != nil {
		logger.Error("updateAlertmanager", "error decoding message", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	urnkey := c.Param("urnkey")
	if urnkey == "" {
		return c.JSON(http.StatusUnauthorized, "Missing urnkey")
	}

	if !checkUrnKey(urnkey) {
		return c.JSON(http.StatusUnauthorized, "Check your urnkey")
	}

	//printDebugAlertManager(message)
	err := parseAlertmanagerHookMessage(message)
	if err != nil {
		logger.Error("fromAlertmanager", "error", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, "")
}

func parseAlertmanagerHookMessage(message HookMessage) error {

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	var post_service_state ServiceState

	post_service_state.UpdateKey = "valid"
	post_service_state.Severity = strings.ToLower(message.CommonLabels["severity"])
	if message.CommonLabels["severity"] == "" {
		err := fmt.Errorf("missing CommonLabels[Severity]")
		logger.Error("parseAlertmanagerHookMessage", "CommonLabels", err)
		post_service_state.Severity = "error"
	}

	status := post_service_state.Severity
	if message.Status == "resolved" {
		status = "ok"
		post_service_state.Ttl = 10800 //cleanup after 3 hours
	}

	post_service_state.Status = status
	post_service_state.Host = message.CommonLabels["prometheus_cluster"]
	if post_service_state.Host == "" {
		post_service_state.Host = message.CommonLabels["cluster"]
	}
	if post_service_state.Host == "" {
		post_service_state.Host = message.CommonLabels["prometheus"]
	}
	if post_service_state.Host == "" {
		err := fmt.Errorf("missing CommonLabels['prometheus_cluster'], CommonLabels['cluster'] or CommonLabels['prometheus']")
		return err
	}
	post_service_state.Host = strings.ToLower(post_service_state.Host)

	post_service_state.Service = message.CommonLabels["namespace"]
	if message.CommonLabels["namespace"] == "" {
		logger.Info("parseAlertmanagerHookMessage", "missing", "CommonLabels['namespace']")
	}
	post_service_state.From = append(post_service_state.From, post_service_state.Host)

	this_is_now := time.Now().Unix()
	for _, alert := range message.Alerts { //parsing the alerts

		post_service_state, err := parseAlertmanagerAlert(alert, post_service_state)
		if err != nil {
			return err
		}

		host_service := post_service_state.Host + post_service_state.Service

		change := iSnewState(post_service_state) // Informs abount state change
		if change {
			post_service_state.Change = this_is_now
		} else {
			post_service_state.Change = ss.serviceStateList[host_service].Change
			logger.Debug("No change recorded")
		}

		post_service_state, err = filterUpdate(post_service_state)
		if err != nil {
			return err
		}
		post_service_state = runDependOn(post_service_state)

		ss.serviceStateList[host_service] = post_service_state

		go updateBuddy(post_service_state, "")
	}

	return nil
}

func parseAlertmanagerAlert(alert Alert, service_state ServiceState) (ServiceState, error) {

	service_state.Message = alert.Labels["alertname"] + " - " + alert.Annotations["summary"] + " - " + alert.Labels["container"]
	if alert.Annotations["summary"] == "" {
		err := fmt.Errorf("missing alert Annotations['summary']")
		logger.Error("parseAlertmanagerAlert", "ServiceState.message", err)
	}

	//Have found namespace in CommonLabels
	if service_state.Service == "" {
		service_state.Service = alert.Labels["namespace"]
	}
	if service_state.Service == "" {
		service_state.Service = alert.Labels["container"]
	}
	if service_state.Service == "" {
		service_state.Service = alert.Labels["alertname"]
	}
	if service_state.Service == "" {
		logger.Error("parseAlertmanagerAlert", "service", "Cant find namespace or container", "alert object", alert.Labels)
	} else {
		service_state.Service = strings.ToLower(service_state.Service)
	}
	return service_state, nil

}

func printDebugAlertManager(message HookMessage) {

	logger.Info("updateAlertmanager", "[Body]", &message)

	fmt.Println("-start Hook-")

	fmt.Println(message.Version)
	fmt.Println(message.GroupKey)
	fmt.Println(message.Status)
	fmt.Println(message.Receiver)
	fmt.Println(message.GroupLabels)
	fmt.Println(message.CommonLabels)
	fmt.Println(message.CommonAnnotations)
	fmt.Println(message.ExternalURL)

	fmt.Println("-end Hook-")

	for _, alert := range message.Alerts {
		fmt.Println("-start alert-")
		fmt.Println("-Labels-")
		fmt.Println(alert.Labels)
		fmt.Println("-Annotations-")
		fmt.Println(alert.Annotations)
		fmt.Println(" -end- ")
	}

}
