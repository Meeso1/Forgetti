package commands

import (
	"time"
	"fmt"
	"Forgetti/io"
	"Forgetti/encryption"
	"Forgetti/interaction"
	"Forgetti/models"
)

type EncryptInput struct {
	InputPath string
	OutputPath string
	Password string
	Expiration time.Time
	ServerAddress string
	Overwrite bool
}

func CreateEncryptInput(
	inputPath string, 
	outputPath string, 
	password string, 
	expiresIn string, 
	serverAddress string, 
	overwrite bool,
) (*EncryptInput, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}

	if !io.FileExists(inputPath) {
		return nil, fmt.Errorf("input file does not exist: '%s'", inputPath)
	}
	
	if outputPath == "" {
		outputPath = inputPath + ".forgetti"
	}

	if io.FileExists(outputPath) && !overwrite {
		return nil, fmt.Errorf("output file already exists: '%s'", outputPath)
	}
	
	expiration, err := parseExpiration(expiresIn)
	if err != nil {
		return nil, err
	}

	if password == "" {
		// TODO: Generate password + allow interactive input
		return nil, fmt.Errorf("password is required")
	}

	if serverAddress == "" {
		// TODO: Get default from config and env
		return nil, fmt.Errorf("server address is required")
	}
	
	return &EncryptInput{
		InputPath: inputPath,
		OutputPath: outputPath,
		Password: password,
		Expiration: expiration,
		ServerAddress: serverAddress,
		Overwrite: overwrite,
	}, nil
}

func parseExpiration(expiresIn string) (time.Time, error) {
	// TODO: Implement
	return time.Now().Add(time.Duration(1) * time.Hour), nil
}

// TODO: Print something sometimes
func Encrypt(input EncryptInput) error {
	content, err := io.ReadFile(input.InputPath)
	if err != nil {
		return err
	}

	interactionResult, err := interaction.GenerateKeyAndEncrypt(input.ServerAddress, input.Password, input.Expiration)
	if err != nil {
		return err
	}

	key, err := encryption.CreateKey(input.Password, interactionResult.EncryptedKeyHash)
	if err != nil {
		return err
	}

	encryptedContent, err := encryption.Encrypt(content, key)
	if err != nil {
		return err
	}

	contentWithMetadata := models.FileContentWithMetadata{
		FileContent: encryptedContent,
		Metadata: interactionResult.Metadata,
	}
	
	if err := io.WriteMetadataToFile(input.OutputPath, input.Overwrite, &contentWithMetadata); err != nil {
		return err
	}

	return nil
}
