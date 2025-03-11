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
        this.hoveredLink = null;
        this.tooltip = null;
        this.raycaster = new THREE.Raycaster();
        this.mouse = new THREE.Vector2();
        this.forceSimulation = null;
        this.nodePositions = new Map();
        this.nodeInfoPanel = document.getElementById('node-info-panel');
        this.nodeInfoTitle = document.getElementById('node-info-title');
        this.nodeInfoContent = document.getElementById('node-info-content');
        this.nodeConnections = document.getElementById('node-connections');
        this.closeInfoPanelBtn = document.getElementById('close-info-panel');
        this.stars = [];
        
        // Initialize the visualizer
        this.init();
        console.log("GraphVisualizer initialized with container:", containerId);
        
        // Add event listener for close button
        this.closeInfoPanelBtn.addEventListener('click', () => {
            this.hideNodeInfoPanel();
        });
    }

    /**
     * Initialize the Three.js scene
     */
    init() {
        console.log("Initializing Three.js scene");
        
        // Create scene
        this.scene = new THREE.Scene();
        this.scene.background = new THREE.Color(0x050510);
        
        // Create camera
        this.camera = new THREE.PerspectiveCamera(75, this.width / this.height, 0.1, 2000);
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
        this.controls.rotateSpeed = 0.5;
        this.controls.zoomSpeed = 1.2;
        
        // Add lights
        const ambientLight = new THREE.AmbientLight(0x404080, 0.6);
        this.scene.add(ambientLight);
        
        const directionalLight = new THREE.DirectionalLight(0x6080ff, 0.8);
        directionalLight.position.set(1, 1, 1);
        this.scene.add(directionalLight);
        
        // Add hemisphere light for better 3D effect
        const hemisphereLight = new THREE.HemisphereLight(0x4060ff, 0x404080, 0.6);
        this.scene.add(hemisphereLight);
        
        // Add stars to the background
        this.createStars();
        
        // Create tooltip
        this.createTooltip();
        
        // Add event listeners
        this.addEventListeners();
        
        console.log("Three.js scene initialized");
    }
    
    /**
     * Create stars in the background
     */
    createStars() {
        const starsGeometry = new THREE.BufferGeometry();
        const starsMaterial = new THREE.PointsMaterial({
            color: 0xffffff,
            size: 1,
            transparent: true,
            opacity: 0.8,
            sizeAttenuation: true
        });
        
        const starsVertices = [];
        for (let i = 0; i < 5000; i++) {
            const x = (Math.random() - 0.5) * 2000;
            const y = (Math.random() - 0.5) * 2000;
            const z = (Math.random() - 0.5) * 2000;
            starsVertices.push(x, y, z);
        }
        
        starsGeometry.setAttribute('position', new THREE.Float32BufferAttribute(starsVertices, 3));
        
        const stars = new THREE.Points(starsGeometry, starsMaterial);
        this.scene.add(stars);
        this.stars.push(stars);
    }

    /**
     * Create tooltip element
     */
    createTooltip() {
        this.tooltip = document.createElement('div');
        this.tooltip.className = 'tooltip';
        this.tooltip.style.display = 'none';
        this.tooltip.style.position = 'absolute';
        this.tooltip.style.backgroundColor = 'rgba(10, 10, 30, 0.9)';
        this.tooltip.style.color = 'white';
        this.tooltip.style.padding = '10px 15px';
        this.tooltip.style.borderRadius = '6px';
        this.tooltip.style.pointerEvents = 'none';
        this.tooltip.style.zIndex = '1000';
        this.tooltip.style.fontSize = '14px';
        this.tooltip.style.maxWidth = '200px';
        this.tooltip.style.boxShadow = '0 0 15px rgba(0, 0, 255, 0.3)';
        this.tooltip.style.border = '1px solid rgba(100, 100, 255, 0.3)';
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
        
        // First check for node intersections
        const nodeIntersects = this.raycaster.intersectObjects(Array.from(this.nodeObjects.values()));
        
        if (nodeIntersects.length > 0) {
            // Get the first intersected object
            const object = nodeIntersects[0].object;
            
            // Set hovered node
            this.hoveredNode = object.userData.node;
            this.hoveredLink = null;
            
            // Update tooltip
            this.tooltip.textContent = this.hoveredNode.name;
            this.tooltip.style.display = 'block';
            this.tooltip.style.left = `${event.clientX + 10}px`;
            this.tooltip.style.top = `${event.clientY + 10}px`;
            
            // Change cursor
            this.renderer.domElement.style.cursor = 'pointer';
        } else {
            // If no nodes are intersected, check for link intersections
            // We need to set a threshold for the raycaster to detect thin lines
            this.raycaster.params.Line.threshold = 5; // Increased threshold for easier selection
            const linkIntersects = this.raycaster.intersectObjects(Array.from(this.linkObjects.values()));
            
            if (linkIntersects.length > 0) {
                // Get the first intersected link
                const object = linkIntersects[0].object;
                
                // Set hovered link
                this.hoveredLink = object.userData.link;
                this.hoveredNode = null;
                
                // Find source and target node names
                const sourceNode = this.nodes.find(node => node.id === this.hoveredLink.source);
                const targetNode = this.nodes.find(node => node.id === this.hoveredLink.target);
                
                // Update tooltip with relationship information
                if (sourceNode && targetNode) {
                    this.tooltip.innerHTML = `<strong>${sourceNode.name}</strong><br>
                                            <span style="color: #64b5f6;">-[${this.hoveredLink.type}]-></span><br>
                                            <strong>${targetNode.name}</strong>`;
                    this.tooltip.style.display = 'block';
                    this.tooltip.style.left = `${event.clientX + 10}px`;
                    this.tooltip.style.top = `${event.clientY + 10}px`;
                }
                
                // Change cursor
                this.renderer.domElement.style.cursor = 'pointer';
            } else {
                // Reset hovered objects
                this.hoveredNode = null;
                this.hoveredLink = null;
                
                // Hide tooltip
                this.tooltip.style.display = 'none';
                
                // Reset cursor
                this.renderer.domElement.style.cursor = 'auto';
            }
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
            
            // Show node info panel
            this.showNodeInfoPanel(this.selectedNode);
            
            console.log('Selected node:', this.selectedNode);
        } else {
            // If clicking on empty space, hide the node info panel
            this.hideNodeInfoPanel();
        }
    }

    /**
     * Show node information panel
     * @param {Object} node - The node to display information for
     */
    showNodeInfoPanel(node) {
        if (!node) return;
        
        // Set node title
        this.nodeInfoTitle.textContent = node.name;
        
        // Set node content (if any additional properties)
        let contentHTML = '';
        for (const [key, value] of Object.entries(node)) {
            if (key !== 'id' && key !== 'name' && key !== 'size' && key !== 'x' && key !== 'y' && key !== 'z') {
                contentHTML += `<p><strong>${key}:</strong> ${value}</p>`;
            }
        }
        this.nodeInfoContent.innerHTML = contentHTML || '<p>No additional properties</p>';
        
        // Find connections (1-hop neighbors)
        const connections = [];
        
        // Debug log to check node ID format
        console.log('Selected node ID:', node.id);
        console.log('All links:', this.links);
        
        this.links.forEach(link => {
            // Handle both string IDs and object references
            const sourceId = typeof link.source === 'object' ? link.source.id : link.source;
            const targetId = typeof link.target === 'object' ? link.target.id : link.target;
            
            console.log(`Checking link: ${sourceId} -> ${targetId}, node.id: ${node.id}`);
            
            if (sourceId === node.id) {
                // Find target node
                const targetNode = this.nodes.find(n => n.id === targetId);
                if (targetNode) {
                    connections.push({
                        node: targetNode,
                        type: link.type,
                        direction: 'outgoing'
                    });
                    console.log(`Found outgoing connection to ${targetNode.name}`);
                }
            } else if (targetId === node.id) {
                // Find source node
                const sourceNode = this.nodes.find(n => n.id === sourceId);
                if (sourceNode) {
                    connections.push({
                        node: sourceNode,
                        type: link.type,
                        direction: 'incoming'
                    });
                    console.log(`Found incoming connection from ${sourceNode.name}`);
                }
            }
        });
        
        // Display connections
        let connectionsHTML = '';
        if (connections.length > 0) {
            connections.forEach(conn => {
                connectionsHTML += `
                    <div class="connection-item">
                        <span class="connection-type">${conn.direction === 'outgoing' ? '' : '← '}${conn.type}${conn.direction === 'outgoing' ? ' →' : ''}</span>
                        <div class="connection-name">${conn.node.name}</div>
                    </div>
                `;
            });
        } else {
            connectionsHTML = '<p>No connections found.</p>';
        }
        this.nodeConnections.innerHTML = connectionsHTML;
        
        // Show the panel
        this.nodeInfoPanel.style.display = 'block';
    }

    /**
     * Hide node information panel
     */
    hideNodeInfoPanel() {
        this.nodeInfoPanel.style.display = 'none';
    }

    /**
     * Highlight a node and its connections
     * @param {Object} node - The node to highlight
     */
    highlightNode(node) {
        // Reset all nodes and links
        this.nodeObjects.forEach((object, id) => {
            object.material.color.set(0x4080ff);
            object.material.opacity = 0.8;
            object.scale.set(1, 1, 1);
        });
        
        this.linkObjects.forEach((object) => {
            object.material.color.set(0x6080c0);
            object.material.opacity = 0.3;
        });
        
        // Highlight selected node
        const nodeObject = this.nodeObjects.get(node.id);
        if (nodeObject) {
            nodeObject.material.color.set(0xff4080);
            nodeObject.material.opacity = 1.0;
            nodeObject.scale.set(1.2, 1.2, 1.2);
        }
        
        // Highlight connected nodes and links
        this.links.forEach(link => {
            if (link.source === node.id || link.target === node.id) {
                // Highlight link
                const linkObject = this.linkObjects.get(`${link.source}-${link.target}`);
                if (linkObject) {
                    linkObject.material.color.set(0xff4080);
                    linkObject.material.opacity = 0.8;
                }
                
                // Highlight connected node
                const connectedId = link.source === node.id ? link.target : link.source;
                const connectedObject = this.nodeObjects.get(connectedId);
                if (connectedObject) {
                    connectedObject.material.color.set(0x40ffff);
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
        const nodeGeometry = new THREE.SphereGeometry(1, 24, 24);
        
        this.nodes.forEach(node => {
            // Create material
            const nodeMaterial = new THREE.MeshPhongMaterial({
                color: 0x4080ff,
                transparent: true,
                opacity: 0.8,
                shininess: 50,
                emissive: 0x102040,
                emissiveIntensity: 0.3
            });
            
            // Create mesh
            const nodeMesh = new THREE.Mesh(nodeGeometry, nodeMaterial);
            
            // Scale node based on size
            const scale = node.size || 5;
            nodeMesh.scale.set(scale, scale, scale);
            
            // Set initial position with true 3D coordinates
            const position = this.nodePositions.get(node.id) || {
                x: (Math.random() - 0.5) * 400,
                y: (Math.random() - 0.5) * 400,
                z: (Math.random() - 0.5) * 400
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
                // Create material with thicker lines
                const linkMaterial = new THREE.LineBasicMaterial({
                    color: 0x6080c0,
                    transparent: true,
                    opacity: 0.5,
                    linewidth: 3 // Note: linewidth only works in certain browsers/GPUs
                });
                
                // Create geometry
                const linkGeometry = new THREE.BufferGeometry();
                
                // Create points for the line
                const points = [
                    sourceNode.position.clone(),
                    targetNode.position.clone()
                ];
                
                // Set positions
                linkGeometry.setFromPoints(points);
                
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
        
        // Create force simulation with 3D forces
        this.forceSimulation = d3.forceSimulation()
            .nodes(this.nodes)
            .force('link', d3.forceLink().id(d => d.id).links(this.links).distance(100))
            .force('charge', d3.forceManyBody().strength(-100)) // Reduced strength
            .force('center', d3.forceCenter(0, 0).strength(0.05)) // Reduced center force
            .force('x', d3.forceX().strength(0.02)) // Reduced x force
            .force('y', d3.forceY().strength(0.02)) // Reduced y force
            .force('z', d3.forceZ().strength(0.01)) // Reduced z force for more stability
            .force('collision', d3.forceCollide().radius(d => (d.size || 5) * 1.8))
            .on('tick', this.updatePositions.bind(this));
        
        // Set alpha target to keep simulation running but at a lower value
        this.forceSimulation.alphaTarget(0.05).restart();
        
        // After 10 seconds, gradually slow down the simulation
        setTimeout(() => {
            this.forceSimulation.alphaTarget(0.01);
            
            // After another 5 seconds, stop the simulation completely
            setTimeout(() => {
                this.forceSimulation.alphaTarget(0);
            }, 5000);
        }, 10000);
    }

    /**
     * Update positions based on force simulation
     */
    updatePositions() {
        // Update node positions
        this.nodes.forEach(node => {
            const nodeObject = this.nodeObjects.get(node.id);
            if (nodeObject) {
                // Apply damping to make movements smoother
                const damping = 0.8;
                
                // Get current position
                const currentPos = this.nodePositions.get(node.id) || {
                    x: nodeObject.position.x,
                    y: nodeObject.position.y,
                    z: nodeObject.position.z
                };
                
                // Calculate new position with damping
                const newX = currentPos.x + (node.x - currentPos.x) * damping;
                const newY = currentPos.y + (node.y - currentPos.y) * damping;
                
                // For z, maintain the current z position if it exists, otherwise use a small random offset
                // This prevents the "flying through" effect
                const newZ = currentPos.z || (Math.random() - 0.5) * 100;
                
                // Update position
                nodeObject.position.x = newX;
                nodeObject.position.y = newY;
                nodeObject.position.z = newZ;
                
                // Store position
                this.nodePositions.set(node.id, {
                    x: newX,
                    y: newY,
                    z: newZ
                });
                
                // Update node z position in the data
                node.z = newZ;
            }
        });
        
        // Update link positions
        this.links.forEach(link => {
            const linkObject = this.linkObjects.get(`${link.source.id || link.source}-${link.target.id || link.target}`);
            const sourceNode = this.nodeObjects.get(link.source.id || link.source);
            const targetNode = this.nodeObjects.get(link.target.id || link.target);
            
            if (linkObject && sourceNode && targetNode) {
                // Update line geometry
                const points = [
                    sourceNode.position.clone(),
                    targetNode.position.clone()
                ];
                
                // Update geometry
                linkObject.geometry.dispose();
                linkObject.geometry = new THREE.BufferGeometry().setFromPoints(points);
            }
        });
        
        // Rotate stars slowly for a space-like effect
        this.stars.forEach(stars => {
            stars.rotation.x += 0.0001;
            stars.rotation.y += 0.0001;
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