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

	dg "github.com/iobear/dashgoat/common"
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

	if len(config.PagerdutyConfig.PdServiceMaps) == 0 {
		logger.Info("no pagerdutyservicemaps, setting pagerdutymode off")
		config.PagerdutyConfig.PdMode = "off"
		return result
	}

	// Default PdMode
	pdClient.config.PdMode = "push"

	// Default timeout value
	if config.PagerdutyConfig.Timeout == 0 {
		config.PagerdutyConfig.Timeout = 10
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
	trigger_level := indexOf(severitys, config.PagerdutyConfig.TriggerLevel)
	to_check := indexOf(severitys, severity_to_check)

	return to_check >= trigger_level

}

func (c *PdClient) CompilePdEvent(fromstate string, dgss dg.ServiceState) {
	var pdevent PagerDutyEvent

	logger.Info("pdevent", "Severity", dgss.Severity)

	pdkey, _ := findKey(dgss)
	if pdkey == "" {
		logger.Info("No key found " + dgss.Host)
		return
	}

	action := "resolve"

	if shouldPagerDutyTrigger(dgss.Status) {
		action = "trigger"
	}

	changetimeTime := time.Unix(dgss.Probe, 0).UTC()
	formattedTime := changetimeTime.Format("2006-01-02T15:04:05.000+0000")

	pdevent.Payload.Timestamp = formattedTime
	pdevent.Payload.Severity = dgss.Severity
	pdevent.Payload.Source = readHostFacts().DashName
	pdevent.Payload.Summary = dgss.Host + " " + dgss.Service + " " + dgss.Message
	pdevent.Payload.Component = dgss.Service
	pdevent.EventAction = action
	pdevent.DedupKey = dgss.Host + dgss.Service + "dashgoat"
	pdevent.RoutingKey = pdkey
	pdevent.Client = "dashgoat"

	err := pdClient.TellPagerDuty(pdevent)
	if err != nil {
		logger.Error("Error sending to PagerDuty:", err)

	}
}

func findKey(dgss dg.ServiceState) (pdkey string, pdmatch string) {

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
func (c *PdClient) TellPagerDuty(pdevent PagerDutyEvent) error {

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
		logger.Error("PagerDuty POST failed", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.Error("PagerDuty error client", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("Failed reading PagerDuty response", err)
		return err
	}

	logger.Info("PagerDuty", "response", string(body), "statuscode", res.StatusCode)

	return nil
}
