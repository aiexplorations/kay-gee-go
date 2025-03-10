#!/bin/bash

# Exit on error
set -e

# Parse command-line arguments
SEED_CONCEPT=""
MAX_NODES=""
TIMEOUT=""
RANDOM_RELATIONSHIPS=""
CONCURRENCY=""

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --seed=*)
      SEED_CONCEPT="${1#*=}"
      shift
      ;;
    --max-nodes=*)
      MAX_NODES="${1#*=}"
      shift
      ;;
    --timeout=*)
      TIMEOUT="${1#*=}"
      shift
      ;;
    --random-relationships=*)
      RANDOM_RELATIONSHIPS="${1#*=}"
      shift
      ;;
    --concurrency=*)
      CONCURRENCY="${1#*=}"
      shift
      ;;
    *)
      echo "Unknown parameter: $1"
      echo "Usage: $0 [--seed=CONCEPT] [--max-nodes=N] [--timeout=MINUTES] [--random-relationships=N] [--concurrency=N]"
      echo "If not specified, default values will be used:"
      echo "  --seed: Artificial Intelligence"
      echo "  --max-nodes: 100"
      echo "  --timeout: 30 minutes"
      echo "  --random-relationships: 50"
      echo "  --concurrency: 5"
      exit 1
      ;;
  esac
done

echo "Starting Knowledge Graph Builder..."
echo "Default configuration:"
echo "  Seed concept: Artificial Intelligence"
echo "  Max nodes: 100"
echo "  Timeout: 30 minutes"
echo "  Random relationships: 50"
echo "  Concurrency: 5"

# Change to the project directory
cd "$(dirname "$0")/kg-builder"

# Copy the config.yaml file if it exists in the parent directory
if [ -f "../config.yaml" ]; then
    echo "Using configuration from config.yaml"
    cp ../config.yaml .
fi

# Show version information
echo "Version: $(cat ../VERSION)"

# Run tests first
echo "Running tests..."
./run-tests.sh

# Check if tests passed
if [ $? -ne 0 ]; then
    echo "Tests failed. Aborting startup."
    exit 1
fi

echo "Tests passed. Starting application..."

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
    echo "Docker and/or Docker Compose are not installed. Please install them first."
    exit 1
fi

# Build the command with any provided arguments
CMD="/kg-builder"
if [ -n "$SEED_CONCEPT" ]; then
    CMD="$CMD --seed='$SEED_CONCEPT'"
    echo "Using seed concept: $SEED_CONCEPT"
fi
if [ -n "$MAX_NODES" ]; then
    CMD="$CMD --max-nodes=$MAX_NODES"
    echo "Using max nodes: $MAX_NODES"
fi
if [ -n "$TIMEOUT" ]; then
    CMD="$CMD --timeout=$TIMEOUT"
    echo "Using timeout: $TIMEOUT minutes"
fi
if [ -n "$RANDOM_RELATIONSHIPS" ]; then
    CMD="$CMD --random-relationships=$RANDOM_RELATIONSHIPS"
    echo "Using random relationships: $RANDOM_RELATIONSHIPS"
fi
if [ -n "$CONCURRENCY" ]; then
    CMD="$CMD --concurrency=$CONCURRENCY"
    echo "Using concurrency: $CONCURRENCY"
fi

# Start the application with Docker Compose
if [ -n "$SEED_CONCEPT" ] || [ -n "$MAX_NODES" ] || [ -n "$TIMEOUT" ] || [ -n "$RANDOM_RELATIONSHIPS" ] || [ -n "$CONCURRENCY" ]; then
    # If any arguments were provided, run the container with the custom command
    echo "Starting with custom configuration: $CMD"
    docker-compose up -d --build
    docker-compose exec -d kg-builder sh -c "$CMD"
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

echo "Knowledge Graph Builder is now running."
echo "You can view the logs with 'docker-compose logs -f'."
echo "You can access the Neo4j browser at http://localhost:7474."
echo "To stop the application, run './stop.sh'." 