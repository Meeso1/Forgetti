package commands

import (
	"Forgetti/encryption"
	"Forgetti/interaction"
	"Forgetti/io"
	"Forgetti/models"
	"fmt"
	"time"
)

type EncryptInput struct {
	InputPath     string
	OutputPath    string
	Password      string
	Expiration    time.Time
	ServerAddress string
	Overwrite     bool
	LogLevel      LogLevel
}

func CreateEncryptInput(
	inputPath string,
	outputPath string,
	password string,
	expiresIn string,
	serverAddress string,
	overwrite bool,
	verbose bool,
	quiet bool,
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

	logLevel := LogLevelInfo
	if verbose {
		logLevel = LogLevelVerbose
	}
	if quiet {
		logLevel = LogLevelError
	}

	return &EncryptInput{
		InputPath:     inputPath,
		OutputPath:    outputPath,
		Password:      password,
		Expiration:    expiration,
		ServerAddress: serverAddress,
		Overwrite:     overwrite,
		LogLevel:      logLevel,
	}, nil
}

func parseExpiration(expiresIn string) (time.Time, error) {
	// TODO: Implement
	return time.Now().Add(time.Duration(1) * time.Hour), nil
}

func Encrypt(input EncryptInput) error {
	logger := MakeLogger(input.LogLevel)

	logger.Verbose("Reading input file '%s'", input.InputPath)
	content, err := io.ReadFile(input.InputPath)
	if err != nil {
		return err
	}
	logger.Verbose("Read %d bytes from input file", len(content))

	logger.Verbose("Creating remote key, using server '%s' and expiration '%s'", input.ServerAddress, input.Expiration.String())
	interactionResult, err := interaction.GenerateKeyAndEncrypt(input.ServerAddress, input.Password, input.Expiration)
	if err != nil {
		return err
	}
	logger.Verbose("Created remote key '%s' with expiration '%s'", interactionResult.Metadata.KeyId, interactionResult.Metadata.Expiration.String())

	logger.Verbose("Creating symmetric key")
	key, err := encryption.CreateKey(input.Password, interactionResult.EncryptedKeyHash)
	if err != nil {
		return err
	}
	logger.Verbose("Created symmetric key")

	logger.Verbose("Encrypting content")
	encryptedContent, err := encryption.Encrypt(content, key)
	if err != nil {
		return err
	}
	logger.Verbose("Encrypted content")

	contentWithMetadata := models.FileContentWithMetadata{
		FileContent: encryptedContent,
		Metadata:    interactionResult.Metadata,
	}

	logger.Verbose("Writing encnrypted content to file '%s' (%d bytes, overwrite: %t)", input.OutputPath, len(encryptedContent), input.Overwrite)
	if err := io.WriteContentWithMetadataToFile(input.OutputPath, input.Overwrite, &contentWithMetadata); err != nil {
		return err
	}

	logger.Info("Output:         %s (%d bytes)", input.OutputPath, len(encryptedContent))
	logger.Info("Key ID:         %s", interactionResult.Metadata.KeyId)
	logger.Info("Expires at:     %s (in %s)", interactionResult.Metadata.Expiration.String(), time.Until(interactionResult.Metadata.Expiration).String())
	logger.Info("Server Address: %s", interactionResult.Metadata.ServerAddress)

	return nil
}
