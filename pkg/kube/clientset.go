package kube

import (
	"errors"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

func New(cfg *Config) (*kubernetes.Clientset, error) {
	if cfg == nil {
		return nil, errors.New("kube -> config is nil")
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("kube -> failed to validate config -> %w", err)
	}

	restCfg, err := buildRESTConfig(cfg)
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("kube -> create clientset -> %w", err)
	}

	return cs, nil
}
