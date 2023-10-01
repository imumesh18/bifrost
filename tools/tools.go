//go:build tools

// Package tools includes the list of tools used in the project.
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // golangci-lint is a linters aggregator
	_ "go.uber.org/mock/mockgen"                            // mockgen generates mocks for Go interfaces
)
