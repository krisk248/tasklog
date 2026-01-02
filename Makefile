.PHONY: build run test clean install release

# Binary name
BINARY=nexus

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Main package path
MAIN=./cmd/tasklog

# Build the application
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY) $(MAIN)

# Run the application
run:
	$(GORUN) $(MAIN)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY)
	rm -rf dist/

# Install to GOPATH/bin
install:
	$(GOCMD) install $(MAIN)

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Update dependencies
update:
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Build for all platforms
build-all: clean
	mkdir -p dist
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 $(MAIN)
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY)-linux-arm64 $(MAIN)
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 $(MAIN)
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64 $(MAIN)
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe $(MAIN)

# Release with goreleaser
release:
	goreleaser release --clean

# Snapshot release (for testing)
snapshot:
	goreleaser release --snapshot --clean

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code
lint:
	golangci-lint run

# Default target
all: tidy fmt build
