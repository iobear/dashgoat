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
)

type (
	Backlog struct {
		mutex        sync.RWMutex
		buddyBacklog map[string][]string
		StateDown    map[string]int64 //host, timestamp
	}
)

// setStateDown on buddy backlog
func setStateDown(host string, down bool) {
	var empty_timestamp int64
	var down_at int64

	if host == "" {
		logger.Info("setStateDown", "error", "no host")
		return
	}

	if down {
		down_at = time.Now().Unix()
	} else {
		down_at = empty_timestamp
	}

	backlog.mutex.Lock()
	defer backlog.mutex.Unlock()
	backlog.StateDown[host] = down_at
}

// getStateDown on buddy backlog
func getStateDown() map[string]int64 {
	backlog.mutex.RLock()
	defer backlog.mutex.RUnlock()

	copy := make(map[string]int64)
	for k, v := range backlog.StateDown {
		copy[k] = v
	}
	return copy
}

// setBacklog on buddy backlog
func setBacklog(host string, data []string) {
	if host == "" {
		logger.Info("setBacklog", "error", "no host")
		return
	}

	backlog.mutex.Lock()
	defer backlog.mutex.Unlock()
	backlog.buddyBacklog[host] = data
}

// Update Buddies with newly received msg
func updateBuddy(event ServiceState, delete string) {
	to_update := listBuddies()

	if len(to_update) < 1 {
		return //No buddy to tell
	}

	buddyDown := getStateDown()

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

// talkToBuddyApiDelete removes a service from Buddies
func talkToBuddyApiDelete(hostURL, serviceName string) error {
	url := hostURL + "/service/" + serviceName

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func talkToBuddyApi(event ServiceState, host Buddy, delete string) {
	logger.Info("talkToBuddyApi", "to Buddy", host.Name)
	my_hostnames := readHostFacts().Hostnames

	my_name := strings.ToLower(config.DashName)

	if host.NSconfig != "" {
		my_name = my_hostnames[0]
	}

	if my_name == host.Name {
		//I'm the sender of this message, I'm not telling my self
		logger.Info("talkToBuddyApi", "msg", "Should not happen, Not sending buddy msg to my self")
		return
	}

	if contains(event.From, my_name) {
		//I have already send this message once, don't repeat
		return
	}

	if delete != "" {
		err := talkToBuddyApiDelete(host.Url, delete)
		if err != nil {
			logger.Error("talkToBuddyApi", "delete", err)
		}
		return
	}

	event.From = append(event.From, my_name)
	event.UpdateKey = host.Key

	jsonMapAsStringFormat, err := json.Marshal(event)
	if err != nil {
		logger.Error("talkToBuddyApi", "JSON marshall", err)
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
		logger.Error("talkToBuddyApi", "Problems talking to "+host.Url, err)
		tellBuddyState(host.Name, false, event.Host)
		return
	}

	defer res.Body.Close()

}

func findBuddy(buddyConfig []Buddy) {

	initBuddyConf(buddyConfig)
	buddy_amount := len(listBuddies())

	if buddy_amount < 1 {
		setDashGoatReady(true)
		//logger.Info("findBuddy", "msg", "0 Buddy found")
	}

	firstRound := true

	buddy_txt := "Buddy"
	if buddy_amount > 1 {
		buddy_txt = "Buddies"
	}

	logger.Info("findBuddy", "buddy name", buddy_txt, "count", buddy_amount)

	wait_for := 3
	if config.CheckBuddyIntervalSec > 1 {
		wait_for = config.CheckBuddyIntervalSec
	}

	for {
		for _, bhost := range listBuddies() {
			if !contains(readHostFacts().Hostnames, bhost.Name) {
				healthy := askHealth(bhost)
				if healthy && firstRound {
					tellBuddyState(bhost.Name, false, "")
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

		time.Sleep(time.Duration(wait_for) * time.Second)
		firstRound = false
	}

}

// report back to UI, statusList
func tellBuddyState(host string, up bool, host_service string) {
	var empty_slice []string

	if _, ok := backlog.StateDown[host]; !ok {
		setStateDown(host, false)
	}

	if up {
		if getStateDown()[host] != 0 {
			tellServiceListAboutBuddy(host, up)
		}
		setStateDown(host, false)
		deliverBacklog(host, backlog.buddyBacklog[host])
		setBacklog(host, empty_slice) //empty backlog for host
	} else {
		if host_service != "" {
			backlog_tmp := append(backlog.buddyBacklog[host], host_service)
			setBacklog(host, backlog_tmp)
		}
		if backlog.StateDown[host] == 0 {
			tellServiceListAboutBuddy(host, up)
			setStateDown(host, true)
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
	var buddy_to_use Buddy

	for _, bhost := range listBuddies() {
		if bhost.Name == hostname {
			buddy_to_use = bhost
		}
	}

	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	for _, host_service := range backlog {
		if _, ok := ss.serviceStateList[host_service]; ok {
			talkToBuddyApi(ss.serviceStateList[host_service], buddy_to_use, "")
		} else {
			err := talkToBuddyApiDelete(buddy_to_use.Url, host_service)
			if err != nil {
				logger.Error("deliverBacklog", "delete", err)
			}
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

// UpdateFromBuddy Fetch statuslist from buddy
func UpdateFromBuddy(bhost Buddy) error {
	err := AskApiFullStatusList(bhost)
	if err != nil {
		return err
	}

	setDashGoatReady(true)
	return nil
}

func AskApiFullStatusList(bhost Buddy) error {

	resultMap := make(map[string]ServiceState)

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

	for host_service, status := range resultMap {
		if status.Service != "buddy" {
			ss.mutex.Lock()
			ss.serviceStateList[host_service] = status
			ss.mutex.Unlock()
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

func tellServiceListAboutBuddy(buddy_name string, up bool) {
	var result ServiceState

	if buddy_name == readHostFacts().DashName { //do not report my self
		return
	}

	time_now := time.Now().Unix()

	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	host_service := buddy_name + "buddy"

	result.Service = "buddy"
	result.Host = buddy_name

	if ss.serviceStateList[host_service].Status != result.Status {
		result.Change = time_now
	} else if ss.serviceStateList[host_service].Change == 0 {
		result.Change = time_now
	}

	result.Probe = time_now

	result.From = append(result.From, config.DashName)
	result.UpdateKey = "valid"
	if up {
		result.Status = "ok"
		result.Message = "buddy up"
		result.Severity = "info"
	} else {
		result.Status = strings.ToLower(config.BuddyDownStatus)
		result.Severity = result.Status
		result.Message = "buddy is down"
	}

	result, err := filterUpdate(result)
	if err != nil {
		logger.Error("tellServiceListAboutBuddy", "error", err)
	}

	iSnewState(result)

	ss.serviceStateList[host_service] = result

}
