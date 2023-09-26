package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestValidate_Logger(t *testing.T) {
	config := LoggerConfig{
		Level: slog.LevelInfo,
	}

	tests := []struct {
		description string
		config      LoggerConfig
		changeFn    func(LoggerConfig) LoggerConfig
		wantErr     bool
	}{
		{
			description: "valid config",
			config:      config,
			changeFn:    func(lc LoggerConfig) LoggerConfig { return lc },
			wantErr:     false,
		},
		{
			description: "invalid level",
			config:      config,
			changeFn: func(lc LoggerConfig) LoggerConfig {
				lc.Level = 6
				return lc
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			config := tt.changeFn(tt.config)
			err := config.validate()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
