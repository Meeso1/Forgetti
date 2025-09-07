package io

import (
	"Forgetti/models"
	"encoding/json"
	"fmt"
	"forgetti-common/io"
)

const delimiterByte = 0x00 // Must be invalid in JSON

func FileExists(path string) bool {
	return io.FileExists(path)
}

func ReadFile(path string) ([]byte, error) {
	return io.ReadFile(path)
}

func WriteFile(path string, overwrite bool, data []byte) error {
	return io.WriteFile(path, overwrite, data)
}

func GetRelativePathFromBin(path string) (string, error) {
	return io.GetRelativePathFromBin(path)
}

func WriteContentWithMetadataToFile(path string, overwrite bool, data *models.FileContentWithMetadata) error {
	// Marshal only the metadata to JSON
	metadataJson, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Create the final file content: metadata JSON + null byte delimiter + raw encrypted bytes
	// Using null byte (0x00) as delimiter since it cannot appear in valid JSON
	fileContent := append(metadataJson, delimiterByte)
	fileContent = append(fileContent, data.FileContent...)

	return io.WriteFile(path, overwrite, fileContent)
}

func ReadContentWithMetadataFromFile(path string) (*models.FileContentWithMetadata, error) {
	fileData, err := io.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Find the first null byte that separates metadata JSON from encrypted bytes
	delimiterIndex := -1
	for i, b := range fileData {
		if b == delimiterByte {
			delimiterIndex = i
			break
		}
	}

	if delimiterIndex == -1 {
		return nil, fmt.Errorf("invalid file format: no delimiter found between metadata and encrypted content")
	}

	// Parse metadata JSON from the beginning up to the delimiter
	metadataJson := fileData[:delimiterIndex]
	var metadata models.Metadata
	if err := json.Unmarshal(metadataJson, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	// Use the remaining bytes after the delimiter as encrypted content
	encryptedContent := fileData[delimiterIndex+1:]

	return &models.FileContentWithMetadata{
		FileContent: encryptedContent,
		Metadata:    metadata,
	}, nil
}
