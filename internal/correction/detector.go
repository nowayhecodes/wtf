package correction

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nowayhecodes/wtf/internal/config"
)

type Detector struct {
	cfg         *config.Config
	commonCmds  map[string]struct{} // Common commands cache
	initialized bool
}

func NewDetector(cfg *config.Config) *Detector {
	return &Detector{
		cfg:        cfg,
		commonCmds: make(map[string]struct{}),
	}
}

// initializeCommonCommands builds a cache of common commands
func (d *Detector) initializeCommonCommands() {
	if d.initialized {
		return
	}

	// Pre-populate with very common commands
	commonCommands := []string{
		"git", "ls", "cd", "pwd", "grep", "find", "docker", "kubectl",
		"vim", "nano", "cat", "echo", "mkdir", "rm", "cp", "mv",
	}

	for _, cmd := range commonCommands {
		d.commonCmds[cmd] = struct{}{}
	}

	// Add commands from PATH
	pathDirs := filepath.SplitList(os.Getenv("PATH"))
	for _, dir := range pathDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			d.commonCmds[entry.Name()] = struct{}{}
		}
	}

	d.initialized = true
}

func (d *Detector) HasError(cmd string) bool {
	if cmd == "" {
		return false
	}

	cmdParts := strings.Fields(cmd)
	if len(cmdParts) == 0 {
		return false
	}

	_, err := exec.LookPath(cmdParts[0])
	return err != nil
}

func (d *Detector) Suggest(cmd string) string {
	d.initializeCommonCommands()

	// Check custom rules first (they take priority)
	if suggestion, ok := d.cfg.CustomRules[cmd]; ok {
		return suggestion
	}

	cmdParts := strings.Fields(cmd)
	if len(cmdParts) == 0 {
		return ""
	}

	wrongCmd := cmdParts[0]
	bestMatch := d.findBestMatch(wrongCmd)
	if bestMatch == "" {
		return ""
	}

	// Replace the wrong command with the suggestion
	cmdParts[0] = bestMatch
	return strings.Join(cmdParts, " ")
}

func (d *Detector) findBestMatch(cmd string) string {
	var bestMatch string
	minDist := d.cfg.LevenThreshold + 1

	// First, try to match against common commands
	for commonCmd := range d.commonCmds {
		// Quick check for common typos
		if isCommonTypo(cmd, commonCmd) {
			return commonCmd
		}

		dist := levenshteinDistance(cmd, commonCmd)
		if dist < minDist {
			// Prefer shorter commands when the distance is the same
			if dist == minDist && len(commonCmd) > len(bestMatch) {
				continue
			}
			minDist = dist
			bestMatch = commonCmd
		}
	}

	if minDist <= d.cfg.LevenThreshold {
		return bestMatch
	}

	return ""
}

// isCommonTypo checks for common typing mistakes
func isCommonTypo(typed, actual string) bool {
	commonTypos := map[string]string{
		"gti":   "git",
		"sl":    "ls",
		"cd..":  "cd ..",
		"grpe":  "grep",
		"mkidr": "mkdir",
	}

	if correct, ok := commonTypos[typed]; ok {
		return correct == actual
	}

	return false
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
