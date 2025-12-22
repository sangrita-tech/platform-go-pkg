package leaderelection

import (
	"errors"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
)

type Elector struct {
	cfg       Config
	cb        Callbacks
	identity  string
	clientset *kubernetes.Clientset
}

func New(cfg Config, cb Callbacks, clientset *kubernetes.Clientset) (*Elector, error) {
	if clientset == nil {
		return nil, errors.New("leaderelection: clientset is nil")
	}

	if cfg.LeaseName == "" {
		cfg.LeaseName = "app-leader"
	}
	if cfg.LeaseNamespace == "" {
		cfg.LeaseNamespace = "default"
	}

	id := cfg.Identity
	if id == "" {
		if h, err := os.Hostname(); err == nil && h != "" {
			id = h
		} else {
			id = "unknown"
		}
	}

	if cfg.LeaseName == "" || cfg.LeaseNamespace == "" {
		return nil, errors.New("leaderelection: lease name/namespace must be set")
	}
	if cfg.LeaseDuration <= 0 || cfg.RenewDeadline <= 0 || cfg.RetryPeriod <= 0 {
		return nil, errors.New("leaderelection: durations must be > 0")
	}
	if !(cfg.LeaseDuration > cfg.RenewDeadline && cfg.RenewDeadline > cfg.RetryPeriod) {
		return nil, fmt.Errorf(
			"leaderelection: expected LeaseDuration(%s) > RenewDeadline(%s) > RetryPeriod(%s)",
			cfg.LeaseDuration, cfg.RenewDeadline, cfg.RetryPeriod,
		)
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
