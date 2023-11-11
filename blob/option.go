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

package blob

import (
	"io"
	"time"
)

// SignedURLOption is a function type that takes a pointer to a SignedURLConfig struct as its argument.
// It is used to configure the SignedURLConfig struct.
type SignedURLOption func(*SignedURLConfig)

// SignedURLConfig represents the configuration for generating a signed URL for a blob object.
type SignedURLConfig struct {
	// Expires is the duration for which the URL is valid.
	Expires time.Time
	// Method is the HTTP method for which the URL is valid.
	Method string
	// Bucket is the bucket name.
	Bucket string
	// Object is the object name.
	Object string
}

func WithExpiry(expiry time.Time) SignedURLOption {
	return func(c *SignedURLConfig) {
		c.Expires = expiry
	}
}

func WithMethod(method string) SignedURLOption {
	return func(c *SignedURLConfig) {
		c.Method = method
	}
}

// WithSignedURLBucket sets the bucket for the SignedURLConfig.
func WithSignedURLBucket(bucket string) SignedURLOption {
	return func(c *SignedURLConfig) {
		c.Bucket = bucket
	}
}

func WithSignedURLObject(object string) SignedURLOption {
	return func(c *SignedURLConfig) {
		c.Object = object
	}
}

func (c *SignedURLConfig) Apply(opts ...SignedURLOption) {
	for _, opt := range opts {
		opt(c)
	}
}

type UploadOption func(*UploadConfig)

// UploadConfig represents the configuration options for uploading an object to a bucket.
type UploadConfig struct {
	// Body is the object body.
	Body io.Reader
	// Metadata is the object metadata.
	Metadata map[string]string
	// Bucket is the bucket name.
	Bucket string
	// Object is the object name.
	Object string
	// ContentType is the object content type.
	ContentType string
	// ContentDisposition is the object content disposition.
	ContentDisposition string
	// ContentEncoding is the object content encoding.
	ContentEncoding string
	// ContentLanguage is the object content language.
	ContentLanguage string
}

func WithUploadBucket(bucket string) UploadOption {
	return func(c *UploadConfig) {
		c.Bucket = bucket
	}
}

func WithUploadObject(object string) UploadOption {
	return func(c *UploadConfig) {
		c.Object = object
	}
}

func WithUploadBody(body io.Reader) UploadOption {
	return func(c *UploadConfig) {
		c.Body = body
	}
}

func WithContentType(contentType string) UploadOption {
	return func(c *UploadConfig) {
		c.ContentType = contentType
	}
}

func (c *UploadConfig) Apply(opts ...UploadOption) {
	for _, opt := range opts {
		opt(c)
	}
}

// DownloadOption is a function type that takes a pointer to a DownloadConfig struct as its argument.
// It is used as an option in the Download function to modify the DownloadConfig struct.
type DownloadOption func(*DownloadConfig)

// DownloadConfig describes the configuration options for downloading an object from a bucket.
type DownloadConfig struct {
	// Bucket is the name of the bucket containing the object to download.
	Bucket string
	// Object is the name of the object to download.
	Object string
}

func WithDownloadBucket(bucket string) DownloadOption {
	return func(c *DownloadConfig) {
		c.Bucket = bucket
	}
}

func WithDownloadObject(object string) DownloadOption {
	return func(c *DownloadConfig) {
		c.Object = object
	}
}

func (c *DownloadConfig) Apply(opts ...DownloadOption) {
	for _, opt := range opts {
		opt(c)
	}
}

// Object represents a blob object with its metadata and content.
type Object struct {
	// LastModified is the time the object was last modified.
	LastModified time.Time

	// Body is a ReadCloser that contains the object's content.
	Body io.ReadCloser

	// Metadata is a map of user-defined metadata for the object.
	Metadata map[string]string

	// Bucket is the name of the bucket containing the object.
	Bucket string

	// Name is the name of the object.
	Name string

	// ContentType is the MIME type of the object's content.
	ContentType string

	// URL is the URL of the object.
	URL string

	// ContentDisposition specifies the presentation style of the object's content.
	ContentDisposition string

	// ContentEncoding specifies the encoding of the object's content.
	ContentEncoding string

	// ContentLanguage specifies the natural language of the object's content.
	ContentLanguage string

	// Size is the size of the object's content in bytes.
	Size int64
}
