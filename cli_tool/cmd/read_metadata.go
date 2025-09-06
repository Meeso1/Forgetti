package cmd

import (
	"Forgetti/commands"
	"fmt"

	"github.com/spf13/cobra"
)

var readMetadata_inputPath string

func init() {
	readMetadataCmd.Flags().StringVarP(&readMetadata_inputPath, "input", "i", "", "The path to the encrypted file")

	rootCmd.AddCommand(readMetadataCmd)
}

var readMetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Read metadata from an encrypted file",
	Long:  `Read metadata from an encrypted file.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := commands.CreateReadMetadataInput(readMetadata_inputPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = commands.ReadMetadata(*input)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}