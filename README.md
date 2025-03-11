# Knowledge Graph Builder

This project builds a knowledge graph using Neo4j and an LLM service. It is designed to be run in a Docker container.  
Large language models are used to retrieve related concepts and mine relationships between concepts. The purpose of this project is to build a knowledge graph that can be used for further analysis or as a foundation for a semantic search engine.

Version: 0.5.2 - See [CHANGELOG.md](CHANGELOG.md) for details.

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

## Components

The project consists of two main components:

1. **Knowledge Graph Builder**: Builds the initial knowledge graph from a seed concept.
   - Starts with a seed concept and expands outward by discovering related concepts
   - Creates nodes and relationships in the Neo4j database
   - Uses concurrent workers to speed up the graph building process
   - Implements timeout and maximum node count limits
   - Mines random relationships between existing concepts

2. **Knowledge Graph Enricher**: Enriches the existing knowledge graph by finding and adding relationships between random pairs of concepts.
   - Runs as a continuous service or in one-shot mode
   - Selects random pairs of concepts from the database
   - Uses LLM to determine if relationships exist between concepts
   - Creates new relationships in the Neo4j database
   - Provides detailed statistics on the enrichment process
   - Implements configurable batch processing with intervals

## Project Structure

```
.
├── cmd/                    # Command-line applications
│   ├── builder/            # Knowledge Graph Builder application
│   ├── enricher/           # Knowledge Graph Enricher application
│   └── frontend/           # Knowledge Graph Frontend application
├── internal/               # Internal packages
│   ├── builder/            # Builder-specific code
│   ├── enricher/           # Enricher-specific code
│   ├── frontend/           # Frontend-specific code
│   └── common/             # Shared code
│       ├── config/         # Configuration handling
│       ├── errors/         # Error handling
│       ├── llm/            # LLM service interactions
│       ├── models/         # Data models
│       └── neo4j/          # Neo4j connection and operations
├── build/                  # Build-related files
│   ├── builder/            # Builder Dockerfile
│   ├── enricher/           # Enricher Dockerfile
│   └── frontend/           # Frontend Dockerfile
├── config/                 # Configuration files
│   ├── builder.yaml        # Builder configuration
│   └── enricher.yaml       # Enricher configuration
├── cache/                  # Cache directory
│   ├── builder/            # Builder cache
│   └── enricher/           # Enricher cache
├── scripts/                # Scripts for running the application
│   ├── run.sh              # Script to run the application
│   ├── stop.sh             # Script to stop the application
│   └── status.sh           # Script to show the status of the application
├── public/                 # Frontend static files
├── docker-compose.yml      # Docker Compose configuration
├── Makefile                # Makefile for building and running the application
├── go.mod                  # Go module file
├── go.sum                  # Go module checksum file
├── run.sh                  # Symbolic link to scripts/run.sh
├── stop.sh                 # Symbolic link to scripts/stop.sh
└── status.sh               # Symbolic link to scripts/status.sh
```

## Setup

1. Ensure Docker and Docker Compose are installed on your system.
2. Navigate to the project directory.
3. Run `./run.sh` to start the application.

## Usage

### Running the Application

You can run the application using the provided scripts:

```bash
# Start the application
./run.sh

# Stop the application
./stop.sh

# Show the status of the application
./status.sh
```

### Command-line Arguments

You can customize the behavior of the application using command-line arguments:

```bash
# Start the application with custom settings
./run.sh --seed="Machine Learning" --max-nodes=200 --timeout=60

# Start only specific components
./run.sh --skip-neo4j --skip-frontend

# Run the enricher once and exit
./run.sh --skip-builder --skip-frontend --run-once --count=20
```

### Using the Makefile

You can also use the Makefile to build and run the application:

```bash
# Build all components
make build

# Run all components
make run

# Run only specific components
make builder
make enricher
make frontend

# Stop the application
make stop

# Show the status of the application
make status

# Clean up build artifacts
make clean

# Run tests
make test

# Show help
make help
```

## Configuration

The application can be configured using YAML configuration files or environment variables.

### Configuration Files

- `config/builder.yaml`: Configuration for the Knowledge Graph Builder
- `config/enricher.yaml`: Configuration for the Knowledge Graph Enricher

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

## Persistence

The Neo4j database data is persisted in Docker volumes, so the graph will be preserved even if the containers are stopped or removed.

### LLM Cache

Responses from the LLM service are cached in the `cache` directory, which is mounted as a volume in the Docker container. This allows for:

1. Reduced API calls to the LLM service
2. Faster graph building by reusing previously retrieved concepts and relationships
3. Offline access to previously mined knowledge

To clear the cache, run:
```bash
rm -rf cache/builder/*.json
rm -rf cache/enricher/*.json
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






