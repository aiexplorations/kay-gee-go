#!/bin/bash

# Exit on error
set -e

echo "Running tests for kg-enricher..."

# Create a temporary file to store the list of packages to test
packages_file=$(mktemp)
go list ./internal/... | grep -v "kg-enricher/internal/enricher" > "$packages_file"

# Run tests for each package
while read -r package; do
  echo "Testing $package..."
  go test -v -cover "$package"
done < "$packages_file"

# Clean up
rm "$packages_file"

echo "Tests completed successfully!" 