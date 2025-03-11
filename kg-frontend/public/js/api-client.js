/**
 * API Client for interacting with the Knowledge Graph services
 */
class ApiClient {
    constructor() {
        this.neo4jUrl = '/db/data/transaction/commit';
        this.neo4jAuth = btoa('neo4j:password'); // Base64 encode neo4j:password
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
            alert('Builder functionality is not available in this version. Please use the command line tools.');
            return { success: false, message: 'Builder functionality is not available in this version.' };
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
            alert('Builder functionality is not available in this version. Please use the command line tools.');
            return { success: false, message: 'Builder functionality is not available in this version.' };
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
            alert('Enricher functionality is not available in this version. Please use the command line tools.');
            return { success: false, message: 'Enricher functionality is not available in this version.' };
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
            alert('Enricher functionality is not available in this version. Please use the command line tools.');
            return { success: false, message: 'Enricher functionality is not available in this version.' };
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
            const nodeCountQuery = 'MATCH (n) RETURN count(n) as nodeCount';
            const relationshipCountQuery = 'MATCH ()-[r]->() RETURN count(r) as relationshipCount';
            const nodeTypesQuery = 'MATCH (n) RETURN labels(n) as type, count(*) as count';
            const relationshipTypesQuery = 'MATCH ()-[r]->() RETURN type(r) as type, count(*) as count';
            
            const response = await fetch(this.neo4jUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Basic ${this.neo4jAuth}`,
                    'Accept': 'application/json; charset=UTF-8'
                },
                body: JSON.stringify({
                    statements: [
                        { statement: nodeCountQuery },
                        { statement: relationshipCountQuery },
                        { statement: nodeTypesQuery },
                        { statement: relationshipTypesQuery }
                    ]
                })
            });
            
            if (!response.ok) {
                throw new Error(`Failed to fetch statistics: ${response.statusText}`);
            }
            
            const data = await response.json();
            
            const statistics = {
                nodeCount: 0,
                relationshipCount: 0,
                nodeTypes: [],
                relationshipTypes: []
            };
            
            if (data.results && data.results.length >= 4) {
                // Node count
                if (data.results[0].data && data.results[0].data.length > 0) {
                    statistics.nodeCount = data.results[0].data[0].row[0];
                }
                
                // Relationship count
                if (data.results[1].data && data.results[1].data.length > 0) {
                    statistics.relationshipCount = data.results[1].data[0].row[0];
                }
                
                // Node types
                if (data.results[2].data) {
                    data.results[2].data.forEach(row => {
                        if (row.row && row.row.length >= 2) {
                            statistics.nodeTypes.push({
                                type: row.row[0].join(':'),
                                count: row.row[1]
                            });
                        }
                    });
                }
                
                // Relationship types
                if (data.results[3].data) {
                    data.results[3].data.forEach(row => {
                        if (row.row && row.row.length >= 2) {
                            statistics.relationshipTypes.push({
                                type: row.row[0],
                                count: row.row[1]
                            });
                        }
                    });
                }
            }
            
            return statistics;
        } catch (error) {
            console.error('Error fetching statistics:', error);
            throw error;
        }
    }
} 