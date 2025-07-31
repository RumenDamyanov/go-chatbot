# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Binary names
BINARY_NAME=go-chatbot
BINARY_UNIX=$(BINARY_NAME)_unix

# Version
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT?=$(shell git rev-parse HEAD)

# Linker flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

.PHONY: all build clean test coverage deps fmt lint help

all: test build

## build: Build the binary file
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

## clean: Remove build related file
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

## test: Run the tests
test:
	$(GOTEST) -v ./...

## test-race: Run tests with race condition detection
test-race:
	$(GOTEST) -race -short ./...

## coverage: Run tests with coverage
coverage:
	$(GOTEST) -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## deps: Get the dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## fmt: Format the Go code
fmt:
	$(GOFMT) -s -w .

## lint: Lint the Go code
lint:
	$(GOLINT) run

## lint-fix: Lint and fix the Go code
lint-fix:
	$(GOLINT) run --fix

## vet: Run go vet
vet:
	$(GOCMD) vet ./...

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test

## build-linux: Cross compilation for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v

## docker-build: Build docker image
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

## security: Run security scan
security:
	gosec ./...

## mod-update: Update all dependencies
mod-update:
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

## install-tools: Install development tools
install-tools:
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

## help: Show this help message
help: Makefile
	@echo "Choose a command to run:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
