package gcs

import (
	"bytes"
	"context"
	"io"

	"github.com/imumesh18/bifrost/blob"
)

func (s *Storage) Download(ctx context.Context, opts ...blob.DownloadOption) (*blob.Object, error) {
	c := blob.DownloadConfig{}
	c.Apply(opts...)

	bucket := s.client.Bucket(c.Bucket)
	rc, err := bucket.Object(c.Object).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return &blob.Object{
		Body:         io.NopCloser(bytes.NewReader(body)),
		ContentType:  rc.Attrs.ContentType,
		Size:         rc.Attrs.Size,
		Bucket:       c.Bucket,
		Name:         c.Object,
		LastModified: rc.Attrs.LastModified,
		URL:          string("https://storage.googleapis.com/" + c.Bucket + "/" + c.Object),
	}, nil
}
