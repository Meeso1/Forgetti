package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

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

// generateRandomPassword creates a random password with specified length
func generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"
	password := make([]byte, length)

	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random password: %v", err)
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

// promptForChoice asks user to choose between entering password or generating random one
func promptForChoice() (bool, error) {
	fmt.Print("Do you want to (e)nter a password or (g)enerate a random one? [e/g]: ")

	var choice string
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return false, fmt.Errorf("failed to read choice: %v", err)
	}

	choice = strings.ToLower(strings.TrimSpace(choice))
	switch choice {
	case "e", "enter":
		return false, nil // false means enter password
	case "g", "generate":
		return true, nil // true means generate password
	default:
		return false, fmt.Errorf("invalid choice '%s', please enter 'e' for enter or 'g' for generate", choice)
	}
}

func promptForEncryptPasswordIfEmpty(password *string, nonInteractive bool) error {
	if *password == "" {
		if nonInteractive {
			// In non-interactive mode, always generate a random password
			randomPassword, err := generateRandomPassword(16)
			if err != nil {
				return err
			}
			*password = randomPassword
			fmt.Printf("Generated random password: %s\n", randomPassword)
			return nil
		}

		// Interactive mode: ask user what they want to do
		generateRandom, err := promptForChoice()
		if err != nil {
			return err
		}

		if generateRandom {
			randomPassword, err := generateRandomPassword(16)
			if err != nil {
				return err
			}
			*password = randomPassword
			fmt.Printf("Generated random password: %s\n", randomPassword)
		} else {
			// User wants to enter password manually
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
	}
	return nil
}

// promptForDecryptPasswordIfEmpty handles password prompting for decryption (no generation option)
func promptForDecryptPasswordIfEmpty(password *string, nonInteractive bool) error {
	if *password == "" {
		if nonInteractive {
			return fmt.Errorf("password is required for decryption but non-interactive mode is enabled")
		}

		// For decryption, only prompt for password entry (no generation option)
		var err error
		*password, err = promptForPassword("Enter password: ")
		if err != nil {
			return err
		}

		if *password == "" {
			return fmt.Errorf("password cannot be empty")
		}
	}
	return nil
}
