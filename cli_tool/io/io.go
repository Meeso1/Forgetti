package io

import (
	"Forgetti/models"
	"forgetti-common/io"
	"encoding/json"
	"fmt"
)

func FileExists(path string) bool {
	return io.FileExists(path)
}

func ReadFile(path string) ([]byte, error) {
	return io.ReadFile(path)
}

func WriteFile(path string, overwrite bool, data []byte) error {
	return io.WriteFile(path, overwrite, data)
}

// TODO: Using just JSON is stupid - write metadata and then raw bytes
func WriteContentWithMetadataToFile(path string, overwrite bool, data *models.FileContentWithMetadata) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return io.WriteFile(path, overwrite, jsonData)
}

func ReadContentWithMetadataFromFile(path string) (*models.FileContentWithMetadata, error) {
	jsonData, err := io.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data models.FileContentWithMetadata
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &data, nil
}
