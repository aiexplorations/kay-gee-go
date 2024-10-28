# Knowledge Graph Builder

This project builds a knowledge graph using Neo4j and an LLM service. It is designed to be run in a Docker container.  
Large language models are used to retrieve related concepts and mine relationships between concepts. The purpose of this project is to build a knowledge graph that can be used for further analysis or as a foundation for a semantic search engine. 

## Setup

1. Ensure Docker and Docker Compose are installed on your system.
2. Navigate to the project directory `kg-builder`.
2. Run `docker-compose up --build` to start the application and Neo4j.

## Usage

The application will automatically start building the knowledge graph from the seed concept "Artificial Intelligence".

## Project Structure

- `cmd/kg-builder/`: Main application entry point
- `internal/neo4j/`: Neo4j connection and operations
- `internal/llm/`: LLM service interactions
- `internal/graph/`: Graph operations and data structures

## File Descriptions

### `internal/graph/graph.go`
This file contains the implementation of the `GraphBuilder` struct, which is responsible for building a knowledge graph using concepts and their relationships. 

- **GraphBuilder Struct**: Holds the Neo4j driver, functions for retrieving related concepts and mining relationships, a map of processed concepts, a node count, and a mutex for thread safety.
  
- **NewGraphBuilder**: A constructor function that initializes a new `GraphBuilder` instance with the provided Neo4j driver and functions.

- **BuildGraph**: The main method that builds the knowledge graph starting from a seed concept. It uses goroutines to process concepts concurrently, managing a queue of concepts to explore. It logs the progress and handles timeouts.

- **worker**: A method that processes concepts from the queue, retrieves related concepts, and creates relationships in the Neo4j database. It ensures that the number of processed concepts does not exceed a specified limit.

- **MineRandomRelationships**: This method mines relationships between random pairs of concepts, using concurrency to speed up the process.

- **getRandomPair**: A helper method that retrieves a random pair of processed concepts for relationship mining.

### `internal/llm/llm.go`
This file contains functions that interact with a language model (LLM) service to retrieve related concepts and mine relationships between concepts.

- **GetRelatedConcepts**: Sends a request to the LLM service with a prompt to get related concepts for a given concept. It expects a JSON response containing related concepts and their relationships. This is one of the main functions that is used to mine relationships between concepts. There is a prompt template that is used to generate the prompt for the LLM service. This can be modified to change the behavior of the LLM service. 

- **MineRelationship**: Similar to `GetRelatedConcepts`, this function sends a request to the LLM service to determine if there is a relationship between two concepts. It returns the relationship details if found. The idea is that this will be used to mine relationships between concepts that have already been added to the graph. 

### `internal/models/models.go`
This file defines the `Concept` struct, which represents a concept in the knowledge graph.

- **Concept Struct**: Contains three fields: `Name`, `Relation`, and `RelatedTo`, which are used to store the name of the concept, the type of relationship, and the concept it is related to, respectively. The struct is annotated for JSON serialization.

### `internal/neo4j/neo4j.go`
This file handles the connection to the Neo4j database and provides functions to create relationships between concepts.

- **SetupNeo4jConnection**: Establishes a connection to the Neo4j database with retry logic to handle connection failures.

- **CreateRelationship**: A function that creates a relationship between two concepts in the Neo4j database using a Cypher query. It ensures that the concepts are created if they do not already exist.

- **connectToNeo4jWithRetry**: A helper function that attempts to connect to the Neo4j database multiple times, logging the attempts and errors. It validates the connection parameters before attempting to connect.


# Build process

1. Run `docker-compose up --build` to start the application and Neo4j.
2. The application will automatically start building the knowledge graph from the seed concept "Artificial Intelligence".
3. If you want to change the seed concept, you can do so by modifying the `seedConcept` variable in the `main` function in `cmd/kg-builder/main.go`. Option to provide the seed concept via a command line argument is in development.






