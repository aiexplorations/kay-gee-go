# Knowledge Graph Builder

This project builds a knowledge graph using Neo4j and an LLM service. It is designed to be run in a Docker container.  
Large language models are used to retrieve related concepts and mine relationships between concepts. The purpose of this project is to build a knowledge graph that can be used for further analysis or as a foundation for a semantic search engine.

Version: 0.2.0 - See [CHANGELOG.md](CHANGELOG.md) for details.

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

## Setup

1. Ensure Docker and Docker Compose are installed on your system.
2. Navigate to the project directory.
3. You can run the components separately or together:
   - Run `./start.sh` to run tests and start the Knowledge Graph Builder only.
   - Run `./start-enricher.sh` to start the Knowledge Graph Enricher only.
   - Run `./start-all.sh` to start both components together.

## Usage

### Integrated Deployment

You can run both the Knowledge Graph Builder and Enricher together using the integrated deployment:

```bash
./start-all.sh
```

This will start both services with their default configurations. You can also pass command-line arguments to customize the behavior of both services:

```bash
./start-all.sh --seed="Machine Learning" --max-nodes=200 --run-once --count=20
```

To stop both services:

```bash
./stop-all.sh
```

### Knowledge Graph Builder

The Knowledge Graph Builder will automatically start building the knowledge graph from the seed concept "Artificial Intelligence" (default).

#### Scripts

The following scripts are provided for convenience:

- `start.sh`: Runs tests and starts the application
  - Accepts command-line arguments: `./start.sh --seed="Machine Learning" --max-nodes=200 --timeout=60`
- `stop.sh`: Stops the application
- `status.sh`: Checks the status of the application and shows statistics
- `update-model.sh`: Updates the LLM model in the configuration file
  - Usage: `./update-model.sh <model_name>`
  - Example: `./update-model.sh llama3.1:latest`

#### Command-line Arguments

You can customize the behavior of the Knowledge Graph Builder using the following command-line arguments:

- `--seed`: Seed concept for graph building (default: "Artificial Intelligence")
- `--max-nodes`: Maximum number of nodes to build (default: 100)
- `--timeout`: Timeout in minutes for graph building (default: 30)
- `--random-relationships`: Number of random relationships to mine (default: 50)
- `--concurrency`: Number of concurrent workers for mining random relationships (default: 5)
- `--stats-only`: Only show statistics without building the graph
- `--version`: Show version information and exit

Example:
```bash
docker-compose run kg-builder /kg-builder --seed "Machine Learning" --max-nodes 200 --timeout 60
```

Or using the start.sh script:
```bash
./start.sh --seed="Machine Learning" --max-nodes=200 --timeout=60
```

### Knowledge Graph Enricher

The Knowledge Graph Enricher will automatically start enriching the knowledge graph based on the configuration in `config.yaml`.

#### Scripts

The following scripts are provided for convenience:

- `start-enricher.sh`: Starts the enricher application
  - Accepts command-line arguments: `./start-enricher.sh --run-once --count=20`
- `stop-enricher.sh`: Stops the enricher application

#### Command-line Arguments

You can customize the behavior of the Knowledge Graph Enricher using the following command-line arguments:

- `--run-once`: Run once and exit
- `--count`: Number of relationships to mine when running once (default: 10)
- `--stats`: Show statistics and exit
- `--version`: Show version information and exit

Example:
```bash
./start-enricher.sh --run-once --count=20
```

## Configuration

Both components can be configured using environment variables or configuration files.

### Environment Variables

#### Knowledge Graph Builder

- `NEO4J_URI`: URI for the Neo4j database (default: "bolt://neo4j:7687")
- `NEO4J_USER`: Username for the Neo4j database (default: "neo4j")
- `NEO4J_PASSWORD`: Password for the Neo4j database (default: "password")
- `LLM_URL`: URL for the LLM service (default: "http://host.docker.internal:11434/api/generate")
- `LLM_MODEL`: Model to use for the LLM service (default: "qwen2.5:3b")
- `CONFIG_FILE`: Path to the configuration file (default: "config.yaml")

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

### Configuration Files

Both components can also be configured using YAML configuration files. By default, they look for files named `config.yaml` in their respective directories.

#### Knowledge Graph Builder

```yaml
# Neo4j Configuration
neo4j:
  uri: "bolt://neo4j:7687"
  user: "neo4j"
  password: "password"
  max_retries: 5
  retry_interval_seconds: 5

# LLM Configuration
llm:
  url: "http://host.docker.internal:11434/api/generate"
  model: "qwen2.5:3b"
  cache_dir: "./cache/llm"

# Graph Configuration
graph:
  seed_concept: "Artificial Intelligence"
  max_nodes: 100
  timeout_minutes: 30
  worker_count: 10
  random_relationships: 50
  concurrency: 5
```

#### Knowledge Graph Enricher

```yaml
# Neo4j Configuration
neo4j:
  uri: "bolt://neo4j:7687"
  user: "neo4j"
  password: "password"
  max_retries: 5
  retry_interval_seconds: 5

# LLM Configuration
llm:
  url: "http://host.docker.internal:11434/api/generate"
  model: "qwen2.5:3b"
  cache_dir: "./cache/llm"

# Enricher Configuration
enricher:
  batch_size: 10
  interval_seconds: 60
  max_relationships: 100
  concurrency: 5
```

## Persistence

### Neo4j Data

The Neo4j database data is persisted in Docker volumes, so the graph will be preserved even if the containers are stopped or removed.

### LLM Cache

Responses from the LLM service are cached in the `cache` directory, which is mounted as a volume in the Docker container. This allows for:

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

You can run the tests using the provided Makefile or the run-tests.sh script:

```bash
# Using the run-tests.sh script (automatically chooses local or Docker based on Go installation)
./kg-builder/run-tests.sh
./kg-enricher/run-tests.sh

# Test application build and basic functionality
cd kg-builder
./test-build.sh
# or
make test-build

# Using Makefile (requires Go installation)
cd kg-builder
make test

# Run tests in short mode (skip integration tests)
make test-short

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run tests in Docker (no Go installation required)
make docker-test

# Clean Docker test containers
make docker-test-clean
```

### Test Coverage

The test suite covers the following areas:

1. **Error Handling**: Tests for the custom error types and retry mechanisms
2. **LLM Service**: Tests for the LLM service integration and caching
3. **Neo4j Integration**: Tests for the Neo4j database connection and operations
4. **Graph Building**: Tests for the graph building logic
5. **Configuration**: Tests for environment variables and command-line arguments
6. **Enricher**: Tests for the enricher functionality and statistics
7. **End-to-End**: Integration tests for the complete workflow

## Project Structure

### Knowledge Graph Builder
- `cmd/kg-builder/`: Main application entry point
- `internal/neo4j/`: Neo4j connection and operations
- `internal/llm/`: LLM service interactions
- `internal/graph/`: Graph operations and data structures
- `internal/models/`: Data models
- `internal/errors/`: Error handling
- `internal/config/`: Configuration handling
- `internal/tests/`: Test suite
- `cache/`: Cached LLM responses

### Knowledge Graph Enricher
- `cmd/kg-enricher/`: Main application entry point
- `internal/neo4j/`: Neo4j connection and operations
- `internal/llm/`: LLM service interactions
- `internal/enricher/`: Enricher operations
- `internal/models/`: Data models
- `internal/config/`: Configuration handling
- `cache/`: Cached LLM responses

## File Descriptions

### Knowledge Graph Builder

#### `internal/graph/graph.go`
This file contains the implementation of the `GraphBuilder` struct, which is responsible for building a knowledge graph using concepts and their relationships. 

- **GraphBuilder Struct**: Holds the Neo4j driver, functions for retrieving related concepts and mining relationships, a map of processed concepts, a node count, and a mutex for thread safety.
  
- **BuildGraph**: The main method that builds the knowledge graph starting from a seed concept. It uses goroutines to process concepts concurrently, managing a queue of concepts to explore.

- **MineRandomRelationships**: This method mines relationships between random pairs of concepts, using concurrency to speed up the process.

#### `internal/llm/llm.go`
This file contains functions that interact with a language model (LLM) service to retrieve related concepts and mine relationships between concepts.

- **GetRelatedConcepts**: Sends a request to the LLM service with a prompt to get related concepts for a given concept.

- **MineRelationship**: Similar to `GetRelatedConcepts`, this function sends a request to the LLM service to determine if there is a relationship between two concepts.

- **Caching**: Implements in-memory and file-based caching of LLM responses for improved performance and offline access.

#### `internal/neo4j/neo4j.go`
This file handles the connection to the Neo4j database and provides functions to create relationships between concepts.

- **SetupNeo4jConnection**: Establishes a connection to the Neo4j database with retry logic to handle connection failures.

- **CreateRelationship**: A function that creates a relationship between two concepts in the Neo4j database using a Cypher query.

- **QueryConcepts**: Retrieves all concepts from the Neo4j database.

- **QueryRelationships**: Retrieves all relationships from the Neo4j database.

### Knowledge Graph Enricher

#### `internal/enricher/enricher.go`
This file contains the implementation of the `Enricher` struct, which is responsible for enriching the knowledge graph with new relationships.

- **Enricher Struct**: Holds the Neo4j driver, configuration, and statistics about the enrichment process.

- **Start**: Starts the enricher service, which runs in a loop processing batches of concept pairs.

- **RunOnce**: Runs the enricher once, processing a specified number of concept pairs.

- **processBatch**: Processes a batch of concept pairs, finding and creating relationships between them.

#### `internal/llm/llm.go`
This file contains functions that interact with a language model (LLM) service to find relationships between concepts.

- **FindRelationship**: Sends a request to the LLM service to determine if there is a relationship between two concepts.

- **Caching**: Implements in-memory and file-based caching of LLM responses for improved performance and offline access.

#### `internal/neo4j/neo4j.go`
This file handles the connection to the Neo4j database and provides functions to query and create relationships.

- **SetupNeo4jConnection**: Establishes a connection to the Neo4j database with retry logic to handle connection failures.

- **QueryRandomConceptPairs**: Retrieves random pairs of concepts from the Neo4j database.

- **CreateRelationship**: Creates a relationship between two concepts in the Neo4j database.

## Build process

1. Run `docker-compose up --build` to start the application and Neo4j.
2. The application will automatically start building the knowledge graph from the seed concept "Artificial Intelligence" (or the one specified via command-line arguments).
3. The graph will be persisted in the Neo4j database, and LLM responses will be cached in the `cache` directory.






