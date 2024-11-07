package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	CustomRules    map[string]string `json:"customRules"`
	HistoryFile    string            `json:"historyFile"`
	ShellType      string            `json:"shellType"`
	MaxSuggestions int               `json:"maxSuggestions"`
	LevenThreshold int               `json:"levenThreshold"`
}

func Load(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	cfg := &Config{
		CustomRules:    make(map[string]string),
		ShellType:      "bash",
		MaxSuggestions: 3,
		LevenThreshold: 2,
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".wtf.json"
	}
	return filepath.Join(homeDir, ".wtf.json")
}
