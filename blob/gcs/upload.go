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
	"net/http"

	"github.com/imumesh18/bifrost/blob"
)

// Upload uploads the contents of a reader to a new object in the specified bucket.
// It returns the newly created object or an error if the operation failed.
// The function takes a context and a variadic list of UploadOption.
// The UploadConfig is created from the options and the contents of the reader are read into memory.
// The function sets the content type of the object to the provided content type or detects it from the contents of the reader.
// It also sets the metadata, content encoding, content language, and content disposition of the object if provided.
// The function returns the newly created object with its metadata, size, and URL.
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
	if c.ContentEncoding != "" {
		wc.ContentEncoding = c.ContentEncoding
	}

	if c.ContentLanguage != "" {
		wc.ContentLanguage = c.ContentLanguage
	}

	if c.ContentDisposition != "" {
		wc.ContentDisposition = c.ContentDisposition
	}

	if _, err := io.Copy(wc, bytes.NewReader(body)); err != nil {
		return nil, err
	}

	if err := wc.Close(); err != nil {
		return nil, err
	}

	return &blob.Object{
		Body:               io.NopCloser(bytes.NewReader(body)),
		ContentType:        contentType,
		Size:               wc.Size,
		Bucket:             c.Bucket,
		Name:               c.Object,
		LastModified:       wc.Attrs().Updated,
		Metadata:           c.Metadata,
		URL:                "https://storage.googleapis.com/" + c.Bucket + "/" + c.Object,
		ContentDisposition: c.ContentDisposition,
		ContentEncoding:    c.ContentEncoding,
		ContentLanguage:    c.ContentLanguage,
	}, nil
}
