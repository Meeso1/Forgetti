package commands

import (
	"Forgetti/io"
	"fmt"
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
	logger := MakeLogger(LogLevelInfo)

	logger.Info("File: '%s'", input.InputPath)
	contentWithMetadata, err := io.ReadContentWithMetadataFromFile(input.InputPath)
	if err != nil {
		return err
	}

	logger.Info("%s", contentWithMetadata.String())

	return nil
}
