package backoff

import (
	"math/rand"
	"time"
)

// Backoff provides an interface for different backoff implementations
type Backoff interface {
	// Backoff returns the duration to wait before retrying the operation,
	// or backoff.
	//
	// Example usage:
	//
	// 	for retries := 0; retries < 10; retries++{
	// 		// your code
	// 		if success {
	// 			break
	// 		}
	// 		time.Sleep(exp.Backoff(retries))
	// 	}
	Backoff(retries int) time.Duration
}

// DefaultConfig is the default config for backoff.
var DefaultConfig = Config{
	BaseDelay:  1.0 * time.Second,
	Multiplier: 1.6,
	Jitter:     0.2,
	MaxDelay:   120 * time.Second,
}

// Config defines the configuration options for backoff.
// All unfilled field will be replaced by the default config values.
type Config struct {
	// BaseDelay is the amount of time to backoff after the first failure.
	BaseDelay time.Duration
	// Multiplier is the factor with which to multiply backoffs after a
	// failed retry. Should ideally be greater than 1.
	Multiplier float64
	// Jitter is the factor with which backoffs are randomized.
	Jitter float64
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration
}

func fillDefaultConfig(cfg *Config) *Config {
	if cfg == nil {
		cloned := DefaultConfig
		return &cloned
	}
	if cfg.BaseDelay == 0 {
		cfg.BaseDelay = DefaultConfig.BaseDelay
	}
	if cfg.Jitter == 0 {
		cfg.Jitter = DefaultConfig.Jitter
	}
	if cfg.MaxDelay == 0 {
		cfg.MaxDelay = DefaultConfig.MaxDelay
	}
	if cfg.Multiplier == 0 {
		cfg.Multiplier = DefaultConfig.Multiplier
	}
	return cfg
}

// DefaultExponential is the default exponential backoff with default config
var DefaultExponential = &Exponential{Config: &DefaultConfig}

// Exponential implements exponential backoff algorithm.
type Exponential struct {
	// Config contains all options to configure the backoff algorithm.
	Config *Config
}

// NewExponential creates a exponential backoff algorithm based on given config
func NewExponential(cfg *Config) Backoff {
	cfg = fillDefaultConfig(cfg)
	return &Exponential{Config: &DefaultConfig}
}

// Backoff returns the amount of time to wait before the next retry given the
// number of retries.
func (bc *Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return bc.Config.BaseDelay
	}
	backoff, max := float64(bc.Config.BaseDelay), float64(bc.Config.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.Config.Multiplier
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so that if a cluster of requests start at
	// the same time, they won't operate in lockstep.
	backoff *= 1 + bc.Config.Jitter*(rand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}
