package leaderelection

import (
	"errors"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
)

type Elector struct {
	cfg       *Config
	cb        Callbacks
	identity  string
	clientset *kubernetes.Clientset
}

func New(cfg *Config, cb Callbacks, clientset *kubernetes.Clientset) (*Elector, error) {
	if cfg == nil {
		return nil, errors.New("leaderelection -> config is nil")
	}

	if clientset == nil {
		return nil, errors.New("leaderelection -> clientset is nil")
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("leaderelection -> failed to validate config -> %w", err)
	}

	id := cfg.Identity
	if id == "" {
		if h, err := os.Hostname(); err == nil && h != "" {
			id = h
		} else {
			id = "unknown"
		}
	}

	return &Elector{
		cfg:       cfg,
		cb:        cb,
		identity:  id,
		clientset: clientset,
	}, nil
}

func (e *Elector) Identity() string {
	return e.identity
}
