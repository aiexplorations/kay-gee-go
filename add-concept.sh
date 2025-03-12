#!/bin/bash

# add-concept.sh - Script to add a new concept to the knowledge graph
# This script launches a new builder container for a specific concept

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
MAX_NODES=50
TIMEOUT=15
RANDOM_RELATIONSHIPS=10
CONCURRENCY=3

# Function to show help
show_help() {
    echo "Usage: $0 CONCEPT [options]"
    echo ""
    echo "Arguments:"
    echo "  CONCEPT                    The concept to add to the knowledge graph"
    echo ""
    echo "Options:"
    echo "  --max-nodes=N             Maximum number of nodes to build (default: 50)"
    echo "  --timeout=N               Timeout in minutes for graph building (default: 15)"
    echo "  --random-relationships=N  Number of random relationships to mine (default: 10)"
    echo "  --concurrency=N           Number of concurrent workers (default: 3)"
    echo ""
    echo "Examples:"
    echo "  $0 \"Machine Learning\"      Add Machine Learning concept with default settings"
    echo "  $0 \"Blockchain\" --max-nodes=30 --timeout=10"
}

# Check if a concept was provided
if [ $# -eq 0 ]; then
    print_red "Error: No concept specified."
    show_help
    exit 1
fi

# Get the concept from the first argument
CONCEPT="$1"
shift

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --max-nodes=*)
            MAX_NODES="${1#*=}"
            shift
            ;;
        --timeout=*)
            TIMEOUT="${1#*=}"
            shift
            ;;
        --random-relationships=*)
            RANDOM_RELATIONSHIPS="${1#*=}"
            shift
            ;;
        --concurrency=*)
            CONCURRENCY="${1#*=}"
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            print_red "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Generate a unique container name based on the concept and timestamp
TIMESTAMP=$(date +%s)
SANITIZED_CONCEPT=$(echo "$CONCEPT" | tr ' ' '_' | tr -cd '[:alnum:]_-')
CONTAINER_NAME="kaygeego-builder-${SANITIZED_CONCEPT}-${TIMESTAMP}"

print_yellow "Launching new builder container for concept: \"$CONCEPT\""
print_yellow "Container name: $CONTAINER_NAME"
print_yellow "Max nodes: $MAX_NODES"
print_yellow "Timeout: $TIMEOUT minutes"
print_yellow "Random relationships: $RANDOM_RELATIONSHIPS"
print_yellow "Concurrency: $CONCURRENCY"

# Launch a new builder container with the specified parameters
docker run -d \
    --name "$CONTAINER_NAME" \
    --network kay-gee-go_kaygeego-network \
    -e NEO4J_URI=bolt://neo4j:7687 \
    -e NEO4J_USER=neo4j \
    -e NEO4J_PASSWORD=password \
    -e LLM_URL=http://host.docker.internal:11434/api/generate \
    -e LLM_MODEL=phi4:latest \
    -e SEED_CONCEPT="$CONCEPT" \
    -e MAX_NODES="$MAX_NODES" \
    -e TIMEOUT_MINUTES="$TIMEOUT" \
    -e RANDOM_RELATIONSHIPS="$RANDOM_RELATIONSHIPS" \
    -e CONCURRENCY="$CONCURRENCY" \
    -v "$(pwd)/kg-builder/cache:/app/cache" \
    -v "$(pwd)/kg-builder/config.yaml:/app/config.yaml" \
    --add-host=host.docker.internal:host-gateway \
    kaygeego-builder \
    /kg-builder

if [ $? -eq 0 ]; then
    print_green "Successfully launched builder container for concept: \"$CONCEPT\""
    print_yellow "You can monitor the progress with:"
    print_yellow "  docker logs -f $CONTAINER_NAME"
else
    print_red "Failed to launch builder container for concept: \"$CONCEPT\""
    exit 1
fi

exit 0 