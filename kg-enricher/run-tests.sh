#!/bin/bash

# Exit on error
set -e

echo "Running tests for kg-enricher..."

# Run tests with coverage
go test -v -cover ./internal/...

echo "Tests completed successfully!" 