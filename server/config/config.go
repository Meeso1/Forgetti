package config

import (
	"fmt"
	"forgetti-common/io"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Host string `json:"host" env:"SERVER_HOST" env-default:"localhost" validate:"required"`
		Port int    `json:"port" env:"SERVER_PORT" env-default:"8080" validate:"min=1,max=65535"`
		Mode string `json:"mode" env:"GIN_MODE" env-default:"release" validate:"oneof=debug release test"`
	} `json:"server"`

	KeyStore struct {
		RecentlyExpiredDurationHours int `json:"recently_expired_duration" env:"KEYSTORE_RECENTLY_EXPIRED_DURATION" env-default:"24" validate:"min=1,max=168"`
	} `json:"keystore"`

	Database struct {
		Path            string `json:"path" env:"DB_PATH" env-default:"./forgetti.db" validate:"required"`
		MaxOpenConns    int    `json:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"25" validate:"min=1,max=100"`
		MaxIdleConns    int    `json:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"5" validate:"min=1,max=25"`
		ConnMaxLifetime int    `json:"conn_max_lifetime_minutes" env:"DB_CONN_MAX_LIFETIME" env-default:"60" validate:"min=1,max=1440"`
	} `json:"database"`

	Logging struct {
		Level  string `json:"level" env:"LOG_LEVEL" env-default:"info" validate:"oneof=debug info warn error"`
		Format string `json:"format" env:"LOG_FORMAT" env-default:"json" validate:"oneof=json text"`
	} `json:"logging"`
}

func (c *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadFromFile(filename string) (*Config, error) {
	path, err := io.GetRelativePathFromBin(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path from bin: %w", err)
	}

	if !io.FileExists(path) {
		return Load()
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from file %s: %w", filename, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
