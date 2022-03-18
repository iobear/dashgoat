package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type (
	// Backlog map
	Backlog struct {
		mutex        sync.RWMutex
		buddyBacklog map[string][]string
		StateDown    map[string]int64
	}
)

func updateBuddy(event ServiceState, delete string) {
	bb.mutex.RLock()
	buddyDown := bb.StateDown
	bb.mutex.RUnlock()

	for _, bhost := range config.BuddyHosts {
		if bhost.Name != dashName && !bhost.Ignore && !contains(event.From, bhost.Name) {
			if buddyDown[bhost.Name] > 0 { //node down, move to backlog
				tellBuddyState(bhost.Name, false, event.Host+event.Service)
			} else {
				talkToBuddy(event, bhost, delete)
			}
		}
	}

}

func talkToBuddyDelete(hosturl string, delete string) {
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

func talkToBuddy(event ServiceState, host Buddy, delete string) {

	if delete != "" {
		talkToBuddyDelete(host.Url, delete)
		return
	}

	if contains(event.From, config.DashName) {
		//I have already sendt this message once, dont repeat
		return
	}

	event.From = append(event.From, config.DashName)
	event.UpdateKey = host.Key

	jsonMapAsStringFormat, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
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
		tellBuddyState(host.Name, false, event.Host)
		return
	}

	defer res.Body.Close()

}

func findBuddy() {
	firstRound := true

	if !config.EnableBuddy {
		dashgoat_ready = true
		fmt.Println("Buddy not enabled")
		return
	}

	waitfor := 10
	if config.CheckBuddyIntervalSec > 1 {
		waitfor = config.CheckBuddyIntervalSec
	}

	for {
		for _, bhost := range config.BuddyHosts {
			if bhost.Name != config.DashName && !bhost.Ignore {
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

		if !dashgoat_ready {
			dashgoat_ready = true
		}
		time.Sleep(time.Duration(waitfor) * time.Second)
		firstRound = false
	}

}

func tellBuddyState(host string, up bool, servicehost string) {

	now := time.Now()
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	if _, ok := bb.StateDown[host]; !ok {
		bb.StateDown[host] = 0
	}

	if up {
		if bb.StateDown[host] != 0 {
			tellDashgoatAboutBuddy(host, up)
		}
		bb.StateDown[host] = 0
		emptyBacklog(host, bb.buddyBacklog[host])
		bb.buddyBacklog[host] = nil
	} else {
		if servicehost != "" {
			bb.buddyBacklog[host] = append(bb.buddyBacklog[host], servicehost)
		}
		if bb.StateDown[host] == 0 {
			tellDashgoatAboutBuddy(host, up)
			bb.StateDown[host] = now.Unix()
		}
	}
}

func emptyBacklog(hostname string, backlog []string) {
	var hostToUse Buddy

	for _, bhost := range config.BuddyHosts {
		if bhost.Name == hostname {
			hostToUse = bhost
		}
	}

	ss.mutex.RLock()
	defer ss.mutex.RUnlock()

	for _, hoststate := range backlog {
		if _, ok := ss.serviceStateList[hoststate]; ok {
			talkToBuddy(ss.serviceStateList[hoststate], hostToUse, "")
		} else {
			talkToBuddyDelete(hostToUse.Url, hoststate)
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

func UpdateFromBuddy(bhost Buddy) error {
	err := AskFullStatusList(bhost)
	if err != nil {
		return err
	}

	dashgoat_ready = true
	return nil
}

func AskFullStatusList(bhost Buddy) error {

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
	body, err := ioutil.ReadAll(res.Body)
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

func tellDashgoatAboutBuddy(buddyName string, up bool) {
	var result ServiceState

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
		result.Status = strings.ToLower(config.BuddyDown)
		result.Message = "My buddy is down"
	}

	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.serviceStateList[serviceName] = result

}
