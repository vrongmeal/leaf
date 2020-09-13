.PHONY: build devtools format help lint test

.DEFAULT_GOAL := build

GO := go
GOPATH := $(shell go env GOPATH)
GOPATH_BIN := $(GOPATH)/bin
GOLANGCI_LINT := $(GOPATH_BIN)/golangci-lint
BUILD_OUTPUT := ./target/leaf
BUILD_INPUT := cmd/leaf/main.go
GO_PACKAGES := $(shell go list ./... | grep -v vendor)

build:
	@echo "Building..."
	@test -d target || mkdir target
	@$(GO) build -o $(BUILD_OUTPUT) $(BUILD_INPUT)
	@echo "Built as $(BUILD_OUTPUT)"

devtools:
	@echo "Installing golangci-lint..."
	@curl -sSfL \
	 	https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
	 	sh -s -- -b $(GOPATH_BIN) v1.24.0
	@echo "Installed successfully"

format:
	@echo "Formatting..."
	@$(GO) fmt $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run --fix --issues-exit-code 0 > /dev/null 2>&1
	@echo "Code formatted"

help:
	@echo "make [command]"
	@echo "build    - Build command line tool"
	@echo "devtools - Install required development tools"
	@echo "format   - Format code using golangci-lint"
	@echo "help     - Prints help message"
	@echo "lint     - Lint code using golangci-lint"
	@echo "test     - Runs all go tests"

lint:
	@echo "Linting..."
	@$(GO) vet $(GO_PACKAGES)
	@$(GOLANGCI_LINT) run
	@echo "No errors found"

test:
	@echo "Testing..."
	@$(GO) test -v -count 1 ./...
	@echo "Tests ran successfully!"
