package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_ServerGRPC(t *testing.T) {
	config := ServerConfig{
		Host: "127.0.0.1",
		Port: 50051,
	}

	tests := []struct {
		description string
		config      ServerConfig
		changeFn    func(ServerConfig) ServerConfig
		wantErr     bool
	}{
		{
			description: "valid config",
			config:      config,
			changeFn:    func(sc ServerConfig) ServerConfig { return sc },
			wantErr:     false,
		},
		{
			description: "invalid host",
			config:      config,
			changeFn: func(sc ServerConfig) ServerConfig {
				sc.Host = ""
				return sc
			},
			wantErr: true,
		},
		{
			description: "invalid port",
			config:      config,
			changeFn: func(sc ServerConfig) ServerConfig {
				sc.Port = 11000000
				return sc
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
