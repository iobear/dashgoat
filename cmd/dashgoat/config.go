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

	if configPath == "" {
		configPath = "dashgoat.yaml"
	}

	fileExists := isExists(configPath, "file")
	if !fileExists {
		return fmt.Errorf("Cant read Config file " + configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return err
	}

	if dashName == "dashGoat" {
		dashName = conf.DashName
	}

	if conf.BuddyDown == "" {
		conf.BuddyDown = "error"
	}

	return nil
}
