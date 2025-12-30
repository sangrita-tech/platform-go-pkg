package kube

import (
	"errors"
	"time"
)

type Config struct {
	KubeConfigPath string        `yaml:"kubeConfigPath" env:"KUBE_CONFIG_PATH"`
	KubeContext    string        `yaml:"kubeContext" env:"KUBE_CONTEXT"`
	QPS            float32       `yaml:"qps" env:"QPS" env-default:"20"`
	Burst          int           `yaml:"burst" env:"BURST" env-default:"40"`
	Timeout        time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"30s"`
}

func (c Config) validate() error {
	if c.QPS < 0 {
		return errors.New("qps must not be negative")
	}

	if c.Burst < 0 {
		return errors.New("burst must not be negative")
	}

	if c.Timeout < 0 {
		return errors.New("timeout must not be negative")
	}

	return nil
}
