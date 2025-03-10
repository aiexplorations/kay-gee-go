# Cache Directory

This directory is used to store cached responses from the LLM service. The cache helps to:

1. Reduce the number of API calls to the LLM service
2. Speed up graph building by reusing previously retrieved concepts and relationships
3. Provide offline access to previously mined knowledge

## Structure

The cache directory contains two types of files:

- `concept_*.json`: Cached related concepts for a given concept
- `rel_*.json`: Cached relationships between two concepts

## Persistence

The cache directory is mounted as a volume in the Docker container, so the cached data persists even when the container is stopped or removed.

## Clearing the Cache

To clear the cache, simply delete the files in this directory:

```bash
rm -rf kg-builder/cache/*.json
```

Note: Do not delete this README.md file. 