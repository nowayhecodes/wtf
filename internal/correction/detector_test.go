package correction

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nowayhecodes/wtf/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetector_HasError(t *testing.T) {
	cfg := &config.Config{LevenThreshold: 2}
	d := NewDetector(cfg)

	tests := []struct {
		name    string
		cmd     string
		want    bool
		setup   func() error
		cleanup func() error
	}{
		{
			name: "valid command",
			cmd:  "ls -la",
			want: false,
		},
		{
			name: "invalid command",
			cmd:  "invalidcmd",
			want: true,
		},
		{
			name: "empty command",
			cmd:  "",
			want: false,
		},
		{
			name: "valid file path",
			cmd:  "cat ./testfile",
			want: false,
			setup: func() error {
				return os.WriteFile("testfile", []byte("test"), 0644)
			},
			cleanup: func() error {
				return os.Remove("testfile")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				require.NoError(t, tt.setup())
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			got := d.HasError(tt.cmd)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetector_Suggest(t *testing.T) {
	cfg := &config.Config{
		LevenThreshold: 2,
		CustomRules: map[string]string{
			"gti": "git",
		},
	}
	d := NewDetector(cfg)

	// Create test directory structure
	tmpDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644))

	tests := []struct {
		name string
		cmd  string
		want string
	}{
		{
			name: "custom rule match",
			cmd:  "gti status",
			want: "git status",
		},
		{
			name: "similar file path",
			cmd:  filepath.Join(tmpDir, "test.tx"),
			want: filepath.Join(tmpDir, "test.txt"),
		},
		{
			name: "empty command",
			cmd:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.Suggest(tt.cmd)
			assert.Equal(t, tt.want, got)
		})
	}
}
