package cmd

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func promptForPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	return string(password), nil
}

func promptForPasswordIfEmpty(password *string) error {
	if *password == "" {
		var err error
		*password, err = promptForPassword("Enter password: ")
		if err != nil {
			return err
		}

		confirmPassword, err := promptForPassword("Confirm password: ")
		if err != nil {
			return err
		}

		if *password != confirmPassword {
			return fmt.Errorf("passwords do not match")
		}
	}
	return nil
}
