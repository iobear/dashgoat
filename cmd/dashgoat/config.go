package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Buddy struct
type Buddy struct {
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
	Ignore bool   `yaml:"ignore"`
}

type Configer struct {
	DashName              string  `yaml:"dashName"`
	IPport                string  `yaml:"ipport"`
	WebLog                string  `yaml:"weblog"`
	WebPath               string  `yaml:"webpath"`
	UpdateKey             string  `yaml:"updatekey"`
	EnableBuddy           bool    `yaml:"enableBuddy"`
	CheckBuddyIntervalSec int     `yaml:"checkBuddyIntervalSec"`
	BuddyDown             string  `yaml:"buddyDown"`
	BuddyHosts            []Buddy `yaml:"buddy"`
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
		if buddyCli.Url != "" && buddyCli.Url != "0" {
			conf.BuddyHosts = append(conf.BuddyHosts, buddyCli)
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

	if conf.BuddyDown == "" {
		conf.BuddyDown = "error"
	}

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
