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
	"sync"
	"time"

	dg "github.com/iobear/dashgoat/common"
)

type (
	Backlog struct {
		mutex        sync.RWMutex
		buddyBacklog map[string][]string
		StateDown    map[string]int64
	}
)

// Update Buddies with newly recieved msg
func updateBuddy(event dg.ServiceState, delete string) {
	to_update := listBuddies()

	if len(to_update) < 1 {
		return //No buddy to tell
	}

	backlog.mutex.RLock()
	buddyDown := backlog.StateDown
	backlog.mutex.RUnlock()

	for _, bhost := range to_update {
		if !contains(event.From, bhost.Name) {
			if buddyDown[bhost.Name] > 0 { //node down, move to backlog
				tellBuddyState(bhost.Name, false, event.Host+event.Service)
			} else {
				talkToBuddyApi(event, bhost, delete)
			}
		}
	}

}

// Remove service from Buddies
func talkToBuddyApiDelete(hosturl string, delete string) {
	url := hosturl + "/service/" + delete

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

}

func talkToBuddyApi(event dg.ServiceState, host Buddy, delete string) {
	my_hostnames := readHostFacts().Hostnames

	my_name := config.DashName

	if host.NSconfig != "" {
		my_name = my_hostnames[0]
	}

	if delete != "" {
		talkToBuddyApiDelete(host.Url, delete)
		return
	}

	if contains(event.From, my_name) {
		//I have already sendt this message once, dont repeat
		return
	}

	event.From = append(event.From, my_name)
	event.UpdateKey = host.Key

	jsonMapAsStringFormat, err := json.Marshal(event)
	if err != nil {
		logger.Error("talkToBuddyApi json marshall", err)
		return
	}

	payload := strings.NewReader(string(jsonMapAsStringFormat))

	req, err := http.NewRequest("POST", host.Url+"/update", payload)
	if err != nil {
		tellBuddyState(host.Name, false, event.Host)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("talkToBuddyApi problems talking to "+host.Url, err)
		tellBuddyState(host.Name, false, event.Host)
		return
	}

	defer res.Body.Close()

}

func findBuddy(buddyConfig []Buddy) {

	initBuddyConf(buddyConfig)
	buddyAmount := len(buddyRunningConfig.Buddies)

	if buddyAmount < 1 {
		setDashGoatReady(true)
		logger.Info("Buddy not found")
	}

	firstRound := true

	buddyTxt := "Buddy"
	if buddyAmount > 1 {
		buddyTxt = "Buddies"
	}

	buddyWelcome := fmt.Sprintf("%d %s ", buddyAmount, buddyTxt)
	logger.Info(buddyWelcome)

	waitfor := 10
	if config.CheckBuddyIntervalSec > 1 {
		waitfor = config.CheckBuddyIntervalSec
	}

	for {
		for _, bhost := range listBuddies() {
			if !contains(readHostFacts().Hostnames, bhost.Name) {
				tellBuddyState(bhost.Name, false, "")
				healthy := askHealth(bhost)
				if healthy && firstRound {
					firstRound = false
					err := UpdateFromBuddy(bhost)
					if err != nil {
						firstRound = true
					}
				}
				tellBuddyState(bhost.Name, healthy, "")
			}
		}

		if !isDashGoatReady() {
			setDashGoatReady(true)
		}

		time.Sleep(time.Duration(waitfor) * time.Second)
		firstRound = false
	}

}

// report back to UI, stausList
func tellBuddyState(host string, up bool, servicehost string) {

	now := time.Now()
	backlog.mutex.Lock()
	defer backlog.mutex.Unlock()

	if _, ok := backlog.StateDown[host]; !ok {
		backlog.StateDown[host] = 0
	}

	if up {
		if backlog.StateDown[host] != 0 {
			tellServiceListAboutBuddy(host, up)
		}
		backlog.StateDown[host] = 0
		deliverBacklog(host, backlog.buddyBacklog[host])
		backlog.buddyBacklog[host] = nil
	} else {
		if servicehost != "" {
			backlog.buddyBacklog[host] = append(backlog.buddyBacklog[host], servicehost)
		}
		if backlog.StateDown[host] == 0 {
			tellServiceListAboutBuddy(host, up)
			backlog.StateDown[host] = now.Unix()
		}
	}
}

func deleteBuddyBacklog(valid_buddies []Buddy) {
	backlog.mutex.Lock()
	defer backlog.mutex.Unlock()
	var buddy_names []string

	for _, ok_buddy := range valid_buddies {
		buddy_names = append(buddy_names, ok_buddy.Name)
	}

	for name := range backlog.buddyBacklog {
		if !contains(buddy_names, name) {
			delete(backlog.buddyBacklog, name)
		}
	}

	for name := range backlog.StateDown {
		if !contains(buddy_names, name) {
			delete(backlog.StateDown, name)
		}
	}
}

// Update buddy with messages that could not be delivered
func deliverBacklog(hostname string, backlog []string) {
	var hostToUse Buddy

	for _, bhost := range listBuddies() {
		if bhost.Name == hostname {
			hostToUse = bhost
		}
	}

	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	for _, hoststate := range backlog {
		if _, ok := ss.serviceStateList[hoststate]; ok {
			talkToBuddyApi(ss.serviceStateList[hoststate], hostToUse, "")
		} else {
			talkToBuddyApiDelete(hostToUse.Url, hoststate)
		}
	}
}

func askHealth(bhost Buddy) bool {
	healthy := true

	req, err := http.NewRequest("GET", bhost.Url+"/health", nil)
	if err != nil {
		healthy = false
		return healthy
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		healthy = false
		return healthy
	}

	if res.StatusCode != 200 {
		healthy = false
	}

	defer res.Body.Close()
	return healthy
}

// Fetch statuslist from buddy
func UpdateFromBuddy(bhost Buddy) error {
	err := AskApiFullStatusList(bhost)
	if err != nil {
		return err
	}

	setDashGoatReady(true)
	return nil
}

func AskApiFullStatusList(bhost Buddy) error {

	resultMap := make(map[string]dg.ServiceState)

	req, err := http.NewRequest("GET", bhost.Url+"/status/list", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dashGoat")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 204 {
		err = fmt.Errorf("no content")
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &resultMap)
	if err != nil {
		return err
	}

	for servicehost, status := range resultMap {
		if status.Service != "Buddy" {
			ss.serviceStateList[servicehost] = status
		}
	}

	return nil
}

func serviceListDeleteBuddy(ok_buddies []Buddy) {
	var buddy_names []string

	for _, buddy := range ok_buddies {
		buddy_names = append(buddy_names, buddy.Name)
	}

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	for key, state := range ss.serviceStateList {
		if strings.ToLower(state.Service) == "buddy" {
			if !contains(buddy_names, state.Host) {
				delete(ss.serviceStateList, key)
			}
		}
	}

}

func tellServiceListAboutBuddy(buddyName string, up bool) {
	var result dg.ServiceState

	serviceName := buddyName + "buddy"

	result.Service = "Buddy"
	result.Host = buddyName
	result.Severity = "error"
	result.Probe = 0
	result.Change = time.Now().Unix()
	result.From = append(result.From, config.DashName)
	result.UpdateKey = "valid"
	if up {
		result.Status = "ok"
		result.Message = ""
	} else {
		result.Status = strings.ToLower(config.BuddyDownStatusMsg)
		result.Message = "My buddy is down"
	}

	iSnewState(result)

	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.serviceStateList[serviceName] = result

}
