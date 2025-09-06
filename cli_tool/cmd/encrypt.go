package cmd

import (
	"Forgetti/commands"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	encryptCmd.Flags().StringVarP(&password, "password", "p", "", "The password to encrypt the file with")
	encryptCmd.Flags().StringVarP(&expiresIn, "expires-in", "e", "1d", "The time after which the encrypted file will expire (format: 1y/2/mo/3w/4d/5h/6min)")
	encryptCmd.Flags().StringVarP(&serverAddress, "server-address", "s", "", "The address of the server to encrypt the file with")
	encryptCmd.Flags().StringVarP(&inputPath, "input-path", "i", "", "The path to the input file")
	encryptCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "The path to the output file")
	encryptCmd.Flags().BoolVarP(&overwrite, "overwrite", "w", false, "Overwrite the output file if it already exists")

	rootCmd.AddCommand(encryptCmd)
}

var password string
var expiresIn string
var serverAddress string
var inputPath string
var outputPath string
var overwrite bool

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file",
	Long:  `Encrypt contents of a given file, writing the output to another specified file.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := commands.CreateEncryptInput(inputPath, outputPath, password, expiresIn, serverAddress, overwrite)
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