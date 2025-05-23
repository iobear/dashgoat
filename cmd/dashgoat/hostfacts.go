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
	"runtime"
	"strings"
	"sync"
	"time"
)

type (
	HostFact struct {
		Hostnames       []string
		DashName        string
		Ready           bool
		UpAt            time.Time
		UpAtEpoch       int64
		DashGoatVersion string
		GoVersion       string
		BuildDate       string
		MetricsHistory  bool
		Prometheus      bool
		Commit          string
	}

	HostFacts struct {
		Items HostFact
		mutex sync.RWMutex
	}
)

var host_facts HostFacts

func generateHostFacts() {
	host_facts.mutex.Lock()
	defer host_facts.mutex.Unlock()

	host_facts.Items.DashName = strings.ToLower(config.DashName)
	host_facts.Items.UpAtEpoch = time.Now().Unix()
	host_facts.Items.UpAt = time.Now()
	host_facts.Items.DashGoatVersion = Version
	host_facts.Items.GoVersion = runtime.Version()
	host_facts.Items.BuildDate = BuildDate
	host_facts.Items.Commit = Commit

	hostname, _ := os.Hostname()
	IPhost := ""

	for _, ip := range getHostIPs() {
		IPhost = hostname + "_" + ip
		host_facts.Items.Hostnames = append(host_facts.Items.Hostnames, IPhost)
	}

	if len(host_facts.Items.Hostnames) == 0 {
		logger.Error("Cant find an IP address, check ignorePrefix config")
		os.Exit(1)
	}

	host_facts.Items.Hostnames = append(host_facts.Items.Hostnames, config.DashName)
	logger.Debug("Welcome", "Hostnames found", host_facts.Items.Hostnames)

	if config.DisableMetrics {
		host_facts.Items.MetricsHistory = false
		logger.Debug("HostFacts", "MetricsHistory", "off")
	} else {
		host_facts.Items.MetricsHistory = true
		logger.Debug("HostFacts", "MetricsHistory", "on")
	}

	if config.Prometheusurl != "" {
		host_facts.Items.Prometheus = true
	}

	if config.UrnKey == "" {
		logger.Warn("No UrnKey")
	}
}

func getHostIPs() []string {
	var result []string

	// Get the list of IP addresses associated with the host
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
	}

	ip_addr := ""
	for _, addr := range addrs {
		ip_addr = addr.String()
		if !ignorePrefix(ip_addr) { //no localhost addr
			result = append(result, strings.Split(addr.String(), "/")[0])
		}
	}

	return result
}

func ignorePrefix(ip_addr string) bool {
	ignore := config.IgnorePrefix

	if len(ignore) == 0 {
		ignore = []string{"8", "64", "128"}
	}

	for _, ignoreStr := range ignore {
		if strings.HasSuffix(ip_addr, ignoreStr) {
			return true
		}
	}
	return false
}

func readHostFacts() HostFact {
	host_facts.mutex.RLock()
	defer host_facts.mutex.RUnlock()
	return host_facts.Items
}

func isDashGoatReady() bool {
	host_facts.mutex.RLock()
	defer host_facts.mutex.RUnlock()
	return host_facts.Items.Ready
}

func setDashGoatReady(ready bool) {
	host_facts.mutex.Lock()
	defer host_facts.mutex.Unlock()
	host_facts.Items.Ready = ready
}
