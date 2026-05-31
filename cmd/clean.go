package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up your sessions longer than a week.",
	Long: `Delete all of the sessions older than a certain amount of days.

		Specify the amount of days to keep sessions. The default is seven.

		Provide a soft flag to move the sessions to your trash folder instead of deleting.

		The root flag will clean out all of the files not in a session, which would be directly under the root of
			the AIDrop directory. If the day's flag is provided with the root flag files with a date-modified timestamp
			older than the specified amount of days will be cleaned. Otherwise, all files directly under the root of the
			AIDrop directory that are not in a session will be cleaned regardless of their date-modified timestamp.

		For example:
		aidrop clean -d 30
			deletes all sessions older than 30 days
		aidrop clean -s
			moves all sessions older than 7 days to the trash folder instead of deleting them
		aidrop clean -r
			deletes all files directly under the root of the AIDrop directory that are not in a session.
		aidrop clean -r -d 30
			deletes all files directly under the root of the AIDrop directory that are not in a session and have a 
			date-modified timestamp older than 30 days.`,

	Run: clean,
}

func clean(cmd *cobra.Command, args []string) {
	fmt.Println("clean called")
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cleanCmd.Flags().IntP("days", "d", 7, "specify the number of days to keep sessions for")
	cleanCmd.Flags().BoolP("soft", "s", false, "move sessions to trash instead of deleting them")
	cleanCmd.Flags().BoolP("root", "r", false, "clean out all files directly under the root of the AIDrop directory that are not in a session")
}
