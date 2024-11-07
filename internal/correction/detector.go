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
	// Get all directories in PATH
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	bestMatch := ""
	minDist := d.cfg.LevenThreshold + 1

	// Search through each directory in PATH
	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		// Compare each executable with the input command
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// Check if file is executable
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.Mode()&0111 == 0 { // Check if executable bit is set
				continue
			}

			dist := levenshteinDistance(cmd, entry.Name())
			if dist < minDist {
				minDist = dist
				bestMatch = entry.Name()
			}
		}
	}

	if minDist <= d.cfg.LevenThreshold {
		return bestMatch
	}

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

// levenshteinDistance calculates the minimum number of single-character edits
// required to change one string into another
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix of size (len(s1)+1) x (len(s2)+1)
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill in the rest of the matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
