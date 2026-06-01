package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List staged files organized by project and session",
	Long: `Print a tree of all files under the AIDrop directory, organized by project and session.

Use --project to narrow the output to a single project.

Examples:
  aidrop ls
    Prints the full AIDrop tree across all projects and sessions.

  aidrop ls -p federation-service
    Prints only the files under ~/AIDrop/federation-service/.`,
	RunE: ls,
}

func ls(cmd *cobra.Command, args []string) error {
	project, _ := cmd.Flags().GetString("project")

	dropDir, err := getAIDropDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dropDir); os.IsNotExist(err) {
		fmt.Println("AIDrop directory does not exist yet. Run `aidrop add` to get started.")
		return nil
	}

	fmt.Printf("AIDrop  %s\n", dropDir)

	if project != "" {
		return printProject(dropDir, project, "")
	}

	entries, err := os.ReadDir(dropDir)
	if err != nil {
		return fmt.Errorf("could not read AIDrop directory: %w", err)
	}

	var projects []string
	for _, e := range entries {
		if e.IsDir() {
			projects = append(projects, e.Name())
		}
	}
	sort.Strings(projects)

	for i, proj := range projects {
		isLast := i == len(projects)-1
		connector, childPrefix := treeChars(isLast)
		fmt.Printf("%s%s/\n", connector, proj)
		if err := printProject(dropDir, proj, childPrefix); err != nil {
			return err
		}
	}
	return nil
}

func printProject(dropDir, project, prefix string) error {
	projDir := filepath.Join(dropDir, project)
	entries, err := os.ReadDir(projDir)
	if err != nil {
		return fmt.Errorf("could not read project %q: %w", project, err)
	}

	var sessions, files []string
	for _, e := range entries {
		if e.IsDir() {
			sessions = append(sessions, e.Name())
		} else {
			files = append(files, e.Name())
		}
	}
	sort.Strings(sessions)
	sort.Strings(files)

	// Sessions are printed first, then loose files.
	type entry struct {
		name      string
		isSession bool
	}
	var all []entry
	for _, s := range sessions {
		all = append(all, entry{s, true})
	}
	for _, f := range files {
		all = append(all, entry{f, false})
	}

	for i, item := range all {
		isLast := i == len(all)-1
		connector, childPrefix := treeChars(isLast)
		if item.isSession {
			fmt.Printf("%s%s%s/\n", prefix, connector, item.name)
			printSession(filepath.Join(projDir, item.name), prefix+childPrefix)
		} else {
			fmt.Printf("%s%s%s\n", prefix, connector, item.name)
		}
	}
	return nil
}

func printSession(sessionDir, prefix string) {
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	for i, name := range names {
		isLast := i == len(names)-1
		connector, _ := treeChars(isLast)
		fmt.Printf("%s%s%s\n", prefix, connector, name)
	}
}

// treeChars returns the branch connector and child prefix for tree rendering.
func treeChars(isLast bool) (connector, childPrefix string) {
	if isLast {
		return "└── ", "    "
	}
	return "├── ", "│   "
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().StringP("project", "p", "", "Restrict output to a specific project")
}
