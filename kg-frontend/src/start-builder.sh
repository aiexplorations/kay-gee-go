#!/bin/bash

# Script to start the knowledge graph builder

# Default values
SEED="Artificial Intelligence"
MAX_NODES=100
TIMEOUT=30
RANDOM_RELATIONSHIPS=50
CONCURRENCY=5

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --seed)
      SEED="$2"
      shift 2
      ;;
    --max-nodes)
      MAX_NODES="$2"
      shift 2
      ;;
    --timeout)
      TIMEOUT="$2"
      shift 2
      ;;
    --random-relationships)
      RANDOM_RELATIONSHIPS="$2"
      shift 2
      ;;
    --concurrency)
      CONCURRENCY="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Check if the builder is already running
if docker ps | grep -q kaygeego-builder; then
  echo "Knowledge Graph Builder is already running"
else
  # Start the builder
  echo "Starting Knowledge Graph Builder with the following parameters:"
  echo "  Seed concept: $SEED"
  echo "  Max nodes: $MAX_NODES"
  echo "  Timeout: $TIMEOUT minutes"
  echo "  Random relationships: $RANDOM_RELATIONSHIPS"
  echo "  Concurrency: $CONCURRENCY"

  # Use docker-compose to start the builder
  cd "$PROJECT_ROOT" && docker-compose up -d builder
fi

# Set environment variables for the builder
docker exec kaygeego-builder /bin/sh -c "export SEED_CONCEPT=\"$SEED\" && \
  export MAX_NODES=$MAX_NODES && \
  export TIMEOUT=$TIMEOUT && \
  export RANDOM_RELATIONSHIPS=$RANDOM_RELATIONSHIPS && \
  export CONCURRENCY=$CONCURRENCY && \
  echo 'Builder configured with seed: $SEED, max nodes: $MAX_NODES'"

echo "Knowledge Graph Builder started successfully" 