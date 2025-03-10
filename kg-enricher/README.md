# Knowledge Graph Enricher

This microservice enriches a knowledge graph by finding and adding relationships between random pairs of concepts. It uses a large language model (LLM) to determine if there is a relationship between two concepts, and if so, what type of relationship it is.

## Features

- Periodically selects random pairs of concepts from the knowledge graph
- Uses an LLM to determine if there is a relationship between the concepts
- Adds the relationship to the knowledge graph if one is found
- Configurable batch size, interval, and concurrency
- Caches LLM responses for faster processing and reduced API calls
- Provides statistics about the enrichment process

## Setup

1. Ensure Docker and Docker Compose are installed on your system.
2. Navigate to the project directory.
3. Run `docker-compose up -d` to start the application.

## Usage

The application will automatically start enriching the knowledge graph based on the configuration in `config.yaml`.

### Command-line Arguments

You can customize the behavior of the application using the following command-line arguments:

- `--run-once`: Run once and exit
- `--count`: Number of relationships to mine when running once (default: 10)
- `--stats`: Show statistics and exit
- `--version`: Show version information and exit

Example:
```bash
docker-compose run kg-enricher --run-once --count 20
```

### Configuration

The application can be configured using environment variables or a configuration file (`config.yaml`).

#### Environment Variables

- `NEO4J_URI`: URI for the Neo4j database (default: "bolt://neo4j:7687")
- `NEO4J_USER`: Username for the Neo4j database (default: "neo4j")
- `NEO4J_PASSWORD`: Password for the Neo4j database (default: "password")
- `LLM_URL`: URL for the LLM service (default: "http://host.docker.internal:11434/api/generate")
- `LLM_MODEL`: Model to use for the LLM service (default: "qwen2.5:3b")
- `ENRICHER_BATCH_SIZE`: Number of pairs to process in each batch (default: 10)
- `ENRICHER_INTERVAL_SECONDS`: Interval between batches in seconds (default: 60)
- `ENRICHER_MAX_RELATIONSHIPS`: Maximum number of relationships to create (default: 100)
- `ENRICHER_CONCURRENCY`: Number of concurrent workers for mining relationships (default: 5)

#### Configuration File

The application can also be configured using a YAML configuration file. By default, the application looks for a file named `config.yaml` in the current directory.

Example configuration file:

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

## How It Works

1. The enricher selects a batch of random nodes from the Neo4j database.
2. It creates pairs of nodes and checks if a relationship already exists between them.
3. For each pair without an existing relationship, it uses the LLM to determine if there is a relationship between the concepts.
4. If a relationship is found, it is added to the Neo4j database.
5. The process repeats at the configured interval.

## Caching

Responses from the LLM service are cached to reduce API calls and improve performance. The cache is stored in the `cache/llm` directory, which is mounted as a volume in the Docker container.

## Statistics

You can view statistics about the enrichment process by running:

```bash
docker-compose run kg-enricher --stats
```

This will show:
- Total concepts in the graph
- Total relationships in the graph
- Total batches processed
- Total pairs processed
- Total relationships found
- Total relationships created
- Last batch time
- Running time

## Development

### Building the Application

```bash
go build -o kg-enricher ./cmd/kg-enricher
```

### Running Tests

```bash
go test ./...
```

### Building the Docker Image

```bash
docker build -t kg-enricher .
``` 