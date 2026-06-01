package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove sessions older than a specified number of days",
	Long: `Delete session directories from the AIDrop staging area that are older than --days days (default: 7).

Only dated session directories (those whose names begin with YYYY-MM-DD) are
considered for automatic age-based removal.

Loose files (files added directly to a project folder without a session) are left
untouched unless --loose is specified. When --loose is set alongside --days, only
loose files whose modification time predates the cutoff are removed. When --loose
is set without an explicit --days value, all loose files are removed regardless of age.

Use --soft to move items to the system trash instead of permanently deleting them.
Use --dry-run to preview what would be removed without making any changes.

Examples:
  aidrop clean
    Removes all sessions older than 7 days.

  aidrop clean -d 30
    Removes all sessions older than 30 days.

  aidrop clean -s
    Moves sessions older than 7 days to the system trash.

  aidrop clean --loose
    Removes all sessions older than 7 days and all loose project files.

  aidrop clean --loose -d 30
    Removes sessions older than 30 days and loose files not modified in the last 30 days.

  aidrop clean --dry-run
    Previews what would be removed without deleting anything.`,
	RunE: clean,
}

func clean(cmd *cobra.Command, args []string) error {
	days, _ := cmd.Flags().GetInt("days")
	soft, _ := cmd.Flags().GetBool("soft")
	loose, _ := cmd.Flags().GetBool("loose")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	daysExplicit := cmd.Flags().Changed("days")

	cutoff := time.Now().AddDate(0, 0, -days)

	dropDir, err := getAIDropDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dropDir); os.IsNotExist(err) {
		fmt.Println("AIDrop directory does not exist. Nothing to clean.")
		return nil
	}

	projects, err := os.ReadDir(dropDir)
	if err != nil {
		return fmt.Errorf("could not read AIDrop directory: %w", err)
	}

	removed := 0

	for _, proj := range projects {
		if !proj.IsDir() {
			continue
		}
		projDir := filepath.Join(dropDir, proj.Name())
		entries, err := os.ReadDir(projDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not read project %s: %v\n", proj.Name(), err)
			continue
		}

		for _, entry := range entries {
			entryPath := filepath.Join(projDir, entry.Name())

			if entry.IsDir() {
				// Only remove session directories whose names start with a parseable date.
				name := entry.Name()
				if len(name) < 10 {
					continue
				}
				sessionDate, err := time.Parse("2006-01-02", name[:10])
				if err != nil {
					continue
				}
				if !sessionDate.Before(cutoff) {
					continue
				}
				if dryRun {
					fmt.Printf("[dry-run] would remove session: %s\n", entryPath)
				} else {
					if err := removeItem(entryPath, soft); err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
						continue
					}
					verb := "removed"
					if soft {
						verb = "trashed"
					}
					fmt.Printf("%s session: %s\n", verb, entryPath)
				}
				removed++

			} else if loose {
				// Remove loose files, optionally gated on modification time.
				if daysExplicit {
					info, err := entry.Info()
					if err != nil || !info.ModTime().Before(cutoff) {
						continue
					}
				}
				if dryRun {
					fmt.Printf("[dry-run] would remove loose file: %s\n", entryPath)
				} else {
					if err := removeItem(entryPath, soft); err != nil {
						fmt.Fprintf(os.Stderr, "error: %v\n", err)
						continue
					}
					verb := "removed"
					if soft {
						verb = "trashed"
					}
					fmt.Printf("%s loose file: %s\n", verb, entryPath)
				}
				removed++
			}
		}
	}

	switch {
	case removed == 0:
		fmt.Println("Nothing to clean.")
	case dryRun:
		fmt.Printf("[dry-run] %d item(s) would be removed.\n", removed)
	default:
		fmt.Printf("Cleaned %d item(s).\n", removed)
	}
	return nil
}

// removeItem deletes a file or directory. If soft is true it moves the item to
// the system trash (~/.Trash on macOS, ~/.local/share/Trash/files on Linux)
// instead of permanently deleting it.
func removeItem(path string, soft bool) error {
	if !soft {
		return os.RemoveAll(path)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	// Prefer macOS Trash; fall back to XDG Trash.
	trash := filepath.Join(home, ".Trash")
	if _, err := os.Stat(trash); os.IsNotExist(err) {
		trash = filepath.Join(home, ".local", "share", "Trash", "files")
		if err := os.MkdirAll(trash, 0755); err != nil {
			return fmt.Errorf("could not locate or create trash directory: %w", err)
		}
	}

	dest := resolveConflict(filepath.Join(trash, filepath.Base(path)))
	if err := os.Rename(path, dest); err != nil {
		// Rename fails across filesystems; give a clear error rather than silently losing data.
		return fmt.Errorf("could not move %s to trash (cross-device?): %w — use without --soft to delete directly", path, err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().IntP("days", "d", 7, "Remove sessions older than this many days")
	cleanCmd.Flags().BoolP("soft", "s", false, "Move items to the system trash instead of permanently deleting them")
	cleanCmd.Flags().Bool("loose", false, "Also remove loose files at the project root (all if --days not set, or older than --days if set)")
	cleanCmd.Flags().Bool("dry-run", false, "Preview what would be removed without making any changes")
}
