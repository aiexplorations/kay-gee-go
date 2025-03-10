#!/bin/bash

# Check if a model name was provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <model_name>"
    echo "Example: $0 qwen2.5:3b"
    exit 1
fi

MODEL_NAME="$1"
CONFIG_FILE="kg-builder/config.yaml"

# Check if the config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Config file not found: $CONFIG_FILE"
    echo "Creating a new config file..."
    
    # Create the config file with the new model
    cat > "$CONFIG_FILE" << EOL
# Knowledge Graph Builder Configuration

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
  model: "$MODEL_NAME"
  cache_dir: "./cache/llm"

# Graph Configuration
graph:
  seed_concept: "Artificial Intelligence"
  max_nodes: 100
  timeout_minutes: 30
  worker_count: 10
  random_relationships: 50
  concurrency: 5
EOL
    echo "Created new config file with model: $MODEL_NAME"
else
    # Update the existing config file
    if grep -q "model:" "$CONFIG_FILE"; then
        # Replace the model line
        sed -i.bak "s/model: \"[^\"]*\"/model: \"$MODEL_NAME\"/" "$CONFIG_FILE"
        rm -f "${CONFIG_FILE}.bak"
        echo "Updated config file with model: $MODEL_NAME"
    else
        echo "Could not find model configuration in $CONFIG_FILE"
        exit 1
    fi
fi

echo "To apply the changes, restart the application with './stop.sh' and then './start.sh'" 