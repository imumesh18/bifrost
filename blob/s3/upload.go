package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/imumesh18/bifrost/blob"
)

func (s *Storage) Upload(ctx context.Context, opts ...blob.UploadOption) (*blob.Object, error) {
	c := &blob.UploadConfig{}
	c.Apply(opts...)

	body, err := io.ReadAll(c.Body)
	if err != nil {
		return nil, err
	}

	if c.ContentType == "" {
		c.ContentType = http.DetectContentType(body)
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.Bucket),
		Key:         aws.String(c.Object),
		Body:        bytes.NewReader(body),
		Metadata:    c.Metadata,
		ContentType: aws.String(c.ContentType),
	})
	if err != nil {
		return nil, err
	}

	obj := &blob.Object{
		Body:         io.NopCloser(bytes.NewReader(body)),
		ContentType:  c.ContentType,
		Size:         int64(len(body)),
		Bucket:       c.Bucket,
		Name:         c.Object,
		LastModified: time.Now(),
		Metadata:     c.Metadata,
		URL:          fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.Bucket, c.Object),
	}

	return obj, nil
}
