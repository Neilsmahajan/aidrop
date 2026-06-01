package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <file> [files...]",
	Short: "Add files to the AIDrop staging area",
	Long: `Copy (or move) one or more files into the AIDrop staging area at ~/AIDrop/<project>/[<date>-<session>/].

Project resolution order:
  1. The value of --project if provided.
  2. The name of the current git repository root (automatic inference).
  3. "default" if no git repository is detected.

If --session is specified, files are placed inside a date-prefixed subdirectory
(YYYY-MM-DD-<session>) within the project folder. Without --session, files land
directly at the project root.

Filename conflicts are resolved automatically by appending a numeric suffix
(e.g., file-2.txt, file-3.txt).

Symbolic links are followed — the resolved target is copied, not the link itself.

Examples:
  aidrop add -p federation-service README.md internal/models.go
    Copies README.md and models.go to ~/AIDrop/federation-service/

  aidrop add -p snake-game -s add-animation animate.go
    Copies animate.go to ~/AIDrop/snake-game/2026-05-31-add-animation/animate.go

  aidrop add -s stack-overflow-issue -m output.log
    Moves output.log to ~/AIDrop/<git-repo>/2026-05-31-stack-overflow-issue/output.log`,
	Args: cobra.MinimumNArgs(1),
	RunE: add,
}

func add(cmd *cobra.Command, args []string) error {
	project, _ := cmd.Flags().GetString("project")
	session, _ := cmd.Flags().GetString("session")
	move, _ := cmd.Flags().GetBool("move")

	// Infer project name if not explicitly provided.
	if project == "" {
		project = getGitRepoName()
	}
	if project == "" {
		project = "default"
	}

	dropDir, err := getAIDropDir()
	if err != nil {
		return err
	}

	destDir := filepath.Join(dropDir, project)
	if session != "" {
		datePrefix := time.Now().Format("2006-01-02")
		destDir = filepath.Join(destDir, datePrefix+"-"+session)
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("could not create destination directory %s: %w", destDir, err)
	}

	for _, arg := range args {
		// Resolve symlinks so we always copy the real file content.
		resolved, err := filepath.EvalSymlinks(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", arg, err)
			continue
		}

		info, err := os.Stat(resolved)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", arg, err)
			continue
		}

		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "warning: skipping %s: directories are not supported; use a glob to expand individual files\n", arg)
			continue
		}

		destPath := resolveConflict(filepath.Join(destDir, filepath.Base(resolved)))

		if move {
			if err := os.Rename(resolved, destPath); err != nil {
				// os.Rename fails across filesystem boundaries; fall back to copy+remove.
				if err2 := copyFile(resolved, destPath); err2 != nil {
					fmt.Fprintf(os.Stderr, "error: could not move %s: %v\n", arg, err2)
					continue
				}
				if err2 := os.Remove(resolved); err2 != nil {
					fmt.Fprintf(os.Stderr, "warning: file copied but could not remove source %s: %v\n", arg, err2)
				}
			}
			fmt.Printf("moved  %s  →  %s\n", arg, destPath)
		} else {
			if err := copyFile(resolved, destPath); err != nil {
				fmt.Fprintf(os.Stderr, "error: could not copy %s: %v\n", arg, err)
				continue
			}
			fmt.Printf("added  %s  →  %s\n", arg, destPath)
		}
	}
	return nil
}

// copyFile copies the file at src to dst, preserving content exactly.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("project", "p", "", "Project name (defaults to git repository name, then \"default\")")
	addCmd.Flags().StringP("session", "s", "", "Session name; places files in a YYYY-MM-DD-<session> subdirectory")
	addCmd.Flags().BoolP("move", "m", false, "Move files instead of copying them")
}
