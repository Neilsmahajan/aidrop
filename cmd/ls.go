package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List out all your files, projects, and sessions.",
	Long: `Print a tree of all of your files under the AIDrop directory organized by project and session.

	Specify a project to only list the files under that project.

	For example:
	aidrop ls
		prints out all your files under the AIDrop directory organized by project and session.
	aidrop ls -p federation-service
		prints out all your files under the AIDrop/federation-service directory organized by session.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ls called")
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	lsCmd.Flags().StringP("project", "p", "", "The project to list the files under")
}
