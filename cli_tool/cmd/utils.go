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

func promptForPasswordIfEmpty(password *string, prompt string) error {
	if *password == "" {
		var err error
		*password, err = promptForPassword(prompt)
		if err != nil {
			return err
		}
	}
	return nil
}
