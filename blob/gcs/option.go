// Copyright (C) 2023 Umesh Yadav
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
