/**
 * Tests for the ApiClient component
 */

// Mock fetch
global.fetch = jest.fn();

// Import the ApiClient class
const ApiClient = require('../api-client');

describe('ApiClient', () => {
  let apiClient;
  
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Create a new instance
    apiClient = new ApiClient();
    
    // Mock successful fetch response
    global.fetch.mockImplementation(() => 
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ data: 'test' })
      })
    );
  });
  
  test('should initialize with correct base URL', () => {
    expect(apiClient.baseUrl).toBe('/api');
  });
  
  test('should fetch graph data', async () => {
    // Mock the Neo4j responses
    global.fetch.mockImplementationOnce(() => 
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ 
          results: [{ 
            data: [{ row: [1, 'Concept 1'] }] 
          }] 
        })
      })
    );
    
    global.fetch.mockImplementationOnce(() => 
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ 
          results: [{ 
            data: [{ row: [1, 2, 'RELATES_TO'] }] 
          }] 
        })
      })
    );
    
    // Execute
    const result = await apiClient.getGraphData();
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith(apiClient.neo4jUrl, expect.any(Object));
    expect(result).toHaveProperty('nodes');
    expect(result).toHaveProperty('links');
  });
  
  test('should handle error when fetching graph data', async () => {
    // Setup
    global.fetch.mockImplementationOnce(() => 
      Promise.resolve({
        ok: false,
        statusText: 'Not Found'
      })
    );
    
    // Execute and verify
    await expect(apiClient.getGraphData()).rejects.toThrow('Failed to fetch nodes: Not Found');
  });
  
  test('should start builder with parameters', async () => {
    // Setup
    const params = {
      seedConcept: 'AI',
      maxNodes: 100,
      timeout: 30,
      randomRelationships: 50,
      concurrency: 5
    };
    
    // Execute
    const result = await apiClient.startBuilder(params);
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith('/api/builder/start', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(params),
    });
    expect(result).toEqual({ data: 'test' });
  });
  
  test('should stop builder', async () => {
    // Execute
    const result = await apiClient.stopBuilder();
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith('/api/builder/stop', {
      method: 'POST',
    });
    expect(result).toEqual({ data: 'test' });
  });
  
  test('should start enricher with parameters', async () => {
    // Setup
    const params = {
      batchSize: 10,
      interval: 60,
      maxRelationships: 100,
      concurrency: 5
    };
    
    // Execute
    const result = await apiClient.startEnricher(params);
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith('/api/enricher/start', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(params),
    });
    expect(result).toEqual({ data: 'test' });
  });
  
  test('should stop enricher', async () => {
    // Execute
    const result = await apiClient.stopEnricher();
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith('/api/enricher/stop', {
      method: 'POST',
    });
    expect(result).toEqual({ data: 'test' });
  });
  
  test('should search concepts', async () => {
    // Setup
    global.fetch.mockImplementationOnce(() => 
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ 
          results: [{ 
            data: [{ row: [{ id: 1, name: 'AI' }] }] 
          }] 
        })
      })
    );
    
    // Execute
    const result = await apiClient.searchConcepts('AI');
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith(apiClient.neo4jUrl, expect.any(Object));
    expect(result).toEqual([{ id: 1, name: 'AI', properties: { id: 1, name: 'AI' } }]);
  });
  
  test('should create relationship', async () => {
    // Setup
    const source = '1';
    const target = '2';
    const type = 'RELATES_TO';
    
    // Mock the response for createRelationship
    global.fetch.mockImplementationOnce(() => 
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ 
          success: true,
          message: 'Relationship created successfully'
        })
      })
    );
    
    // Execute
    const result = await apiClient.createRelationship(source, target, type);
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith(apiClient.neo4jUrl, expect.any(Object));
    expect(result).toEqual({ 
      success: true,
      message: 'Relationship created successfully'
    });
  });
  
  test('should get statistics', async () => {
    // Execute
    const result = await apiClient.getStatistics();
    
    // Verify
    expect(global.fetch).toHaveBeenCalledWith('/api/statistics');
    expect(result).toEqual({ data: 'test' });
  });
}); 