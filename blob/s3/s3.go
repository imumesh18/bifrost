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
	url             string
}

type Storage struct {
	client *s3.Client
}

func New(ctx context.Context, opts ...Option) (*Storage, error) {
	c, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(c.region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: c.url,
			}, nil
		})),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     c.accessKeyID,
				SecretAccessKey: c.secretAccessKey,
				Source:          "client provided credentials",
			},
		}))
	if err != nil {
		return nil, fmt.Errorf("[s3.New] failed to load config from provided credentials: %w", err)
	}

	return &Storage{client: s3.NewFromConfig(cfg)}, nil
}
