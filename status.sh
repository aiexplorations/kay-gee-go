#!/bin/bash

echo "Checking Knowledge Graph Builder status..."

# Change to the project directory
cd "$(dirname "$0")/kg-builder"

# Show version information
echo "Version: $(cat ../VERSION)"

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose is not installed. Please install it first."
    exit 1
fi

# Check if the containers are running
if [ "$(docker-compose ps -q | wc -l)" -eq 0 ]; then
    echo "Knowledge Graph Builder is not running."
    exit 0
fi

# Show the status of the containers
echo "Knowledge Graph Builder is running."
echo "Container status:"
docker-compose ps

# Show Neo4j statistics if the Neo4j container is running
if docker-compose ps | grep -q "neo4j.*Up"; then
    echo "Neo4j is running. Checking graph statistics..."
    
    # Show current configuration
    echo "Current configuration:"
    docker-compose exec kg-builder /kg-builder --version
    
    # Run a command to get graph statistics
    echo "Running query to get graph statistics..."
    docker-compose exec kg-builder /kg-builder --stats-only || echo "Failed to get statistics."
fi 