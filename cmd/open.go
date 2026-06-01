package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [project [session]]",
	Short: "Open the AIDrop directory, a project, or a session in the system file manager",
	Long: `Open a path within the AIDrop staging area using the system file manager (Finder on
macOS, xdg-open on Linux). With no arguments the top-level AIDrop directory is
opened. Provide a project name to open that project folder, and optionally a
session name to open a specific session inside that project.

Examples:
  aidrop open
    Opens ~/AIDrop/ in Finder / file manager.

  aidrop open federation-service
    Opens ~/AIDrop/federation-service/.

  aidrop open federation-service 2026-05-31-auth-bug
    Opens ~/AIDrop/federation-service/2026-05-31-auth-bug/.`,
	Args: cobra.MaximumNArgs(2),
	RunE: openPath,
}

func openPath(cmd *cobra.Command, args []string) error {
	dropDir, err := getAIDropDir()
	if err != nil {
		return err
	}

	target := dropDir
	if len(args) >= 1 {
		target = filepath.Join(dropDir, args[0])
	}
	if len(args) >= 2 {
		target = filepath.Join(target, args[1])
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", target)
	}

	var openCmd string
	switch runtime.GOOS {
	case "darwin":
		openCmd = "open"
	case "linux":
		openCmd = "xdg-open"
	default:
		return fmt.Errorf("aidrop open is not supported on %s", runtime.GOOS)
	}

	c := exec.Command(openCmd, target)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("could not open %s: %w", target, err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(openCmd)
}
