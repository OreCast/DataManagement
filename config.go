package main

// config module
//
// Copyright (c) 2023 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Configuration stores server configuration parameters
type Configuration struct {
	// web server parts
	Base    string `json:"base"`     // base URL
	LogFile string `json:"log_file"` // server log file
	Port    int    `json:"port"`     // server port number
	Verbose int    `json:"verbose"`  // verbose output

	// server parts
	RootCAs       string   `json:"rootCAs"`      // server Root CAs path
	ServerCrt     string   `json:"server_cert"`  // server certificate
	ServerKey     string   `json:"server_key"`   // server certificate
	DomainNames   []string `json:"domain_names"` // LetsEncrypt domain names
	LimiterPeriod string   `json:"rate"`         // limiter rate value

	// OreCast parts
	DiscoveryPassword string `json:"discovery_secret"` // data-discovery password
	DiscoveryCipher   string `json:"discovery_cipher"` // data-discovery cipher
	DiscoveryURL      string `json:"discovery_url"`    // data-discovery URL
}

// Config variable represents configuration object
var Config Configuration

// helper function to parse server configuration file
func parseConfig(configFile string) error {
	data, err := os.ReadFile(filepath.Clean(configFile))
	if err != nil {
		log.Println("WARNING: Unable to read", err)
	} else {
		err = json.Unmarshal(data, &Config)
		if err != nil {
			log.Println("ERROR: Unable to parse", err)
			return err
		}
	}

	// default values
	if Config.Port == 0 {
		Config.Port = 8340
	}
	return nil
}
