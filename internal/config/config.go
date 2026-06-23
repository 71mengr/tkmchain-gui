package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	NodeURL     string `json:"node_url"`
	DataDir     string `json:"data_dir"`
	AutoRefresh int    `json:"auto_refresh_seconds"`
	Theme       string `json:"theme"`
}

func Load() *Config {
	cfg := &Config{
		NodeURL:     "http://localhost:8545",
		DataDir:     getDefaultDataDir(),
		AutoRefresh: 30,
		Theme:       "dark",
	}
	
	configPath := filepath.Join(cfg.DataDir, "config.json")
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, cfg)
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
