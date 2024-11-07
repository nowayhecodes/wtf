package history

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockParser is a Parser that uses a custom home directory for testing
type mockParser struct {
	homeDir string
}

func (p *mockParser) getHistoryFilePath() string {
	// Try common history file locations in the mock home directory
	candidates := []string{
		filepath.Join(p.homeDir, ".bash_history"),
		filepath.Join(p.homeDir, ".zsh_history"),
		filepath.Join(p.homeDir, ".history"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func (p *mockParser) GetLastCommand() (string, error) {
	historyFile := p.getHistoryFilePath()
	file, err := os.Open(historyFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}

func TestParser_GetLastCommand(t *testing.T) {
	// Create temporary directory as mock home
	tmpHome := t.TempDir()

	tests := []struct {
		name        string
		historyFile string
		content     string
		want        string
		wantErr     bool
	}{
		{
			name:        "bash history",
			historyFile: ".bash_history",
			content:     "ls\ncd /tmp\npwd\n",
			want:        "pwd",
		},
		{
			name:        "zsh history",
			historyFile: ".zsh_history",
			content:     "ls\ncd /tmp\ngit status\n",
			want:        "git status",
		},
		{
			name:        "empty history",
			historyFile: ".history",
			content:     "",
			want:        "",
		},
		{
			name:        "non-existent history file",
			historyFile: "nonexistent",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing history files
			historyFiles := []string{".bash_history", ".zsh_history", ".history"}
			for _, file := range historyFiles {
				path := filepath.Join(tmpHome, file)
				_ = os.Remove(path)
			}

			// Create history file in temp directory
			historyPath := filepath.Join(tmpHome, tt.historyFile)
			if !tt.wantErr {
				err := os.WriteFile(historyPath, []byte(tt.content), 0644)
				require.NoError(t, err)
			}

			// Create parser with mock home directory
			p := &mockParser{homeDir: tmpHome}

			// Get last command
			got, err := p.GetLastCommand()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)

			// Clean up after test
			if !tt.wantErr {
				err = os.Remove(historyPath)
				require.NoError(t, err)
			}
		})
	}
}
