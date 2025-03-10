.PHONY: test test-short test-verbose test-coverage build run clean docker-test docker-test-clean check-go test-build

# Default target
all: check-go test build

# Check if Go is installed
check-go:
	@which go > /dev/null || (echo "Go is not installed or not in PATH. Please install Go or use Docker tests with 'make docker-test'." && exit 1)

# Build the application
build: check-go
	mkdir -p bin
	go build -o bin/kg-builder ./cmd/kg-builder

# Run the application
run: build
	./bin/kg-builder

# Test application build and basic functionality
test-build:
	./test-build.sh

# Run all tests
test: check-go
	go test ./internal/... ./cmd/...

# Run tests in short mode (skip integration tests)
test-short: check-go
	go test -short ./internal/... ./cmd/...

# Run tests with verbose output
test-verbose: check-go
	go test -v ./internal/... ./cmd/...

# Run tests with coverage report
test-coverage: check-go
	go test -coverprofile=coverage.out ./internal/... ./cmd/...
	go tool cover -html=coverage.out -o coverage.html

# Run tests in Docker (no Go installation required)
docker-test:
	mkdir -p test-cache
	docker-compose -f docker-compose.test.yml up --build

# Clean Docker test containers
docker-test-clean:
	docker-compose -f docker-compose.test.yml down -v

# Clean build artifacts
clean: docker-test-clean
	rm -rf bin
	rm -rf test-cache
	rm -f coverage.out coverage.html 