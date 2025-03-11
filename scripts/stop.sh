#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_green() {
    echo -e "${GREEN}$1${NC}"
}

print_yellow() {
    echo -e "${YELLOW}$1${NC}"
}

print_red() {
    echo -e "${RED}$1${NC}"
}

# Default values
SKIP_NEO4J=false
SKIP_BUILDER=false
SKIP_ENRICHER=false
SKIP_FRONTEND=false

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-neo4j)
            SKIP_NEO4J=true
            shift
            ;;
        --skip-builder)
            SKIP_BUILDER=true
            shift
            ;;
        --skip-enricher)
            SKIP_ENRICHER=true
            shift
            ;;
        --skip-frontend)
            SKIP_FRONTEND=true
            shift
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --skip-neo4j              Skip stopping Neo4j"
            echo "  --skip-builder            Skip stopping the Knowledge Graph Builder"
            echo "  --skip-enricher           Skip stopping the Knowledge Graph Enricher"
            echo "  --skip-frontend           Skip stopping the Knowledge Graph Frontend"
            echo "  --help                    Show this help message"
            exit 0
            ;;
        *)
            print_red "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Stop the application
if [ "$SKIP_FRONTEND" = false ]; then
    print_yellow "Stopping Knowledge Graph Frontend..."
    docker-compose stop kg-frontend
    print_green "Knowledge Graph Frontend stopped."
fi

if [ "$SKIP_ENRICHER" = false ]; then
    print_yellow "Stopping Knowledge Graph Enricher..."
    docker-compose stop kg-enricher
    print_green "Knowledge Graph Enricher stopped."
fi

if [ "$SKIP_BUILDER" = false ]; then
    print_yellow "Stopping Knowledge Graph Builder..."
    docker-compose stop kg-builder
    print_green "Knowledge Graph Builder stopped."
fi

if [ "$SKIP_NEO4J" = false ]; then
    print_yellow "Stopping Neo4j..."
    docker-compose stop neo4j
    print_green "Neo4j stopped."
fi

print_green "Application stopped successfully." 