package commands

import (
	"Forgetti/encryption"
	"Forgetti/interaction"
	"Forgetti/io"
	"fmt"
	"strings"
	"time"
)

type DecryptInput struct {
	InputPath     string
	OutputPath    string
	Password      string
	ServerAddress string
	Overwrite     bool
	LogLevel      LogLevel
}

func CreateDecryptInput(
	inputPath string,
	outputPath string,
	password string,
	serverAddress string,
	overwrite bool,
	verbose bool,
	quiet bool,
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

	logLevel := LogLevelInfo
	if verbose {
		logLevel = LogLevelVerbose
	}
	if quiet {
		logLevel = LogLevelError
	}

	return &DecryptInput{
		InputPath:     inputPath,
		OutputPath:    outputPath,
		Password:      password,
		ServerAddress: serverAddress,
		Overwrite:     overwrite,
		LogLevel:      logLevel,
	}, nil
}

func Decrypt(input DecryptInput) error {
	logger := MakeLogger(input.LogLevel)

	logger.Verbose("Reading file '%s'", input.InputPath)
	contentWithMetadata, err := io.ReadContentWithMetadataFromFile(input.InputPath)
	if err != nil {
		return err
	}
	logger.Verbose("Read encrypted content from file")
	logger.Info("\n%s", contentWithMetadata.String())

	if contentWithMetadata.Metadata.Expiration.Before(time.Now()) {
		return fmt.Errorf("key has expired at %s (%s ago)", contentWithMetadata.Metadata.Expiration.String(), time.Since(contentWithMetadata.Metadata.Expiration).String())
	}

	serverAddress := input.ServerAddress
	if serverAddress == "" {
		logger.Verbose("Server address not provided, using server address from metadata: '%s'", contentWithMetadata.Metadata.ServerAddress)
		serverAddress = contentWithMetadata.Metadata.ServerAddress
	}

	logger.Verbose("Getting remote key '%s', using server '%s'", contentWithMetadata.Metadata.KeyId, serverAddress)
	encryptedKeyHash, err := interaction.EncryptWithExistingKey(serverAddress, input.Password, &contentWithMetadata.Metadata)
	if err != nil {
		return err
	}
	logger.Verbose("Got remote key '%s'", contentWithMetadata.Metadata.KeyId)

	logger.Verbose("Creating symmetric key")
	key, err := encryption.CreateKey(input.Password, encryptedKeyHash)
	if err != nil {
		return err
	}
	logger.Verbose("Created symmetric key")

	logger.Verbose("Decrypting content")
	decryptedContent, err := encryption.Decrypt(contentWithMetadata.FileContent, key)
	if err != nil {
		return err
	}
	logger.Verbose("Decrypted content")

	logger.Verbose("Writing content to file '%s' (%d bytes, overwrite: %t)", input.OutputPath, len(decryptedContent), input.Overwrite)
	if err := io.WriteFile(input.OutputPath, input.Overwrite, decryptedContent); err != nil {
		return err
	}

	logger.Info("Output: '%s' (%d bytes)", input.OutputPath, len(decryptedContent))

	return nil
}
