package gcs

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/imumesh18/bifrost/blob"
)

func (s *Storage) Upload(ctx context.Context, opts ...blob.UploadOption) (*blob.Object, error) {
	c := &blob.UploadConfig{}
	c.Apply(opts...)

	bucket := s.client.Bucket(c.Bucket)
	wc := bucket.Object(c.Object).NewWriter(ctx)
	body, err := io.ReadAll(c.Body)
	if err != nil {
		return nil, err
	}

	contentType := c.ContentType
	if c.ContentType == "" {
		contentType = http.DetectContentType(body)
	}
	wc.ContentType = contentType
	wc.Metadata = c.Metadata

	if _, err := io.Copy(wc, bytes.NewReader(body)); err != nil {
		return nil, err
	}

	if err := wc.Close(); err != nil {
		return nil, err
	}

	return &blob.Object{
		Body:         io.NopCloser(bytes.NewReader(body)),
		ContentType:  contentType,
		Size:         wc.Size,
		Bucket:       c.Bucket,
		Name:         c.Object,
		LastModified: wc.Attrs().Updated,
		Metadata:     c.Metadata,
		URL:          string("https://storage.googleapis.com/" + c.Bucket + "/" + c.Object),
	}, nil
}
