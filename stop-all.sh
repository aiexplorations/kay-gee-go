#!/bin/bash

# Exit on error
set -e

echo "Stopping Knowledge Graph Builder and Enricher..."

# Check if Docker and Docker Compose are installed
if ! command -v docker &> /dev/null || ! command -v docker-compose &> /dev/null; then
    echo "Docker and/or Docker Compose are not installed. Please install them first."
    exit 1
fi

# Stop the services
docker-compose down

echo "Knowledge Graph Builder and Enricher have been stopped." 