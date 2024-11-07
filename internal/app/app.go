package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nowayhecodes/wtf/internal/config"
	"github.com/nowayhecodes/wtf/internal/correction"
	"github.com/nowayhecodes/wtf/internal/history"
	"github.com/nowayhecodes/wtf/internal/shell"
)

type App struct {
	cfg        *config.Config
	detector   *correction.Detector
	executor   *shell.Executor
	histParser *history.Parser
}

func New(cfg *config.Config) *App {
	return &App{
		cfg:        cfg,
		detector:   correction.NewDetector(cfg),
		executor:   shell.NewExecutor(),
		histParser: history.NewParser(),
	}
}

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
