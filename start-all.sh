#!/bin/bash

# Exit on error
set -e

# Parse command-line arguments for kg-builder
SEED_CONCEPT=""
MAX_NODES=""
TIMEOUT=""
RANDOM_RELATIONSHIPS=""
CONCURRENCY=""

# Parse command-line arguments for kg-enricher
RUN_ONCE=""
COUNT=""
SHOW_STATS=""

# Flag to skip tests
SKIP_TESTS=false

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
    --skip-tests)
      SKIP_TESTS=true
      shift
      ;;
    --version)
      echo "Knowledge Graph Builder and Enricher v0.1.0"
      exit 0
      ;;
    *)
      echo "Unknown parameter: $1"
      echo "Usage: $0 [--seed=CONCEPT] [--max-nodes=N] [--timeout=MINUTES] [--random-relationships=N] [--concurrency=N] [--run-once] [--count=N] [--stats] [--skip-tests]"
      exit 1
      ;;
  esac
done

echo "Starting Knowledge Graph Builder and Enricher..."
echo "Default configuration:"
echo "Knowledge Graph Builder:"
echo "  Seed concept: Artificial Intelligence"
echo "  Max nodes: 100"
echo "  Timeout: 30 minutes"
echo "  Random relationships: 50"
echo "  Concurrency: 5"
echo "Knowledge Graph Enricher:"
echo "  Batch size: 10"
echo "  Interval: 60 seconds"
echo "  Max relationships: 100"
echo "  Concurrency: 5"
echo "  LLM model: qwen2.5:3b"

# Show version information
echo "Version: $(cat VERSION)"

# Run tests if not skipped
if [ "$SKIP_TESTS" = false ]; then
  echo "Running tests for both components..."
  
  # Run kg-builder tests
  echo "Running kg-builder tests..."
  cd kg-builder
  
  echo "Running unit tests..."
  ./run-tests.sh
  
  # Check if unit tests passed
  if [ $? -ne 0 ]; then
    echo "kg-builder unit tests failed. Aborting startup."
    exit 1
  fi
  
  echo "Running integration tests..."
  go test -v ./internal/tests/...
  
  # Check if integration tests passed
  if [ $? -ne 0 ]; then
    echo "kg-builder integration tests failed. Aborting startup."
    exit 1
  fi
  
  echo "Running end-to-end tests..."
  go test -v ./internal/tests/end_to_end_test.go
  
  # Check if end-to-end tests passed
  if [ $? -ne 0 ]; then
    echo "kg-builder end-to-end tests failed. Aborting startup."
    exit 1
  fi
  
  cd ..
  
  # Run kg-enricher tests
  echo "Running kg-enricher tests..."
  cd kg-enricher
  
  echo "Running unit tests..."
  ./run-tests.sh
  
  # Check if unit tests passed
  if [ $? -ne 0 ]; then
    echo "kg-enricher unit tests failed. Aborting startup."
    exit 1
  fi
  
  echo "Running integration tests..."
  go test -v ./internal/tests/...
  
  # Check if integration tests passed
  if [ $? -ne 0 ]; then
    echo "kg-enricher integration tests failed. Aborting startup."
    exit 1
  fi
  
  cd ..
  
  echo "All tests passed successfully!"
else
  echo "Skipping tests as requested."
fi

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
    echo "Docker and/or Docker Compose are not installed. Please install them first."
    exit 1
fi

# Start the services with Docker Compose
echo "Starting services with Docker Compose..."
docker-compose up -d --build

# Wait for the services to start

echo "Waiting for the services to start..."
sleep 5

# Check if the containers are running
if [ "$(docker-compose ps -q | wc -l)" -eq 0 ]; then
    echo "Failed to start the services. Check the logs with 'docker-compose logs'."
    exit 1
fi

# Build the command for kg-builder with any provided arguments
if [ -n "$SEED_CONCEPT" ] || [ -n "$MAX_NODES" ] || [ -n "$TIMEOUT" ] || [ -n "$RANDOM_RELATIONSHIPS" ] || [ -n "$CONCURRENCY" ]; then
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
    
    # Run the kg-builder with custom arguments
    echo "Starting kg-builder with custom configuration: $CMD"
    docker-compose exec -d kg-builder sh -c "$CMD"
fi

# Build the command for kg-enricher with any provided arguments
if [ -n "$RUN_ONCE" ] || [ -n "$COUNT" ] || [ -n "$SHOW_STATS" ]; then
    CMD="/app/enricher"
    if [ -n "$RUN_ONCE" ]; then
        CMD="$CMD $RUN_ONCE"
        echo "Running kg-enricher once"
    fi
    if [ -n "$COUNT" ]; then
        CMD="$CMD $COUNT"
        echo "Using count: ${COUNT#*=}"
    fi
    if [ -n "$SHOW_STATS" ]; then
        CMD="$CMD $SHOW_STATS"
        echo "Showing kg-enricher stats"
    fi
    
    # Run the kg-enricher with custom arguments
    echo "Starting kg-enricher with custom configuration: $CMD"
    docker-compose exec -d kg-enricher sh -c "$CMD"
fi

echo "Knowledge Graph Builder and Enricher are now running."
echo "You can view the logs with 'docker-compose logs -f'."
echo "You can access the Neo4j browser at http://localhost:7474."
echo "To stop the services, run './stop-all.sh'." 