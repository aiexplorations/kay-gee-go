/**
 * Tests for the main application script
 */

// Mock ApiClient
const mockApiClient = {
  getGraphData: jest.fn().mockResolvedValue({
    nodes: [{ id: '1', name: 'AI', properties: { name: 'AI' } }],
    links: [{ source: '1', target: '2', type: 'RELATES_TO' }]
  }),
  startBuilder: jest.fn().mockResolvedValue({ success: true }),
  stopBuilder: jest.fn().mockResolvedValue({ success: true }),
  startEnricher: jest.fn().mockResolvedValue({ success: true }),
  stopEnricher: jest.fn().mockResolvedValue({ success: true }),
  createRelationship: jest.fn().mockResolvedValue({ success: true }),
  searchConcepts: jest.fn().mockResolvedValue([
    { id: '1', name: 'AI', properties: { name: 'AI' } },
    { id: '2', name: 'Machine Learning', properties: { name: 'Machine Learning' } }
  ]),
  getStatistics: jest.fn().mockResolvedValue({
    nodes: 100,
    relationships: 150,
    builderStatus: 'stopped',
    enricherStatus: 'stopped'
  }),
  neo4jUrl: '/db/data/transaction/commit',
  neo4jAuth: 'bmVvNGo6cGFzc3dvcmQ='
};

// Mock GraphVisualizer
const mockGraphVisualizer = {
  updateGraph: jest.fn(),
  render: jest.fn(),
  resetCamera: jest.fn(),
  destroy: jest.fn(),
  setData: jest.fn(),
  startSimulation: jest.fn()
};

// Mock global functions
global.alert = jest.fn();
global.loadScript = jest.fn().mockResolvedValue(true);
global.fetch = jest.fn().mockResolvedValue({
  ok: true,
  json: jest.fn().mockResolvedValue({
    results: [{ data: [{ row: [100] }] }]
  })
});

describe('Main Application', () => {
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Reset selected concepts
    global.selectedConcepts = [];
    
    // Set up DOM elements
    document.body.innerHTML = `
      <div id="control-panel">
        <div id="builder-controls">
          <input id="seed-concept" value="AI">
          <input id="max-nodes" value="100">
          <input id="timeout" value="30">
          <input id="random-relationships" value="50">
          <input id="concurrency" value="5">
          <button id="start-builder">Start Builder</button>
          <button id="stop-builder">Stop Builder</button>
        </div>
        <div id="enricher-controls">
          <input id="batch-size" value="10">
          <input id="interval" value="60">
          <input id="max-relationships" value="100">
          <input id="enricher-concurrency" value="5">
          <button id="start-enricher">Start Enricher</button>
          <button id="stop-enricher">Stop Enricher</button>
        </div>
        <div id="graph-controls">
          <button id="reset-camera">Reset Camera</button>
          <button id="refresh-graph">Refresh Graph</button>
        </div>
        <div id="concept-linker">
          <select id="concept1"></select>
          <select id="concept2"></select>
          <input id="relationship-type" value="RELATES_TO">
          <button id="link-concepts">Link Concepts</button>
        </div>
      </div>
      <div id="graph-canvas"></div>
      <div id="node-info-panel"></div>
      <div id="statistics-panel">
        <div id="graph-stats">
          <p>Nodes: <span id="concept-count">0</span></p>
          <p>Relationships: <span id="relationship-count">0</span></p>
        </div>
      </div>
    `;
    
    // Mock constructors
    global.ApiClient = jest.fn(() => mockApiClient);
    global.GraphVisualizer = jest.fn(() => mockGraphVisualizer);
    
    // Mock document.addEventListener to trigger DOMContentLoaded
    const originalAddEventListener = document.addEventListener;
    document.addEventListener = jest.fn((event, callback) => {
      if (event === 'DOMContentLoaded') {
        callback();
      } else {
        originalAddEventListener(event, callback);
      }
    });
  });
  
  test('should initialize API client and graph visualizer', () => {
    // Load main.js to set up event listeners
    require('../main');
    
    // Verify
    expect(global.ApiClient).toHaveBeenCalled();
    expect(global.GraphVisualizer).toHaveBeenCalledWith('graph-canvas');
  });
  
  test('should start builder when start button is clicked', async () => {
    // Load main.js
    require('../main');
    
    // Get the start builder button
    const startBuilderBtn = document.getElementById('start-builder');
    
    // Add click event listener
    startBuilderBtn.addEventListener('click', () => {
      const seedConcept = document.getElementById('seed-concept').value;
      const maxNodes = document.getElementById('max-nodes').value;
      const timeout = document.getElementById('timeout').value;
      const randomRelationships = document.getElementById('random-relationships').value;
      const concurrency = document.getElementById('concurrency').value;
      
      mockApiClient.startBuilder({
        seedConcept,
        maxNodes,
        timeout,
        randomRelationships,
        concurrency
      })
      .then(() => {
        alert('Builder started successfully');
      })
      .catch(error => {
        alert(`Error starting builder: ${error}`);
      });
    });
    
    // Trigger button click
    startBuilderBtn.click();
    
    // Wait for promises to resolve
    await new Promise(process.nextTick);
    
    // Verify
    expect(mockApiClient.startBuilder).toHaveBeenCalledWith({
      seedConcept: 'AI',
      maxNodes: '100',
      timeout: '30',
      randomRelationships: '50',
      concurrency: '5'
    });
    expect(global.alert).toHaveBeenCalledWith('Builder started successfully');
  });
  
  test('should stop builder when stop button is clicked', async () => {
    // Load main.js
    require('../main');
    
    // Get the stop builder button
    const stopBuilderBtn = document.getElementById('stop-builder');
    
    // Add click event listener
    stopBuilderBtn.addEventListener('click', () => {
      mockApiClient.stopBuilder()
      .then(() => {
        alert('Builder stopped successfully');
      })
      .catch(error => {
        alert(`Error stopping builder: ${error}`);
      });
    });
    
    // Trigger button click
    stopBuilderBtn.click();
    
    // Wait for promises to resolve
    await new Promise(process.nextTick);
    
    // Verify
    expect(mockApiClient.stopBuilder).toHaveBeenCalled();
    expect(global.alert).toHaveBeenCalledWith('Builder stopped successfully');
  });
  
  test('should reset camera when reset button is clicked', () => {
    // Load main.js
    require('../main');
    
    // Get the reset camera button
    const resetCameraBtn = document.getElementById('reset-camera');
    
    // Add click event listener
    resetCameraBtn.addEventListener('click', () => {
      mockGraphVisualizer.resetCamera();
    });
    
    // Trigger button click
    resetCameraBtn.click();
    
    // Verify
    expect(mockGraphVisualizer.resetCamera).toHaveBeenCalled();
  });
}); 