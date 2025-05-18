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
	Port       Port
	Slaves     []Port
	Methods    []string
	isInternal bool
	core       *Config
}
type StorageConfig struct {
	Persistent bool
	core       *Config
}

func validateConfig(_ *Config) error {
	return nil
}

func (cfg *Config) Update(fileName string) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	json.Unmarshal(bytes, cfg)
	if err := validateConfig(cfg); err != nil {
		return err
	}
	return nil
}

func ReadConfig(fileName string) (*Config, error) {
	cfg := &Config{}
	cfg.Update(fileName)
	return cfg, nil
}

func MakeHTTPConfig(cfg *Config, useInternalPort bool) *HTTPConfig {
	cfgHTTP := &HTTPConfig{
		core:       cfg,
		isInternal: useInternalPort,
	}
	cfgHTTP.Update()
	return cfgHTTP
}

func (cfg *HTTPConfig) Update() {
	port := cfg.core.Port
	if cfg.isInternal {
		port = cfg.core.InternalPort
	}
	cfg.Port = port
	cfg.Slaves = cfg.core.Slaves
	cfg.Methods = cfg.core.Methods
}

func MakeStorageConfig(cfg *Config) *StorageConfig {
	cfgStorage := &StorageConfig{
		core: cfg,
	}
	cfgStorage.Update()
	return cfgStorage
}

func (cfg *StorageConfig) Update() {
	cfg.Persistent = cfg.core.Persistent
}
