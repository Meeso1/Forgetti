package io

import (
	"os"
	"path/filepath"
	"Forgetti/models"
	"fmt"
	"encoding/json"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// TODO: Instead of reading everything to memory, implement encryption chunk by chunk
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, overwrite bool, data []byte) error {
	if FileExists(path) && !overwrite {
		return fmt.Errorf("file already exists: '%s'", path)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directories: '%s'", path)
	}

	return os.WriteFile(path, data, 0644)
}

// TODO: Using just JSON is stupid - write metadata and then raw bytes
func WriteMetadataToFile(path string, overwrite bool, data *models.FileContentWithMetadata) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return WriteFile(path, overwrite, jsonData)
}

func ReadMetadataFromFile(path string) (*models.FileContentWithMetadata, error) {
	jsonData, err := ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data models.FileContentWithMetadata
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &data, nil
}
