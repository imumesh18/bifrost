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

// Package gcs provides an interface for uploading, downloading, and generating signed URLs for objects in Google Cloud Storage.
package gcs

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/imumesh18/bifrost/blob"
)

var _ blob.Storage = (*Storage)(nil)

type Config struct {
	projectID       string
	credentialsFile string
	credentialsJSON []byte
}

type Storage struct {
	client *storage.Client
}

func New(ctx context.Context, opts ...Option) (*Storage, error) {
	c, err := newConfig(opts...)
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(c.credentialsJSON), option.WithCredentialsFile(c.credentialsFile))
	if err != nil {
		return nil, err
	}

	return &Storage{client: client}, nil
}
