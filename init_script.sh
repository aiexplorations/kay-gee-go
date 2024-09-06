#!/bin/bash

# Initialize Knowledge Graph Builder Project

# Create project root directory
mkdir -p kg-builder
cd kg-builder

# Create project structure
mkdir -p cmd/kg-builder internal/neo4j internal/llm internal/graph

# Initialize Go module
go mod init kg-builder

# Create main.go
cat > cmd/kg-builder/main.go << EOL
package main

import (
	"log"
	"os"
	"time"
	"kg-builder/internal/graph"
	"kg-builder/internal/neo4j"
	"kg-builder/internal/llm"
)

func main() {
	log.Println("Starting Knowledge Graph Builder")

	neo4jDriver, err := neo4j.SetupNeo4jConnection()
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer neo4jDriver.Close()

	graphBuilder := graph.NewGraphBuilder(neo4jDriver, llm.GetRelatedConcepts)

	seedConcept := "Artificial Intelligence"
	maxNodes := 100
	timeout := 20 * time.Minute

	err = graphBuilder.BuildGraph(seedConcept, maxNodes, timeout)
	if err != nil {
		log.Fatalf("Failed to build graph: %v", err)
	}

	log.Println("Knowledge Graph Builder completed")
}
EOL

# Create neo4j.go
cat > internal/neo4j/neo4j.go << EOL
package neo4j

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func SetupNeo4jConnection() (neo4j.Driver, error) {
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "password"

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return driver, nil
}

func CreateRelationship(driver neo4j.Driver, from, to, relation string) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := \`
			MERGE (a:Concept {name: $from})
			MERGE (b:Concept {name: $to})
			MERGE (a)-[:$relation]->(b)
		\`
		params := map[string]interface{}{
			"from":     from,
			"to":       to,
			"relation": relation,
		}
		_, err := tx.Run(query, params)
		return nil, err
	})

	return err
}
EOL

# Create llm.go
cat > internal/llm/llm.go << EOL
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Concept struct {
	Name      string \`json:"name"\`
	Relation  string \`json:"relation"\`
	RelatedTo string \`json:"relatedTo"\`
}

func GetRelatedConcepts(concept string) ([]Concept, error) {
	url := "http://localhost:11434/api/generate"
	prompt := fmt.Sprintf("Given the concept '%s', provide 5 related concepts. For each, specify the relationship type. Format as JSON array with 'name', 'relation', and 'relatedTo' keys.", concept)

	requestBody, err := json.Marshal(map[string]string{
		"model":  "llama3.1",
		"prompt": prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Response string \`json:"response"\`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var concepts []Concept
	err = json.Unmarshal([]byte(result.Response), &concepts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal concepts: %w", err)
	}

	return concepts, nil
}
EOL

# Create graph.go
cat > internal/graph/graph.go << EOL
package graph

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	kgneo4j "kg-builder/internal/neo4j"
)

type Concept struct {
	Name      string
	Relation  string
	RelatedTo string
}

type GraphBuilder struct {
	driver              neo4j.Driver
	getRelatedConcepts  func(string) ([]Concept, error)
	processedConcepts   map[string]bool
	mutex               sync.Mutex
}

func NewGraphBuilder(driver neo4j.Driver, getRelatedConcepts func(string) ([]Concept, error)) *GraphBuilder {
	return &GraphBuilder{
		driver:             driver,
		getRelatedConcepts: getRelatedConcepts,
		processedConcepts:  make(map[string]bool),
	}
}

func (gb *GraphBuilder) BuildGraph(seedConcept string, maxNodes int, timeout time.Duration) error {
	queue := []string{seedConcept}
	startTime := time.Now()

	for len(gb.processedConcepts) < maxNodes && len(queue) > 0 {
		concept := queue[0]
		queue = queue[1:]

		gb.mutex.Lock()
		if gb.processedConcepts[concept] {
			gb.mutex.Unlock()
			continue
		}
		gb.processedConcepts[concept] = true
		gb.mutex.Unlock()

		relatedConcepts, err := gb.getRelatedConcepts(concept)
		if err != nil {
			log.Printf("Error getting related concepts for %s: %v", concept, err)
			continue
		}

		for _, rc := range relatedConcepts {
			err := kgneo4j.CreateRelationship(gb.driver, concept, rc.Name, rc.Relation)
			if err != nil {
				log.Printf("Error creating relationship: %v", err)
				continue
			}

			gb.mutex.Lock()
			if !gb.processedConcepts[rc.Name] {
				queue = append(queue, rc.Name)
			}
			gb.mutex.Unlock()
		}

		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout reached after processing %d concepts", len(gb.processedConcepts))
		}
	}

	return nil
}
EOL

# Create Dockerfile
cat > Dockerfile << EOL
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./

RUN go mod download

COPY . .

RUN go build -o /kg-builder ./cmd/kg-builder

CMD ["/kg-builder"]
EOL

# Create docker-compose.yml
cat > docker-compose.yml << EOL
version: '3'
services:
  kg-builder:
    build: .
    depends_on:
      - neo4j
    environment:
      - NEO4J_URI=bolt://neo4j:7687
      - NEO4J_USER=neo4j
      - NEO4J_PASSWORD=password
      - LLM_URL=http://host.docker.internal:11434/api/generate
    network_mode: "host"

  neo4j:
    image: neo4j:4.4
    ports:
      - "7474:7474"
      - "7687:7687"
    environment:
      - NEO4J_AUTH=neo4j/password
EOL

# Create README.md
cat > README.md << EOL
# Knowledge Graph Builder

This project builds a knowledge graph using Neo4j and an LLM service.

## Setup

1. Ensure Docker and Docker Compose are installed on your system.
2. Run \`docker-compose up --build\` to start the application and Neo4j.

## Usage

The application will automatically start building the knowledge graph from the seed concept "Artificial Intelligence".

## Project Structure

- \`cmd/kg-builder/\`: Main application entry point
- \`internal/neo4j/\`: Neo4j connection and operations
- \`internal/llm/\`: LLM service interactions
- \`internal/graph/\`: Graph operations and data structures
EOL

# Initialize Git repository
git init
echo -e "\n# Ignore Go build artifacts\n*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n\n# Ignore Go test artifacts\n*.test\n\n# Ignore Go coverage artifacts\n*.out\n\n# Ignore Go workspace file\ngo.work" > .gitignore

# Download Go dependencies
go mod tidy

echo "Project initialized successfully!"
echo "Navigate to the project directory: cd kg-builder"
echo "Build and run the project: docker-compose up --build"