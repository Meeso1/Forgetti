package io

import (
	"fmt"
	"os"
	"path/filepath"
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

func GetRelativePathFromBin(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	
	binDir := filepath.Dir(execPath)
	return filepath.Join(binDir, path), nil
}
