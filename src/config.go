package main

import (
	"encoding/json"
	"os"
)

type Port int
type Config struct {
	Port    Port     `json:"port"`
	LogFile string   `json:"log_file"`
	Slaves  []Port   `json:"slaves"`
	Methods []string `json:"methods"`
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

func MakeHTTPConfig(cfg *Config) *HTTPConfig {
	cfgHTTP := &HTTPConfig{
		Port:    cfg.Port,
		Slaves:  cfg.Slaves,
		Methods: cfg.Methods,
	}
	return cfgHTTP
}
