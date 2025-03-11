.PHONY: build test clean run stop status

# Build all components
build:
	@echo "Building all components..."
	@docker-compose build

# Run tests for all components
test:
	@echo "Running tests for all components..."
	@go test -v ./internal/...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf */cache/*.json
	@docker-compose down -v

# Run the application
run:
	@echo "Starting the application..."
	@docker-compose up -d

# Stop the application
stop:
	@echo "Stopping the application..."
	@docker-compose down

# Show application status
status:
	@echo "Application status:"
	@docker-compose ps

# Build and run the builder component only
builder:
	@echo "Building and running the builder component..."
	@docker-compose up -d neo4j
	@docker-compose up -d kg-builder

# Build and run the enricher component only
enricher:
	@echo "Building and running the enricher component..."
	@docker-compose up -d neo4j
	@docker-compose up -d kg-enricher

# Build and run the frontend component only
frontend:
	@echo "Building and running the frontend component..."
	@docker-compose up -d neo4j
	@docker-compose up -d kg-frontend

# Show help
help:
	@echo "Available targets:"
	@echo "  build     - Build all components"
	@echo "  test      - Run tests for all components"
	@echo "  clean     - Clean up build artifacts"
	@echo "  run       - Run the application"
	@echo "  stop      - Stop the application"
	@echo "  status    - Show application status"
	@echo "  builder   - Build and run the builder component only"
	@echo "  enricher  - Build and run the enricher component only"
	@echo "  frontend  - Build and run the frontend component only"
	@echo "  help      - Show this help message" 