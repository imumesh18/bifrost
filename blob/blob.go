package blob

import (
	"context"
)

type Uploader interface {
	Upload(ctx context.Context, opts ...UploadOption) (*Object, error)
}

type Downloader interface {
	Download(ctx context.Context, opts ...DownloadOption) (*Object, error)
}

type SignedURLGenerator interface {
	GenerateSignedURL(ctx context.Context, opts ...SignedURLOption) (string, error)
}

type Storage interface {
	Uploader
	Downloader
	SignedURLGenerator
}
