# KayGeeGo - Knowledge Graph Generator

This project builds a knowledge graph using Neo4j and an LLM service. It is designed to be run in a Docker container.  
Large language models are used to retrieve related concepts and mine relationships between concepts. The purpose of this project is to build a knowledge graph that can be used for further analysis or as a foundation for a semantic search engine.

Version: 0.6.1 - See [CHANGELOG.md](CHANGELOG.md) for details.

## Features

- Builds a knowledge graph starting from a seed concept
- Uses an LLM service to retrieve related concepts and mine relationships
- Persists the graph in a Neo4j database
- Caches LLM responses for offline access and faster rebuilding
- Provides robust error handling with retry mechanisms
- Configurable via command-line arguments
- Supports continuous enrichment of the knowledge graph
- Provides detailed statistics on graph building and enrichment processes
- Includes comprehensive test coverage for all components
- Supports low connectivity concept seeding to enhance graph diversity

## Components

The project consists of three main components:

1. **Builder**: Builds the initial knowledge graph from a seed concept.
   - Starts with a seed concept and expands outward by discovering related concepts
   - Creates nodes and relationships in the Neo4j database
   - Uses concurrent workers to speed up the graph building process
   - Implements timeout and maximum node count limits
   - Mines random relationships between existing concepts
   - Can use low connectivity concepts as seeds for more balanced graph expansion

2. **Enricher**: Enriches the existing knowledge graph by finding and adding relationships between random pairs of concepts.
   - Runs as a continuous service or in one-shot mode
   - Selects random pairs of concepts from the database
   - Uses LLM to determine if relationships exist between concepts
   - Creates new relationships in the Neo4j database
   - Provides detailed statistics on the enrichment process
   - Implements configurable batch processing with intervals

3. **Frontend**: Visualizes the knowledge graph in an interactive 3D interface.
   - Displays nodes and relationships in a 3D space
   - Allows for interactive exploration of the graph
   - Provides statistics on the graph structure
   - Supports searching and filtering of concepts

## Project Structure

```
.
├── kg-builder/            # Knowledge Graph Builder component
├── kg-enricher/           # Knowledge Graph Enricher component
├── kg-frontend/           # Knowledge Graph Frontend component
├── scripts/               # Scripts for running the application
│   ├── run.sh             # Script to run the application
│   ├── stop.sh            # Script to stop the application
│   └── status.sh          # Script to show the status of the application
├── docker-compose.yml     # Docker Compose configuration
├── Makefile               # Makefile for building and running the application
├── run.sh                 # Symbolic link to scripts/run.sh
├── stop.sh                # Symbolic link to scripts/stop.sh
└── status.sh              # Symbolic link to scripts/status.sh
```

## Setup

1. Ensure Docker and Docker Compose are installed on your system.
2. Navigate to the project directory.
3. Run `./run.sh` to start the application.

## Usage

The application can be managed using the unified `kg.sh` script:

```bash
# Start the application with default settings
./kg.sh start

# Start with custom settings
./kg.sh start --seed="Machine Learning" --max-nodes=200

# Stop the application
./kg.sh stop

# Restart the application
./kg.sh restart

# Show application status
./kg.sh status

# Run tests
./kg.sh test

# View logs
./kg.sh logs
./kg.sh logs --service=builder --follow

# Show help
./kg.sh help
```

For more details, run `./kg.sh help`.

## Configuration

The application can be configured using YAML configuration files or environment variables.

### Configuration Files

- `kg-builder/config.yaml`: Configuration for the Knowledge Graph Builder
- `kg-enricher/config.yaml`: Configuration for the Knowledge Graph Enricher

### Environment Variables

#### Knowledge Graph Builder

- `NEO4J_URI`: URI for the Neo4j database (default: "bolt://neo4j:7687")
- `NEO4J_USER`: Username for the Neo4j database (default: "neo4j")
- `NEO4J_PASSWORD`: Password for the Neo4j database (default: "password")
- `LLM_URL`: URL for the LLM service (default: "http://host.docker.internal:11434/api/generate")
- `LLM_MODEL`: Model to use for the LLM service (default: "qwen2.5:3b")

#### Knowledge Graph Enricher

- `NEO4J_URI`: URI for the Neo4j database (default: "bolt://neo4j:7687")
- `NEO4J_USER`: Username for the Neo4j database (default: "neo4j")
- `NEO4J_PASSWORD`: Password for the Neo4j database (default: "password")
- `LLM_URL`: URL for the LLM service (default: "http://host.docker.internal:11434/api/generate")
- `LLM_MODEL`: Model to use for the LLM service (default: "qwen2.5:3b")
- `ENRICHER_BATCH_SIZE`: Number of pairs to process in each batch (default: 10)
- `ENRICHER_INTERVAL_SECONDS`: Interval between batches in seconds (default: 60)
- `ENRICHER_MAX_RELATIONSHIPS`: Maximum number of relationships to create (default: 100)
- `ENRICHER_CONCURRENCY`: Number of concurrent workers for mining relationships (default: 5)

## Low Connectivity Feature

The low connectivity feature enhances graph diversity by targeting concepts with fewer connections:

- Uses `--use-low-connectivity` flag to enable the feature
- Identifies concepts with the lowest number of connections in the graph
- Uses these low connectivity concepts as seeds for subsequent graph building
- Creates a more balanced and comprehensive knowledge graph
- Prevents isolated clusters by creating pathways between distant concepts
- Improves overall graph connectivity and exploration potential

## Persistence

The Neo4j database data is persisted in Docker volumes, so the graph will be preserved even if the containers are stopped or removed.

### LLM Cache

Responses from the LLM service are cached in the cache directories, which are mounted as volumes in the Docker containers. This allows for:

1. Reduced API calls to the LLM service
2. Faster graph building by reusing previously retrieved concepts and relationships
3. Offline access to previously mined knowledge

To clear the cache, run:
```bash
rm -rf kg-builder/cache/*.json
rm -rf kg-enricher/cache/*.json
```

## Testing

The project includes a comprehensive test suite to ensure the application is correctly configured and functioning as expected.

### Running Tests

You can run the tests using the provided Makefile:

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage
```

## Troubleshooting

### Debugging Build Issues

If you encounter issues when starting the system, you can run the start script with the debug flag:

```bash
./run.sh --debug
```

This will provide more detailed error messages and show the Docker logs for any containers that fail to start.






