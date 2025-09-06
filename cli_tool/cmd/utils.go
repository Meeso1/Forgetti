package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

	"golang.org/x/term"
)

const passwordEnv = "FORGETTI_PASSWORD"
const generatedPasswordLength = 16

func promptForPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	return string(password), nil
}

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

func promptForChoice() (bool, error) {
	fmt.Print("Do you want to (p)rovide a password or (g)enerate a random one? [p/g]: ")

	var choice string
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return false, fmt.Errorf("failed to read choice: %v", err)
	}

	choice = strings.ToLower(strings.TrimSpace(choice))
	switch choice {
	case "p", "provide":
		return false, nil // false means provide password
	case "g", "generate":
		return true, nil // true means generate password
	default:
		return false, fmt.Errorf("invalid choice '%s', please enter 'p' for provide or 'g' for generate", choice)
	}
}

func getFomEnv(password *string) {
	if envPassword := os.Getenv(passwordEnv); envPassword != "" {
		fmt.Printf("Using password from environment variable %s\n", passwordEnv)
		*password = envPassword
	}
}

func promptForEncryptPasswordIfEmpty(password *string, nonInteractive bool) error {
	if *password != "" {
		return nil
	}

	if getFomEnv(password); *password != "" {
		return nil
	}

	// Ask user what they want to do - assume generation if non-interactive
	var generateRandom bool
	var err error
	if nonInteractive {
		generateRandom = true
	} else {	
		generateRandom, err = promptForChoice()
		if err != nil {
			return err
		}
	}

	if generateRandom {
		randomPassword, err := generateRandomPassword(generatedPasswordLength)
		if err != nil {
			return err
		}
		*password = randomPassword
		fmt.Printf("Generated random password: %s\n", randomPassword)
	} else {
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

func promptForDecryptPasswordIfEmpty(password *string, nonInteractive bool) error {
	if *password != "" {
		return nil
	}

	if getFomEnv(password); *password != "" {
		return nil
	}

	if nonInteractive {
		return fmt.Errorf("password is not provided and non-interactive mode is enabled")
	}

	var err error
	*password, err = promptForPassword("Enter password: ")
	if err != nil {
		return err
	}

	if *password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}
