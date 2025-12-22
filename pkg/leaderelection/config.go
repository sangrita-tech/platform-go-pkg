package leaderelection

import "time"

type Config struct {
	LeaseName      string        `yaml:"leaseName" env:"LEASE_NAME" env-default:"app-leader"`
	LeaseNamespace string        `yaml:"leaseNamespace" env:"LEASE_NAMESPACE" env-default:"default"`
	Identity       string        `yaml:"identity" env:"IDENTITY"`
	LeaseDuration  time.Duration `yaml:"leaseDuration" env:"LEASE_DURATION" env-default:"60s"`
	RenewDeadline  time.Duration `yaml:"renewDeadline" env:"RENEW_DEADLINE" env-default:"20s"`
	RetryPeriod    time.Duration `yaml:"retryPeriod" env:"RETRY_PERIOD" env-default:"5s"`
}
