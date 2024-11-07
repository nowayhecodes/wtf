package app

import (
	"os"
	"testing"

	"github.com/nowayhecodes/wtf/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHistoryParser struct {
	lastCmd string
	err     error
}

func (m *mockHistoryParser) GetLastCommand() (string, error) {
	return m.lastCmd, m.err
}

type mockDetector struct {
	hasError bool
	suggest  string
}

func (m *mockDetector) HasError(cmd string) bool {
	return m.hasError
}

func (m *mockDetector) Suggest(cmd string) string {
	return m.suggest
}

type mockExecutor struct {
	executed string
	err      error
}

func (m *mockExecutor) Execute(cmd string) error {
	m.executed = cmd
	return m.err
}

func TestApp_Run(t *testing.T) {
	tests := []struct {
		name        string
		lastCmd     string
		hasError    bool
		suggestion  string
		input       string
		wantExecute string
		wantErr     bool
	}{
		{
			name:        "correct command execution",
			lastCmd:     "gti status",
			hasError:    true,
			suggestion:  "git status",
			input:       "y\n",
			wantExecute: "git status",
		},
		{
			name:     "no error in command",
			lastCmd:  "ls -la",
			hasError: false,
		},
		{
			name:       "no suggestion available",
			lastCmd:    "invalid",
			hasError:   true,
			suggestion: "",
		},
		{
			name:       "user rejects suggestion",
			lastCmd:    "gti status",
			hasError:   true,
			suggestion: "git status",
			input:      "n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			histParser := &mockHistoryParser{lastCmd: tt.lastCmd}
			detector := &mockDetector{hasError: tt.hasError, suggest: tt.suggestion}
			executor := &mockExecutor{}

			// Create app with mocks
			app := &App{
				cfg:        &config.Config{},
				histParser: histParser,
				detector:   detector,
				executor:   executor,
			}

			// Mock stdin if input is provided
			if tt.input != "" {
				oldStdin := os.Stdin
				defer func() { os.Stdin = oldStdin }()

				r, w, err := os.Pipe()
				require.NoError(t, err)
				os.Stdin = r

				go func() {
					w.Write([]byte(tt.input))
					w.Close()
				}()
			}

			// Run the app
			err := app.Run()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantExecute, executor.executed)
		})
	}
}
