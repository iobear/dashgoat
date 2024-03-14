/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"fmt"
	"os"
	"strings"

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
)

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

	// buddy settings
	if configPath == "" {
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

	// Default status when Buddy is down
	if conf.BuddyDownStatusMsg == "" {
		conf.BuddyDownStatusMsg = "warning"
	}

	// Default TTL bahaviour
	if conf.TtlBehavior == "" {
		conf.TtlBehavior = "promotetook"
	} else {
		conf.TtlBehavior = strings.ToLower(conf.TtlBehavior)
	}

	// Default delete time on resolved TTL
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
