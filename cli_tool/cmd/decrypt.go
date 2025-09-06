package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file",
	Long:  `Decrypt contents of a given file, writing the output to another specified file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Decrypting file...")
	},
}