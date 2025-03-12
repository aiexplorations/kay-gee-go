# KayGeeGo - Knowledge Graph Generator

KayGeeGo is a comprehensive knowledge graph generation and visualization system that uses Neo4j and Large Language Models (LLMs) to build, enrich, and visualize knowledge graphs. The system automatically discovers related concepts, mines relationships between them, and presents the resulting graph in an interactive 3D visualization.

**Version:** 0.6.1 - See [CHANGELOG.md](CHANGELOG.md) for details.

## Overview

KayGeeGo builds a knowledge graph starting from a seed concept, using LLMs to discover related concepts and relationships. The graph is stored in a Neo4j database and can be visualized through an interactive 3D interface. The system is designed to run in Docker containers and is managed through a unified shell script.

<div align="center">
  <img src="img/kg_1.png" alt="Knowledge Graph Visualization" width="1024">
</div>

## Features
- **Automated Knowledge Graph Construction**: Builds a knowledge graph starting from a seed concept
- **LLM-Powered Relationship Mining**: Uses LLMs to discover relationships between concepts
- **Interactive 3D Visualization**: Visualizes the graph in an interactive 3D interface
- **Continuous Enrichment**: Automatically enriches the graph by finding new relationships
- **Low Connectivity Seeding**: Enhances graph diversity by targeting concepts with fewer connections
- **Caching System**: Caches LLM responses for offline access and faster rebuilding
- **Comprehensive Configuration**: Configurable via command-line arguments and configuration files
- **Unified Management**: Single script (`kg.sh`) for all operations
- **Concurrent Processing**: Implements multi-worker architecture for faster graph building
- **Statistics Dashboard**: Provides real-time metrics on graph structure and growth
- **Search Functionality**: Enables finding specific concepts within the knowledge graph
- **Manual Relationship Creation**: Allows users to manually define relationships between concepts
- **Configurable LLM Integration**: Supports different LLM providers and models
- **Docker Containerization**: Runs all components in isolated Docker containers
- **Neo4j Database Backend**: Leverages the power of graph databases for efficient storage and querying

## System Requirements
- **Docker** and **Docker Compose**
- **LLM Service**: Local or remote LLM service (default configuration uses Ollama with qwen2.5:3b model)
- **Disk Space**: At least 4GB of free disk space for Docker images, volumes, and cached LLM responses
- **Memory**: At least 6GB of RAM (12GB recommended for optimal performance with larger graphs)
- **CPU**: Multi-core processor recommended for concurrent processing operations
- **Network**: Internet connection required for remote LLM services (not needed for local Ollama setup)
- **Browser**: Modern web browser with WebGL support for 3D visualization

## Components

The project consists of three main components:
1. **Builder (`kg-builder`)**: Builds the initial knowledge graph from a seed concept
   - Discovers related concepts and creates nodes and relationships
   - Implements concurrent workers for faster graph building
   - Supports timeout and maximum node count limits
   - Mines random relationships between existing concepts
   - Uses low connectivity concepts as seeds for balanced graph expansion
   - Caches LLM responses for faster rebuilding and offline operation
   - Provides detailed logging of the building process

2. **Enricher (`kg-enricher`)**: Enriches the existing knowledge graph
   - Runs as a continuous service or in one-shot mode
   - Selects random pairs of concepts and finds relationships between them
   - Provides detailed statistics on the enrichment process
   - Implements configurable batch processing with intervals
   - Focuses on low-connectivity nodes to improve graph cohesion
   - Supports different enrichment strategies (random, targeted, similarity-based)

3. **Frontend (`kg-frontend`)**: Visualizes the knowledge graph
   - Displays nodes and relationships in an interactive 3D space
   - Supports zooming, panning, and rotating the visualization
   - Provides real-time statistics on the graph structure
   - Enables searching and filtering of concepts
   - Allows manual creation of relationships between concepts
   - Features a dark/light theme toggle for better visibility
   - Displays node details and relationship information on hover/click
   - Includes a mini-map for easier navigation of large graphs

## Quick Start

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/kay-gee-go.git
   cd kay-gee-go
   ```

2. Make sure the script is executable:
   ```bash
   chmod +x kg.sh
   ```

### Starting the Application

Start the application with default settings:
```bash
./kg.sh start
```

Start with custom settings:
```bash
./kg.sh start --seed="Machine Learning" --max-nodes=200 --timeout=60
```

### Accessing the Application

- **Neo4j Browser**: http://localhost:7474 (username: neo4j, password: password)
- **Knowledge Graph Frontend**: http://localhost:8080

### Stopping the Application

```bash
./kg.sh stop
```

## Detailed Usage

The application is managed using the unified `kg.sh` script:

```bash
# Start the application with default settings
./kg.sh start

# Start with custom settings
./kg.sh start --seed="Machine Learning" --max-nodes=200 --timeout=60 --random-relationships=100 --concurrency=10

# Start only specific components
./kg.sh start --skip-builder --skip-enricher  # Only start Neo4j and frontend

# Start in stats-only mode (no graph building)
./kg.sh start --stats-only

# Run the enricher once with a specific number of relationships
./kg.sh start --run-once --count=50

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

# Optimize the repository (clean up cache files)
./kg.sh optimize
./kg.sh optimize --aggressive  # Remove all cache files

# Show help
./kg.sh help
```

## Configuration Options

### Command-Line Options

When starting the application with `./kg.sh start`, you can use the following options:

| Option | Description | Default |
|--------|-------------|---------|
| `--seed=<concept>` | Seed concept to start with | "Artificial Intelligence" |
| `--max-nodes=<number>` | Maximum number of nodes to create | 100 |
| `--timeout=<minutes>` | Timeout in minutes | 30 |
| `--random-relationships=<number>` | Number of random relationships to create | 50 |
| `--concurrency=<number>` | Number of concurrent operations | 5 |
| `--stats-only` | Only show statistics, don't build graph | false |
| `--skip-neo4j` | Skip starting Neo4j | false |
| `--skip-builder` | Skip starting the builder | false |
| `--skip-enricher` | Skip starting the enricher | false |
| `--skip-frontend` | Skip starting the frontend | false |
| `--run-once` | Run the enricher once and exit | false |
| `--count=<number>` | Number of relationships to create when using --run-once | 10 |
| `--skip-optimization` | Skip repository optimization | false |
| `--use-low-connectivity` | Use low connectivity concepts as seeds | false |

### Configuration Files

The application can also be configured using YAML configuration files:

- `config/builder.yaml`: Configuration for the Knowledge Graph Builder
- `config/enricher.yaml`: Configuration for the Knowledge Graph Enricher

### Environment Variables

The Docker Compose configuration uses the following environment variables:

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

## Data Persistence

### Neo4j Database

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

Alternatively, use the optimize command with the aggressive flag:
```bash
./kg.sh optimize --aggressive
```

## Repository Optimization

The repository includes an optimization feature to manage cache files and keep the repository size under control:

```bash
# Standard optimization (keeps recent cache files)
./kg.sh optimize

# Aggressive optimization (removes all cache files)
./kg.sh optimize --aggressive

# Keep a specific number of example cache files
./kg.sh optimize --keep-examples=20
```

## Testing

The project includes a comprehensive test suite to ensure the application is correctly configured and functioning as expected.

### Running Tests

You can run the tests using the provided script:

```bash
# Run all tests
./kg.sh test
```

## Troubleshooting

### Common Issues

1. **Neo4j Connection Issues**:
   - Ensure Neo4j is running: `docker ps | grep neo4j`
   - Check Neo4j logs: `./kg.sh logs --service=neo4j`
   - Verify Neo4j credentials in configuration files

2. **LLM Service Issues**:
   - Ensure your LLM service is running and accessible
   - Check the LLM URL in configuration files
   - Verify the model name is correct

3. **Docker Issues**:
   - Ensure Docker and Docker Compose are installed and running
   - Check Docker logs: `docker logs kaygeego-neo4j`
   - Restart Docker if necessary

### Debugging

If you encounter issues when starting the system, you can check the logs:

```bash
# View all logs
./kg.sh logs

# View logs for a specific service
./kg.sh logs --service=builder

# Follow logs in real-time
./kg.sh logs --follow
```

## License

MIT






