package commands

import (
	"Forgetti/io"
	"fmt"
	"forgetti-common/logging"
)

type ReadMetadataInput struct {
	InputPath string
}

func CreateReadMetadataInput(path string) (*ReadMetadataInput, error) {
	if path == "" {
		return nil, fmt.Errorf("input path is required")
	}

	return &ReadMetadataInput{
		InputPath: path,
	}, nil
}

func ReadMetadata(input ReadMetadataInput) error {
	logging.SetGlobalConfig(logging.Config{
		LogLevel: logging.LogLevelInfo,
		LogFile:  "", // CLI tool logs only to console
	})
	logger := logging.MakeLogger("read_metadata")

	logger.Info("File: '%s'", input.InputPath)
	contentWithMetadata, err := io.ReadContentWithMetadataFromFile(input.InputPath)
	if err != nil {
		return err
	}

	logger.Info("%s", contentWithMetadata.String())

	return nil
}
