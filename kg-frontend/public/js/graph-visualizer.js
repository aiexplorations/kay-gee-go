/**
 * Graph Visualizer component for rendering the knowledge graph in 3D
 */
class GraphVisualizer {
    constructor(containerId) {
        this.containerId = containerId;
        this.container = document.getElementById(containerId);
        
        // Initialize Three.js scene
        this.scene = new THREE.Scene();
        this.camera = new THREE.PerspectiveCamera(60, window.innerWidth / window.innerHeight, 0.1, 2000);
        this.renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
        
        // Set renderer properties
        this.renderer.setSize(this.container.clientWidth, this.container.clientHeight);
        this.renderer.setClearColor(0x000000, 0);
        this.renderer.setPixelRatio(window.devicePixelRatio);
        this.container.appendChild(this.renderer.domElement);
        
        // Set camera position
        this.camera.position.set(0, 0, 300);
        
        // Add orbit controls
        this.controls = new THREE.OrbitControls(this.camera, this.renderer.domElement);
        this.controls.enableDamping = true;
        this.controls.dampingFactor = 0.25;
        
        // Add ambient light
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.5);
        this.scene.add(ambientLight);
        
        // Add directional light
        const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
        directionalLight.position.set(1, 1, 1);
        this.scene.add(directionalLight);
        
        // Initialize node info panel
        this.nodeInfoPanel = document.getElementById('node-info-panel');
        this.nodeInfoTitle = document.getElementById('node-info-title');
        this.nodeInfoContent = document.getElementById('node-info-content');
        this.nodeConnections = document.getElementById('node-connections');
        this.closeInfoPanelBtn = document.getElementById('close-info-panel');
        
        // Close info panel when close button is clicked
        this.closeInfoPanelBtn.addEventListener('click', () => {
            this.nodeInfoPanel.style.display = 'none';
        });
        
        // Initialize data structures
        this.nodes = [];
        this.links = [];
        this.nodeObjects = new Map();
        this.linkObjects = new Map();
        
        // Initialize CSS variables
        this.initCssVariables();
        
        // Initialize raycaster for node selection
        this.raycaster = new THREE.Raycaster();
        this.mouse = new THREE.Vector2();
        this.selectedNode = null;
        
        // Initialize simulation
        this.simulation = null;
        
        // Add event listeners
        this.addEventListeners();
    }
    
    /**
     * Initialize CSS variables from the :root element
     */
    initCssVariables() {
        // Node colors
        this.nodeBaseColor = this.getCssVariable('--node-base-color', '#0088ff');
        this.nodeHighlightColor = this.getCssVariable('--node-highlight-color', '#00ccff');
        this.nodeSelectedColor = this.getCssVariable('--node-selected-color', '#ff3333');
        this.nodeEmissiveIntensityNormal = parseFloat(this.getCssVariable('--node-emissive-intensity-normal', '0.3'));
        this.nodeEmissiveIntensityHover = parseFloat(this.getCssVariable('--node-emissive-intensity-hover', '0.8'));
        this.nodeEmissiveIntensitySelected = parseFloat(this.getCssVariable('--node-emissive-intensity-selected', '1.0'));
        this.nodeShininess = parseFloat(this.getCssVariable('--node-shininess', '50'));
        
        // Edge colors
        this.edgeBaseColor = this.getCssVariable('--edge-base-color', '#88aaff');
        this.edgeSelectedColor = this.getCssVariable('--edge-selected-color', '#ff5555');
        this.edgeOpacity = parseFloat(this.getCssVariable('--edge-opacity', '0.7'));
        this.edgeWidth = parseFloat(this.getCssVariable('--edge-width', '1'));
        
        // Node sizes
        this.nodeSizeMin = parseFloat(this.getCssVariable('--node-size-min', '4'));
        this.nodeSizeMax = parseFloat(this.getCssVariable('--node-size-max', '12'));
        this.nodeSizeMultiplier = parseFloat(this.getCssVariable('--node-size-multiplier', '0.7'));
        
        // Force simulation parameters
        this.forceChargeStrength = parseFloat(this.getCssVariable('--force-charge-strength', '-200'));
        this.forceLinkStrength = parseFloat(this.getCssVariable('--force-link-strength', '0.7'));
        this.forceCenterStrength = parseFloat(this.getCssVariable('--force-center-strength', '1'));
        this.forcePositionStrength = parseFloat(this.getCssVariable('--force-position-strength', '0.05'));
        this.forceCollisionStrength = parseFloat(this.getCssVariable('--force-collision-strength', '0.8'));
        this.forceClusterStrength = parseFloat(this.getCssVariable('--force-cluster-strength', '0.3'));
        
        // Animation parameters
        this.cameraAnimationDuration = parseFloat(this.getCssVariable('--camera-animation-duration', '1000'));
        this.alphaDecay = parseFloat(this.getCssVariable('--alpha-decay', '0.02'));
        this.velocityDecay = parseFloat(this.getCssVariable('--velocity-decay', '0.3'));
    }
    
    /**
     * Get CSS variable value from :root
     * @param {string} name - CSS variable name
     * @param {string} fallback - Fallback value if variable is not defined
     * @returns {string} CSS variable value or fallback
     */
    getCssVariable(name, fallback) {
        const value = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
        return value || fallback;
    }
    
    /**
     * Set the graph data
     * @param {Array} nodes - Array of node objects
     * @param {Array} links - Array of link objects
     */
    setData(nodes, links) {
        this.nodes = nodes;
        this.links = links;
        
        // Process links to ensure they have proper source/target references
        this.links = this.links.map(link => {
            // If source/target are strings (IDs), convert to references to the actual node objects
            if (typeof link.source === 'string' || typeof link.source === 'number') {
                const sourceNode = this.nodes.find(node => node.id === link.source);
                const targetNode = this.nodes.find(node => node.id === link.target);
                
                if (sourceNode && targetNode) {
                    return {
                        ...link,
                        source: sourceNode,
                        target: targetNode
                    };
                }
            }
            return link;
        });
        
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
        // Remove all node objects
        this.nodeObjects.forEach(nodeObj => {
            this.scene.remove(nodeObj);
        });
        
        // Remove all link objects
        this.linkObjects.forEach(linkObj => {
            this.scene.remove(linkObj);
        });
        
        // Clear maps
        this.nodeObjects.clear();
        this.linkObjects.clear();
    }
    
    /**
     * Create node objects
     */
    createNodeObjects() {
        // Create node objects
        this.nodes.forEach(node => {
            // Calculate color based on node properties or connections
            const connections = this.links.filter(link => 
                (link.source.id === node.id || link.target.id === node.id)
            ).length;
            
            // Use CSS variables for node color
            let color;
            if (node.color) {
                // Use node's color if specified
                color = new THREE.Color(node.color);
            } else {
                // Use base color from CSS
                color = new THREE.Color(this.nodeBaseColor);
            }
            
            // Create geometry with appropriate size
            const size = node.size || Math.max(
                this.nodeSizeMin, 
                Math.min(this.nodeSizeMax, this.nodeSizeMin + connections * this.nodeSizeMultiplier)
            );
            
            const geometry = new THREE.SphereGeometry(size, 32, 32);
            
            // Create material with glow effect using CSS variables
            const material = new THREE.MeshPhongMaterial({ 
                color: color,
                emissive: color.clone().multiplyScalar(0.3),
                specular: 0xffffff,
                shininess: this.nodeShininess
            });
            
            const mesh = new THREE.Mesh(geometry, material);
            
            // Set position
            if (node.x !== undefined && node.y !== undefined && node.z !== undefined) {
                mesh.position.set(node.x, node.y, node.z);
            } else {
                // Initial position in a sphere formation
                const phi = Math.acos(-1 + (2 * Math.random()));
                const theta = Math.random() * Math.PI * 2;
                const radius = 100 + Math.random() * 50;
                
                mesh.position.set(
                    radius * Math.sin(phi) * Math.cos(theta),
                    radius * Math.sin(phi) * Math.sin(theta),
                    radius * Math.cos(phi)
                );
                
                // Update node position
                node.x = mesh.position.x;
                node.y = mesh.position.y;
                node.z = mesh.position.z;
            }
            
            // Store reference to node data
            mesh.userData = node;
            
            // Add to scene
            this.scene.add(mesh);
            
            // Store in map
            this.nodeObjects.set(node.id, mesh);
        });
    }
    
    /**
     * Create link objects
     */
    createLinkObjects() {
        // Create link objects
        this.links.forEach(link => {
            const sourceId = typeof link.source === 'object' ? link.source.id : link.source;
            const targetId = typeof link.target === 'object' ? link.target.id : link.target;
            
            const sourceNode = this.nodeObjects.get(sourceId);
            const targetNode = this.nodeObjects.get(targetId);
            
            if (sourceNode && targetNode) {
                // Create curved line for better visualization
                const curvePoints = this.createCurvedLine(
                    sourceNode.position.clone(),
                    targetNode.position.clone()
                );
                
                // Create line geometry
                const geometry = new THREE.BufferGeometry().setFromPoints(curvePoints);
                
                // Create line material with custom color based on relationship type
                let color;
                if (link.color) {
                    // Use link's color if specified
                    color = new THREE.Color(link.color);
                } else if (link.type) {
                    // Generate consistent color based on relationship type
                    const hash = this.hashString(link.type);
                    color = new THREE.Color(`hsl(${hash % 360}, 70%, 70%)`);
                } else {
                    // Use base color from CSS
                    color = new THREE.Color(this.edgeBaseColor);
                }
                
                const material = new THREE.LineBasicMaterial({ 
                    color: color,
                    transparent: true,
                    opacity: this.edgeOpacity,
                    linewidth: this.edgeWidth
                });
                
                const line = new THREE.Line(geometry, material);
                
                // Store reference to link data
                line.userData = link;
                
                // Add to scene
                this.scene.add(line);
                
                // Store in map
                const linkId = `${sourceId}-${targetId}`;
                this.linkObjects.set(linkId, line);
            }
        });
    }
    
    /**
     * Create curved line between two points
     * @param {THREE.Vector3} start - Start position
     * @param {THREE.Vector3} end - End position
     * @returns {Array} Array of points
     */
    createCurvedLine(start, end) {
        // Calculate midpoint
        const mid = start.clone().add(end).multiplyScalar(0.5);
        
        // Add slight curve
        const direction = end.clone().sub(start);
        const perpendicular = new THREE.Vector3(-direction.y, direction.x, direction.z).normalize();
        const distance = start.distanceTo(end);
        
        // Curve more for longer distances
        const curveAmount = Math.min(distance * 0.1, 10);
        mid.add(perpendicular.multiplyScalar(curveAmount));
        
        // Create curve
        const curve = new THREE.QuadraticBezierCurve3(start, mid, end);
        
        // Get points along curve
        return curve.getPoints(10);
    }
    
    /**
     * Simple string hash function
     * @param {string} str - String to hash
     * @returns {number} Hash value
     */
    hashString(str) {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            hash = ((hash << 5) - hash) + str.charCodeAt(i);
            hash |= 0; // Convert to 32bit integer
        }
        return Math.abs(hash);
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
     * Handle mouse move event
     * @param {Event} event - Mouse event
     */
    onMouseMove(event) {
        // Calculate mouse position in normalized device coordinates
        const rect = this.renderer.domElement.getBoundingClientRect();
        this.mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
        this.mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
        
        // Update the picking ray with the camera and mouse position
        this.raycaster.setFromCamera(this.mouse, this.camera);
        
        // Calculate objects intersecting the picking ray
        const intersects = this.raycaster.intersectObjects(Array.from(this.nodeObjects.values()));
        
        // Reset cursor and node highlighting
        this.renderer.domElement.style.cursor = 'auto';
        this.nodeObjects.forEach(node => {
            if (node !== this.selectedNode) {
                node.material.emissiveIntensity = this.nodeEmissiveIntensityNormal;
                // Reset color if not the selected node
                node.material.color = new THREE.Color(this.nodeBaseColor);
                node.material.emissive = node.material.color.clone().multiplyScalar(0.3);
            }
        });
        
        if (intersects.length > 0) {
            // Change cursor to pointer
            this.renderer.domElement.style.cursor = 'pointer';
            
            // Highlight hovered node
            const hoveredNode = intersects[0].object;
            
            // Only change emissive intensity if not the selected node
            if (hoveredNode !== this.selectedNode) {
                hoveredNode.material.emissiveIntensity = this.nodeEmissiveIntensityHover;
            }
        }
    }
    
    /**
     * Handle click event
     * @param {Event} event - Mouse event
     */
    onClick(event) {
        // Calculate mouse position in normalized device coordinates
        const rect = this.renderer.domElement.getBoundingClientRect();
        this.mouse.x = ((event.clientX - rect.left) / rect.width) * 2 - 1;
        this.mouse.y = -((event.clientY - rect.top) / rect.height) * 2 + 1;
        
        // Update the picking ray with the camera and mouse position
        this.raycaster.setFromCamera(this.mouse, this.camera);
        
        // Calculate objects intersecting the picking ray
        const intersects = this.raycaster.intersectObjects(Array.from(this.nodeObjects.values()));
        
        // Reset all edges to base color
        this.resetEdgeColors();
        
        if (intersects.length > 0) {
            const selectedObject = intersects[0].object;
            const selectedNodeId = selectedObject.userData.id;
            
            // Reset previous selection
            if (this.selectedNode && this.selectedNode !== selectedObject) {
                // Reset node color and emissive intensity
                this.selectedNode.material.emissiveIntensity = this.nodeEmissiveIntensityNormal;
                this.selectedNode.material.color = new THREE.Color(this.nodeBaseColor);
                this.selectedNode.material.emissive = this.selectedNode.material.color.clone().multiplyScalar(0.3);
            }
            
            this.selectedNode = selectedObject;
            
            // Set selected node color and emissive intensity
            this.selectedNode.material.color = new THREE.Color(this.nodeSelectedColor);
            this.selectedNode.material.emissive = this.selectedNode.material.color.clone().multiplyScalar(0.3);
            this.selectedNode.material.emissiveIntensity = this.nodeEmissiveIntensitySelected;
            
            // Highlight connected edges
            this.highlightConnectedEdges(selectedNodeId);
            
            // Show node info
            this.showNodeInfo(selectedObject.userData);
            
            // Focus camera on selected node
            this.focusOnNode(selectedObject);
        } else {
            // Reset selection
            if (this.selectedNode) {
                // Reset node color and emissive intensity
                this.selectedNode.material.emissiveIntensity = this.nodeEmissiveIntensityNormal;
                this.selectedNode.material.color = new THREE.Color(this.nodeBaseColor);
                this.selectedNode.material.emissive = this.selectedNode.material.color.clone().multiplyScalar(0.3);
                this.selectedNode = null;
            }
            
            // Hide node info
            this.nodeInfoPanel.style.display = 'none';
        }
    }
    
    /**
     * Reset all edge colors to base color
     */
    resetEdgeColors() {
        this.linkObjects.forEach(linkObj => {
            linkObj.material.color = new THREE.Color(this.edgeBaseColor);
        });
    }
    
    /**
     * Highlight edges connected to the selected node
     * @param {string} nodeId - ID of the selected node
     */
    highlightConnectedEdges(nodeId) {
        this.links.forEach(link => {
            const sourceId = typeof link.source === 'object' ? link.source.id : link.source;
            const targetId = typeof link.target === 'object' ? link.target.id : link.target;
            
            if (sourceId === nodeId || targetId === nodeId) {
                const linkId = `${sourceId}-${targetId}`;
                const linkObj = this.linkObjects.get(linkId);
                
                if (linkObj) {
                    linkObj.material.color = new THREE.Color(this.edgeSelectedColor);
                }
            }
        });
    }
    
    /**
     * Focus camera on a node
     * @param {THREE.Mesh} node - Node to focus on
     */
    focusOnNode(node) {
        const position = node.position.clone();
        const distance = this.camera.position.distanceTo(position);
        const targetDistance = Math.max(100, distance * 0.8);
        
        // Animate camera movement
        const startPosition = this.camera.position.clone();
        const endPosition = position.clone().add(
            new THREE.Vector3(0, 0, targetDistance)
        );
        
        const duration = this.cameraAnimationDuration; // Use CSS variable
        const startTime = Date.now();
        
        const animate = () => {
            const elapsed = Date.now() - startTime;
            const progress = Math.min(elapsed / duration, 1);
            
            // Ease in-out function
            const easeProgress = progress < 0.5 
                ? 2 * progress * progress 
                : -1 + (4 - 2 * progress) * progress;
            
            // Interpolate position
            const newPosition = startPosition.clone().lerp(endPosition, easeProgress);
            this.camera.position.copy(newPosition);
            
            // Look at node
            this.camera.lookAt(position);
            this.controls.target.copy(position);
            
            if (progress < 1) {
                requestAnimationFrame(animate);
            }
        };
        
        animate();
    }
    
    /**
     * Show node info
     * @param {Object} node - Node data
     */
    showNodeInfo(node) {
        // Set node info
        this.nodeInfoTitle.textContent = node.name;
        
        // Set node content
        let content = '';
        for (const key in node) {
            if (key !== 'name' && key !== 'id' && key !== 'x' && key !== 'y' && key !== 'z' && key !== 'size') {
                content += `<p><strong>${key}:</strong> ${node[key]}</p>`;
            }
        }
        
        this.nodeInfoContent.innerHTML = content || '<p>No additional information available.</p>';
        
        // Set node connections
        this.nodeConnections.innerHTML = '';
        
        // Find connections
        const connections = this.links.filter(link => {
            const sourceId = typeof link.source === 'object' ? link.source.id : link.source;
            const targetId = typeof link.target === 'object' ? link.target.id : link.target;
            return sourceId === node.id || targetId === node.id;
        });
        
        if (connections.length > 0) {
            connections.forEach(connection => {
                const sourceId = typeof connection.source === 'object' ? connection.source.id : connection.source;
                const targetId = typeof connection.target === 'object' ? connection.target.id : connection.target;
                
                const isSource = sourceId === node.id;
                const otherNodeId = isSource ? targetId : sourceId;
                const otherNode = this.nodes.find(n => n.id === otherNodeId);
                
                if (otherNode) {
                    const direction = isSource ? 'outgoing' : 'incoming';
                    const connectionItem = document.createElement('div');
                    connectionItem.className = 'connection-item';
                    connectionItem.innerHTML = `
                        <span class="connection-direction ${direction}">${isSource ? '→' : '←'}</span>
                        <span class="connection-type">${connection.type || 'RELATED_TO'}</span>
                        <span class="connection-node">${otherNode.name}</span>
                    `;
                    
                    // Add click event to focus on connected node
                    connectionItem.addEventListener('click', () => {
                        const connectedNodeObj = this.nodeObjects.get(otherNodeId);
                        if (connectedNodeObj) {
                            // Reset all edges to base color
                            this.resetEdgeColors();
                            
                            // Reset previous selection
                            if (this.selectedNode) {
                                // Reset node color and emissive intensity
                                this.selectedNode.material.emissiveIntensity = this.nodeEmissiveIntensityNormal;
                                this.selectedNode.material.color = new THREE.Color(this.nodeBaseColor);
                                this.selectedNode.material.emissive = this.selectedNode.material.color.clone().multiplyScalar(0.3);
                            }
                            
                            this.selectedNode = connectedNodeObj;
                            
                            // Set selected node color and emissive intensity
                            this.selectedNode.material.color = new THREE.Color(this.nodeSelectedColor);
                            this.selectedNode.material.emissive = this.selectedNode.material.color.clone().multiplyScalar(0.3);
                            this.selectedNode.material.emissiveIntensity = this.nodeEmissiveIntensitySelected;
                            
                            // Highlight connected edges
                            this.highlightConnectedEdges(otherNodeId);
                            
                            // Show node info
                            this.showNodeInfo(connectedNodeObj.userData);
                            
                            // Focus camera on selected node
                            this.focusOnNode(connectedNodeObj);
                        }
                    });
                    
                    this.nodeConnections.appendChild(connectionItem);
                }
            });
        } else {
            this.nodeConnections.innerHTML = '<p>No connections.</p>';
        }
        
        // Show panel
        this.nodeInfoPanel.style.display = 'block';
    }
    
    /**
     * Handle window resize event
     */
    onWindowResize() {
        // Update camera
        this.camera.aspect = this.container.clientWidth / this.container.clientHeight;
        this.camera.updateProjectionMatrix();
        
        // Update renderer
        this.renderer.setSize(this.container.clientWidth, this.container.clientHeight);
    }
    
    /**
     * Start force simulation
     */
    startSimulation() {
        // Create simulation with improved parameters
        this.simulation = d3.forceSimulation(this.nodes)
            // Repulsion between nodes (stronger and with a minimum distance)
            .force('charge', d3.forceManyBody()
                .strength(this.forceChargeStrength)
                .distanceMin(20)
                .distanceMax(300)
            )
            // Links with variable distance based on node size
            .force('link', d3.forceLink(this.links)
                .id(d => d.id)
                .distance(link => {
                    const sourceSize = link.source.size || this.nodeSizeMin;
                    const targetSize = link.target.size || this.nodeSizeMin;
                    return 50 + sourceSize + targetSize;
                })
                .strength(this.forceLinkStrength)
            )
            // Center force
            .force('center', d3.forceCenter(0, 0, 0).strength(this.forceCenterStrength))
            // 3D positioning forces
            .force('x', d3.forceX().strength(this.forcePositionStrength))
            .force('y', d3.forceY().strength(this.forcePositionStrength))
            .force('z', d3.forceZ().strength(this.forcePositionStrength))
            // Collision detection to prevent overlap
            .force('collision', d3.forceCollide().radius(d => (d.size || this.nodeSizeMin) + 2).strength(this.forceCollisionStrength))
            // Cluster by relationship type
            .force('cluster', this.forceCluster())
            .alphaDecay(this.alphaDecay) // Use CSS variable
            .velocityDecay(this.velocityDecay) // Use CSS variable
            .on('tick', this.updatePositions.bind(this));
            
        // Run simulation for a few ticks to get a better initial layout
        for (let i = 0; i < 100; i++) {
            this.simulation.tick();
        }
    }
    
    /**
     * Custom force to cluster nodes by relationship type
     */
    forceCluster() {
        // Group nodes by their most common relationship type
        const nodeGroups = new Map();
        
        this.links.forEach(link => {
            const sourceId = typeof link.source === 'object' ? link.source.id : link.source;
            const targetId = typeof link.target === 'object' ? link.target.id : link.target;
            const type = link.type || 'RELATED_TO';
            
            // Add to source node's types
            if (!nodeGroups.has(sourceId)) {
                nodeGroups.set(sourceId, new Map());
            }
            const sourceTypes = nodeGroups.get(sourceId);
            sourceTypes.set(type, (sourceTypes.get(type) || 0) + 1);
            
            // Add to target node's types
            if (!nodeGroups.has(targetId)) {
                nodeGroups.set(targetId, new Map());
            }
            const targetTypes = nodeGroups.get(targetId);
            targetTypes.set(type, (targetTypes.get(type) || 0) + 1);
        });
        
        // Assign each node to its most common group
        const nodeToGroup = new Map();
        nodeGroups.forEach((types, nodeId) => {
            let maxCount = 0;
            let maxType = null;
            
            types.forEach((count, type) => {
                if (count > maxCount) {
                    maxCount = count;
                    maxType = type;
                }
            });
            
            nodeToGroup.set(nodeId, maxType);
        });
        
        // Calculate group centers
        const groupCenters = new Map();
        const groupCounts = new Map();
        
        return function(alpha) {
            // Update group centers
            groupCenters.clear();
            groupCounts.clear();
            
            // Calculate current group centers
            this.nodes.forEach(node => {
                const group = nodeToGroup.get(node.id) || 'default';
                
                if (!groupCenters.has(group)) {
                    groupCenters.set(group, { x: 0, y: 0, z: 0 });
                    groupCounts.set(group, 0);
                }
                
                const center = groupCenters.get(group);
                center.x += node.x || 0;
                center.y += node.y || 0;
                center.z += node.z || 0;
                groupCounts.set(group, groupCounts.get(group) + 1);
            });
            
            // Normalize centers
            groupCenters.forEach((center, group) => {
                const count = groupCounts.get(group);
                if (count > 0) {
                    center.x /= count;
                    center.y /= count;
                    center.z /= count;
                }
            });
            
            // Apply clustering force
            const clusterStrength = this.forceClusterStrength * alpha;
            
            this.nodes.forEach(node => {
                const group = nodeToGroup.get(node.id) || 'default';
                const center = groupCenters.get(group);
                
                if (center) {
                    node.vx = (node.vx || 0) + (center.x - node.x) * clusterStrength;
                    node.vy = (node.vy || 0) + (center.y - node.y) * clusterStrength;
                    node.vz = (node.vz || 0) + (center.z - node.z) * clusterStrength;
                }
            });
        }.bind(this);
    }
    
    /**
     * Update positions based on simulation
     */
    updatePositions() {
        // Update node positions
        this.nodes.forEach(node => {
            const nodeObj = this.nodeObjects.get(node.id);
            if (nodeObj) {
                nodeObj.position.set(node.x, node.y, node.z);
            }
        });
        
        // Update link positions
        this.links.forEach(link => {
            const sourceId = typeof link.source === 'object' ? link.source.id : link.source;
            const targetId = typeof link.target === 'object' ? link.target.id : link.target;
            
            const sourceNode = this.nodeObjects.get(sourceId);
            const targetNode = this.nodeObjects.get(targetId);
            
            if (sourceNode && targetNode) {
                const linkId = `${sourceId}-${targetId}`;
                const linkObj = this.linkObjects.get(linkId);
                
                if (linkObj) {
                    // Update curved line
                    const curvePoints = this.createCurvedLine(
                        sourceNode.position.clone(),
                        targetNode.position.clone()
                    );
                    
                    // Update geometry
                    linkObj.geometry.dispose();
                    linkObj.geometry = new THREE.BufferGeometry().setFromPoints(curvePoints);
                }
            }
        });
    }
    
    /**
     * Reset camera position
     */
    resetCamera() {
        // Animate camera movement
        const startPosition = this.camera.position.clone();
        const endPosition = new THREE.Vector3(0, 0, 300);
        
        const duration = this.cameraAnimationDuration; // Use CSS variable
        const startTime = Date.now();
        
        const animate = () => {
            const elapsed = Date.now() - startTime;
            const progress = Math.min(elapsed / duration, 1);
            
            // Ease in-out function
            const easeProgress = progress < 0.5 
                ? 2 * progress * progress 
                : -1 + (4 - 2 * progress) * progress;
            
            // Interpolate position
            const newPosition = startPosition.clone().lerp(endPosition, easeProgress);
            this.camera.position.copy(newPosition);
            
            // Reset target
            this.controls.target.set(0, 0, 0);
            
            if (progress < 1) {
                requestAnimationFrame(animate);
            }
        };
        
        animate();
    }
    
    /**
     * Render the scene
     */
    render() {
        this.controls.update();
        this.renderer.render(this.scene, this.camera);
    }
}

// Export for use in tests
if (typeof module !== 'undefined') {
    module.exports = GraphVisualizer;
} 