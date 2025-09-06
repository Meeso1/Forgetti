package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(encryptCmd)
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file",
	Long:  `Encrypt contents of a given file, writing the output to another specified file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Encrypting file...")
	},
}