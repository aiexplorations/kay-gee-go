#!/bin/bash

# Check if Go is installed
if command -v go &> /dev/null; then
    echo "Go is installed. Running tests locally..."
    cd "$(dirname "$0")"
    go test ./internal/... ./cmd/...
else
    echo "Go is not installed. Running tests in Docker..."
    cd "$(dirname "$0")"
    mkdir -p test-cache
    docker-compose -f docker-compose.test.yml up --build
fi 