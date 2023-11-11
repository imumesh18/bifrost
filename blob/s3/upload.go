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

// Upload uploads the contents of the given reader to the specified S3 bucket and returns the uploaded object.
// The function takes a context and a variadic list of UploadOption which can be used to configure the upload.
// The function returns a pointer to the uploaded object and an error if any.
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
		Bucket:             aws.String(c.Bucket),
		Key:                aws.String(c.Object),
		Body:               bytes.NewReader(body),
		Metadata:           c.Metadata,
		ContentType:        aws.String(c.ContentType),
		ContentDisposition: aws.String(c.ContentDisposition),
		ContentEncoding:    aws.String(c.ContentEncoding),
		ContentLanguage:    aws.String(c.ContentLanguage),
	})
	if err != nil {
		return nil, err
	}

	obj := &blob.Object{
		Body:               io.NopCloser(bytes.NewReader(body)),
		ContentType:        c.ContentType,
		Size:               int64(len(body)),
		Bucket:             c.Bucket,
		Name:               c.Object,
		Metadata:           c.Metadata,
		URL:                fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.Bucket, c.Object),
		ContentDisposition: c.ContentDisposition,
		ContentEncoding:    c.ContentEncoding,
		ContentLanguage:    c.ContentLanguage,
		LastModified:       time.Now(),
	}

	return obj, nil
}
