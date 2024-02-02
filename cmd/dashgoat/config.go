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

	"gopkg.in/yaml.v2"
)

type (
	Configer struct {
		DashName              string   `yaml:"dashName"`
		IPport                string   `yaml:"ipport"`
		WebLog                string   `yaml:"weblog"`
		WebPath               string   `yaml:"webpath"`
		UpdateKey             string   `yaml:"updatekey"`
		CheckBuddyIntervalSec int      `yaml:"checkBuddyIntervalSec"`
		BuddyDownStatusMsg    string   `yaml:"buddyDownStatusMsg"`
		BuddyHosts            []Buddy  `yaml:"buddy"`
		IgnorePrefix          []string `yaml:"ignorePrefix"`
		TtlBehavior           string   `yaml:"ttlbehavior"`
		TtlOkDelete           int      `yaml:"ttlokdelete"`
		DisableDependOn       bool     `yaml:"disableDependOn"`
		DisableMetrics        bool     `yaml:"disableMetrics"`
		Prometheusurl         string   `yaml:"prometheusurl"`
	}
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

func (conf *Configer) ReadEnv() {
	var tmp_buddy Buddy

	if os.Getenv("DASHNAME") != "" {
		config.DashName = os.Getenv("DASHNAME")
	}
	if os.Getenv("IPPORT") != "" {
		config.IPport = os.Getenv("IPPORT")
	}
	if os.Getenv("WEBLOG") != "" {
		conf.WebLog = os.Getenv("WEBLOG")
	}
	if os.Getenv("WEBPATH") != "" {
		conf.WebPath = os.Getenv("WEBPATH")
	}
	if os.Getenv("UPDATEKEY") != "" {
		conf.UpdateKey = os.Getenv("UPDATEKEY")
	}
	if os.Getenv("CHECKBUDDYINTERVALSEC") != "" {
		conf.CheckBuddyIntervalSec = str2int(os.Getenv("CHECKBUDDYINTERVALSEC"))
	}
	if os.Getenv("BUDDYDOWNSTATUSMSG") != "" {
		conf.BuddyDownStatusMsg = os.Getenv("BUDDYDOWNSTATUSMSG")
	}
	if os.Getenv("BUDDYNAME") != "" && os.Getenv("BUDDYURL") != "" {
		tmp_buddy.Name = os.Getenv("BUDDYNAME")
		tmp_buddy.Url = os.Getenv("BUDDYURL")
		if os.Getenv("BUDDYKEY") != "" {
			tmp_buddy.Key = os.Getenv("BUDDYKEY")
		}
		conf.BuddyHosts = append(conf.BuddyHosts, tmp_buddy)
	}
	if os.Getenv("IGNOREPREFIX") != "" {
		conf.IgnorePrefix = append(conf.IgnorePrefix, os.Getenv("IGNOREPREFIX"))
	}
	if os.Getenv("NSCONFIG") != "" {
		nsconfig = os.Getenv("NSCONFIG")
	}
	if os.Getenv("TTLBEHAVIOR") != "" {
		conf.TtlBehavior = os.Getenv("TTLBEHAVIOR")
	}
	if os.Getenv("TTLOKDELETE") != "" {
		conf.TtlOkDelete = str2int(os.Getenv("TTLOKDELETE"))
	}
	if os.Getenv("DISABLEDEPENDSON") != "" {
		conf.DisableDependOn = str2bool(os.Getenv("DISABLEDEPENDSON"))
	}
	if os.Getenv("DISABLEMETRICS") != "" {
		conf.DisableMetrics = str2bool(os.Getenv("DISABLEMETRICS"))
	}
	if os.Getenv("PROMETHEUSURL") != "" {
		conf.Prometheusurl = os.Getenv("PROMETHEUSURL")
	}
}

// InitConfig initiates a new decoded Config struct Alex style
func (conf *Configer) InitConfig(configPath string) error {
	var result error

	if configPath == "" {
		configPath = "dashgoat.yaml"
	}

	fileExists := isExists(configPath, "file")
	if !fileExists {
		result = fmt.Errorf("Cant find Config file " + configPath + ", moving on")
		configPath = ""
	}

	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return err
		}

		defer file.Close()

		d := yaml.NewDecoder(file)

		if err := d.Decode(&config); err != nil {
			return err
		}
		fmt.Println("Using settings from " + configPath + " ignoring cli args")
	}

	if conf.DashName == "" {
		conf.DashName = "dashGoat"
	}

	if configPath == "" { // buddy settings
		if buddy_cli.Url != "" && buddy_cli.Url != "0" {
			conf.BuddyHosts = append(conf.BuddyHosts, buddy_cli)
		}

		conf.CheckBuddyIntervalSec = 11

		if len(conf.BuddyHosts) > 0 {
			err := validateBuddyConf()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if conf.BuddyDownStatusMsg == "" {
		conf.BuddyDownStatusMsg = "error"
	}

	if conf.TtlBehavior == "" {
		conf.TtlBehavior = "promotetook"
	} else {
		conf.TtlBehavior = strings.ToLower(conf.TtlBehavior)
	}
	if conf.TtlOkDelete == 0 {
		conf.TtlOkDelete = 3600
	}
	generateHostFacts()
	return result
}

func validateBuddyConf() error {

	var message error

	for idx, buddy := range config.BuddyHosts {
		if buddy.Name == "" {
			message = fmt.Errorf("Missing buddyname, for " + buddy.Url)
			return message
		}

		if buddy.Key == "" {
			config.BuddyHosts[idx].Key = config.UpdateKey
		}
	}

	return nil
}

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
		fmt.Println("Cant find an IP address, check ignorePrefix config")
		os.Exit(1)
	}
	host_facts.Items.Hostnames = append(host_facts.Items.Hostnames, config.DashName)
	fmt.Print("Hostnames found: ")
	fmt.Println(host_facts.Items.Hostnames)

	if config.DisableMetrics || config.Prometheusurl == "" {
		host_facts.Items.MetricsHistory = false
		fmt.Println("MetricsHistory: off")
	} else {
		host_facts.Items.MetricsHistory = true
		fmt.Println("MetricsHistory: on")
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

func dashGoatReady() bool {
	host_facts.mutex.RLock()
	defer host_facts.mutex.RUnlock()
	return host_facts.Items.Ready
}

func setDashGoatReady(ready bool) {
	host_facts.mutex.Lock()
	defer host_facts.mutex.Unlock()
	host_facts.Items.Ready = ready
}
