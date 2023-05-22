package gcs

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/imumesh18/bifrost/blob"
)

var _ blob.Storage = (*Storage)(nil)

type Config struct {
	projectID       string
	credentialsFile string
	credentialsJSON []byte
}

type Storage struct {
	client *storage.Client
}

func New(ctx context.Context, opts ...Option) (*Storage, error) {
	c, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(c.credentialsJSON), option.WithCredentialsFile(c.credentialsFile))
	if err != nil {
		return nil, err
	}

	return &Storage{client: client}, nil
}
