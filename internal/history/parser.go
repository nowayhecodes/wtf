package history

import (
	"bufio"
	"os"
	"path/filepath"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) GetLastCommand() (string, error) {
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

func (p *Parser) getHistoryFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Try common history file locations
	candidates := []string{
		filepath.Join(homeDir, ".bash_history"),
		filepath.Join(homeDir, ".zsh_history"),
		filepath.Join(homeDir, ".history"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
} 