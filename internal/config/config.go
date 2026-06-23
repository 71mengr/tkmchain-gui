package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	NodeURL              string        `json:"node_url"`
	DataDir              string        `json:"data_dir"`
	AutoRefresh          int           `json:"auto_refresh_seconds"`
	Theme                string        `json:"theme"`
	DialTimeoutSeconds   int           `json:"dial_timeout_seconds"`    // seconds for initial dial / verification
	RetryCount           int           `json:"retry_count"`             // number of dial retries
	RetryIntervalSeconds int           `json:"retry_interval_seconds"`  // seconds between retries
}

func Load() *Config {
	cfg := &Config{
		NodeURL:              "http://localhost:8545",
		DataDir:              getDefaultDataDir(),
		AutoRefresh:          30,
		Theme:                "dark",
		DialTimeoutSeconds:   10,
		RetryCount:           5,
		RetryIntervalSeconds: 2,
	}

	configPath := filepath.Join(cfg.DataDir, "config.json")
	if data, err := os.ReadFile(configPath); err == nil {
		// ignore unmarshal errors but keep defaults
		_ = json.Unmarshal(data, cfg)
	}

	return cfg
}

func getDefaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".gtkm-wallet")
}
