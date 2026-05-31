package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add files to your AI drop directory",
	Long: `Copy a file or files by providing the relative path to your AIDrop directory.

	Specify a project to move the files into the project name subdirectory under the AIDrop directory. 
	If no project is specified, the git repository name will be used. Otherwise, "default-project" will be used.

	Provide a session name to move the files inside a subdirectory within the project subdirectory. 
	The session directory name will have a day in the format of year-month-day prepended to the directory name.
	If no session is provided, the files will be added to the root of the project directory under the AIDrop directory. 

	Use the move flag to move the files instead of copying them to the AIDrop directory. 

	Note: if a file name conflict occurs, a number will be prepended to the file name. If main.go already exists, 
		main-2.go will be used instead.
	
	Note: for symlinks if a symlink is provided, the target will be copied instead of the symlink. 

	For example:
	aidrop add -p federation-service README.md internal/models.go
		copies ./README.md and ./internal/models.go to ~/AIDrop/federation-service/README.md and ~/AIDrop/federation-service/models.go
	aidrop add -p snake-game -s add-animation animate.go
		copies ./animate.go to ~/AIDrop/snake-game/2024-06-14-add-animation/animate.go
	aidrop add -p simulator -s stack-overflow-issue -m output.log
		moves ./output.log to ~/AIDrop/simulator/2026-05-30-stack-overflow-issue -m output.log`,
	Run: add,
}

func add(cmd *cobra.Command, args []string) {
	fmt.Println("add called")
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addCmd.Flags().StringP("project", "p", "default-project", "specify a project name")
	addCmd.Flags().StringP("session", "s", "", "specify a session name")
	addCmd.Flags().BoolP("move", "m", false, "move file to project instead of copying")
}
