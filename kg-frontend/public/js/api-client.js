/**
 * API Client for interacting with the Knowledge Graph services
 */
class ApiClient {
    constructor() {
        this.neo4jUrl = '/db/data/transaction/commit';
        this.neo4jAuth = btoa('neo4j:password'); // Base64 encode neo4j:password
        this.baseUrl = '/api'; // Base URL for REST API endpoints
    }

    /**
     * Get the current graph data from Neo4j
     * @returns {Promise<Object>} The graph data with nodes and relationships
     */
    async getGraphData() {
        try {
            const query = `
                MATCH (n:Concept)
                RETURN id(n) AS id, n.name AS name
            `;
            
            const relationshipQuery = `
                MATCH (a:Concept)-[r:RELATED_TO]->(b:Concept)
                RETURN id(a) AS source, id(b) AS target, r.type AS type
            `;
            
            // Get nodes
            const nodeResponse = await fetch(this.neo4jUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${this.neo4jAuth}`,
                    'Accept': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify({
                    statements: [{ statement: query }]
                })
            });
            
            if (!nodeResponse.ok) {
                throw new Error(`Failed to fetch nodes: ${nodeResponse.statusText}`);
            }
            
            const nodeData = await nodeResponse.json();
            
            // Get relationships
            const linkResponse = await fetch(this.neo4jUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${this.neo4jAuth}`,
                    'Accept': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify({
                    statements: [{ statement: relationshipQuery }]
                })
            });
            
            if (!linkResponse.ok) {
                throw new Error(`Failed to fetch relationships: ${linkResponse.statusText}`);
            }
            
            const linkData = await linkResponse.json();
            
            // Transform Neo4j response to the format expected by the frontend
            const nodes = [];
            const links = [];
            
            // Process nodes
            if (nodeData.results && nodeData.results.length > 0 && nodeData.results[0].data) {
                nodeData.results[0].data.forEach(row => {
                    if (row.row && row.row.length >= 2) {
                        const id = row.row[0];
                        const name = row.row[1];
                        
                        nodes.push({
                            id: id.toString(),
                            name: name,
                            size: 5 // Default size
                        });
                    }
                });
            }
            
            // Process relationships
            if (linkData.results && linkData.results.length > 0 && linkData.results[0].data) {
                linkData.results[0].data.forEach(row => {
                    if (row.row && row.row.length >= 3) {
                        const source = row.row[0];
                        const target = row.row[1];
                        const type = row.row[2];
                        
                        links.push({
                            source: source.toString(),
                            target: target.toString(),
                            type: type
                        });
                    }
                });
            }
            
            console.log(`Processed ${nodes.length} nodes and ${links.length} links`);
            
            return {
                nodes: nodes,
                links: links
            };
        } catch (error) {
            console.error('Error fetching graph data:', error);
            throw error;
        }
    }

    /**
     * Start the knowledge graph builder
     * @param {Object} params - Parameters for the builder
     * @param {string} params.seedConcept - The seed concept to start with
     * @param {number} params.maxNodes - Maximum number of nodes to create
     * @param {number} params.timeout - Timeout in minutes
     * @param {number} params.randomRelationships - Number of random relationships to create
     * @param {number} params.concurrency - Number of concurrent operations
     * @returns {Promise<Object>} The response from the builder
     */
    async startBuilder(params) {
        try {
            const response = await fetch(`${this.baseUrl}/builder/start`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(params),
            });
            
            if (!response.ok) {
                throw new Error(`Failed to start builder: ${response.statusText}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error starting builder:', error);
            throw error;
        }
    }

    /**
     * Stop the knowledge graph builder
     * @returns {Promise<Object>} The response from the builder
     */
    async stopBuilder() {
        try {
            const response = await fetch(`${this.baseUrl}/builder/stop`, {
                method: 'POST',
            });
            
            if (!response.ok) {
                throw new Error(`Failed to stop builder: ${response.statusText}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error stopping builder:', error);
            throw error;
        }
    }

    /**
     * Start the knowledge graph enricher
     * @param {Object} params - Parameters for the enricher
     * @param {number} params.batchSize - Number of pairs to process in each batch
     * @param {number} params.interval - Interval between batches in seconds
     * @param {number} params.maxRelationships - Maximum number of relationships to create
     * @param {number} params.concurrency - Number of concurrent operations
     * @returns {Promise<Object>} The response from the enricher
     */
    async startEnricher(params) {
        try {
            const response = await fetch(`${this.baseUrl}/enricher/start`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(params),
            });
            
            if (!response.ok) {
                throw new Error(`Failed to start enricher: ${response.statusText}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error starting enricher:', error);
            throw error;
        }
    }

    /**
     * Stop the knowledge graph enricher
     * @returns {Promise<Object>} The response from the enricher
     */
    async stopEnricher() {
        try {
            const response = await fetch(`${this.baseUrl}/enricher/stop`, {
                method: 'POST',
            });
            
            if (!response.ok) {
                throw new Error(`Failed to stop enricher: ${response.statusText}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error stopping enricher:', error);
            throw error;
        }
    }

    /**
     * Search for concepts in the knowledge graph
     * @param {string} query - The search query
     * @returns {Promise<Array>} The matching concepts
     */
    async searchConcepts(query) {
        try {
            const cypher = `
                MATCH (n)
                WHERE n.name CONTAINS $query
                RETURN n
                LIMIT 10
            `;
            
            const response = await fetch(this.neo4jUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${this.neo4jAuth}`,
                    'Accept': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify({
                    statements: [{ 
                        statement: cypher,
                        parameters: { query: query }
                    }]
                })
            });
            
            if (!response.ok) {
                throw new Error(`Failed to search concepts: ${response.statusText}`);
            }
            
            const data = await response.json();
            
            const concepts = [];
            if (data.results && data.results.length > 0 && data.results[0].data) {
                data.results[0].data.forEach(row => {
                    if (row.row && row.row.length > 0) {
                        const concept = row.row[0];
                        concepts.push({
                            id: concept.id,
                            name: concept.name || concept.id,
                            properties: concept
                        });
                    }
                });
            }
            
            return concepts;
        } catch (error) {
            console.error('Error searching concepts:', error);
            throw error;
        }
    }

    /**
     * Create a relationship between two concepts
     * @param {string} source - The source concept
     * @param {string} target - The target concept
     * @param {string} type - The relationship type
     * @returns {Promise<Object>} The created relationship
     */
    async createRelationship(source, target, type) {
        try {
            const cypher = `
                MATCH (a), (b)
                WHERE a.name = $source AND b.name = $target
                CREATE (a)-[r:${type}]->(b)
                RETURN r
            `;
            
            const response = await fetch(this.neo4jUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${this.neo4jAuth}`,
                    'Accept': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify({
                    statements: [{ 
                        statement: cypher,
                        parameters: { 
                            source: source,
                            target: target
                        }
                    }]
                })
            });
            
            if (!response.ok) {
                throw new Error(`Failed to create relationship: ${response.statusText}`);
            }
            
            const data = await response.json();
            
            if (data.errors && data.errors.length > 0) {
                throw new Error(`Neo4j error: ${data.errors[0].message}`);
            }
            
            return { success: true, message: 'Relationship created successfully' };
        } catch (error) {
            console.error('Error creating relationship:', error);
            throw error;
        }
    }

    /**
     * Get statistics about the knowledge graph
     * @returns {Promise<Object>} The statistics
     */
    async getStatistics() {
        try {
            const response = await fetch(`${this.baseUrl}/statistics`);
            
            if (!response.ok) {
                throw new Error(`Failed to fetch statistics: ${response.statusText}`);
            }
            
            return await response.json();
        } catch (error) {
            console.error('Error fetching statistics:', error);
            throw error;
        }
    }

    /**
     * Clean up orphan relationships and nodes in the graph
     * @returns {Promise<Object>} The cleanup result
     */
    async cleanupGraph() {
        try {
            // First, clean up orphan relationships
            const relationshipQuery = `
                MATCH ()-[r:RELATED_TO]->()
                WHERE NOT EXISTS(r.type) OR r.type = ""
                DELETE r
                RETURN count(r) AS count
            `;
            
            const relationshipResponse = await fetch(this.neo4jUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${this.neo4jAuth}`,
                    'Accept': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify({
                    statements: [{ statement: relationshipQuery }]
                })
            });
            
            if (!relationshipResponse.ok) {
                throw new Error(`Failed to clean up orphan relationships: ${relationshipResponse.statusText}`);
            }
            
            const relationshipData = await relationshipResponse.json();
            const relationshipsRemoved = relationshipData.results[0].data[0].row[0];
            
            // Then, clean up orphan nodes
            const nodeQuery = `
                MATCH (n:Concept)
                WHERE NOT (n)-[]-() 
                DELETE n
                RETURN count(n) AS count
            `;
            
            const nodeResponse = await fetch(this.neo4jUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${this.neo4jAuth}`,
                    'Accept': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify({
                    statements: [{ statement: nodeQuery }]
                })
            });
            
            if (!nodeResponse.ok) {
                throw new Error(`Failed to clean up orphan nodes: ${nodeResponse.statusText}`);
            }
            
            const nodeData = await nodeResponse.json();
            const nodesRemoved = nodeData.results[0].data[0].row[0];
            
            return {
                relationshipsRemoved,
                nodesRemoved
            };
        } catch (error) {
            console.error('Error cleaning up graph:', error);
            throw error;
        }
    }
}

// Export for use in tests
if (typeof module !== 'undefined') {
    module.exports = ApiClient;
}