#!/bin/bash

# kg.sh - Unified script for Kay-Gee-Go application management
# This script provides commands to start, stop, and manage the Kay-Gee-Go application

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
SKIP_OPTIMIZATION=false

# Function to show help
show_help() {
    echo "Usage: $0 COMMAND [options]"
    echo ""
    echo "Commands:"
    echo "  start       Start the Kay-Gee-Go application"
    echo "  stop        Stop the Kay-Gee-Go application"
    echo "  restart     Restart the Kay-Gee-Go application"
    echo "  status      Show the status of the Kay-Gee-Go application"
    echo "  test        Run tests for the Kay-Gee-Go application"
    echo "  logs        Show logs for the Kay-Gee-Go application"
    echo "  optimize    Optimize the repository size"
    echo ""
    echo "Options for 'start' command:"
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
    echo "  --skip-optimization       Skip the automatic optimization step"
    echo ""
    echo "Options for 'logs' command:"
    echo "  --service=NAME            Show logs for a specific service (neo4j, builder, enricher, frontend)"
    echo "  --follow                  Follow log output"
    echo ""
    echo "Options for 'optimize' command:"
    echo "  --aggressive              Run aggressive optimization (removes all cache files)"
    echo "  --keep-examples=N         Number of example cache files to keep (default: 5)"
    echo ""
    echo "Examples:"
    echo "  $0 start                  Start all services with default settings"
    echo "  $0 start --seed=\"Machine Learning\" --max-nodes=200"
    echo "  $0 stop                   Stop all services"
    echo "  $0 status                 Show status of all services"
    echo "  $0 test                   Run all tests"
    echo "  $0 logs --service=builder --follow"
    echo "  $0 optimize               Optimize the repository size"
}

# Function to optimize the repository
optimize_repo() {
    AGGRESSIVE=false
    KEEP_EXAMPLES=5
    
    # Parse command-line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --aggressive)
                AGGRESSIVE=true
                shift
                ;;
            --keep-examples=*)
                KEEP_EXAMPLES="${1#*=}"
                shift
                ;;
            *)
                print_red "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_yellow "Running repository optimization..."
    
    # Clean up cache files
    if [ "$AGGRESSIVE" = true ]; then
        print_yellow "Performing aggressive optimization (removing all cache files)..."
        
        # Remove all cache files
        find ./kg-builder/cache -name "*.json" -delete
        find ./kg-enricher/cache -name "*.json" -delete
    else
        print_yellow "Performing standard optimization..."
        
        # Count cache files
        BUILDER_CACHE_COUNT=$(find ./kg-builder/cache -name "*.json" | wc -l)
        ENRICHER_CACHE_COUNT=$(find ./kg-enricher/cache -name "*.json" | wc -l)
        
        # Only optimize if there are more than 100 cache files
        if [ "$BUILDER_CACHE_COUNT" -gt 100 ]; then
            print_yellow "Optimizing builder cache ($BUILDER_CACHE_COUNT files found)..."
            
            # Create backup directory for important cache files
            mkdir -p cache/builder_examples
            
            # Keep a few example cache files
            find ./kg-builder/cache -name "*.json" | head -$KEEP_EXAMPLES | xargs -I{} cp {} cache/builder_examples/
            
            # Remove excess cache files (keep the 50 most recent)
            find ./kg-builder/cache -name "*.json" | sort -r | tail -n +51 | xargs rm -f
        fi
        
        if [ "$ENRICHER_CACHE_COUNT" -gt 100 ]; then
            print_yellow "Optimizing enricher cache ($ENRICHER_CACHE_COUNT files found)..."
            
            # Create backup directory for important cache files
            mkdir -p cache/enricher_examples
            
            # Keep a few example cache files
            find ./kg-enricher/cache -name "*.json" | head -$KEEP_EXAMPLES | xargs -I{} cp {} cache/enricher_examples/
            
            # Remove excess cache files (keep the 50 most recent)
            find ./kg-enricher/cache -name "*.json" | sort -r | tail -n +51 | xargs rm -f
        fi
    fi
    
    # Remove compiled binaries
    print_yellow "Removing compiled binaries..."
    rm -f ./kg-enricher/enricher
    rm -f ./kg-builder/builder
    
    # Add binaries to .gitignore if not already there
    if ! grep -q "kg-enricher/enricher" .gitignore; then
        echo "" >> .gitignore
        echo "# Binaries" >> .gitignore
        echo "kg-enricher/enricher" >> .gitignore
        echo "kg-builder/builder" >> .gitignore
        echo "*/*/enricher" >> .gitignore
        echo "*/*/builder" >> .gitignore
    fi
    
    print_green "Repository optimization completed!"
}

# Function to start the application
start_app() {
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
            --skip-optimization)
                SKIP_OPTIMIZATION=true
                shift
                ;;
            *)
                print_red "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # Run optimization before starting (unless skipped)
    if [ "$SKIP_OPTIMIZATION" = false ]; then
        optimize_repo
    fi

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
        docker-compose up -d builder
        print_green "Knowledge Graph Builder started."
    fi

    if [ "$SKIP_ENRICHER" = false ]; then
        print_yellow "Starting Knowledge Graph Enricher..."
        docker-compose up -d enricher
        print_green "Knowledge Graph Enricher started."
    fi

    if [ "$SKIP_FRONTEND" = false ]; then
        print_yellow "Starting Knowledge Graph Frontend..."
        docker-compose up -d frontend
        print_green "Knowledge Graph Frontend started."
    fi

    print_green "Application started successfully."
    print_yellow "You can access the Neo4j browser at http://localhost:7474"
    print_yellow "You can access the Knowledge Graph Frontend at http://localhost:8080"

    # Show application status
    print_yellow "Application status:"
    docker-compose ps
}

# Function to stop the application
stop_app() {
    print_yellow "Stopping Kay-Gee-Go application..."

    # Stop all services defined in docker-compose.yml
    print_yellow "Stopping docker-compose services..."
    docker-compose down

    # Check if any individual containers are still running and stop them
    CONTAINERS=("kaygeego-builder" "kaygeego-enricher" "kaygeego-neo4j" "kaygeego-frontend" "kaygeego-wait-for-neo4j")

    for container in "${CONTAINERS[@]}"; do
        if [ "$(docker ps -q -f name=$container)" ]; then
            print_yellow "Stopping $container container..."
            docker stop $container
        fi
    done

    # Remove any stopped containers that might still exist
    for container in "${CONTAINERS[@]}"; do
        if [ "$(docker ps -a -q -f name=$container)" ]; then
            print_yellow "Removing $container container..."
            docker rm $container
        fi
    done

    print_green "Kay-Gee-Go application has been stopped successfully."
}

# Function to show application status
show_status() {
    print_yellow "Kay-Gee-Go application status:"
    docker-compose ps
}

# Function to run tests
run_tests() {
    print_yellow "Running tests..."
    go test -v ./internal/...

    # Check if tests passed
    if [ $? -eq 0 ]; then
        print_green "All tests passed!"
        return 0
    else
        print_red "Tests failed!"
        return 1
    fi
}

# Function to show logs
show_logs() {
    SERVICE=""
    FOLLOW=""
    
    # Parse command-line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --service=*)
                SERVICE="${1#*=}"
                shift
                ;;
            --follow)
                FOLLOW="--follow"
                shift
                ;;
            *)
                print_red "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    if [ -z "$SERVICE" ]; then
        print_yellow "Showing logs for all services..."
        docker-compose logs $FOLLOW
    else
        print_yellow "Showing logs for $SERVICE service..."
        docker-compose logs $FOLLOW $SERVICE
    fi
}

# Main script logic
if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

COMMAND=$1
shift

case $COMMAND in
    start)
        start_app "$@"
        ;;
    stop)
        stop_app
        ;;
    restart)
        stop_app
        start_app "$@"
        ;;
    status)
        show_status
        ;;
    test)
        run_tests
        exit $?
        ;;
    logs)
        show_logs "$@"
        ;;
    optimize)
        optimize_repo "$@"
        ;;
    help)
        show_help
        ;;
    *)
        print_red "Unknown command: $COMMAND"
        show_help
        exit 1
        ;;
esac

exit 0 