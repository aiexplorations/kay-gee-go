# Changelog

All notable changes to this project will be documented in this file.

## [0.6.5] - 2025-03-13

### Added
- Added CSS variable integration to the graph visualizer for consistent styling
- Implemented a new feature to highlight selected nodes and their connected edges in red
- Added new CSS variable `--edge-selected-color` for styling selected edges
- Added helper methods for edge color management: `resetEdgeColors()` and `highlightConnectedEdges()`

### Changed
- Modified the graph visualizer to use CSS variables instead of hardcoded styles
- Updated node and edge styling to use values from CSS variables
- Changed selected node color to red for better visibility
- Improved node highlighting behavior to maintain consistent styling during hover and selection
- Enhanced the node info panel to properly highlight connected nodes when clicked

### Fixed
- Fixed inconsistent styling issues in the graph visualizer
- Resolved issues with node highlighting during hover and selection
- Fixed edge coloring to properly highlight connections to selected nodes
- Improved color consistency across the visualization

## [0.6.3] - 2025-03-12

### Fixed
- Fixed node highlighting in the graph visualization to properly display selected nodes in bright red
- Resolved issue with nodes being scaled down to small points when clicked
- Enhanced the `highlightNode` method in GraphVisualizer to preserve original node sizes
- Improved handling of node and link highlighting with proper size scaling
- Added support for both string IDs and object references in link source/target handling
- Enhanced debugging with detailed console logs for node sizing and highlighting
- Fixed edge coloring to ensure highlighted edges are clearly visible

### Changed
- Updated node highlighting color from pink/red (0xff4080) to bright red (0xff0000) for better visibility
- Improved node scaling logic to scale based on original node size rather than fixed values
- Enhanced connected node highlighting with better size preservation
- Added more robust error handling for node and link interactions

## [0.6.1] - 2025-03-12

### Added
- Integrated automatic repository optimization into the main script
- Added new `optimize` command to the unified script for manual optimization
- Added `--skip-optimization` flag to bypass automatic optimization
- Created smart caching strategy that preserves recent and important cache files
- Added options for aggressive optimization with `--aggressive` flag
- Added customizable example preservation with `--keep-examples=N` parameter

### Changed
- Consolidated all management scripts into a single unified `kg.sh` script
- Improved script organization with command-based interface
- Enhanced help documentation with detailed examples
- Optimized repository size by removing redundant files and binaries
- Updated .gitignore to properly exclude cache files and binaries

### Fixed
- Fixed issue with redundant script files cluttering the repository
- Resolved Git tracking of cache files that should be excluded
- Fixed repository bloat by implementing automatic optimization
- Addressed potential performance issues from excessive cache files

## [0.6.0] - 2025-03-12

### Added
- Implemented low connectivity concept seeding to enhance graph diversity
- Added `--use-low-connectivity` flag to enable building the graph using low connectivity concepts as seeds
- Created new Neo4j functions to support low connectivity features:
  - `GetLowConnectivityConcepts`: Retrieves concepts with the least number of connections
  - `GetRandomLowConnectivityConcept`: Retrieves a random concept from low connectivity concepts
- Added new graph building method `BuildGraphWithLowConnectivitySeeds` to utilize low connectivity concepts
- Implemented comprehensive testing for the new low connectivity features

### Changed
- Renamed the entire deployment from "kg-builder" to "kaygeego" to avoid term overloading
- Restructured Docker Compose services with clearer naming:
  - Renamed services from "kg-builder", "kg-enricher", and "kg-frontend" to "builder", "enricher", and "frontend"
  - Updated container names to use the "kaygeego-" prefix (kaygeego-neo4j, kaygeego-builder, etc.)
  - Changed network name from "kg-network" to "kaygeego-network"
- Updated all scripts and configuration files to use the new naming convention
- Improved graph building algorithm to create more balanced and comprehensive knowledge graphs
- Enhanced the diversity of the knowledge graph by focusing on expanding low-connectivity areas

### Fixed
- Fixed issues with isolated nodes by targeting low connectivity concepts for expansion
- Improved graph balance by prioritizing concepts with fewer connections
- Enhanced overall graph connectivity by creating more pathways between distant concepts
- Fixed potential bottlenecks in graph expansion by diversifying seed concepts

## [0.5.3] - 2025-03-11

### Fixed
- Fixed LLM client to properly handle Ollama API responses by extracting JSON from markdown-formatted responses
- Added support for multiline JSON parsing in LLM responses with improved regex patterns
- Updated LLM prompts to explicitly request JSON in the expected format for better response consistency
- Fixed frontend graph visualization to directly query Neo4j database through the proxy
- Implemented proper parsing of Neo4j responses in the frontend API client
- Updated statistics loading to use direct Neo4j queries for accurate data display
- Enhanced error handling and logging in the LLM client for better debugging
- Fixed the graph visualization to properly display nodes and relationships from the Neo4j database
- Improved node positioning in the 3D visualization with a spherical formation for better visual appeal

### Changed
- Modified the LLM client to use non-streaming responses for more reliable JSON parsing
- Updated the frontend API client to bypass the Go backend API and directly access Neo4j
- Enhanced the graph visualization with better node sizing based on connection counts
- Improved the frontend statistics display with more accurate data from Neo4j

## [0.5.2] - 2025-03-11

### Fixed
- Fixed frontend implementation to correctly display the space-like knowledge graph visualization
- Resolved issue with static file serving by replacing the old frontend with a dedicated Nginx-based solution
- Corrected CSS and JavaScript file paths in the frontend HTML to ensure proper loading of resources
- Removed conflicting old frontend implementation to prevent interference with the new visualization
- Ensured proper serving of Three.js-based graph visualization with correct file references

### Changed
- Migrated frontend to use Nginx for static file serving instead of the Go-based server
- Updated docker-compose.yml to use the new frontend implementation from kg-frontend directory
- Simplified frontend container by removing unnecessary Go dependencies

## [0.5.1] - 2025-03-11

### Fixed
- Fixed duplicate node handling in Neo4j database by implementing a robust detection and removal system
- Added proper error handling for duplicate 'Complexity Theory' nodes that were preventing constraint creation
- Enhanced the `InitializeSchema` method to gracefully handle existing constraints and duplicate nodes
- Fixed null value handling in `GetRandomConcepts` and `GetRandomConceptPairs` methods to prevent frontend crashes
- Improved frontend JavaScript code to better handle errors and display graph data more effectively
- Added proper error handling in the frontend to display meaningful error messages
- Enhanced D3.js graph visualization with zoom functionality, node coloring, and tooltips
- Fixed frontend template creation to ensure HTML, CSS, and JS files are properly generated
- Added proper directory creation in the frontend container to ensure all required directories exist
- Implemented a more robust graph data processing pipeline in the frontend JavaScript

### Changed
- Refactored Neo4j client to handle null values safely with proper type checking
- Updated the graph visualization to use a force-directed layout with links between nodes
- Enhanced the frontend UI with better error handling and user feedback
- Improved the builder service to handle database constraints more gracefully

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

## [0.4.1] - 2025-03-11

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

## [0.4.0] - 2025-03-11

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
