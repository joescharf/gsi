# Makefile for gsi
BINARY_NAME := gsi
MODULE := $(shell head -1 go.mod | awk '{print $$2}')
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X $(MODULE)/cmd.version=$(VERSION) -X $(MODULE)/cmd.commit=$(COMMIT) -X $(MODULE)/cmd.date=$(BUILD_DATE)"

.DEFAULT_GOAL := build

##@ App
.PHONY: build install run clean tidy test test-cover lint vet fmt

build: ## Build the Go binary
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

install: ## Install the binary to $GOPATH/bin
	go install $(LDFLAGS) .

run: build ## Build and run the binary
	./bin/$(BINARY_NAME)

clean: ## Remove build artifacts
	rm -rf bin/
	rm -f coverage.out

tidy: ## Run go mod tidy
	go mod tidy

test: ## Run tests
	go test -v -race -count=1 ./...

test-cover: ## Run tests with coverage
	go test -v -race -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: vet ## Run golangci-lint
	@which golangci-lint > /dev/null 2>&1 || { echo "Install golangci-lint: https://golangci-lint.run/welcome/install/"; exit 1; }
	golangci-lint run ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Run gofmt
	gofmt -s -w .

##@ Help
.PHONY: help

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)
