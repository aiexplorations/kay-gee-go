#!/bin/bash

# Script to stop the knowledge graph enricher

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Check if the enricher is running
if ! docker ps | grep -q kaygeego-enricher; then
  echo "Knowledge Graph Enricher is not running"
  exit 0
fi

# Stop the enricher
echo "Stopping Knowledge Graph Enricher..."

# Use docker-compose to stop the enricher
cd "$PROJECT_ROOT" && docker-compose stop enricher

echo "Knowledge Graph Enricher stopped successfully" 