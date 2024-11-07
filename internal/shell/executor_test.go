package shell

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutor_Execute(t *testing.T) {
	e := NewExecutor()
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	tests := []struct {
		name    string
		cmd     string
		setup   func() error
		verify  func() error
		wantErr bool
	}{
		{
			name: "simple command",
			cmd:  "touch " + testFile,
			verify: func() error {
				_, err := os.Stat(testFile)
				return err
			},
		},
		{
			name:    "invalid command",
			cmd:     "invalidcommand",
			wantErr: true,
		},
		{
			name: "empty command",
			cmd:  "",
		},
		{
			name: "command with arguments",
			cmd:  "echo 'test' > " + testFile,
			verify: func() error {
				data, err := os.ReadFile(testFile)
				if err != nil {
					return err
				}
				if string(data) != "test\n" {
					return os.ErrInvalid
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				require.NoError(t, tt.setup())
			}

			err := e.Execute(tt.cmd)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.verify != nil {
				assert.NoError(t, tt.verify())
			}
		})
	}
}
