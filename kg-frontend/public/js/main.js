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
    loadScript('https://d3js.org/d3.v7.min.js')
        .then(() => {
            console.log("D3.js loaded successfully");
            // Load initial graph data
            refreshGraph();
            // Load initial statistics
            loadStatistics();
        })
        .catch(error => {
            console.error('Error loading D3.js:', error);
            alert('Failed to load D3.js. The graph visualization may not work properly.');
        });
    
    // Add ambient light to the scene
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.5);
    graphVisualizer.scene.add(ambientLight);
    
    // Add directional light to the scene
    const directionalLight = new THREE.DirectionalLight(0xffffff, 0.8);
    directionalLight.position.set(0, 1, 1);
    graphVisualizer.scene.add(directionalLight);
    
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
    
    // Function to load graph data
    function refreshGraph() {
        console.log("Refreshing graph data");
        apiClient.getGraphData()
            .then(data => {
                console.log(`Graph data loaded: ${data.nodes.length} nodes, ${data.links.length} links`);
                
                // Process nodes to ensure they have size property
                data.nodes.forEach(node => {
                    if (!node.size) {
                        node.size = 5; // Default size
                    }
                });
                
                // Set data to visualizer
                graphVisualizer.setData(data.nodes, data.links);
                
                // Start force simulation
                graphVisualizer.startSimulation();
            })
            .catch(error => {
                console.error('Error loading graph data:', error);
                alert('Failed to load graph data. Please check the console for details.');
            });
    }
    
    // Function to load statistics
    function loadStatistics() {
        console.log("Loading statistics");
        apiClient.getStatistics()
            .then(data => {
                console.log(`Statistics loaded: ${data.conceptCount} concepts, ${data.relationshipCount} relationships`);
                document.getElementById('concept-count').textContent = data.conceptCount;
                document.getElementById('relationship-count').textContent = data.relationshipCount;
            })
            .catch(error => {
                console.error('Error loading statistics:', error);
                alert('Failed to load statistics. Please check the console for details.');
            });
    }
    
    // Function to start the builder
    function startBuilder() {
        const seedConcept = document.getElementById('seed-concept').value;
        const maxNodes = parseInt(document.getElementById('max-nodes').value);
        const timeout = parseInt(document.getElementById('timeout').value);
        const randomRelationships = parseInt(document.getElementById('random-relationships').value);
        const concurrency = parseInt(document.getElementById('concurrency').value);
        
        const params = {
            seedConcept,
            maxNodes,
            timeout,
            randomRelationships,
            concurrency
        };
        
        console.log("Starting builder with params:", params);
        apiClient.startBuilder(params)
            .then(response => {
                console.log('Builder started:', response);
                alert('Knowledge Graph Builder started successfully.');
            })
            .catch(error => {
                console.error('Error starting builder:', error);
                alert('Error starting Knowledge Graph Builder. Check the console for details.');
            });
    }
    
    // Function to stop the builder
    function stopBuilder() {
        console.log("Stopping builder");
        apiClient.stopBuilder()
            .then(response => {
                console.log('Builder stopped:', response);
                alert('Knowledge Graph Builder stopped successfully.');
                refreshGraph();
                loadStatistics();
            })
            .catch(error => {
                console.error('Error stopping builder:', error);
                alert('Error stopping Knowledge Graph Builder. Check the console for details.');
            });
    }
    
    // Function to start the enricher
    function startEnricher() {
        const batchSize = parseInt(document.getElementById('batch-size').value);
        const interval = parseInt(document.getElementById('interval').value);
        const maxRelationships = parseInt(document.getElementById('max-relationships').value);
        const concurrency = parseInt(document.getElementById('enricher-concurrency').value);
        
        const params = {
            batchSize,
            interval,
            maxRelationships,
            concurrency
        };
        
        console.log("Starting enricher with params:", params);
        apiClient.startEnricher(params)
            .then(response => {
                console.log('Enricher started:', response);
                alert('Knowledge Graph Enricher started successfully.');
            })
            .catch(error => {
                console.error('Error starting enricher:', error);
                alert('Error starting Knowledge Graph Enricher. Check the console for details.');
            });
    }
    
    // Function to stop the enricher
    function stopEnricher() {
        console.log("Stopping enricher");
        apiClient.stopEnricher()
            .then(response => {
                console.log('Enricher stopped:', response);
                alert('Knowledge Graph Enricher stopped successfully.');
                refreshGraph();
                loadStatistics();
            })
            .catch(error => {
                console.error('Error stopping enricher:', error);
                alert('Error stopping Knowledge Graph Enricher. Check the console for details.');
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