package leaderelection

import (
	"errors"
	"fmt"
	"time"
)

type Config struct {
	LeaseName      string        `yaml:"leaseName" env:"LEASE_NAME" env-default:"app-leader"`
	LeaseNamespace string        `yaml:"leaseNamespace" env:"LEASE_NAMESPACE" env-default:"default"`
	Identity       string        `yaml:"identity" env:"IDENTITY"`
	LeaseDuration  time.Duration `yaml:"leaseDuration" env:"LEASE_DURATION" env-default:"60s"`
	RenewDeadline  time.Duration `yaml:"renewDeadline" env:"RENEW_DEADLINE" env-default:"20s"`
	RetryPeriod    time.Duration `yaml:"retryPeriod" env:"RETRY_PERIOD" env-default:"5s"`
}

func (c Config) validate() error {
	if c.LeaseName == "" {
		return errors.New("lease name must be set")
	}

	if c.LeaseNamespace == "" {
		return errors.New("lease namespace must be set")
	}

	if c.LeaseDuration <= 0 {
		return errors.New("lease duration must be > 0")
	}

	if c.RenewDeadline <= 0 {
		return errors.New("renew deadline must be > 0")
	}

	if c.RetryPeriod <= 0 {
		return errors.New("retry period must be > 0")
	}

	if !(c.LeaseDuration > c.RenewDeadline && c.RenewDeadline > c.RetryPeriod) {
		return fmt.Errorf(
			"expected LeaseDuration(%s) > RenewDeadline(%s) > RetryPeriod(%s)",
			c.LeaseDuration, c.RenewDeadline, c.RetryPeriod,
		)
	}

	return nil
}
