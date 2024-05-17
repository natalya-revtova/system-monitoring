package config

import (
	"errors"
)

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func (sc ServerConfig) validate() error {
	if len(sc.Host) == 0 {
		return errors.New("invalid host field")
	}
	if sc.Port <= 0 || sc.Port > 65535 {
		return errors.New("invalid port field")
	}
	return nil
}
