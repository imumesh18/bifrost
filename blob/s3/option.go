package s3

import (
	"errors"
)

type Option func(*Config)

func WithRegion(region string) Option {
	return func(c *Config) {
		c.region = region
	}
}

func WithAccessKeyID(accessKeyID string) Option {
	return func(c *Config) {
		c.accessKeyID = accessKeyID
	}
}

func WithSecretAccessKey(secretAccessKey string) Option {
	return func(c *Config) {
		c.secretAccessKey = secretAccessKey
	}
}

func WithURL(url string) Option {
	return func(c *Config) {
		c.url = url
	}
}

func newConfig(opts ...Option) (Config, error) {
	var c Config
	for _, opt := range opts {
		opt(&c)
	}
	if err := c.validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

func (c *Config) validate() error {
	if c.region == "" {
		return errors.New("region is required")
	}
	if c.accessKeyID == "" || c.secretAccessKey == "" {
		return errors.New("accessKeyID and secretAccessKey are required")
	}
	return nil
}
