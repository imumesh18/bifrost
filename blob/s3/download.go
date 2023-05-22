package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/imumesh18/bifrost/blob"
)

// Download downloads an object from the s3 bucket.
func (s *Storage) Download(ctx context.Context, opts ...blob.DownloadOption) (*blob.Object, error) {
	c := &blob.DownloadConfig{}
	c.Apply(opts...)

	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(c.Object),
	})
	if err != nil {
		return nil, err
	}

	return &blob.Object{
		Body:         resp.Body,
		ContentType:  aws.ToString(resp.ContentType),
		Metadata:     resp.Metadata,
		Size:         resp.ContentLength,
		LastModified: aws.ToTime(resp.LastModified),
		Bucket:       c.Bucket,
		Name:         c.Object,
	}, nil
}
