#!/bin/bash

# Test script for MDMA
# Supports running tests in CI mode (no audio hardware) or full mode

set -e

case "${1:-ci}" in
    "ci")
        echo "Running tests in CI mode (no audio hardware required)..."
        go test ./... -short -v
        ;;
    "full")
        echo "Running full test suite (requires audio hardware)..."
        go test ./... -v
        ;;
    "coverage")
        echo "Running tests with coverage..."
        go test ./... -short -v -coverprofile=coverage.out
        go tool cover -html=coverage.out -o coverage.html
        echo "Coverage report generated: coverage.html"
        ;;
    *)
        echo "Usage: $0 [ci|full|coverage]"
        echo "  ci       - Run tests without audio hardware (default)"
        echo "  full     - Run all tests including audio hardware tests"
        echo "  coverage - Run tests and generate coverage report"
        exit 1
        ;;
esac