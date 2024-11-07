package correction

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nowayhecodes/wtf/internal/config"
)

type Detector struct {
	cfg *config.Config
}

func NewDetector(cfg *config.Config) *Detector {
	return &Detector{cfg: cfg}
}

func (d *Detector) HasError(cmd string) bool {
	// Check if command exists in PATH
	cmdParts := strings.Fields(cmd)
	if len(cmdParts) == 0 {
		return false
	}

	_, err := exec.LookPath(cmdParts[0])
	if err == nil {
		return false
	}

	// Check if it's a file path error
	if strings.Contains(cmd, "/") {
		_, err := os.Stat(cmd)
		if err == nil {
			return false
		}
	}

	return true
}

func (d *Detector) Suggest(cmd string) string {
	// Check custom rules first
	if suggestion, ok := d.cfg.CustomRules[cmd]; ok {
		return suggestion
	}

	cmdParts := strings.Fields(cmd)
	if len(cmdParts) == 0 {
		return ""
	}

	// Try to find similar commands in PATH
	suggestion := d.findSimilarCommand(cmdParts[0])
	if suggestion != "" {
		return suggestion + strings.Join(cmdParts[1:], " ")
	}

	// Try to find similar file paths
	if strings.Contains(cmd, "/") {
		suggestion = d.findSimilarPath(cmd)
		if suggestion != "" {
			return suggestion
		}
	}

	return ""
}

func (d *Detector) findSimilarCommand(cmd string) string {
	// Implementation of Levenshtein distance calculation
	// and similar command finding logic
	return ""
}

func (d *Detector) findSimilarPath(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	var bestMatch string
	minDist := d.cfg.LevenThreshold + 1

	for _, entry := range entries {
		dist := levenshteinDistance(base, entry.Name())
		if dist < minDist {
			minDist = dist
			bestMatch = filepath.Join(dir, entry.Name())
		}
	}

	if minDist <= d.cfg.LevenThreshold {
		return bestMatch
	}

	return ""
}
