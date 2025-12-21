package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load[T any](path string) (T, error) {
	var cfg T

	if path == "" {
		return cfg, cleanenv.ReadEnv(&cfg)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			var envCfg T
			return envCfg, cleanenv.ReadEnv(&envCfg)
		}
		var zero T
		return zero, err
	}

	return cfg, cleanenv.ReadEnv(&cfg)
}
