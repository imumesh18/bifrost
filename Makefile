SHELL := /bin/bash
APP_NAME=bifrost
ALL_PACKAGES=$(shell go list ./... | grep -v /vendor | uniq)

TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)
ifeq ($(VERSION),)
	VERSION := 0.0.0
endif
ifneq ($(COMMIT), $(TAG_COMMIT))
	VERSION := $(VERSION)-devel-$(COMMIT)-$(DATE)
endif
ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

.PHONY: all
all: clean setup copy-config test-copy-config ci-copy-config lint test build


.PHONY: download
download: ## Download go mod dependencies
	@go mod download

.PHONY: setup
setup: ## Setup Installs Dependecies for Projects
	@cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: lint
lint: ## Check style mistake
	@golangci-lint run -c .golangci.yml

.PHONY: fix-lint
fix-lint: ## Fix found lint issues (if it's supported by the linter)
	@golangci-lint run --fix -c .golangci.yml

.PHONY: tidy
tidy: ## Add missing and remove unused modules
	@go mod tidy

.PHONY: clean
clean: ## Clean the builds
	@rm -rf out/
	@rm -f coverage*.out

.PHONY: test
test: test-copy-config ## Run test with short coverage
	@env $(shell grep -v '^#' configs/test.env | xargs -d '\n') go test -v -race -covermode=atomic -coverprofile=./out/coverage.out -coverpkg=$(ALL_PACKAGES) $(ALL_PACKAGES)

.PHONY: test.ci
test.ci: ci-copy-config ## Run ci test with short coverage
	@env $(shell grep -v '^#' configs/test.env | xargs -d '\n') go test -v -race -covermode=atomic -coverprofile=./out/coverage.out -coverpkg=$(ALL_PACKAGES) $(ALL_PACKAGES)

.PHONY: help
help: ## Shows help.
	@echo
	@echo 'Usage:'
	@echo '    make <target>'
	@echo
	@echo 'Targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo
