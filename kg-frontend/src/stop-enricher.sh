#!/bin/bash

# Script to stop the knowledge graph enricher

# Check if the enricher is running
if ! docker ps | grep -q kg-enricher; then
  echo "Knowledge Graph Enricher is not running"
  exit 0
fi

# Stop the enricher
echo "Stopping Knowledge Graph Enricher..."

# Use docker-compose to stop the enricher
cd /app && docker-compose stop kg-enricher

echo "Knowledge Graph Enricher stopped successfully" 