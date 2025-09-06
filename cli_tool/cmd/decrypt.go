package cmd

import (
	"Forgetti/commands"
	"fmt"
	"github.com/spf13/cobra"
)

var decrypt_inputPath string
var decrypt_outputPath string
var decrypt_password string
var decrypt_serverAddress string
var decrypt_overwrite bool
var decrypt_verbose bool
var decrypt_quiet bool

func init() {
	decryptCmd.Flags().StringVarP(&decrypt_inputPath, "input-path", "i", "", "The path to the encrypted file")
	decryptCmd.Flags().StringVarP(&decrypt_outputPath, "output-path", "o", "", "The path to the output file")
	decryptCmd.Flags().StringVarP(&decrypt_password, "password", "p", "", "The password to decrypt the file with")
	decryptCmd.Flags().StringVarP(&decrypt_serverAddress, "server-address", "s", "", "The address of the server to decrypt the file with")
	decryptCmd.Flags().BoolVarP(&decrypt_overwrite, "overwrite", "w", false, "Overwrite the output file if it already exists")
	decryptCmd.Flags().BoolVarP(&decrypt_verbose, "verbose", "v", false, "Verbose output")
	decryptCmd.Flags().BoolVarP(&decrypt_quiet, "quiet", "q", false, "Quiet output")
	
	rootCmd.AddCommand(decryptCmd)
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file",
	Long:  `Decrypt contents of a given file, writing the output to another specified file.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := commands.CreateDecryptInput(
			decrypt_inputPath, 
			decrypt_outputPath, 
			decrypt_password, 
			decrypt_serverAddress,
			decrypt_overwrite,
			decrypt_verbose,
			decrypt_quiet,
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = commands.Decrypt(*input)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}