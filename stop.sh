#!/bin/bash

echo "Stopping Knowledge Graph Builder..."

# Change to the project directory
cd "$(dirname "$0")/kg-builder"

# Show version information
echo "Version: $(cat ../VERSION)"

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose is not installed. Please install it first."
    exit 1
fi

# Stop the application with Docker Compose
docker-compose down

# Check if the containers are stopped
if [ "$(docker-compose ps -q | wc -l)" -ne 0 ]; then
    echo "Failed to stop some containers. Forcing stop..."
    docker-compose down -v --remove-orphans
fi

echo "Knowledge Graph Builder has been stopped." 