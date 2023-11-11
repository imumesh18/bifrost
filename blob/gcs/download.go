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

package gcs

import (
	"bytes"
	"context"
	"io"

	"github.com/imumesh18/bifrost/blob"
)

// Download downloads an object from Google Cloud Storage and returns a blob.Object.
// It takes a context.Context and a variadic slice of blob.DownloadOption as input.
// It returns a pointer to a blob.Object and an error.
// The function reads the object from the specified bucket and object name.
// It returns an error if the object does not exist or if there is an error reading the object.
// The returned blob.Object contains the object's body, content type, size, bucket name, object name, last modified time, and URL.
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
		Body:            io.NopCloser(bytes.NewReader(body)),
		ContentType:     rc.Attrs.ContentType,
		Size:            rc.Attrs.Size,
		Bucket:          c.Bucket,
		Name:            c.Object,
		LastModified:    rc.Attrs.LastModified,
		URL:             "https://storage.googleapis.com/" + c.Bucket + "/" + c.Object,
		ContentEncoding: rc.Attrs.ContentEncoding,
	}, nil
}
