# Changelog

All notable changes to this project will be documented in this file.

## [0.5.0] - 2025-03-15

### Added
- Created a clean, modular directory structure with clear separation of concerns
- Added a unified Makefile with targets for building, testing, and running the application
- Created a centralized configuration system with dedicated config directory
- Added symbolic links for convenience (run.sh, stop.sh, status.sh)

### Changed
- Completely reorganized the project structure for better maintainability
  - Moved command-line applications to `cmd/` directory
  - Moved internal packages to `internal/` directory
  - Created `build/` directory for Dockerfiles
  - Created `config/` directory for configuration files
  - Created `cache/` directory for caching LLM responses
  - Created `scripts/` directory for scripts to run the application
- Consolidated multiple scripts into a single set of scripts in the `scripts/` directory
- Created a unified docker-compose.yml file for all components
- Updated the README.md with the new project structure and simplified instructions
- Improved the .gitignore file with more comprehensive patterns

### Removed
- Eliminated redundant scripts and consolidated functionality
- Removed duplicate configuration files
- Removed nested project directories for cleaner organization

## [0.4.1] - 2025-03-12

### Fixed
- Fixed routing conflict in the frontend Go API server that caused the application to crash
- Removed unused import (net/http) in main.go that was causing build failures
- Fixed static file serving configuration to avoid conflicts with API routes
- Updated check-system.sh to correctly identify Docker networks with prefixed names
- Enhanced error detection and reporting in the frontend container
- Fixed Go dependency issues by improving the fix-go-deps.sh script
- Added better error handling in start-all.sh with --debug flag to show detailed logs
- Fixed container startup sequence to ensure proper initialization

### Changed
- Modified static file serving path from root ("/") to "/static" to avoid routing conflicts
- Updated NoRoute handler to properly serve index.html for client-side routing
- Improved system check to provide more accurate status information
- Enhanced Docker network detection in check-system.sh

## [0.4.0] - 2025-03-12

### Added
- Implemented a Three.js frontend with Go backend API server for visualizing the knowledge graph
- Added 3D visualization of the knowledge graph with interactive nodes and relationships
- Created a Go-based API server to serve the frontend and communicate with Neo4j
- Added RESTful API endpoints for retrieving graph data, statistics, and managing concepts
- Implemented unit tests for all API handlers (concepts, graph, statistics, utils, builder, enricher)
- Created a Makefile for building, testing, and running the frontend application
- Added Docker support for the frontend with multi-stage build process
- Added new system management scripts:
  - `check-system.sh`: Checks the status of all components and reports any issues
  - `fix-go-deps.sh`: Fixes Go dependency issues for the frontend
- Enhanced existing scripts to include frontend component:
  - Updated `start-all.sh` to include frontend startup with debug mode
  - Updated `stop-all.sh` to properly stop the frontend
  - Updated `test-all.sh` to run frontend tests
  - Updated `update-all.sh` to update frontend dependencies

### Changed
- Refactored system architecture to include the frontend as a separate container
- Updated docker-compose.yml to include the frontend service
- Enhanced error reporting in all scripts with colored output
- Improved system startup process with better error handling and debugging
- Updated README.md with frontend documentation and troubleshooting information

### Fixed
- Fixed Go dependency issues in the frontend by adding proper dependency verification
- Fixed Docker build process for the frontend to ensure proper dependency resolution
- Fixed issues with go.mod and go.sum synchronization
- Enhanced error detection and reporting in start-all.sh script
- Added dependency checks to prevent build failures
- Improved test coverage to catch dependency issues before deployment
- Fixed container networking to ensure proper communication between services

## [0.3.0] - 2025-03-11

### Added
- Created a Neo4j service interface (`Neo4jService`) to abstract Neo4j operations
- Implemented a real Neo4j service (`RealNeo4jService`) that uses the Neo4j driver
- Created a mock Neo4j service (`MockNeo4jService`) for testing without a real Neo4j connection
- Added dependency injection for Neo4j service in the Enricher
- Added proper timeout handling in tests to prevent hangs

### Changed
- Refactored the Enricher to use the Neo4j service interface instead of direct Neo4j calls
- Updated all tests to use the mock Neo4j service
- Improved test reliability by properly mocking Neo4j operations
- Enhanced the test framework to allow tests to run without a real Neo4j connection
- Increased test coverage for the enricher package to 67.5%

### Fixed
- Fixed interface conversion panic in tests by properly implementing mock interfaces
- Fixed test timeouts in the TestStartAndStop function
- Fixed issues with mock implementations not matching Neo4j driver interfaces
- Fixed test failures due to incomplete mock implementations
- Fixed issues with the RunOnce method in tests

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