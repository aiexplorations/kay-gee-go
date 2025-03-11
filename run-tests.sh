#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_green() {
    echo -e "${GREEN}$1${NC}"
}

print_yellow() {
    echo -e "${YELLOW}$1${NC}"
}

print_red() {
    echo -e "${RED}$1${NC}"
}

# Run tests
print_yellow "Running tests..."
go test -v ./internal/...

# Check if tests passed
if [ $? -eq 0 ]; then
    print_green "All tests passed!"
    exit 0
else
    print_red "Tests failed!"
    exit 1
fi 