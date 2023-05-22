package s3

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/imumesh18/bifrost/blob"
)

func (s *Storage) GenerateSignedURL(ctx context.Context, opts ...blob.SignedURLOption) (string, error) {
	c := blob.SignedURLConfig{}
	c.Apply(opts...)

	var err error
	var presignResult *v4.PresignedHTTPRequest
	presignClient := s3.NewPresignClient(s.client)
	expirationDuration := time.Until(c.Expires)
	switch strings.ToUpper(c.Method) {
	case "GET":
		presignResult, err = presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(c.Object),
		}, s3.WithPresignExpires(expirationDuration))
	case "PUT":
		presignResult, err = presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(c.Object),
		}, s3.WithPresignExpires(expirationDuration))
	case "HEAD":
		presignResult, err = presignClient.PresignHeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(c.Object),
		}, s3.WithPresignExpires(expirationDuration))
	default:
		return "", fmt.Errorf("unsupported method: %s", c.Method)
	}

	return presignResult.URL, err
}
