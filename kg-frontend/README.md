# Knowledge Graph Frontend

This is the frontend for the Knowledge Graph Visualizer application. It provides a user interface for interacting with the knowledge graph, including building, enriching, and visualizing the graph.

## Features

- Graph visualization with 3D force-directed layout
- Graph builder for creating new knowledge graphs
- Graph enricher for adding relationships to existing concepts
- Manual linking of concepts
- Statistics about the knowledge graph

## Prerequisites

- Node.js (v14 or higher)
- npm (v6 or higher)

## Installation

1. Clone the repository
2. Navigate to the `kg-frontend` directory
3. Install dependencies:

```bash
npm install
```

## Running the Application

To start the application, run:

```bash
npm start
```

This will start a local server on port 3000. You can access the application at http://localhost:3000.

## Running Tests

To run the tests, use:

```bash
npm test
```

To run the tests in watch mode (automatically re-run when files change):

```bash
npm run test:watch
```

## Test Coverage

The tests cover the following aspects of the application:

1. API client functionality
   - Connecting to backend endpoints
   - Handling responses and errors

2. Button click handlers
   - Graph Builder buttons (Start/Stop)
   - Graph Enricher buttons (Start/Stop)
   - Manual Linking button
   - Graph controls (Reset Camera, Refresh Graph)

3. Integration tests
   - Verifying that the UI correctly interacts with the backend
   - Validating input before making API calls
   - Handling API responses and errors

## Backend Integration

The frontend communicates with the backend through the following endpoints:

- `/api/graph` - Get graph data
- `/api/builder/start` - Start the graph builder
- `/api/builder/stop` - Stop the graph builder
- `/api/enricher/start` - Start the graph enricher
- `/api/enricher/stop` - Stop the graph enricher
- `/api/concepts/search` - Search for concepts
- `/api/relationships` - Create relationships between concepts
- `/api/statistics` - Get statistics about the knowledge graph

## Known Issues

- The builder and enricher functionality may not be available in all environments
- The graph visualization may be slow for large graphs
- The manual linking functionality requires exact concept names 