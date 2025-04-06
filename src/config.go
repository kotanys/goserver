package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port    int    `json:"port"`
	LogFile string `json:"log_file"`
}

func ReadConfig(fileName string) (*Config, error) {
	config := &Config{}
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, config)
	return config, nil
}
