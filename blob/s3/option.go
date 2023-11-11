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

func WithEndpoint(endpoint string) Option {
	return func(c *Config) {
		c.endpoint = endpoint
	}
}

func WithEnableDebug(enableDebug bool) Option {
	return func(c *Config) {
		c.enableDebug = enableDebug
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
