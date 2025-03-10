#!/bin/bash

# Exit on error
set -e

# Parse command-line arguments
RUN_ONCE=""
COUNT=""
SHOW_STATS=""
SHOW_VERSION=""
SKIP_TESTS=false

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --run-once)
      RUN_ONCE="--run-once"
      shift
      ;;
    --count=*)
      COUNT="--count=${1#*=}"
      shift
      ;;
    --stats)
      SHOW_STATS="--stats"
      shift
      ;;
    --version)
      SHOW_VERSION="--version"
      shift
      ;;
    --skip-tests)
      SKIP_TESTS=true
      shift
      ;;
    *)
      echo "Unknown parameter: $1"
      echo "Usage: $0 [--run-once] [--count=N] [--stats] [--version] [--skip-tests]"
      exit 1
      ;;
  esac
done

echo "Starting Knowledge Graph Enricher..."
echo "Default configuration:"
echo "  Batch size: 10"
echo "  Interval: 60 seconds"
echo "  Max relationships: 100"
echo "  Concurrency: 5"
echo "  LLM model: qwen2.5:3b"

# Change to the project directory
cd "$(dirname "$0")/kg-enricher"

# Show version information
echo "Version: $(cat ../VERSION)"

# Run tests if not skipped
if [ "$SKIP_TESTS" = false ]; then
  echo "Running tests..."
  
  echo "Running unit tests..."
  ./run-tests.sh
  
  # Check if unit tests passed
  if [ $? -ne 0 ]; then
    echo "Unit tests failed. Aborting startup."
    exit 1
  fi
  
  echo "Running integration tests..."
  go test -v ./internal/tests/...
  
  # Check if integration tests passed
  if [ $? -ne 0 ]; then
    echo "Integration tests failed. Aborting startup."
    exit 1
  fi
  
  echo "All tests passed successfully!"
else
  echo "Skipping tests as requested."
fi

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
    echo "Docker and/or Docker Compose are not installed. Please install them first."
    exit 1
fi

# Build the command with any provided arguments
CMD=""
if [ -n "$RUN_ONCE" ]; then
    CMD="$CMD $RUN_ONCE"
fi
if [ -n "$COUNT" ]; then
    CMD="$CMD $COUNT"
fi
if [ -n "$SHOW_STATS" ]; then
    CMD="$CMD $SHOW_STATS"
fi
if [ -n "$SHOW_VERSION" ]; then
    CMD="$CMD $SHOW_VERSION"
fi

# Start the application with Docker Compose
if [ -n "$CMD" ]; then
    # If any arguments were provided, run the container with the custom command
    echo "Starting with custom configuration: $CMD"
    docker-compose up -d --build
    docker-compose exec -d kg-enricher /app/kg-enricher $CMD
else
    # Otherwise, just start the containers normally
    docker-compose up --build -d
fi

# Wait for the application to start
echo "Waiting for the application to start..."
sleep 5

# Check if the containers are running
if [ "$(docker-compose ps -q | wc -l)" -eq 0 ]; then
    echo "Failed to start the application. Check the logs with 'docker-compose logs'."
    exit 1
fi

echo "Knowledge Graph Enricher is now running."
echo "You can view the logs with 'docker-compose logs -f'."
echo "You can access the Neo4j browser at http://localhost:7474."
echo "To stop the application, run './stop-enricher.sh'." 