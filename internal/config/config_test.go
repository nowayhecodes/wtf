package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	configData := `{
		"customRules": {
			"gti": "git",
			"sl": "ls"
		},
		"historyFile": "/custom/history",
		"shellType": "zsh",
		"maxSuggestions": 5,
		"levenThreshold": 3
	}`
	err := os.WriteFile(configPath, []byte(configData), 0644)
	require.NoError(t, err)

	tests := []struct {
		name        string
		configPath  string
		wantConfig  *Config
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid config file",
			configPath: configPath,
			wantConfig: &Config{
				CustomRules: map[string]string{
					"gti": "git",
					"sl":  "ls",
				},
				HistoryFile:    "/custom/history",
				ShellType:      "zsh",
				MaxSuggestions: 5,
				LevenThreshold: 3,
			},
		},
		{
			name:       "empty config path uses default values",
			configPath: "",
			wantConfig: &Config{
				CustomRules:    make(map[string]string),
				ShellType:      "bash",
				MaxSuggestions: 3,
				LevenThreshold: 2,
			},
		},
		{
			name:        "invalid json",
			configPath:  filepath.Join(tmpDir, "invalid.json"),
			wantErr:     true,
			errContains: "no such file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Load(tt.configPath)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantConfig, cfg)
		})
	}
}
