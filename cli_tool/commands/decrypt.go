package commands

import (
	"fmt"
	"Forgetti/io"
	"Forgetti/encryption"
	"Forgetti/interaction"
	"strings"
)

type DecryptInput struct {
	InputPath string
	OutputPath string
	Password string
	ServerAddress string
	Overwrite bool
}

func CreateDecryptInput(
	inputPath string, 
	outputPath string, 
	password string, 
	serverAddress string, 
	overwrite bool,
) (*DecryptInput, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}

	if !io.FileExists(inputPath) {
		return nil, fmt.Errorf("input file does not exist: '%s'", inputPath)
	}
	
	if outputPath == "" {
		if strings.HasSuffix(inputPath, ".forgetti") {
			outputPath = strings.TrimSuffix(inputPath, ".forgetti")
		} else {
			outputPath = inputPath + ".decrypted"
		}
	}

	if io.FileExists(outputPath) && !overwrite {
		return nil, fmt.Errorf("output file already exists: '%s'", outputPath)
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	return &DecryptInput{
		InputPath: inputPath,
		OutputPath: outputPath,
		Password: password,
		ServerAddress: serverAddress,
		Overwrite: overwrite,
	}, nil
}

// TODO: Print stuff + handle HTTP errors in a pretty way
func Decrypt(input DecryptInput) error {
	contentWithMetadata, err := io.ReadMetadataFromFile(input.InputPath)
	if err != nil {
		return err
	}

	serverAddress := input.ServerAddress
	if serverAddress == "" {
		serverAddress = contentWithMetadata.Metadata.ServerAddress
	}

	encryptedKeyHash, err := interaction.EncryptWithExistingKey(serverAddress, input.Password, &contentWithMetadata.Metadata)
	if err != nil {
		return err
	}

	key, err := encryption.CreateKey(input.Password, encryptedKeyHash)
	if err != nil {
		return err
	}

	decryptedContent, err := encryption.Decrypt(contentWithMetadata.FileContent, key)
	if err != nil {
		return err
	}

	if err := io.WriteFile(input.OutputPath, input.Overwrite, decryptedContent); err != nil {
		return err
	}

	return nil
}