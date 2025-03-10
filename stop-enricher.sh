#!/bin/bash

# Exit on error
set -e

echo "Stopping Knowledge Graph Enricher..."

# Change to the project directory
cd "$(dirname "$0")/kg-enricher"

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
    echo "Docker and/or Docker Compose are not installed. Please install them first."
    exit 1
fi

# Stop the application
docker-compose down

echo "Knowledge Graph Enricher has been stopped." 