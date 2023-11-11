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

// Package s3 provides an implementation of the blob.Storage interface for Amazon S3.
package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/imumesh18/bifrost/blob"
)

var _ blob.Storage = (*Storage)(nil)

type Config struct {
	region          string
	accessKeyID     string
	secretAccessKey string
	endpoint        string
	enableDebug     bool
}

type Storage struct {
	client *s3.Client
}

func New(ctx context.Context, opts ...Option) (*Storage, error) {
	c, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	var clientLogMode aws.ClientLogMode
	if c.enableDebug {
		clientLogMode = aws.LogRetries | aws.LogRequest | aws.LogRequestWithBody | aws.LogResponse | aws.LogResponseWithBody | aws.LogDeprecatedUsage | aws.LogRequestEventMessage | aws.LogResponseEventMessage
	}
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: c.endpoint,
				}, nil
			})),
		config.WithClientLogMode(clientLogMode),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     c.accessKeyID,
				SecretAccessKey: c.secretAccessKey,
				Source:          "client provided credentials",
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("[s3.New] failed to load config from provided credentials: %w", err)
	}

	return &Storage{client: s3.NewFromConfig(cfg)}, nil
}
