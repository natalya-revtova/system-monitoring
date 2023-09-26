package config

import (
	"errors"

	"golang.org/x/exp/slog"
)

type LoggerConfig struct {
	Level slog.Level `toml:"level"`
}

func (lc LoggerConfig) validate() error {
	switch lc.Level {
	case slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError:
		return nil
	}
	return errors.New("invalid level field")
}
