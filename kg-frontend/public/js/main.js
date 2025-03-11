/**
 * Main application script
 */
document.addEventListener('DOMContentLoaded', () => {
    console.log("DOM loaded, initializing application");
    
    // Initialize API client
    const apiClient = new ApiClient();
    
    // Initialize graph visualizer
    console.log("Creating graph visualizer");
    const graphVisualizer = new GraphVisualizer('graph-canvas');
    
    // Add D3.js for force simulation
    console.log("Loading D3.js");
    Promise.all([
        loadScript('https://d3js.org/d3.v7.min.js'),
        loadScript('https://unpkg.com/d3-force-3d')
    ])
        .then(() => {
            console.log("D3.js and D3-force-3d loaded successfully");
            // Load initial graph data
            refreshGraph();
            // Load initial statistics
            loadStatistics();
            // Create space particles
            createSpaceParticles();
        })
        .catch(error => {
            console.error('Error loading D3.js or D3-force-3d:', error);
            alert('Failed to load required libraries. The graph visualization may not work properly.');
        });
    
    // Create space particles
    function createSpaceParticles() {
        const container = document.querySelector('.graph-container');
        for (let i = 0; i < 100; i++) {
            const particle = document.createElement('div');
            particle.className = 'space-particle';
            
            // Random size between 1 and 3 pixels
            const size = Math.random() * 2 + 1;
            particle.style.width = `${size}px`;
            particle.style.height = `${size}px`;
            
            // Random position
            const posX = Math.random() * 100;
            const posY = Math.random() * 100;
            particle.style.left = `${posX}%`;
            particle.style.top = `${posY}%`;
            
            // Random opacity
            particle.style.opacity = Math.random() * 0.7 + 0.3;
            
            // Random color (mostly white/blue)
            const hue = Math.random() * 60 + 180; // Blue to purple range
            const saturation = Math.random() * 50 + 50;
            const lightness = Math.random() * 30 + 70;
            particle.style.backgroundColor = `hsl(${hue}, ${saturation}%, ${lightness}%)`;
            
            // Add to container
            container.appendChild(particle);
            
            // Animate particle
            animateParticle(particle);
        }
    }
    
    // Animate space particle
    function animateParticle(particle) {
        // Random duration between 20 and 60 seconds
        const duration = Math.random() * 40000 + 20000;
        
        // Random movement
        const startX = parseFloat(particle.style.left);
        const startY = parseFloat(particle.style.top);
        const endX = Math.random() * 100;
        const endY = Math.random() * 100;
        
        // Start time
        const startTime = Date.now();
        
        // Animation function
        function update() {
            const elapsed = Date.now() - startTime;
            const progress = elapsed / duration;
            
            if (progress < 1) {
                const x = startX + (endX - startX) * progress;
                const y = startY + (endY - startY) * progress;
                
                particle.style.left = `${x}%`;
                particle.style.top = `${y}%`;
                
                // Twinkle effect
                particle.style.opacity = 0.3 + Math.sin(elapsed / 1000) * 0.3;
                
                requestAnimationFrame(update);
            } else {
                // Reset animation
                particle.style.left = `${Math.random() * 100}%`;
                particle.style.top = `${Math.random() * 100}%`;
                animateParticle(particle);
            }
        }
        
        // Start animation
        update();
    }
    
    // Function to load graph data
    function refreshGraph() {
        console.log("Refreshing graph data");
        
        // Show loading indicator
        const loadingIndicator = document.createElement('div');
        loadingIndicator.className = 'loading';
        loadingIndicator.textContent = 'Loading graph data...';
        document.querySelector('.graph-container').appendChild(loadingIndicator);
        
        apiClient.getGraphData()
            .then(data => {
                console.log(`Graph data loaded: ${data.nodes.length} nodes, ${data.links.length} links`);
                
                // Process nodes to ensure they have size property and 3D positions
                data.nodes.forEach(node => {
                    // Vary node size based on connections
                    const connections = data.links.filter(link => 
                        link.source === node.id || link.target === node.id
                    ).length;
                    
                    // Scale size based on connections (min 3, max 10)
                    node.size = Math.max(3, Math.min(10, 3 + connections * 0.5));
                    
                    // Create initial positions in a sphere formation rather than a cube
                    // This helps prevent the "flying through" effect
                    const phi = Math.acos(-1 + (2 * Math.random()));
                    const theta = Math.random() * Math.PI * 2;
                    const radius = 150 + Math.random() * 50; // Radius between 150-200
                    
                    node.x = radius * Math.sin(phi) * Math.cos(theta);
                    node.y = radius * Math.sin(phi) * Math.sin(theta);
                    node.z = radius * Math.cos(phi);
                });
                
                // Set data to visualizer
                graphVisualizer.setData(data.nodes, data.links);
                
                // Start force simulation
                graphVisualizer.startSimulation();
                
                // Remove loading indicator
                loadingIndicator.remove();
            })
            .catch(error => {
                console.error('Error loading graph data:', error);
                alert('Failed to load graph data. Please check the console for details.');
                
                // Remove loading indicator
                loadingIndicator.remove();
            });
    }
    
    // Function to load statistics
    function loadStatistics() {
        console.log("Loading statistics");
        
        // Get node count
        const nodeCountQuery = `
            MATCH (n:Concept) 
            RETURN count(n) as count
        `;
        
        // Get relationship count
        const relationshipCountQuery = `
            MATCH ()-[r:RELATED_TO]->() 
            RETURN count(r) as count
        `;
        
        // Get relationship types
        const relationshipTypesQuery = `
            MATCH ()-[r:RELATED_TO]->() 
            RETURN r.type as type, count(r) as count
            ORDER BY count DESC
            LIMIT 10
        `;
        
        // Execute node count query
        fetch(apiClient.neo4jUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Basic ${apiClient.neo4jAuth}`,
                'Accept': 'application/json; charset=UTF-8'
            },
            body: JSON.stringify({
                statements: [{ statement: nodeCountQuery }]
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.results && data.results.length > 0 && data.results[0].data && data.results[0].data.length > 0) {
                const nodeCount = data.results[0].data[0].row[0];
                document.getElementById('concept-count').textContent = nodeCount;
            }
        })
        .catch(error => {
            console.error('Error fetching node count:', error);
        });
        
        // Execute relationship count query
        fetch(apiClient.neo4jUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Basic ${apiClient.neo4jAuth}`,
                'Accept': 'application/json; charset=UTF-8'
            },
            body: JSON.stringify({
                statements: [{ statement: relationshipCountQuery }]
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.results && data.results.length > 0 && data.results[0].data && data.results[0].data.length > 0) {
                const relationshipCount = data.results[0].data[0].row[0];
                document.getElementById('relationship-count').textContent = relationshipCount;
            }
        })
        .catch(error => {
            console.error('Error fetching relationship count:', error);
        });
        
        // Execute relationship types query
        fetch(apiClient.neo4jUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Basic ${apiClient.neo4jAuth}`,
                'Accept': 'application/json; charset=UTF-8'
            },
            body: JSON.stringify({
                statements: [{ statement: relationshipTypesQuery }]
            })
        })
        .then(response => response.json())
        .then(data => {
            if (data.results && data.results.length > 0 && data.results[0].data) {
                const relationshipTypes = data.results[0].data.map(row => ({
                    type: row.row[0],
                    count: row.row[1]
                }));
                
                const relationshipTypesContainer = document.getElementById('relationship-types');
                relationshipTypesContainer.innerHTML = '';
                
                relationshipTypes.forEach(type => {
                    const typeElement = document.createElement('div');
                    typeElement.className = 'type-item';
                    typeElement.innerHTML = `
                        <span class="type-name">${type.type}</span>
                        <span class="type-count">${type.count}</span>
                    `;
                    relationshipTypesContainer.appendChild(typeElement);
                });
            }
        })
        .catch(error => {
            console.error('Error fetching relationship types:', error);
        });
    }
    
    // Get DOM elements
    const startBuilderBtn = document.getElementById('start-builder');
    const stopBuilderBtn = document.getElementById('stop-builder');
    const startEnricherBtn = document.getElementById('start-enricher');
    const stopEnricherBtn = document.getElementById('stop-enricher');
    const resetCameraBtn = document.getElementById('reset-camera');
    const refreshGraphBtn = document.getElementById('refresh-graph');
    const conceptSearch = document.getElementById('concept-search');
    const searchResults = document.getElementById('search-results');
    const selectedConcept1 = document.getElementById('selected-concept-1');
    const selectedConcept2 = document.getElementById('selected-concept-2');
    const linkConceptsBtn = document.getElementById('link-concepts');
    
    // Selected concepts for manual linking
    let selectedConcepts = [];
    
    // Event listeners for builder controls
    startBuilderBtn.addEventListener('click', startBuilder);
    stopBuilderBtn.addEventListener('click', stopBuilder);
    
    // Event listeners for enricher controls
    startEnricherBtn.addEventListener('click', startEnricher);
    stopEnricherBtn.addEventListener('click', stopEnricher);
    
    // Event listeners for graph controls
    resetCameraBtn.addEventListener('click', () => {
        graphVisualizer.resetCamera();
        console.log("Camera reset");
    });
    
    refreshGraphBtn.addEventListener('click', () => {
        refreshGraph();
        loadStatistics();
        console.log("Graph and statistics refreshed");
    });
    
    // Event listener for concept search
    conceptSearch.addEventListener('input', debounce(searchConcepts, 300));
    
    // Event listener for linking concepts
    linkConceptsBtn.addEventListener('click', linkConcepts);
    
    // Function to start the builder
    function startBuilder() {
        const seedConcept = document.getElementById('seed-concept').value;
        const maxNodes = parseInt(document.getElementById('max-nodes').value);
        const timeout = parseInt(document.getElementById('timeout').value);
        const randomRelationships = parseInt(document.getElementById('random-relationships').value);
        const concurrency = parseInt(document.getElementById('concurrency').value);
        
        const params = {
            seedConcept: seedConcept,
            maxNodes: maxNodes,
            timeout: timeout,
            randomRelationships: randomRelationships,
            concurrency: concurrency
        };
        
        // Disable buttons during operation
        const startBtn = document.getElementById('start-builder');
        startBtn.disabled = true;
        startBtn.textContent = 'Starting...';
        
        console.log("Starting builder with params:", params);
        apiClient.startBuilder(params)
            .then(response => {
                console.log('Builder started:', response);
                if (response.status === 'success') {
                    alert(response.message || 'Knowledge Graph Builder started successfully.');
                } else {
                    alert('Error starting Knowledge Graph Builder: ' + (response.error || 'Unknown error'));
                }
            })
            .catch(error => {
                console.error('Error starting builder:', error);
                alert('Error starting Knowledge Graph Builder: ' + error.message);
            })
            .finally(() => {
                // Re-enable button
                startBtn.disabled = false;
                startBtn.textContent = 'Start Builder';
            });
    }
    
    // Function to stop the builder
    function stopBuilder() {
        // Disable buttons during operation
        const stopBtn = document.getElementById('stop-builder');
        stopBtn.disabled = true;
        stopBtn.textContent = 'Stopping...';
        
        console.log("Stopping builder");
        apiClient.stopBuilder()
            .then(response => {
                console.log('Builder stopped:', response);
                if (response.status === 'success') {
                    alert(response.message || 'Knowledge Graph Builder stopped successfully.');
                    refreshGraph();
                    loadStatistics();
                } else {
                    alert('Error stopping Knowledge Graph Builder: ' + (response.error || 'Unknown error'));
                }
            })
            .catch(error => {
                console.error('Error stopping builder:', error);
                alert('Error stopping Knowledge Graph Builder: ' + error.message);
            })
            .finally(() => {
                // Re-enable button
                stopBtn.disabled = false;
                stopBtn.textContent = 'Stop Builder';
            });
    }
    
    // Function to start the enricher
    function startEnricher() {
        const batchSize = parseInt(document.getElementById('batch-size').value);
        const interval = parseInt(document.getElementById('interval').value);
        const maxRelationships = parseInt(document.getElementById('max-relationships').value);
        const concurrency = parseInt(document.getElementById('enricher-concurrency').value);
        
        const params = {
            batchSize: batchSize,
            interval: interval,
            maxRelationships: maxRelationships,
            concurrency: concurrency
        };
        
        // Disable buttons during operation
        const startBtn = document.getElementById('start-enricher');
        startBtn.disabled = true;
        startBtn.textContent = 'Starting...';
        
        console.log("Starting enricher with params:", params);
        apiClient.startEnricher(params)
            .then(response => {
                console.log('Enricher started:', response);
                if (response.status === 'success') {
                    alert(response.message || 'Knowledge Graph Enricher started successfully.');
                } else {
                    alert('Error starting Knowledge Graph Enricher: ' + (response.error || 'Unknown error'));
                }
            })
            .catch(error => {
                console.error('Error starting enricher:', error);
                alert('Error starting Knowledge Graph Enricher: ' + error.message);
            })
            .finally(() => {
                // Re-enable button
                startBtn.disabled = false;
                startBtn.textContent = 'Start Enricher';
            });
    }
    
    // Function to stop the enricher
    function stopEnricher() {
        // Disable buttons during operation
        const stopBtn = document.getElementById('stop-enricher');
        stopBtn.disabled = true;
        stopBtn.textContent = 'Stopping...';
        
        console.log("Stopping enricher");
        apiClient.stopEnricher()
            .then(response => {
                console.log('Enricher stopped:', response);
                if (response.status === 'success') {
                    alert(response.message || 'Knowledge Graph Enricher stopped successfully.');
                    refreshGraph();
                    loadStatistics();
                } else {
                    alert('Error stopping Knowledge Graph Enricher: ' + (response.error || 'Unknown error'));
                }
            })
            .catch(error => {
                console.error('Error stopping enricher:', error);
                alert('Error stopping Knowledge Graph Enricher: ' + error.message);
            })
            .finally(() => {
                // Re-enable button
                stopBtn.disabled = false;
                stopBtn.textContent = 'Stop Enricher';
            });
    }
    
    // Function for concept search
    function searchConcepts() {
        const query = conceptSearch.value.trim();
        
        if (query.length < 2) {
            searchResults.innerHTML = '';
            return;
        }
        
        console.log(`Searching for concepts: "${query}"`);
        apiClient.searchConcepts(query)
            .then(concepts => {
                console.log(`Found ${concepts.length} concepts matching "${query}"`);
                searchResults.innerHTML = '';
                
                concepts.forEach(concept => {
                    const resultItem = document.createElement('div');
                    resultItem.className = 'search-result-item';
                    resultItem.textContent = concept.name;
                    resultItem.addEventListener('click', () => {
                        if (selectedConcepts.length < 2) {
                            selectedConcepts.push(concept);
                            updateSelectedConcepts();
                        }
                        searchResults.innerHTML = '';
                        conceptSearch.value = '';
                    });
                    
                    searchResults.appendChild(resultItem);
                });
                
                // Show search results
                if (concepts.length > 0) {
                    searchResults.style.display = 'block';
                } else {
                    searchResults.style.display = 'none';
                }
            })
            .catch(error => {
                console.error('Error searching concepts:', error);
                alert('Error searching concepts. Check the console for details.');
            });
    }
    
    // Function to update selected concepts display
    function updateSelectedConcepts() {
        console.log("Updating selected concepts:", selectedConcepts);
        if (selectedConcepts.length > 0) {
            selectedConcept1.textContent = selectedConcepts[0].name;
        } else {
            selectedConcept1.textContent = 'None';
        }
        
        if (selectedConcepts.length > 1) {
            selectedConcept2.textContent = selectedConcepts[1].name;
        } else {
            selectedConcept2.textContent = 'None';
        }
    }
    
    // Function to link concepts
    function linkConcepts() {
        if (selectedConcepts.length !== 2) {
            alert('Please select two concepts to link.');
            return;
        }
        
        const relationshipType = document.getElementById('relationship-type').value.trim();
        
        if (!relationshipType) {
            alert('Please enter a relationship type.');
            return;
        }
        
        console.log(`Creating relationship: ${selectedConcepts[0].name} -[${relationshipType}]-> ${selectedConcepts[1].name}`);
        apiClient.createRelationship(selectedConcepts[0].id, selectedConcepts[1].id, relationshipType)
            .then(response => {
                console.log('Relationship created:', response);
                alert('Relationship created successfully.');
                selectedConcepts = [];
                updateSelectedConcepts();
                refreshGraph();
                loadStatistics();
            })
            .catch(error => {
                console.error('Error creating relationship:', error);
                alert('Error creating relationship. Check the console for details.');
            });
    }
    
    // Helper function to load a script dynamically
    function loadScript(src) {
        return new Promise((resolve, reject) => {
            const script = document.createElement('script');
            script.src = src;
            script.onload = resolve;
            script.onerror = reject;
            document.head.appendChild(script);
        });
    }
    
    // Helper function to debounce function calls
    function debounce(func, wait) {
        let timeout;
        return function() {
            const context = this;
            const args = arguments;
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                func.apply(context, args);
            }, wait);
        };
    }
    
    // Start the animation loop
    function animate() {
        requestAnimationFrame(animate);
        graphVisualizer.render();
    }
    
    // Start the animation loop
    animate();
    
    console.log("Application initialized");
}); 