package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Config ... Configuration of the program
type Config struct {
	Backend struct {
		ExePath    string `yaml:"exePath"`
		ExeArgs    string `yaml:"exeArgs"`
		StatusFile string `yaml:"statusFile"`
		LogoFile   string `yaml:"logoFile"`
	} `yaml:"backend"`
	Userspace struct {
		Enabled             bool  `yaml:"enabled"`
		NotifyOnStepsStart  []int `yaml:"notifyOnStepsStarted"`
		NotifyOnStepsFinish []int `yaml:"notifyOnStepsFinished"`
	} `yaml:"userspace"`
}

// This is the default config file
var configFilename = "abasgo.yml"

// InitializeConfiguration ... Load the configuration from 1. Defaults 2. File
// Also write the file with defaults if it doesn't exist yet
func InitializeConfiguration() *Config {
	cfg := defaults()
	if _, err := os.Stat(configFilename); errors.Is(err, os.ErrNotExist) {
		err := writeFile(&cfg)
		if err != nil {
			log.Println("Could not write default configuration file.", err)
		}
		return &cfg
	}

	// File was found? Then read the values from it
	err := readFile(&cfg)
	if err != nil {
		log.Panic(" PANIC! Could not read configuration file.", err)
	}

	return &cfg
}

// These are the defaults of the configuration
func defaults() Config {
	var cfg Config
	// Initialize defaults for backend settings
	cfg.Backend.ExePath = "abasgui.exe"
	cfg.Backend.ExeArgs = "-lang E"
	cfg.Backend.StatusFile = "buildStatus.txt"
	cfg.Backend.LogoFile = "logo.png"
	// Initialize defaults for user space settings
	cfg.Userspace.Enabled = true
	cfg.Userspace.NotifyOnStepsStart = []int{1, 8}
	cfg.Userspace.NotifyOnStepsFinish = []int{1, 8}

	return cfg
}

// Write the struct to the configuration file
func writeFile(cfg *Config) error {
	// Parse struct to yaml format
	yamlString, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	// Finally write to file
	err = os.WriteFile(configFilename, yamlString, 0644)
	return err
}

// Read struct from yaml file
func readFile(cfg *Config) error {
	yamlString, err := os.ReadFile(configFilename)
	if err != nil {
		return err
	}

	// Merge the file into the existing struct, in case anything is missing from the config file
	err = yaml.Unmarshal(yamlString, cfg)
	if err != nil {
		return err
	}

	return nil
}
