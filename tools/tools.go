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

//go:build tools

// Package tools includes the list of tools used in the project.
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // golangci-lint is a linters aggregator
	_ "go.uber.org/mock/mockgen"                            // mockgen generates mocks for Go interfaces
)
