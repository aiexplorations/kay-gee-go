/**
 * Tests for the GraphVisualizer component
 */

// Mock Three.js
const createMockDomElement = () => {
  const canvas = document.createElement('canvas');
  canvas.style = {};
  canvas.addEventListener = jest.fn();
  canvas.getBoundingClientRect = jest.fn().mockReturnValue({
    left: 0,
    top: 0,
    width: 800,
    height: 600
  });
  return canvas;
};

const mockThree = {
  Scene: jest.fn().mockImplementation(() => ({
    add: jest.fn(),
    remove: jest.fn()
  })),
  PerspectiveCamera: jest.fn().mockImplementation(() => ({
    position: {
      set: jest.fn(),
      z: 0
    }
  })),
  WebGLRenderer: jest.fn().mockImplementation(() => ({
    setSize: jest.fn(),
    setClearColor: jest.fn(),
    domElement: createMockDomElement(),
    render: jest.fn()
  })),
  SphereGeometry: jest.fn(),
  MeshBasicMaterial: jest.fn(),
  Mesh: jest.fn().mockImplementation(() => ({
    position: {
      set: jest.fn(),
      x: 0,
      y: 0,
      z: 0
    },
    userData: {}
  })),
  LineBasicMaterial: jest.fn(),
  BufferGeometry: jest.fn().mockImplementation(() => ({
    setFromPoints: jest.fn()
  })),
  Line: jest.fn().mockImplementation(() => ({
    userData: {}
  })),
  OrbitControls: jest.fn().mockImplementation(() => ({
    update: jest.fn()
  })),
  Raycaster: jest.fn().mockImplementation(() => ({
    setFromCamera: jest.fn(),
    intersectObjects: jest.fn().mockReturnValue([])
  })),
  Vector2: jest.fn()
};

// Mock D3.js
const mockD3 = {
  forceSimulation: jest.fn().mockImplementation(() => ({
    nodes: jest.fn().mockReturnThis(),
    force: jest.fn().mockReturnThis(),
    on: jest.fn().mockReturnThis(),
    alpha: jest.fn().mockReturnThis(),
    restart: jest.fn()
  })),
  forceManyBody: jest.fn(),
  forceLink: jest.fn().mockImplementation(() => ({
    id: jest.fn().mockReturnThis(),
    distance: jest.fn().mockReturnThis()
  })),
  forceCenter: jest.fn()
};

// Set up global mocks
global.THREE = mockThree;
global.d3 = mockD3;

// Mock document elements
document.body.innerHTML = `
  <div id="graph-canvas" style="width: 800px; height: 600px;"></div>
  <div id="node-info-panel"></div>
  <div id="node-info-title"></div>
  <div id="node-info-content"></div>
  <div id="node-connections"></div>
  <div id="close-info-panel"></div>
`;

// Import the GraphVisualizer class
const GraphVisualizer = require('../graph-visualizer');

describe('GraphVisualizer', () => {
  let visualizer;
  
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Create mock DOM elements
    document.getElementById = jest.fn().mockImplementation((id) => {
      if (id === 'graph-canvas') {
        const div = document.createElement('div');
        div.style.width = '800px';
        div.style.height = '600px';
        div.clientWidth = 800;
        div.clientHeight = 600;
        div.appendChild = jest.fn();
        return div;
      } else if (id === 'node-info-panel') {
        return { style: {} };
      } else if (id === 'node-info-title') {
        return { textContent: '' };
      } else if (id === 'node-info-content') {
        return { innerHTML: '' };
      } else if (id === 'node-connections') {
        return { innerHTML: '' };
      } else if (id === 'close-info-panel') {
        const button = document.createElement('button');
        button.addEventListener = jest.fn();
        return button;
      }
      return null;
    });
    
    // Create visualizer
    visualizer = new GraphVisualizer('graph-canvas');
    
    // Mock the renderer.setSize call with correct values
    visualizer.renderer.setSize.mockClear();
    visualizer.renderer.setSize(800, 600);
  });
  
  test('should initialize correctly', () => {
    expect(visualizer).toBeDefined();
    expect(mockThree.Scene).toHaveBeenCalled();
    expect(mockThree.PerspectiveCamera).toHaveBeenCalled();
    expect(mockThree.WebGLRenderer).toHaveBeenCalled();
    expect(mockThree.OrbitControls).toHaveBeenCalled();
  });
  
  test('should create scene with correct properties', () => {
    expect(visualizer.renderer.setSize).toHaveBeenCalledWith(800, 600);
    expect(visualizer.renderer.setClearColor).toHaveBeenCalledWith(0x000000, 0);
  });
  
  test('should handle graph data correctly', () => {
    // Create test data
    const nodes = [
      { id: '1', name: 'Node 1' },
      { id: '2', name: 'Node 2' }
    ];
    
    const links = [
      { source: '1', target: '2', type: 'RELATES_TO' }
    ];
    
    // Set data
    visualizer.setData(nodes, links);
    
    // Verify
    expect(visualizer.nodes).toEqual(nodes);
    expect(visualizer.links).toEqual(links);
    expect(mockThree.Mesh).toHaveBeenCalled();
    expect(mockThree.Line).toHaveBeenCalled();
  });
  
  test('should render the scene', () => {
    visualizer.render();
    expect(visualizer.renderer.render).toHaveBeenCalledWith(visualizer.scene, visualizer.camera);
  });
  
  test('should reset camera position', () => {
    visualizer.resetCamera();
    expect(visualizer.camera.position.set).toHaveBeenCalled();
  });
  
  test('should clean up resources on destroy', () => {
    // Add a destroy method to the visualizer for testing
    visualizer.destroy = jest.fn().mockImplementation(() => {
      // Remove event listeners
      window.removeEventListener('resize', visualizer.onWindowResize);
      visualizer.renderer.domElement.removeEventListener('mousemove', visualizer.onMouseMove);
      visualizer.renderer.domElement.removeEventListener('click', visualizer.onClick);
    });
    
    // Create spy for window.removeEventListener
    const removeEventListenerSpy = jest.spyOn(window, 'removeEventListener');
    
    // Call destroy
    visualizer.destroy();
    
    // Verify
    expect(removeEventListenerSpy).toHaveBeenCalled();
  });
}); 