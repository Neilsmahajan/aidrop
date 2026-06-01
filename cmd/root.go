package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aidrop",
	Short: "A fast CLI staging area for files you are about to send to an AI",
	Long: `aidrop is a command-line tool that maintains a structured staging area at ~/AIDrop/
for files you intend to share with an AI assistant or chat application.

Files are organized by project (inferred from the current git repository or set
explicitly) and optionally grouped into named sessions with automatic date prefixes.

The staging directory can be overridden by setting the AIDROP_DIR environment variable.

Available commands:
  add    Copy or move files (or directories) into the staging area.
  ls     Print a tree of staged files, projects, and sessions.
  open   Open the staging area (or a project/session) in the system file manager.
  rm     Remove a project or session from the staging area.
  clean  Remove session directories older than a specified number of days.

Run "aidrop <command> --help" for detailed usage of any command.`,
}

// Execute is the entry point called by main. It runs the root command and exits
// with a non-zero status code on error.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
