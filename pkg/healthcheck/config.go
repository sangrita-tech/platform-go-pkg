package healthcheck

import (
	"errors"
	"time"
)

type Config struct {
	Addr            string        `yaml:"addr" env:"ADDR" env-default:":8080"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"5s"`
}

func (c *Config) validate() error {
	if c.Addr == "" {
		return errors.New("addr must not be empty")
	}

	if c.ShutdownTimeout <= 0 {
		return errors.New("shutdown timeout must be positive")
	}

	return nil
}
