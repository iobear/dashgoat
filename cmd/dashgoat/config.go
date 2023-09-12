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
		EnableBuddy           bool     `yaml:"enableBuddy"`
		CheckBuddyIntervalSec int      `yaml:"checkBuddyIntervalSec"`
		BuddyDownStatusMsg    string   `yaml:"buddyDown"`
		BuddyHosts            []Buddy  `yaml:"buddy"`
		IgnorePrefix          []string `yaml:"ignorePrefix"`
	}
	HostFact struct {
		Hostnames       []string
		DashName        string
		Ready           bool
		UpAt            time.Time
		UpAtEpoch       int64
		DashGoatVersion string
		GoVersion       string
	}
	HostFacts struct {
		Items HostFact
		mutex sync.RWMutex
	}
)

var host_facts HostFacts

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

		if conf.BuddyHosts != nil {
			conf.EnableBuddy = true
		}

		if conf.EnableBuddy {
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
	host_facts.Items.DashGoatVersion = "1.3.0"
	host_facts.Items.GoVersion = runtime.Version()

	hostname, _ := os.Hostname()
	IPhost := ""

	for _, ip := range getHostIPs() {
		IPhost = hostname + "-" + ip
		host_facts.Items.Hostnames = append(host_facts.Items.Hostnames, IPhost)
	}
	if len(host_facts.Items.Hostnames) == 0 {
		fmt.Println("Cant find an IP address, check ignorePrefix config")
		os.Exit(1)
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
			result = append(result, addr.String())
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
			fmt.Println("ignoring " + ip_addr)
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
