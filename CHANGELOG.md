# Changelog

All notable changes to this project will be documented in this file.

## [0.2.0] - 2025-03-11

### Added
- Comprehensive test coverage for all components
- Unit tests for configuration management
- Unit tests for graph building with concurrency, max nodes, and timeout
- Unit tests for LLM service with caching and filename sanitization
- Unit tests for Neo4j operations
- Integration tests for the complete workflow
- End-to-end tests for the entire system
- Added `--skip-tests` flag to start-all.sh script to optionally skip tests
- Exported SanitizeFilename function for testing purposes

### Changed
- Improved configuration system to properly handle environment variables overriding file values
- Enhanced test scripts to run unit, integration, and end-to-end tests
- Updated Docker configuration to support test execution
- Standardized test naming and structure across all components
- Improved error handling in tests

### Fixed
- Fixed configuration loading to properly respect environment variables
- Fixed LLM cache directory path to ensure consistency
- Fixed Neo4j session handling in tests
- Fixed type assertions in database operations
- Fixed concurrency issues in graph building tests
- Fixed build errors in test files
- Fixed Docker container test execution

## [0.1.0] - 2025-03-10

### Added
- Configuration management system using environment variables and command-line flags
- Comprehensive test suite with unit and integration tests
- Error handling system with custom error types and retry mechanisms
- LLM response caching for offline access and improved performance
- Neo4j data persistence using Docker volumes
- Command-line interface with various configuration options
- Convenience scripts for starting, stopping, and checking status
- `--stats-only` flag for showing graph statistics without building
- Support for passing configuration arguments to the start.sh script
- YAML configuration file support for easier configuration management
- Changed default LLM model from llama3.1:latest to qwen2.5:3b
- Added Knowledge Graph Enricher microservice for finding and adding relationships between random pairs of concepts
- Added scripts for starting and stopping the Knowledge Graph Enricher
- Added integrated deployment with start-all.sh and stop-all.sh scripts to run both components together

### Changed
- Refactored Neo4j connection logic to use configuration system
- Refactored LLM service to use configuration system
- Refactored graph building to use configuration system
- Improved logging with more detailed information
- Updated Docker Compose configuration for better persistence
- Standardized default values for all configuration parameters
- Updated README.md with information about the Knowledge Graph Enricher
- Created a combined docker-compose.yml file for integrated deployment

### Fixed
- Fixed error handling in Neo4j connection
- Fixed error handling in LLM service
- Fixed concurrency issues in graph building
- Fixed compilation errors in the Knowledge Graph Enricher

### Removed
- Removed hardcoded configuration values throughout the codebase 