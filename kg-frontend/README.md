# Knowledge Graph Visualizer

A Three.js-based frontend for visualizing and interacting with the knowledge graph.

## Features

- 3D visualization of the knowledge graph using Three.js
- Controls for starting and stopping the knowledge graph builder
- Controls for starting and stopping the knowledge graph enricher
- Manual concept linking interface
- Real-time statistics about the knowledge graph

## Architecture

The frontend consists of two main components:

1. **Go Backend API Server**: Handles API requests from the frontend, communicates with Neo4j, and manages the knowledge graph builder and enricher.
2. **HTML/CSS/JS Frontend**: Provides the user interface and 3D visualization using Three.js.

## API Endpoints

- `GET /api/graph`: Get the current graph data from Neo4j
- `POST /api/builder/start`: Start the knowledge graph builder
- `POST /api/builder/stop`: Stop the knowledge graph builder
- `POST /api/enricher/start`: Start the knowledge graph enricher
- `POST /api/enricher/stop`: Stop the knowledge graph enricher
- `GET /api/concepts/search`: Search for concepts in the graph
- `POST /api/relationships`: Create a relationship between two concepts
- `GET /api/statistics`: Get statistics about the graph

## Development

### Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- Neo4j 4.4

### Building and Running

1. Build the Docker image:
   ```
   docker-compose build kg-frontend
   ```

2. Run the frontend:
   ```
   docker-compose up -d kg-frontend
   ```

3. Access the frontend at http://localhost:8080

## License

MIT 