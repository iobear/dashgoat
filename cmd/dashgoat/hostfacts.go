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
		MetricsHistory  bool
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

	host_facts.Items.DashName = config.DashName
	host_facts.Items.UpAtEpoch = time.Now().Unix()
	host_facts.Items.UpAt = time.Now()
	host_facts.Items.DashGoatVersion = Version
	host_facts.Items.GoVersion = runtime.Version()

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
	logger.Info("Welcome", "Hostnames found", host_facts.Items.Hostnames)

	if config.DisableMetrics && config.Prometheusurl == "" {
		host_facts.Items.MetricsHistory = false
		logger.Info("HostFacts", "MetricsHistory", "off")
	} else {
		host_facts.Items.MetricsHistory = true
		logger.Info("HostFacts", "MetricsHistory", "on")
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
