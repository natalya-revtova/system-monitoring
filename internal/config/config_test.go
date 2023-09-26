package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		description string
		path        string
		want        Config
		wantErr     bool
	}{
		{
			description: "correct config path & content",
			path:        "./testdata/valid_config.toml",
			want: Config{
				Server: ServerConfig{
					Host: "127.0.0.1",
					Port: 50051,
				},
				Logger: LoggerConfig{
					Level: slog.LevelInfo,
				},
				Metrics: MetricsConfig{
					LoadAvg:  true,
					CPUUsg:   true,
					DiskInfo: true,
				},
			},
			wantErr: false,
		},
		{
			description: "invalid path",
			path:        "./invalid",
			want:        Config{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got, err := NewConfig(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
