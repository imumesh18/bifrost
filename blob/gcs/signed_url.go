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
	"context"

	"cloud.google.com/go/storage"

	"github.com/imumesh18/bifrost/blob"
)

// GenerateSignedURL generates a signed URL for the given object in the specified bucket.
// It takes a context and a variadic list of SignedURLOption as input.
// It returns a string representing the signed URL and an error if any.
func (s *Storage) GenerateSignedURL(_ context.Context, opts ...blob.SignedURLOption) (string, error) {
	c := blob.SignedURLConfig{}
	c.Apply(opts...)

	u, err := s.client.Bucket(c.Bucket).SignedURL(c.Object, &storage.SignedURLOptions{
		Method:  c.Method,
		Expires: c.Expires,
		Scheme:  storage.SigningSchemeV4,
	})
	if err != nil {
		return "", err
	}

	return u, nil
}
