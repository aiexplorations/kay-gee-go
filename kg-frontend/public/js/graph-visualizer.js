/**
 * Graph Visualizer using Three.js
 */
class GraphVisualizer {
    constructor(containerId) {
        this.containerId = containerId;
        this.container = document.getElementById(containerId);
        this.width = this.container.clientWidth;
        this.height = this.container.clientHeight;
        this.nodes = [];
        this.links = [];
        this.nodeObjects = new Map();
        this.linkObjects = new Map();
        this.selectedNode = null;
        this.hoveredNode = null;
        this.tooltip = null;
        this.raycaster = new THREE.Raycaster();
        this.mouse = new THREE.Vector2();
        this.forceSimulation = null;
        this.nodePositions = new Map();
        
        // Initialize the visualizer
        this.init();
        console.log("GraphVisualizer initialized with container:", containerId);
    }

    /**
     * Initialize the Three.js scene
     */
    init() {
        console.log("Initializing Three.js scene");
        
        // Create scene
        this.scene = new THREE.Scene();
        this.scene.background = new THREE.Color(0xf0f0f0);
        
        // Create camera
        this.camera = new THREE.PerspectiveCamera(75, this.width / this.height, 0.1, 1000);
        this.camera.position.z = 200;
        
        // Create renderer
        this.renderer = new THREE.WebGLRenderer({ antialias: true });
        this.renderer.setSize(this.width, this.height);
        this.renderer.setPixelRatio(window.devicePixelRatio);
        this.container.appendChild(this.renderer.domElement);
        
        // Create controls
        this.controls = new THREE.OrbitControls(this.camera, this.renderer.domElement);
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.25;
        
        // Add lights
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.6);
        this.scene.add(ambientLight);
        
        const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
        directionalLight.position.set(1, 1, 1);
        this.scene.add(directionalLight);
        
        // Create tooltip
        this.createTooltip();
        
        // Add event listeners
        this.addEventListeners();
        
        console.log("Three.js scene initialized");
    }

    /**
     * Create tooltip element
     */
    createTooltip() {
        this.tooltip = document.createElement('div');
        this.tooltip.className = 'tooltip';
        this.tooltip.style.display = 'none';
        this.tooltip.style.position = 'absolute';
        this.tooltip.style.backgroundColor = 'rgba(0, 0, 0, 0.7)';
        this.tooltip.style.color = 'white';
        this.tooltip.style.padding = '5px 10px';
        this.tooltip.style.borderRadius = '5px';
        this.tooltip.style.pointerEvents = 'none';
        this.tooltip.style.zIndex = '1000';
        document.body.appendChild(this.tooltip);
    }

    /**
     * Add event listeners
     */
    addEventListeners() {
        this.renderer.domElement.addEventListener('mousemove', this.onMouseMove.bind(this));
        this.renderer.domElement.addEventListener('click', this.onClick.bind(this));
        window.addEventListener('resize', this.onWindowResize.bind(this));
    }

    /**
     * Handle window resize
     */
    onWindowResize() {
        this.width = this.container.clientWidth;
        this.height = this.container.clientHeight;
        
        this.camera.aspect = this.width / this.height;
        this.camera.updateProjectionMatrix();
        
        this.renderer.setSize(this.width, this.height);
    }

    /**
     * Handle mouse move event
     * @param {MouseEvent} event - Mouse event
     */
    onMouseMove(event) {
        // Calculate mouse position in normalized device coordinates
        const rect = this.renderer.domElement.getBoundingClientRect();
        this.mouse.x = ((event.clientX - rect.left) / this.width) * 2 - 1;
        this.mouse.y = -((event.clientY - rect.top) / this.height) * 2 + 1;
        
        // Update the picking ray with the camera and mouse position
        this.raycaster.setFromCamera(this.mouse, this.camera);
        
        // Calculate objects intersecting the picking ray
        const intersects = this.raycaster.intersectObjects(Array.from(this.nodeObjects.values()));
        
        if (intersects.length > 0) {
            // Get the first intersected object
            const object = intersects[0].object;
            
            // Set hovered node
            this.hoveredNode = object.userData.node;
            
            // Update tooltip
            this.tooltip.textContent = this.hoveredNode.name;
            this.tooltip.style.display = 'block';
            this.tooltip.style.left = `${event.clientX + 10}px`;
            this.tooltip.style.top = `${event.clientY + 10}px`;
            
            // Change cursor
            this.renderer.domElement.style.cursor = 'pointer';
        } else {
            // Reset hovered node
            this.hoveredNode = null;
            
            // Hide tooltip
            this.tooltip.style.display = 'none';
            
            // Reset cursor
            this.renderer.domElement.style.cursor = 'auto';
        }
    }

    /**
     * Handle click event
     * @param {MouseEvent} event - Mouse event
     */
    onClick(event) {
        // Calculate mouse position in normalized device coordinates
        const rect = this.renderer.domElement.getBoundingClientRect();
        this.mouse.x = ((event.clientX - rect.left) / this.width) * 2 - 1;
        this.mouse.y = -((event.clientY - rect.top) / this.height) * 2 + 1;
        
        // Update the picking ray with the camera and mouse position
        this.raycaster.setFromCamera(this.mouse, this.camera);
        
        // Calculate objects intersecting the picking ray
        const intersects = this.raycaster.intersectObjects(Array.from(this.nodeObjects.values()));
        
        if (intersects.length > 0) {
            // Get the first intersected object
            const object = intersects[0].object;
            
            // Set selected node
            this.selectedNode = object.userData.node;
            
            // Highlight selected node
            this.highlightNode(this.selectedNode);
            
            console.log('Selected node:', this.selectedNode);
        }
    }

    /**
     * Highlight a node and its connections
     * @param {Object} node - The node to highlight
     */
    highlightNode(node) {
        // Reset all nodes and links
        this.nodeObjects.forEach((object, id) => {
            object.material.color.set(0x3498db);
            object.material.opacity = 0.8;
            object.scale.set(1, 1, 1);
        });
        
        this.linkObjects.forEach((object) => {
            object.material.color.set(0xbdc3c7);
            object.material.opacity = 0.3;
        });
        
        // Highlight selected node
        const nodeObject = this.nodeObjects.get(node.id);
        if (nodeObject) {
            nodeObject.material.color.set(0xe74c3c);
            nodeObject.material.opacity = 1.0;
            nodeObject.scale.set(1.2, 1.2, 1.2);
        }
        
        // Highlight connected nodes and links
        this.links.forEach(link => {
            if (link.source === node.id || link.target === node.id) {
                // Highlight link
                const linkObject = this.linkObjects.get(`${link.source}-${link.target}`);
                if (linkObject) {
                    linkObject.material.color.set(0xe74c3c);
                    linkObject.material.opacity = 0.8;
                }
                
                // Highlight connected node
                const connectedId = link.source === node.id ? link.target : link.source;
                const connectedObject = this.nodeObjects.get(connectedId);
                if (connectedObject) {
                    connectedObject.material.color.set(0xf39c12);
                    connectedObject.material.opacity = 0.9;
                    connectedObject.scale.set(1.1, 1.1, 1.1);
                }
            }
        });
    }

    /**
     * Set graph data
     * @param {Array} nodes - Array of nodes
     * @param {Array} links - Array of links
     */
    setData(nodes, links) {
        console.log(`Setting graph data: ${nodes.length} nodes, ${links.length} links`);
        this.nodes = nodes;
        this.links = links;
        
        // Clear existing objects
        this.clearScene();
        
        // Create node objects
        this.createNodeObjects();
        
        // Create link objects
        this.createLinkObjects();
    }

    /**
     * Clear the scene
     */
    clearScene() {
        // Remove node objects
        this.nodeObjects.forEach(object => {
            this.scene.remove(object);
        });
        this.nodeObjects.clear();
        
        // Remove link objects
        this.linkObjects.forEach(object => {
            this.scene.remove(object);
        });
        this.linkObjects.clear();
        
        // Clear node positions
        this.nodePositions.clear();
    }

    /**
     * Create node objects
     */
    createNodeObjects() {
        const nodeGeometry = new THREE.SphereGeometry(1, 16, 16);
        
        this.nodes.forEach(node => {
            // Create material
            const nodeMaterial = new THREE.MeshPhongMaterial({
                color: 0x3498db,
                transparent: true,
                opacity: 0.8,
                shininess: 30
            });
            
            // Create mesh
            const nodeMesh = new THREE.Mesh(nodeGeometry, nodeMaterial);
            
            // Scale node based on size
            const scale = node.size || 5;
            nodeMesh.scale.set(scale, scale, scale);
            
            // Set initial position
            const position = this.nodePositions.get(node.id) || {
                x: (Math.random() - 0.5) * 300,
                y: (Math.random() - 0.5) * 300,
                z: (Math.random() - 0.5) * 300
            };
            
            nodeMesh.position.set(position.x, position.y, position.z);
            
            // Store position
            this.nodePositions.set(node.id, {
                x: nodeMesh.position.x,
                y: nodeMesh.position.y,
                z: nodeMesh.position.z
            });
            
            // Set user data
            nodeMesh.userData.node = node;
            
            // Add to scene
            this.scene.add(nodeMesh);
            
            // Store in map
            this.nodeObjects.set(node.id, nodeMesh);
        });
        
        console.log(`Created ${this.nodeObjects.size} node objects`);
    }

    /**
     * Create link objects
     */
    createLinkObjects() {
        this.links.forEach(link => {
            // Get source and target nodes
            const sourceNode = this.nodeObjects.get(link.source);
            const targetNode = this.nodeObjects.get(link.target);
            
            if (sourceNode && targetNode) {
                // Create material
                const linkMaterial = new THREE.LineBasicMaterial({
                    color: 0xbdc3c7,
                    transparent: true,
                    opacity: 0.3
                });
                
                // Create geometry
                const linkGeometry = new THREE.BufferGeometry().setFromPoints([
                    sourceNode.position,
                    targetNode.position
                ]);
                
                // Create line
                const linkLine = new THREE.Line(linkGeometry, linkMaterial);
                
                // Set user data
                linkLine.userData.link = link;
                
                // Add to scene
                this.scene.add(linkLine);
                
                // Store in map
                this.linkObjects.set(`${link.source}-${link.target}`, linkLine);
            }
        });
        
        console.log(`Created ${this.linkObjects.size} link objects`);
    }

    /**
     * Start force simulation
     */
    startSimulation() {
        console.log("Starting force simulation");
        
        if (typeof d3 === 'undefined') {
            console.error("D3.js is not loaded. Cannot start force simulation.");
            return;
        }
        
        // Create force simulation
        this.forceSimulation = d3.forceSimulation()
            .nodes(this.nodes)
            .force('link', d3.forceLink().id(d => d.id).links(this.links).distance(50))
            .force('charge', d3.forceManyBody().strength(-100))
            .force('center', d3.forceCenter(0, 0))
            .force('collision', d3.forceCollide().radius(d => (d.size || 5) * 1.5))
            .on('tick', this.updatePositions.bind(this));
        
        // Set alpha target to keep simulation running
        this.forceSimulation.alphaTarget(0.1).restart();
        
        // After 3 seconds, stop the simulation
        setTimeout(() => {
            this.forceSimulation.alphaTarget(0);
        }, 3000);
    }

    /**
     * Update positions based on force simulation
     */
    updatePositions() {
        // Update node positions
        this.nodes.forEach(node => {
            const nodeObject = this.nodeObjects.get(node.id);
            if (nodeObject) {
                nodeObject.position.x = node.x;
                nodeObject.position.y = node.y;
                nodeObject.position.z = node.z || 0;
                
                // Store position
                this.nodePositions.set(node.id, {
                    x: nodeObject.position.x,
                    y: nodeObject.position.y,
                    z: nodeObject.position.z
                });
            }
        });
        
        // Update link positions
        this.links.forEach(link => {
            const linkObject = this.linkObjects.get(`${link.source.id || link.source}-${link.target.id || link.target}`);
            const sourceNode = this.nodeObjects.get(link.source.id || link.source);
            const targetNode = this.nodeObjects.get(link.target.id || link.target);
            
            if (linkObject && sourceNode && targetNode) {
                // Update geometry
                const positions = new Float32Array([
                    sourceNode.position.x, sourceNode.position.y, sourceNode.position.z,
                    targetNode.position.x, targetNode.position.y, targetNode.position.z
                ]);
                
                linkObject.geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
                linkObject.geometry.attributes.position.needsUpdate = true;
            }
        });
    }

    /**
     * Reset camera position
     */
    resetCamera() {
        this.camera.position.set(0, 0, 200);
        this.camera.lookAt(0, 0, 0);
        this.controls.reset();
    }

    /**
     * Render the scene
     */
    render() {
        // Update controls
        this.controls.update();
        
        // Render scene
        this.renderer.render(this.scene, this.camera);
    }
} 