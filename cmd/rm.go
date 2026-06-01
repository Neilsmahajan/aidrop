package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <project> [session]",
	Short: "Remove a project or session from the AIDrop staging area",
	Long: `Permanently remove a project directory or a specific session directory from the
AIDrop staging area.

  aidrop rm <project>
    Removes the entire project directory and everything inside it.

  aidrop rm <project> <session>
    Removes only the named session directory within the project.

Use --soft to move the target to the system trash instead of deleting it.
Use --dry-run to preview what would be removed without making any changes.

Examples:
  aidrop rm federation-service
    Deletes ~/AIDrop/federation-service/ and all its contents.

  aidrop rm federation-service 2026-05-31-auth-bug
    Deletes ~/AIDrop/federation-service/2026-05-31-auth-bug/.

  aidrop rm -s federation-service
    Moves ~/AIDrop/federation-service/ to the system trash.

  aidrop rm --dry-run federation-service
    Shows what would be removed without deleting anything.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: rm,
}

func rm(cmd *cobra.Command, args []string) error {
	soft, _ := cmd.Flags().GetBool("soft")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	dropDir, err := getAIDropDir()
	if err != nil {
		return err
	}

	target := filepath.Join(dropDir, args[0])
	if len(args) == 2 {
		target = filepath.Join(target, args[1])
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", target)
	}

	if dryRun {
		fmt.Printf("[dry-run] would remove: %s\n", target)
		return nil
	}

	if err := removeItem(target, soft); err != nil {
		return err
	}

	verb := "removed"
	if soft {
		verb = "trashed"
	}
	fmt.Printf("%s: %s\n", verb, target)
	return nil
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolP("soft", "s", false, "Move to the system trash instead of permanently deleting")
	rmCmd.Flags().Bool("dry-run", false, "Preview what would be removed without making any changes")
}
