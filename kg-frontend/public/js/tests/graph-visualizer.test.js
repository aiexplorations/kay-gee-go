/**
 * Tests for the GraphVisualizer component
 */

// Mock THREE.js
global.THREE = {
  Scene: jest.fn(() => ({
    add: jest.fn(),
    background: { set: jest.fn() }
  })),
  PerspectiveCamera: jest.fn(() => ({
    position: { z: 0, set: jest.fn() },
    aspect: 1,
    updateProjectionMatrix: jest.fn(),
    lookAt: jest.fn()
  })),
  WebGLRenderer: jest.fn(() => ({
    setSize: jest.fn(),
    setPixelRatio: jest.fn(),
    domElement: {
      addEventListener: jest.fn(),
      style: {}
    },
    render: jest.fn()
  })),
  OrbitControls: jest.fn(() => ({
    enableDamping: false,
    dampingFactor: 0,
    rotateSpeed: 0,
    zoomSpeed: 0,
    update: jest.fn(),
    reset: jest.fn()
  })),
  AmbientLight: jest.fn(() => ({
    position: { set: jest.fn() }
  })),
  DirectionalLight: jest.fn(() => ({
    position: { set: jest.fn() }
  })),
  HemisphereLight: jest.fn(),
  BufferGeometry: jest.fn(() => ({
    setAttribute: jest.fn(),
    setFromPoints: jest.fn()
  })),
  Float32BufferAttribute: jest.fn(),
  PointsMaterial: jest.fn(),
  Points: jest.fn(() => ({
    rotation: { x: 0, y: 0 }
  })),
  SphereGeometry: jest.fn(),
  MeshPhongMaterial: jest.fn(),
  Mesh: jest.fn(() => ({
    position: { x: 0, y: 0, z: 0, copy: jest.fn() },
    scale: { set: jest.fn() },
    userData: {},
    quaternion: { copy: jest.fn() }
  })),
  LineBasicMaterial: jest.fn(),
  Line: jest.fn(() => ({
    userData: {},
    geometry: {
      dispose: jest.fn(),
      attributes: {
        position: { needsUpdate: false }
      }
    }
  })),
  Vector2: jest.fn(() => ({ x: 0, y: 0 })),
  Vector3: jest.fn(() => ({
    x: 0, y: 0, z: 0,
    subVectors: jest.fn(() => ({
      normalize: jest.fn(),
      length: jest.fn(() => 1)
    })),
    normalize: jest.fn(() => ({ x: 0, y: 0, z: 0 })),
    clone: jest.fn(() => ({ x: 0, y: 0, z: 0 }))
  })),
  Quaternion: jest.fn(() => ({
    setFromUnitVectors: jest.fn(),
    copy: jest.fn()
  })),
  Raycaster: jest.fn(() => ({
    setFromCamera: jest.fn(),
    intersectObjects: jest.fn(() => []),
    params: {
      Line: { threshold: 0 }
    }
  })),
  Color: jest.fn()
};

// Mock document
document.getElementById = jest.fn(() => ({
  clientWidth: 800,
  clientHeight: 600,
  appendChild: jest.fn()
}));

document.createElement = jest.fn(() => ({
  className: '',
  style: {},
  addEventListener: jest.fn()
}));

document.body = {
  appendChild: jest.fn()
};

// Import the GraphVisualizer class
const GraphVisualizer = require('../graph-visualizer');

describe('GraphVisualizer', () => {
  let visualizer;
  
  beforeEach(() => {
    // Reset mocks
    jest.clearAllMocks();
    
    // Create a new instance
    visualizer = new GraphVisualizer('graph-canvas');
  });
  
  test('should initialize correctly', () => {
    expect(visualizer).toBeDefined();
    expect(visualizer.scene).toBeDefined();
    expect(visualizer.camera).toBeDefined();
    expect(visualizer.renderer).toBeDefined();
    expect(visualizer.controls).toBeDefined();
    expect(visualizer.nodes).toEqual([]);
    expect(visualizer.links).toEqual([]);
    expect(visualizer.nodeObjects).toBeInstanceOf(Map);
    expect(visualizer.linkObjects).toBeInstanceOf(Map);
  });
  
  test('should create node objects', () => {
    // Setup
    const mockNodes = [
      { id: '1', name: 'Node 1', size: 5 },
      { id: '2', name: 'Node 2', size: 7 }
    ];
    visualizer.nodes = mockNodes;
    
    // Execute
    visualizer.createNodeObjects();
    
    // Verify
    expect(global.THREE.SphereGeometry).toHaveBeenCalled();
    expect(global.THREE.MeshPhongMaterial).toHaveBeenCalledTimes(2);
    expect(global.THREE.Mesh).toHaveBeenCalledTimes(2);
    expect(visualizer.nodeObjects.size).toBe(2);
  });
  
  test('should create link objects', () => {
    // Setup
    const mockNodes = [
      { id: '1', name: 'Node 1', size: 5 },
      { id: '2', name: 'Node 2', size: 7 }
    ];
    const mockLinks = [
      { source: '1', target: '2', type: 'RELATES_TO' }
    ];
    visualizer.nodes = mockNodes;
    visualizer.links = mockLinks;
    
    // Create node objects first
    visualizer.createNodeObjects();
    
    // Execute
    visualizer.createLinkObjects();
    
    // Verify
    expect(global.THREE.LineBasicMaterial).toHaveBeenCalled();
    expect(global.THREE.BufferGeometry).toHaveBeenCalled();
    expect(global.THREE.Line).toHaveBeenCalled();
    expect(visualizer.linkObjects.size).toBe(1);
  });
  
  test('should handle node click', () => {
    // Setup
    const mockNode = { id: '1', name: 'Node 1' };
    const mockMesh = { userData: { node: mockNode } };
    global.THREE.Raycaster.prototype.intersectObjects = jest.fn(() => [{ object: mockMesh }]);
    
    // Mock methods
    visualizer.highlightNode = jest.fn();
    visualizer.showNodeInfoPanel = jest.fn();
    
    // Execute
    visualizer.onClick({ clientX: 100, clientY: 100 });
    
    // Verify
    expect(visualizer.selectedNode).toBe(mockNode);
    expect(visualizer.highlightNode).toHaveBeenCalledWith(mockNode);
    expect(visualizer.showNodeInfoPanel).toHaveBeenCalledWith(mockNode);
  });
  
  test('should handle node hover', () => {
    // Setup
    const mockNode = { id: '1', name: 'Node 1' };
    const mockMesh = { userData: { node: mockNode } };
    global.THREE.Raycaster.prototype.intersectObjects = jest.fn(() => [{ object: mockMesh }]);
    
    // Execute
    visualizer.onMouseMove({ clientX: 100, clientY: 100 });
    
    // Verify
    expect(visualizer.hoveredNode).toBe(mockNode);
    expect(visualizer.tooltip.textContent).toBe(mockNode.name);
    expect(visualizer.tooltip.style.display).toBe('block');
  });
  
  test('should handle link hover', () => {
    // Setup
    const mockLink = { 
      source: '1', 
      target: '2', 
      type: 'RELATES_TO' 
    };
    const mockLine = { userData: { link: mockLink } };
    
    // First return empty array for node intersections, then return link intersection
    global.THREE.Raycaster.prototype.intersectObjects = jest.fn()
      .mockReturnValueOnce([]) // No node intersections
      .mockReturnValueOnce([{ object: mockLine }]); // Link intersection
    
    visualizer.nodes = [
      { id: '1', name: 'Node 1' },
      { id: '2', name: 'Node 2' }
    ];
    
    // Execute
    visualizer.onMouseMove({ clientX: 100, clientY: 100 });
    
    // Verify
    expect(visualizer.hoveredLink).toBe(mockLink);
    expect(visualizer.tooltip.style.display).toBe('block');
  });
  
  test('should reset camera', () => {
    // Execute
    visualizer.resetCamera();
    
    // Verify
    expect(visualizer.camera.position.set).toHaveBeenCalledWith(0, 0, 200);
    expect(visualizer.camera.lookAt).toHaveBeenCalledWith(0, 0, 0);
    expect(visualizer.controls.reset).toHaveBeenCalled();
  });
  
  test('should handle window resize', () => {
    // Setup
    visualizer.container.clientWidth = 1000;
    visualizer.container.clientHeight = 800;
    
    // Execute
    visualizer.onWindowResize();
    
    // Verify
    expect(visualizer.width).toBe(1000);
    expect(visualizer.height).toBe(800);
    expect(visualizer.camera.updateProjectionMatrix).toHaveBeenCalled();
    expect(visualizer.renderer.setSize).toHaveBeenCalledWith(1000, 800);
  });
}); 