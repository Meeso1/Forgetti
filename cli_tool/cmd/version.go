package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version variables that can be set at build time
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `All software has versions. This is Forgetti's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Forgetti version %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
	},
} 