# MDMA Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=mdma
BINARY_PATH=./cmd/mdma

# Build targets
.PHONY: all build clean test test-unit test-integration test-audio test-db test-ui deps run test-coverage test-benchmark

all: clean deps build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Test targets
test: test-unit test-integration
	@echo "All tests completed successfully!"

test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v -short ./internal/audio ./internal/database ./internal/config

test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v ./internal/ui

test-audio:
	@echo "Running audio engine tests..."
	$(GOTEST) -v ./internal/audio

test-db:
	@echo "Running database tests..."
	$(GOTEST) -v ./internal/database

test-ui:
	@echo "Running UI tests..."
	$(GOTEST) -v ./internal/ui

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-benchmark:
	@echo "Running benchmark tests..."
	$(GOTEST) -bench=. -benchmem ./internal/ui

# Create test MP3 if it doesn't exist (for testing)
test-setup:
	@if [ ! -f "/tmp/power.mp3" ]; then \
		echo "Test MP3 not found. Please run: ffmpeg -i 'path/to/sample.m4a' /tmp/power.mp3"; \
	else \
		echo "Test MP3 found at /tmp/power.mp3"; \
	fi

deps:
	$(GOMOD) download
	$(GOMOD) tidy

run:
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)
	./$(BINARY_NAME)

# Development
dev: clean deps run

# Install system dependencies (macOS)
install-deps-macos:
	brew install portaudio

# Install system dependencies (Linux)
install-deps-linux:
	sudo apt-get update
	sudo apt-get install -y portaudio19-dev

.DEFAULT_GOAL := all
