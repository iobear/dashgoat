package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

//Buddy struct
type Buddy struct {
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
	Key    string `yaml:"key"`
	Ignore bool   `yaml:"ignore"`
}

type Configer struct {
	DashName              string  `yaml:"dashName"`
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
	}

	if dashName == "" {
		conf.DashName = "dashGoat"
		dashName = conf.DashName
	}

	if configPath == "" {
		conf.EnableBuddy = false
		conf.BuddyHosts = nil
		conf.CheckBuddyIntervalSec = 30
	}

	if conf.BuddyDown == "" {
		conf.BuddyDown = "error"
	}

	return result
}

func (conf *Configer) OneBuddy(url string, key string, name string) error {
	var onebuddyconf Buddy

	conf.EnableBuddy = true
	onebuddyconf.Ignore = false
	onebuddyconf.Key = key
	onebuddyconf.Url = url
	onebuddyconf.Name = "MyBuddy"

	if name != "" {
		onebuddyconf.Name = name
	}

	conf.BuddyHosts = nil
	conf.BuddyHosts = append(conf.BuddyHosts, onebuddyconf)

	return nil
}
