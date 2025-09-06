package config

import (
	"Forgetti/io"
	"encoding/json"
	"fmt"
	"os"
)

const defaultConfigPath = ".config.json"
const configPathEnvVar = "FORGETTI_CONFIG_PATH"

type Config struct {
	ServerAddress string `json:"server_address"`
}

func GetConfigPath() string {
	envConfigPath := os.Getenv(configPathEnvVar)
	if envConfigPath != "" {
		return envConfigPath
	}

	return defaultConfigPath
}

func DoesConfigExist() bool {
	return io.FileExists(GetConfigPath())
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()
	if !io.FileExists(configPath) {
		return nil, fmt.Errorf("config file does not exist: '%s'", configPath)
	}

	content, err := io.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
