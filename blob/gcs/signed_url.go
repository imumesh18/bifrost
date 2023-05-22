package gcs

import (
	"context"

	"cloud.google.com/go/storage"

	"github.com/imumesh18/bifrost/blob"
)

func (s *Storage) GenerateSignedURL(ctx context.Context, opts ...blob.SignedURLOption) (string, error) {
	c := blob.SignedURLConfig{}
	c.Apply(opts...)

	u, err := s.client.Bucket(c.Bucket).SignedURL(c.Object, &storage.SignedURLOptions{
		Method:  c.Method,
		Expires: c.Expires,
		Scheme:  storage.SigningSchemeV4,
	})
	if err != nil {
		return "", err
	}

	return u, nil
}
