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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/imumesh18/bifrost/blob"
)

// GenerateSignedURL generates a signed URL for the given object in the given bucket.
// The URL can be used to perform the specified HTTP method on the object until the expiration time.
// The expiration time is determined by the Expires option passed in the opts parameter.
// The supported HTTP methods are GET, PUT, and HEAD.
// Returns the signed URL and any error encountered.
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
