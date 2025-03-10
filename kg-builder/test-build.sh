#!/bin/bash

# Exit on error
set -e

echo "Testing application build and basic functionality..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Testing with Docker..."
    
    # Build the Docker image
    docker build -t kg-builder-test .
    
    # Run the Docker container with --help flag to test basic functionality
    docker run --rm kg-builder-test /kg-builder --help
    
    echo "Docker build and basic functionality test passed!"
else
    echo "Go is installed. Testing local build..."
    
    # Build the application
    mkdir -p bin
    go build -o bin/kg-builder ./cmd/kg-builder
    
    # Run the application with --help flag to test basic functionality
    ./bin/kg-builder --help
    
    echo "Local build and basic functionality test passed!"
fi 