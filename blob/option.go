package blob

import (
	"io"
	"time"
)

type SignedURLOption func(*SignedURLConfig)

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

type DownloadOption func(*DownloadConfig)

type DownloadConfig struct {
	// Bucket is the bucket name.
	Bucket string
	// Object is the object name.
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

type Object struct {
	// LastModified is the object last modified time.
	LastModified time.Time
	// Body is the object body.
	Body io.ReadCloser
	// Metadata is the object metadata.
	Metadata map[string]string
	// Bucket is the bucket name.
	Bucket string
	// Name is the object name.
	Name string
	// ContentType is the object content type.
	ContentType string
	// URL is the object URL.
	URL string
	// Size is the object size.
	Size int64
}
