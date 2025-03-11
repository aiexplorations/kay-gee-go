/**
 * API Client for interacting with the Knowledge Graph services
 */
class ApiClient {
    constructor() {
        this.baseUrl = '/api';
    }

    /**
     * Get the current graph data from Neo4j
     * @returns {Promise<Object>} The graph data with nodes and relationships
     */
    async getGraphData() {
        try {
            const response = await fetch(`${this.baseUrl}/graph`);
            if (!response.ok) {
                throw new Error(`Failed to fetch graph data: ${response.statusText}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error fetching graph data:', error);
            throw error;
        }
    }

    /**
     * Start the knowledge graph builder
     * @param {Object} params - Parameters for the builder
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
     * Search for concepts in the graph
     * @param {string} query - The search query
     * @returns {Promise<Array>} The matching concepts
     */
    async searchConcepts(query) {
        try {
            const response = await fetch(`${this.baseUrl}/concepts/search?q=${encodeURIComponent(query)}`);
            if (!response.ok) {
                throw new Error(`Failed to search concepts: ${response.statusText}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error searching concepts:', error);
            throw error;
        }
    }

    /**
     * Create a relationship between two concepts
     * @param {string} source - The source concept ID
     * @param {string} target - The target concept ID
     * @param {string} type - The relationship type
     * @returns {Promise<Object>} The created relationship
     */
    async createRelationship(source, target, type) {
        try {
            const response = await fetch(`${this.baseUrl}/relationships`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    source,
                    target,
                    type,
                }),
            });
            if (!response.ok) {
                throw new Error(`Failed to create relationship: ${response.statusText}`);
            }
            return await response.json();
        } catch (error) {
            console.error('Error creating relationship:', error);
            throw error;
        }
    }

    /**
     * Get statistics about the graph
     * @returns {Promise<Object>} The graph statistics
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
} 