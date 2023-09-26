package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger  LoggerConfig  `toml:"logger"`
	Server  ServerConfig  `toml:"server"`
	Metrics MetricsConfig `toml:"metrics"`
}

func NewConfig(path string) (Config, error) {
	var config Config

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return Config{}, fmt.Errorf("parsing error: %w", err)
	}

	if err := config.validate(); err != nil {
		return Config{}, fmt.Errorf("validation error: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if err := c.Logger.validate(); err != nil {
		return fmt.Errorf("invalid logger definition: %w", err)
	}
	if err := c.Server.validate(); err != nil {
		return fmt.Errorf("invalid server definition: %w", err)
	}
	return nil
}
