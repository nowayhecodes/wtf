package history

import (
	"bytes"
	"io"
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

	// Read last 4KB of the file (usually more than enough for last command)
	const readSize = 4096
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	size := stat.Size()
	if size == 0 {
		return "", nil
	}

	offset := size - readSize
	if offset < 0 {
		offset = 0
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return "", err
	}

	buf := make([]byte, readSize)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	buf = buf[:n]

	// Find the last complete line
	lastNewline := bytes.LastIndex(buf, []byte("\n"))
	if lastNewline == -1 {
		return string(buf), nil
	}

	// If we started reading in the middle of a line and found multiple newlines,
	// take the last complete line
	start := bytes.LastIndex(buf[:lastNewline], []byte("\n"))
	if start == -1 {
		start = 0
	} else {
		start++ // Skip the newline
	}

	return string(buf[start:lastNewline]), nil
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
