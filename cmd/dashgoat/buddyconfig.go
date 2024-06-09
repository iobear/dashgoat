/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	BuddyConfig struct {
		mutex   sync.RWMutex
		Buddies []Buddy
	}

	Buddy struct {
		Name     string `yaml:"name"`
		Url      string `yaml:"url"`
		Key      string `yaml:"key"`
		NSconfig string
		LastSeen int64
		Down     bool
	}
)

var buddyRunningConfig BuddyConfig

func initBuddyConf(rawConfig []Buddy) {

	for _, buddy := range rawConfig {
		buddy.Name = strings.ToLower(buddy.Name)
		if !buddy.Down && buddy.Name != config.DashName {
			logger.Info("initBuddyConf", "Adding buddy", buddy.Name)
			addBuddy(buddy)
		}
	}

	if buddy_nsconfig == "" {
		buddy_nsconfig = config.BuddyNsConfig
	}

	if buddy_nsconfig != "" {
		findDNSBuddies()
	}

	//look for headless service IPs
	if isK8s() && buddy_nsconfig == "" {
		buddy_nsconfig = "dashgoat-headless-svc"
		findDNSBuddies()
	}

	if buddy_nsconfig != "" {
		go loopFindDNSBuddies()
	}
}

// Add buddy to running config
func addBuddy(host Buddy) {

	idx := findBuddyIndex(host.Name)
	buddyRunningConfig.mutex.Lock()
	defer buddyRunningConfig.mutex.Unlock()

	if idx >= 0 {
		buddyRunningConfig.Buddies[idx].LastSeen = host.LastSeen
	} else {
		buddyRunningConfig.Buddies = append(buddyRunningConfig.Buddies, host)
	}
}

func findBuddyIndex(host string) int {
	buddyRunningConfig.mutex.RLock()
	defer buddyRunningConfig.mutex.RUnlock()

	for idx, buddy := range listBuddies() {
		if buddy.Name == host {
			return idx
		}
	}

	return -1
}

// List all buddy running config
func listBuddies() []Buddy {
	buddyRunningConfig.mutex.RLock()
	defer buddyRunningConfig.mutex.RUnlock()
	return buddyRunningConfig.Buddies
}

// Del buddy running config
func delBuddy(del_buddy Buddy) {
	var result []Buddy

	buddyRunningConfig.mutex.Lock()
	defer buddyRunningConfig.mutex.Unlock()

	for _, buddy := range buddyRunningConfig.Buddies {
		if buddy.Name != del_buddy.Name {
			result = append(result, buddy)
		}
	}

	deleteBuddyBacklog(result)
	serviceListDeleteBuddy(result)

	buddyRunningConfig.Buddies = nil
	buddyRunningConfig.Buddies = result

}

// look for IPs in A-records, update list of IPs
func ipLookup() (error, bool) {
	var err error
	var new_ip bool

	ips, err := net.LookupIP(buddy_nsconfig)
	if err != nil {
		fmt.Println(err)
		return err, new_ip
	}

	if len(ips) < 1 {
		return err, new_ip
	}
	new_ip = true
	adjustBuddyConfig(ips)
	return err, new_ip
}

func adjustBuddyConfig(ips []net.IP) {
	time_now := time.Now().Unix()

	//add buddies found via IP slice
	for _, ip := range ips {
		buddy_items, my_self := compileBuddyConfig(ip.String(), time_now)
		if !my_self {
			addBuddy(buddy_items)
		}
	}

	//Remove old nsconfig-Buddies not found last lookup
	for _, buddy := range listBuddies() {
		if buddy.NSconfig == buddy_nsconfig && buddy.LastSeen != time_now {
			delBuddy(buddy)
		}
	}
}

func compileBuddyConfig(ip string, time_now int64) (Buddy, bool) {
	var result Buddy
	var my_self bool

	hostname, _ := os.Hostname()
	result.Name = hostname + "_" + ip
	if contains(readHostFacts().Hostnames, result.Name) {
		my_self = true
	}

	result.LastSeen = time_now
	result.NSconfig = buddy_nsconfig
	result.Key = config.UpdateKey
	result.Url = "http://" + ip + ":" + strings.Split(config.IPport, ":")[1] //TODO - needs some work

	if my_self { //If the host is added earlier
		delBuddy(result)
	}

	return result, my_self
}

func isK8s() bool {
	var result bool

	k8s_service := os.Getenv("KUBERNETES_SERVICE_HOST")
	if len(k8s_service) > 4 {
		result = true
	}

	return result
}

func translateDNSresult() string {
	var result string

	result_len := len(listBuddies())
	if result_len == 1 {
		result = "There is only one "
	} else if result_len == 0 {
		result = "There is no "
	} else if result_len > 1 {
		result = "There is more than one "
	}

	return result + "buddy found"
}

func lookForIPs() bool {
	var new_ip bool

	err, new_ip := ipLookup()

	if err != nil {
		logger.Error("lookForIPs", "error", err)
	}

	return new_ip
}

func findDNSBuddies() {

	new_ip := lookForIPs()
	if new_ip {
		fmt.Println("DNS result - " + translateDNSresult())
	}
}

func loopFindDNSBuddies() {

	for {
		logger.Info("loopFindDNSBuddies", "loop", "sleep 20 Sec")
		time.Sleep(20 * time.Second)
		findDNSBuddies()
	}
}
