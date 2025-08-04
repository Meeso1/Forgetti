package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "forgetti",
	Short: "Forgetti CLI tool",
	Long:  `CLI tool that encrypts and data and sometimes decrypts it too.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run `forgetti --help` for a list of available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
