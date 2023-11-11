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

// Package blob provides an interface for uploading, downloading, and generating signed URLs for objects in a storage service.
package blob

import (
	"context"
)

// Uploader is an interface for uploading objects to a storage service.
type Uploader interface {
	Upload(ctx context.Context, opts ...UploadOption) (*Object, error)
}

// Downloader is an interface for downloading objects.
type Downloader interface {
	Download(ctx context.Context, opts ...DownloadOption) (*Object, error)
}

// SignedURLGenerator is an interface for generating signed URLs.
type SignedURLGenerator interface {
	GenerateSignedURL(ctx context.Context, opts ...SignedURLOption) (string, error)
}

// Storage is an interface that combines the Uploader, Downloader, and SignedURLGenerator interfaces.
// It represents a storage system that can upload, download, and generate signed URLs for files.
type Storage interface {
	Uploader
	Downloader
	SignedURLGenerator
}
