package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nowayhecodes/wtf/internal/config"
)

// HistoryParser interface for command history operations
type HistoryParser interface {
	GetLastCommand() (string, error)
}

// ErrorDetector interface for error detection and suggestion
type ErrorDetector interface {
	HasError(cmd string) bool
	Suggest(cmd string) string
}

// CommandExecutor interface for command execution
type CommandExecutor interface {
	Execute(cmd string) error
}

// App represents the main application
type App struct {
	cfg        *config.Config
	histParser HistoryParser
	detector   ErrorDetector
	executor   CommandExecutor
}

// New creates a new App instance
func New(cfg *config.Config, histParser HistoryParser, detector ErrorDetector, executor CommandExecutor) *App {
	return &App{
		cfg:        cfg,
		histParser: histParser,
		detector:   detector,
		executor:   executor,
	}
}

// Run executes the main application logic
func (a *App) Run() error {
	// Get last command from history
	lastCmd, err := a.histParser.GetLastCommand()
	if err != nil {
		return fmt.Errorf("failed to get last command: %w", err)
	}

	// Check if command had an error
	if !a.detector.HasError(lastCmd) {
		return nil
	}

	// Get correction suggestion
	suggestion := a.detector.Suggest(lastCmd)
	if suggestion == "" {
		return nil
	}

	// Prompt user for correction
	fmt.Printf("Did you mean: %s? [Y/n] ", suggestion)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response == "y" || response == "" {
		return a.executor.Execute(suggestion)
	}

	return nil
}
