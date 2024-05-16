/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type PagerDutyEvent struct {
	Payload struct {
		Summary   string `json:"summary"`
		Severity  string `json:"severity"`
		Source    string `json:"source"`
		Component string `json:"component"`
		Timestamp string `json:"timestamp"`
	} `json:"payload"`
	RoutingKey  string `json:"routing_key"`
	DedupKey    string `json:"dedup_key"`
	EventAction string `json:"event_action"`
	Client      string `json:"client"`
}

type PdConfig struct {
	URL           string         `yaml:"url"`
	Timeout       time.Duration  `yaml:"timeout"`
	PdMode        string         `yaml:"pagerdutymode"`
	TriggerLevel  string         `yaml:"triggerlevel"`
	PdServiceMaps []PdServiceMap `yaml:"pagerdutyservicemaps"`
}

type PdServiceMap struct {
	HostService string `yaml:"hostservice"`
	Tag         string `yaml:"tag"`
	EapiKey     string `yaml:"eapikey"`
}

type PdClient struct {
	config PdConfig
}

var pdClient = &PdClient{}

func validatePagerdutyConf() error {
	var result error

	// Is PagerDuty enabled?
	if config.PagerdutyConfig.PdMode == "off" {
		return result
	}

	// Default PdMode
	config.PagerdutyConfig.PdMode = "push"

	if len(config.PagerdutyConfig.PdServiceMaps) == 0 {
		logger.Debug("no pagerdutyservicemaps, setting pagerdutymode off")
		config.PagerdutyConfig.PdMode = "off"
		return result
	}

	// Default timeout value
	if config.PagerdutyConfig.Timeout < (3 * time.Second) {
		config.PagerdutyConfig.Timeout = 10 * time.Second
	}

	// Default PagerDuty Url US
	if config.PagerdutyConfig.URL == "" {
		config.PagerdutyConfig.URL = "https://events.pagerduty.com/v2/enqueue"
	}

	// Default PagerDuty trigger level
	if config.PagerdutyConfig.TriggerLevel == "" {
		config.PagerdutyConfig.TriggerLevel = "error"
	}

	for key, val := range config.PagerdutyConfig.PdServiceMaps {
		if val.EapiKey == "" {
			return fmt.Errorf("pagerDuty eapikey missing")
		}
		if val.HostService == "" && val.Tag == "" {
			config.PagerdutyConfig.PdServiceMaps[key].HostService = "0catchall0"
		}
	}

	return result
}

func initPagerDuty() {
	logger.Info("PagerDuty mode is " + config.PagerdutyConfig.PdMode)

	if config.PagerdutyConfig.PdMode == "off" {
		return
	}

	pdClient.config = PdConfig{
		URL:           config.PagerdutyConfig.URL,
		Timeout:       config.PagerdutyConfig.Timeout,
		PdServiceMaps: config.PagerdutyConfig.PdServiceMaps,
	}

}

func shouldPagerDutyTrigger(severity_to_check string) bool {
	trigger_level := indexOf(severitys[:], config.PagerdutyConfig.TriggerLevel)
	to_check := indexOf(severitys[:], severity_to_check)

	logger.Info("shouldPagerDutyTrigger", "severity_to_check", to_check, "trigger_level", trigger_level)
	return to_check >= trigger_level

}

func CompilePdEvent(fromstate string, reportss ServiceState) PagerDutyEvent {
	var pdevent PagerDutyEvent

	logger.Info("pdevent", "Severity", reportss.Severity)

	action := "resolve"

	if shouldPagerDutyTrigger(reportss.Status) {
		action = "trigger"
	}

	changetimeTime := time.Unix(reportss.Probe, 0).UTC()
	formattedTime := changetimeTime.Format("2006-01-02T15:04:05.000+0000")

	pdevent.Payload.Timestamp = formattedTime
	pdevent.Payload.Severity = reportss.Severity
	pdevent.Payload.Source = readHostFacts().DashName
	pdevent.Payload.Summary = reportss.Host + " " + reportss.Service + " " + reportss.Message
	pdevent.Payload.Component = reportss.Service
	pdevent.EventAction = action
	pdevent.DedupKey = reportss.Host + reportss.Service + "dashgoat"
	pdevent.RoutingKey = reportss.UpdateKey
	pdevent.Client = "dashgoat"

	return pdevent

}

func (c *PdClient) pagerDutyShipper(fromstate string, reportss ServiceState) {

	pdkey, _ := findKey(reportss)
	if pdkey == "" {
		logger.Debug("pagerDutyShipper", "mgs", "No match found for PagerDuty "+reportss.Host)
		return
	}

	reportss.UpdateKey = pdkey
	pdevent := CompilePdEvent(fromstate, reportss)

	retry := 3
	for retry >= 1 {

		err := pdClient.TellPagerDutyApi(pdevent)
		if err == nil {
			break
		} else {
			logger.Error("pagerDutyShipper", "error", err, "msg", "retrying..")
		}

		time.Sleep(3 * time.Second)
		retry--
	}

	logger.Error("pagerDutyShipper", "msg", "update was not send - giving up")

}

func findKey(dgss ServiceState) (pdkey string, pdmatch string) {

	var result string
	var match string

	for _, item := range config.PagerdutyConfig.PdServiceMaps {
		if item.HostService == dgss.Host+dgss.Service || item.HostService == "0catchall0" { // look for match or catch all token
			return item.EapiKey, item.HostService
		}

		for _, tag := range dgss.Tags {
			if tag == item.Tag {
				return item.EapiKey, item.Tag
			}
		}
	}

	return result, match
}

// TellPagerDuty updates PagerDuty via HTTP
func (c *PdClient) TellPagerDutyApi(pdevent PagerDutyEvent) error {

	client := &http.Client{
		Timeout: c.config.Timeout,
	}

	json_data, err := json.Marshal(pdevent)
	if err != nil {
		logger.Error("Error marshaling to JSON", err)
		return err
	}

	var payload = strings.NewReader(string(json_data))

	req, err := http.NewRequest("POST", c.config.URL, payload)
	if err != nil {
		logger.Error("TellPagerDuty", "POST failed", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("TellPagerDuty", "Do error", err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("TellPagerDuty", "ReadAll PagerDuty", err)
		return err
	}

	logger.Info("PagerDuty", "response", string(body), "statuscode", res.StatusCode)

	return nil
}
