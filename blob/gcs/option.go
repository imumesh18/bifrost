package gcs

import (
	"errors"
)

type Option func(*Config)

func WithProjectID(projectID string) Option {
	return func(c *Config) {
		c.projectID = projectID
	}
}

func WithCredentialsFile(credentialsFile string) Option {
	return func(c *Config) {
		c.credentialsFile = credentialsFile
	}
}

func WithCredentialsJSON(credentialsJSON []byte) Option {
	return func(c *Config) {
		c.credentialsJSON = credentialsJSON
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
	if c.projectID == "" {
		return errors.New("projectID is required")
	}
	if c.credentialsFile == "" && len(c.credentialsJSON) == 0 {
		return errors.New("credentialsFile or credentialsJSON is required")
	}

	if c.credentialsFile != "" && len(c.credentialsJSON) != 0 {
		return errors.New("only one of credentialsFile or credentialsJSON can be specified")
	}

	return nil
}
