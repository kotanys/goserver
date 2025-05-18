package main

import (
	"encoding/json"
	"os"
)

type Port int
type Config struct {
	Port         Port     `json:"port"`
	InternalPort Port     `json:"internal_port"`
	LogFile      string   `json:"log_file"`
	Slaves       []Port   `json:"slaves"`
	Methods      []string `json:"methods"`
	Persistent   bool     `json:"persistent"`
}
type HTTPConfig struct {
	Port    Port
	Slaves  []Port
	Methods []string
}

func validateConfig(_ *Config) error {
	return nil
}

func ReadConfig(fileName string) (*Config, error) {
	config := &Config{}
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, config)
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

func MakeHTTPConfig(cfg *Config, useInternalPort bool) *HTTPConfig {
	port := cfg.Port
	if useInternalPort {
		port = cfg.InternalPort
	}
	cfgHTTP := &HTTPConfig{
		Port:    port,
		Slaves:  cfg.Slaves,
		Methods: cfg.Methods,
	}
	return cfgHTTP
}
