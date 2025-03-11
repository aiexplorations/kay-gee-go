#!/bin/bash

# Script to stop the knowledge graph builder

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Check if the builder is running
if ! docker ps | grep -q kg-builder; then
  echo "Knowledge Graph Builder is not running"
  exit 0
fi

# Stop the builder
echo "Stopping Knowledge Graph Builder..."

# Use docker-compose to stop the builder
cd "$PROJECT_ROOT" && docker-compose stop kg-builder

echo "Knowledge Graph Builder stopped successfully" 