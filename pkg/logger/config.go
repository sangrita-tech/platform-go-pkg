package logger

import (
	"fmt"
	"strings"
)

type Config struct {
	Level      string            `yaml:"level" env:"LEVEL" env-default:"info"`
	Format     string            `json:"format" yaml:"format" env:"FORMAT" env-default:"json"`
	DevMode    bool              `yaml:"devMode" env:"DEV_MODE" env-default:"false"`
	BaseFields map[string]string `yaml:"baseFields" env:"BASE_FIELDS"`
}

func (c *Config) validate() error {
	switch strings.ToLower(c.Level) {
	case "debug", "info", "warn", "warning", "error":
	default:
		return fmt.Errorf("unknown level %q", c.Level)
	}

	switch strings.ToLower(c.Format) {
	case "json", "console":
	default:
		return fmt.Errorf("unknown format %q", c.Format)
	}

	return nil
}
