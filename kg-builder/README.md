# Knowledge Graph Builder

This project builds a knowledge graph using Neo4j and an LLM service.

## Setup

1. Ensure Docker and Docker Compose are installed on your system.
2. Run `docker-compose up --build` to start the application and Neo4j.

## Usage

The application will automatically start building the knowledge graph from the seed concept "Artificial Intelligence".

## Project Structure

- `cmd/kg-builder/`: Main application entry point
- `internal/neo4j/`: Neo4j connection and operations
- `internal/llm/`: LLM service interactions
- `internal/graph/`: Graph operations and data structures
