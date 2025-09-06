package cmd

import (
	"Forgetti/commands"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	encryptCmd.Flags().StringVarP(&encrypt_password, "password", "p", "", "The password to encrypt the file with")
	encryptCmd.Flags().StringVarP(&encrypt_expiresIn, "expires-in", "e", "1d", "The time after which the encrypted file will expire (format: 1y/2/mo/3w/4d/5h/6min)")
	encryptCmd.Flags().StringVarP(&encrypt_serverAddress, "server-address", "s", "", "The address of the server to encrypt the file with")
	encryptCmd.Flags().StringVarP(&encrypt_inputPath, "input", "i", "", "The path to the input file")
	encryptCmd.Flags().StringVarP(&encrypt_outputPath, "output", "o", "", "The path to the output file")
	encryptCmd.Flags().BoolVarP(&encrypt_overwrite, "overwrite", "w", false, "Overwrite the output file if it already exists")
	encryptCmd.Flags().BoolVarP(&encrypt_verbose, "verbose", "v", false, "Verbose output")
	encryptCmd.Flags().BoolVarP(&encrypt_quiet, "quiet", "q", false, "Quiet output")

	rootCmd.AddCommand(encryptCmd)
}

var encrypt_password string
var encrypt_expiresIn string
var encrypt_serverAddress string
var encrypt_inputPath string
var encrypt_outputPath string
var encrypt_overwrite bool
var encrypt_verbose bool
var encrypt_quiet bool

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file",
	Long:  `Encrypt contents of a given file, writing the output to another specified file.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := commands.CreateEncryptInput(
			encrypt_inputPath, 
			encrypt_outputPath,
			encrypt_password,
			encrypt_expiresIn,
			encrypt_serverAddress,
			encrypt_overwrite,
			encrypt_verbose,
			encrypt_quiet,
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = commands.Encrypt(*input)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}