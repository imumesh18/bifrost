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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/imumesh18/bifrost/blob"
)

// Download downloads an object from the specified S3 bucket and returns a blob.Object.
// It takes a context.Context and a variadic list of blob.DownloadOption as input.
// It returns a pointer to blob.Object and an error.
// The function applies the provided options to a blob.DownloadConfig and uses it to get the object from S3.
// If the object is successfully retrieved, it is returned as a blob.Object.
// If there is an error while retrieving the object, the function returns nil and the error.
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
