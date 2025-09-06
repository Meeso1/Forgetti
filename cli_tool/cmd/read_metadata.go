package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(readMetadataCmd)
}

var readMetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Read metadata from an encrypted file",
	Long:  `Read metadata from an encrypted file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reading metadata...")
	},
}