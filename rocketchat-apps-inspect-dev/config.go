package main

import (
	"errors"
	"log"
	"os"
)

type config struct {
	RootUrl  string `yaml:"rootUrl" json:"rootUrl"`
	User     string `yaml:"user" json:"user"`
	ApiKey   string `yaml:"apiKey" json:"apiKey"`
	Password string `yaml:"password" json:"password"`
}

func LoadConfig() (*config, error) {
	f, err := os.Open("config.yaml")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("failed to open config.yaml: %v", err)
	}

	f, err = os.Open("config.json")
	if err != nil {
		log.Fatalf("failed to open config.json: %v", err)
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("no config file found")
		}
	}
}
