package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_GetLastCommand(t *testing.T) {
	// Create temporary history files
	tmpDir := t.TempDir()
	bashHistory := filepath.Join(tmpDir, ".bash_history")
	zshHistory := filepath.Join(tmpDir, ".zsh_history")

	tests := []struct {
		name        string
		historyFile string
		content     string
		want        string
		wantErr     bool
	}{
		{
			name:        "bash history",
			historyFile: bashHistory,
			content:     "ls\ncd /tmp\npwd\n",
			want:        "pwd",
		},
		{
			name:        "zsh history",
			historyFile: zshHistory,
			content:     "ls\ncd /tmp\ngit status\n",
			want:        "git status",
		},
		{
			name:        "empty history",
			historyFile: filepath.Join(tmpDir, "empty"),
			content:     "",
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write test content
			err := os.WriteFile(tt.historyFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Create parser with mock home directory
			p := &Parser{}

			// Get last command
			got, err := p.GetLastCommand()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
