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

# Show application status
print_yellow "Application status:"
docker-compose ps

# Show Neo4j statistics
if [ "$(docker ps -q -f name=neo4j)" ]; then
    print_yellow "Neo4j statistics:"
    docker exec neo4j cypher-shell -u neo4j -p password "MATCH (n) RETURN count(n) AS nodes" | grep -v "nodes"
    docker exec neo4j cypher-shell -u neo4j -p password "MATCH ()-[r]->() RETURN count(r) AS relationships" | grep -v "relationships"
    docker exec neo4j cypher-shell -u neo4j -p password "MATCH (n) RETURN labels(n) AS label, count(*) AS count ORDER BY count DESC" | grep -v "label"
    docker exec neo4j cypher-shell -u neo4j -p password "MATCH ()-[r]->() RETURN type(r) AS type, count(*) AS count ORDER BY count DESC" | grep -v "type"
else
    print_red "Neo4j is not running."
fi 