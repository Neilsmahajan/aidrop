package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getAIDropDir returns the path to the AIDrop staging directory (~/AIDrop).
func getAIDropDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, "AIDrop"), nil
}

// getGitRepoName returns the base name of the current git repository root,
// or an empty string if the working directory is not inside a git repository.
func getGitRepoName() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	return filepath.Base(strings.TrimSpace(string(out)))
}

// resolveConflict returns a non-conflicting destination path. If destPath already
// exists, it appends a numeric suffix before the extension (-2, -3, …) until a
// free slot is found.
func resolveConflict(destPath string) string {
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return destPath
	}
	ext := filepath.Ext(destPath)
	base := strings.TrimSuffix(destPath, ext)
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s-%d%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
