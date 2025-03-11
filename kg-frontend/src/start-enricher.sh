#!/bin/bash

# Script to start the knowledge graph enricher

# Default values
BATCH_SIZE=10
INTERVAL=60
MAX_RELATIONSHIPS=100
CONCURRENCY=5

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --batch-size)
      BATCH_SIZE="$2"
      shift 2
      ;;
    --interval)
      INTERVAL="$2"
      shift 2
      ;;
    --max-relationships)
      MAX_RELATIONSHIPS="$2"
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

# Check if the enricher is already running
if docker ps | grep -q kg-enricher; then
  echo "Knowledge Graph Enricher is already running"
  exit 0
fi

# Start the enricher
echo "Starting Knowledge Graph Enricher with the following parameters:"
echo "  Batch size: $BATCH_SIZE"
echo "  Interval: $INTERVAL seconds"
echo "  Max relationships: $MAX_RELATIONSHIPS"
echo "  Concurrency: $CONCURRENCY"

# Use docker-compose to start the enricher
cd /app && docker-compose up -d kg-enricher

# Set environment variables for the enricher
docker exec kg-enricher /bin/sh -c "export ENRICHER_BATCH_SIZE=$BATCH_SIZE && \
  export ENRICHER_INTERVAL_SECONDS=$INTERVAL && \
  export ENRICHER_MAX_RELATIONSHIPS=$MAX_RELATIONSHIPS && \
  export ENRICHER_CONCURRENCY=$CONCURRENCY && \
  /app/kg-enricher"

echo "Knowledge Graph Enricher started successfully" 