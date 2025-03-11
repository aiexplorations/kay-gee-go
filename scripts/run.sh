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
SEED_CONCEPT="Artificial Intelligence"
MAX_NODES=100
TIMEOUT=30
RANDOM_RELATIONSHIPS=50
CONCURRENCY=5
STATS_ONLY=false
SKIP_NEO4J=false
SKIP_BUILDER=false
SKIP_ENRICHER=false
SKIP_FRONTEND=false
RUN_ONCE=false
COUNT=10

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --seed=*)
            SEED_CONCEPT="${1#*=}"
            shift
            ;;
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
        --stats-only)
            STATS_ONLY=true
            shift
            ;;
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
        --run-once)
            RUN_ONCE=true
            shift
            ;;
        --count=*)
            COUNT="${1#*=}"
            shift
            ;;
        --help)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --seed=CONCEPT            Seed concept for graph building (default: \"Artificial Intelligence\")"
            echo "  --max-nodes=N             Maximum number of nodes to build (default: 100)"
            echo "  --timeout=N               Timeout in minutes for graph building (default: 30)"
            echo "  --random-relationships=N  Number of random relationships to mine (default: 50)"
            echo "  --concurrency=N           Number of concurrent workers (default: 5)"
            echo "  --stats-only              Only show statistics without building the graph"
            echo "  --skip-neo4j              Skip starting Neo4j"
            echo "  --skip-builder            Skip starting the Knowledge Graph Builder"
            echo "  --skip-enricher           Skip starting the Knowledge Graph Enricher"
            echo "  --skip-frontend           Skip starting the Knowledge Graph Frontend"
            echo "  --run-once                Run the enricher once and exit"
            echo "  --count=N                 Number of relationships to mine when running once (default: 10)"
            echo "  --help                    Show this help message"
            exit 0
            ;;
        *)
            print_red "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Update configuration files
if [ "$STATS_ONLY" = false ]; then
    print_yellow "Updating configuration files..."
    
    # Update builder configuration
    sed -i.bak "s/seed_concept:.*/seed_concept: \"$SEED_CONCEPT\"/" config/builder.yaml
    sed -i.bak "s/max_nodes:.*/max_nodes: $MAX_NODES/" config/builder.yaml
    sed -i.bak "s/timeout_minutes:.*/timeout_minutes: $TIMEOUT/" config/builder.yaml
    sed -i.bak "s/random_relationships:.*/random_relationships: $RANDOM_RELATIONSHIPS/" config/builder.yaml
    sed -i.bak "s/concurrency:.*/concurrency: $CONCURRENCY/" config/builder.yaml
    
    # Update enricher configuration
    if [ "$RUN_ONCE" = true ]; then
        export ENRICHER_ARGS="--run-once --count=$COUNT"
    fi
    
    print_green "Configuration updated."
fi

# Start the application
if [ "$SKIP_NEO4J" = false ]; then
    print_yellow "Starting Neo4j..."
    docker-compose up -d neo4j
    print_green "Neo4j started."
    
    # Wait for Neo4j to be ready
    print_yellow "Waiting for Neo4j to be ready..."
    sleep 15
    print_green "Neo4j is ready."
fi

if [ "$SKIP_BUILDER" = false ]; then
    print_yellow "Starting Knowledge Graph Builder..."
    docker-compose up -d kg-builder
    print_green "Knowledge Graph Builder started."
fi

if [ "$SKIP_ENRICHER" = false ]; then
    print_yellow "Starting Knowledge Graph Enricher..."
    docker-compose up -d kg-enricher
    print_green "Knowledge Graph Enricher started."
fi

if [ "$SKIP_FRONTEND" = false ]; then
    print_yellow "Starting Knowledge Graph Frontend..."
    docker-compose up -d kg-frontend
    print_green "Knowledge Graph Frontend started."
fi

print_green "Application started successfully."
print_yellow "You can access the Neo4j browser at http://localhost:7474"
print_yellow "You can access the Knowledge Graph Frontend at http://localhost:8080"

# Show application status
print_yellow "Application status:"
docker-compose ps 