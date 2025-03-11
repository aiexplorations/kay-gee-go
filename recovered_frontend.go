package frontend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/kay-gee-go/internal/common/models"
	"github.com/kay-gee-go/internal/common/neo4j"
)

// Frontend represents a knowledge graph frontend
type Frontend struct {
	config     *config.Neo4jConfig
	neo4jClient *neo4j.Client
	templates  *template.Template
}

// NewFrontend creates a new knowledge graph frontend
func NewFrontend(config *config.Neo4jConfig) (*Frontend, error) {
	// Create Neo4j client
	neo4jClient, err := neo4j.NewClient(*config)
	if err != nil {
		return nil, err
	}

	// Create public directory if it doesn't exist
	if err := os.MkdirAll("public", 0755); err != nil {
		return nil, fmt.Errorf("failed to create public directory: %w", err)
	}

	// Create CSS directory if it doesn't exist
	if err := os.MkdirAll("public/css", 0755); err != nil {
		return nil, fmt.Errorf("failed to create CSS directory: %w", err)
	}

	// Create JS directory if it doesn't exist
	if err := os.MkdirAll("public/js", 0755); err != nil {
		return nil, fmt.Errorf("failed to create JS directory: %w", err)
	}

	// Create index.html if it doesn't exist
	if _, err := os.Stat("public/index.html"); os.IsNotExist(err) {
		content := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Knowledge Graph Viewer</title>
    <link rel="stylesheet" href="/css/style.css">
</head>
<body>
    <header>
        <h1>Knowledge Graph Viewer</h1>
    </header>
    <main>
        <section class="stats">
            <h2>Graph Statistics</h2>
            <div class="stats-container">
                <div class="stat-card">
                    <h3>Nodes</h3>
                    <p id="node-count">{{.Stats.NodeCount}}</p>
                </div>
                <div class="stat-card">
                    <h3>Relationships</h3>
                    <p id="relationship-count">{{.Stats.RelationshipCount}}</p>
                </div>
                <div class="stat-card">
                    <h3>Last Updated</h3>
                    <p id="last-updated">{{.Stats.LastUpdated.Format "Jan 02, 2006 15:04:05"}}</p>
                </div>
            </div>
        </section>
        <section class="graph-container">
            <h2>Graph Visualization</h2>
            <div id="graph"></div>
        </section>
    </main>
    <footer>
        <p>&copy; 2023 Knowledge Graph Builder</p>
    </footer>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <script src="/js/script.js"></script>
</body>
</html>`
		if err := os.WriteFile("public/index.html", []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to create index.html: %w", err)
		}
	}

	// Create style.css if it doesn't exist
	if _, err := os.Stat("public/css/style.css"); os.IsNotExist(err) {
		if err := createStyleCSS(); err != nil {
			return nil, fmt.Errorf("failed to create style.css: %w", err)
		}
	}

	// Create script.js if it doesn't exist
	if _, err := os.Stat("public/js/script.js"); os.IsNotExist(err) {
		if err := createScriptJS(); err != nil {
			return nil, fmt.Errorf("failed to create script.js: %w", err)
		}
	}

	// Parse templates
	templates, err := template.ParseGlob("public/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Frontend{
		config:     config,
		neo4jClient: neo4jClient,
		templates:  templates,
	}, nil
}

// Close closes the frontend
func (f *Frontend) Close() error {
	return f.neo4jClient.Close()
}

// Start starts the frontend server
func (f *Frontend) Start(port int) error {
	// Create public directory if it doesn't exist
	if err := os.MkdirAll("public", 0755); err != nil {
		return fmt.Errorf("failed to create public directory: %w", err)
	}

	// Create CSS directory if it doesn't exist
	if err := os.MkdirAll("public/css", 0755); err != nil {
		return fmt.Errorf("failed to create CSS directory: %w", err)
	}

	// Create JS directory if it doesn't exist
	if err := os.MkdirAll("public/js", 0755); err != nil {
		return fmt.Errorf("failed to create JS directory: %w", err)
	}

	// Create index.html if it doesn't exist
	if _, err := os.Stat("public/index.html"); os.IsNotExist(err) {
		if err := f.createIndexHTML(); err != nil {
			return fmt.Errorf("failed to create index.html: %w", err)
		}
	}

	// Create style.css if it doesn't exist
	if _, err := os.Stat("public/css/style.css"); os.IsNotExist(err) {
		if err := createStyleCSS(); err != nil {
			return fmt.Errorf("failed to create style.css: %w", err)
		}
	}

	// Create script.js if it doesn't exist
	if _, err := os.Stat("public/js/script.js"); os.IsNotExist(err) {
		if err := createScriptJS(); err != nil {
			return fmt.Errorf("failed to create script.js: %w", err)
		}
	}

	// Set up routes
	http.HandleFunc("/", f.handleIndex)
	http.HandleFunc("/api/stats", f.handleStats)
	http.HandleFunc("/api/graph", f.handleGraph)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))

	// Start server
	fmt.Printf("Starting frontend server on port %d...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// handleIndex handles the index page
func (f *Frontend) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get graph stats
	stats, err := f.neo4jClient.GetGraphStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get graph stats: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Render template
	data := struct {
		Stats *models.GraphStats
	}{
		Stats: stats,
	}

	if err := f.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to render template: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// handleStats handles the stats API endpoint
func (f *Frontend) handleStats(w http.ResponseWriter, r *http.Request) {
	// Get graph stats
	stats, err := f.neo4jClient.GetGraphStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get graph stats: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleGraph handles the graph API endpoint
func (f *Frontend) handleGraph(w http.ResponseWriter, r *http.Request) {
	// Get all concepts
	concepts, err := f.neo4jClient.GetRandomConcepts(1000)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get concepts: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(concepts)
}

// createIndexHTML creates the index.html file
func (f *Frontend) createIndexHTML() error {
	content := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Knowledge Graph Viewer</title>
    <link rel="stylesheet" href="/css/style.css">
</head>
<body>
    <header>
        <h1>Knowledge Graph Viewer</h1>
    </header>
    <main>
        <section class="stats">
            <h2>Graph Statistics</h2>
            <div class="stats-container">
                <div class="stat-card">
                    <h3>Nodes</h3>
                    <p id="node-count">{{.Stats.NodeCount}}</p>
                </div>
                <div class="stat-card">
                    <h3>Relationships</h3>
                    <p id="relationship-count">{{.Stats.RelationshipCount}}</p>
                </div>
                <div class="stat-card">
                    <h3>Last Updated</h3>
                    <p id="last-updated">{{.Stats.LastUpdated.Format "Jan 02, 2006 15:04:05"}}</p>
                </div>
            </div>
        </section>
        <section class="graph-container">
            <h2>Graph Visualization</h2>
            <div id="graph"></div>
        </section>
    </main>
    <footer>
        <p>&copy; 2023 Knowledge Graph Builder</p>
    </footer>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <script src="/js/script.js"></script>
</body>
</html>`

	return os.WriteFile("public/index.html", []byte(content), 0644)
}

// createStyleCSS creates the style.css file
func (f *Frontend) createStyleCSS() error {
	content := `/* Reset */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

/* Global styles */
body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    line-height: 1.6;
    color: #333;
    background-color: #f8f9fa;
}

header {
    background-color: #343a40;
    color: white;
    padding: 1rem;
    text-align: center;
}

main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
}

h1, h2, h3 {
    margin-bottom: 1rem;
}

/* Stats section */
.stats-container {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2rem;
}

.stat-card {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 1.5rem;
    width: 30%;
    text-align: center;
}

.stat-card h3 {
    color: #6c757d;
    font-size: 1.2rem;
}

.stat-card p {
    font-size: 2rem;
    font-weight: bold;
    color: #007bff;
}

/* Graph section */
.graph-container {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 1.5rem;
    margin-bottom: 2rem;
}

#graph {
    width: 100%;
    height: 600px;
    border: 1px solid #dee2e6;
    border-radius: 4px;
}

/* Footer */
footer {
    background-color: #343a40;
    color: white;
    text-align: center;
    padding: 1rem;
    margin-top: 2rem;
}`

	return os.WriteFile("public/css/style.css", []byte(content), 0644)
}

// createScriptJS creates the script.js file
func (f *Frontend) createScriptJS() error {
	content := `// Update stats every 5 seconds
setInterval(updateStats, 5000);

// Initial graph load
document.addEventListener('DOMContentLoaded', function() {
    loadGraph();
});

// Update stats function
function updateStats() {
    fetch('/api/stats')
        .then(response => response.json())
        .then(data => {
            document.getElementById('node-count').textContent = data.node_count;
            document.getElementById('relationship-count').textContent = data.relationship_count;
            document.getElementById('last-updated').textContent = new Date(data.last_updated).toLocaleString();
        })
        .catch(error => console.error('Error fetching stats:', error));
}

// Load graph function
function loadGraph() {
    fetch('/api/graph')
        .then(response => response.json())
        .then(data => {
            if (!data || data.length === 0) {
                document.getElementById('graph').innerHTML = '<p class="no-data">No data available. Please build the knowledge graph first.</p>';
                return;
            }
            
            createForceGraph(data);
        })
        .catch(error => {
            console.error('Error fetching graph data:', error);
            document.getElementById('graph').innerHTML = '<p class="error">Error loading graph data. Please check the console for details.</p>';
        });
}

// Create force-directed graph
function createForceGraph(concepts) {
    const width = document.getElementById('graph').clientWidth;
    const height = 600;
    
    // Clear previous graph
    d3.select('#graph').html('');
    
    // Create SVG
    const svg = d3.select('#graph')
        .append('svg')
        .attr('width', width)
        .attr('height', height);
        
    // Add zoom functionality
    const g = svg.append('g');
    svg.call(d3.zoom().on('zoom', (event) => {
        g.attr('transform', event.transform);
    }));
    
    // Process data to create nodes and links
    const nodes = concepts.map(concept => ({
        id: concept.id,
        name: concept.name || 'Unnamed Concept',
        description: concept.description || 'No description available'
    }));
    
    // Create some sample links between nodes
    // In a real application, you would fetch actual relationships from the backend
    const links = [];
    for (let i = 0; i < nodes.length; i++) {
        // Connect to a random node
        const target = Math.floor(Math.random() * nodes.length);
        if (i !== target) {
            links.push({
                source: i,
                target: target
            });
        }
    }
    
    // Create simulation
    const simulation = d3.forceSimulation(nodes)
        .force('link', d3.forceLink(links).id(d => d.id).distance(100))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('collision', d3.forceCollide().radius(50));
    
    // Create links
    const link = g.append('g')
        .attr('class', 'links')
        .selectAll('line')
        .data(links)
        .enter()
        .append('line')
        .attr('stroke', '#999')
        .attr('stroke-opacity', 0.6)
        .attr('stroke-width', 1);
    
    // Create nodes
    const node = g.append('g')
        .attr('class', 'nodes')
        .selectAll('g')
        .data(nodes)
        .enter()
        .append('g')
        .call(d3.drag()
            .on('start', dragstarted)
            .on('drag', dragged)
            .on('end', dragended));
    
    // Add circles to nodes
    node.append('circle')
        .attr('r', 20)
        .attr('fill', d => stringToColor(d.name))
        .append('title')
        .text(d => d.description);
    
    // Add text to nodes
    node.append('text')
        .text(d => d.name)
        .attr('text-anchor', 'middle')
        .attr('dy', 30)
        .attr('font-size', 10);
    
    // Update positions on tick
    simulation.on('tick', () => {
        link
            .attr('x1', d => d.source.x)
            .attr('y1', d => d.source.y)
            .attr('x2', d => d.target.x)
            .attr('y2', d => d.target.y);
            
        node.attr('transform', d => 'translate(' + d.x + ',' + d.y + ')');
    });
    
    // Drag functions
    function dragstarted(event, d) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
    }
    
    function dragged(event, d) {
        d.fx = event.x;
        d.fy = event.y;
    }
    
    function dragended(event, d) {
        if (!event.active) simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
    }
    
    // Helper function to generate colors from strings
    function stringToColor(str) {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            hash = str.charCodeAt(i) + ((hash << 5) - hash);
        }
        let color = '#';
        for (let i = 0; i < 3; i++) {
            const value = (hash >> (i * 8)) & 0xFF;
            color += ('00' + value.toString(16)).substr(-2);
        }
        return color;
    }
}`

	return os.WriteFile("public/js/script.js", []byte(content), 0644)
}

// createStyleCSS creates the style.css file
func createStyleCSS() error {
	content := `/* Reset */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

/* Global styles */
body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    line-height: 1.6;
    color: #333;
    background-color: #f8f9fa;
}

header {
    background-color: #343a40;
    color: white;
    padding: 1rem;
    text-align: center;
}

main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
}

h1, h2, h3 {
    margin-bottom: 1rem;
}

/* Stats section */
.stats-container {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2rem;
}

.stat-card {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 1.5rem;
    width: 30%;
    text-align: center;
}

.stat-card h3 {
    color: #6c757d;
    font-size: 1.2rem;
}

.stat-card p {
    font-size: 2rem;
    font-weight: bold;
    color: #007bff;
}

/* Graph section */
.graph-container {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 1.5rem;
    margin-bottom: 2rem;
}

#graph {
    width: 100%;
    height: 600px;
    border: 1px solid #dee2e6;
    border-radius: 4px;
}

/* Footer */
footer {
    background-color: #343a40;
    color: white;
    text-align: center;
    padding: 1rem;
    margin-top: 2rem;
}`

	return os.WriteFile("public/css/style.css", []byte(content), 0644)
}

// createScriptJS creates the script.js file
func createScriptJS() error {
	content := `// Update stats every 5 seconds
setInterval(updateStats, 5000);

// Initial graph load
document.addEventListener('DOMContentLoaded', function() {
    loadGraph();
});

// Update stats function
function updateStats() {
    fetch('/api/stats')
        .then(response => response.json())
        .then(data => {
            document.getElementById('node-count').textContent = data.node_count;
            document.getElementById('relationship-count').textContent = data.relationship_count;
            document.getElementById('last-updated').textContent = new Date(data.last_updated).toLocaleString();
        })
        .catch(error => console.error('Error fetching stats:', error));
}

// Load graph function
function loadGraph() {
    fetch('/api/graph')
        .then(response => response.json())
        .then(data => {
            if (!data || data.length === 0) {
                document.getElementById('graph').innerHTML = '<p class="no-data">No data available. Please build the knowledge graph first.</p>';
                return;
            }
            
            createForceGraph(data);
        })
        .catch(error => {
            console.error('Error fetching graph data:', error);
            document.getElementById('graph').innerHTML = '<p class="error">Error loading graph data. Please check the console for details.</p>';
        });
}

// Create force-directed graph
function createForceGraph(concepts) {
    const width = document.getElementById('graph').clientWidth;
    const height = 600;
    
    // Clear previous graph
    d3.select('#graph').html('');
    
    // Create SVG
    const svg = d3.select('#graph')
        .append('svg')
        .attr('width', width)
        .attr('height', height);
        
    // Add zoom functionality
    const g = svg.append('g');
    svg.call(d3.zoom().on('zoom', (event) => {
        g.attr('transform', event.transform);
    }));
    
    // Process data to create nodes and links
    const nodes = concepts.map(concept => ({
        id: concept.id,
        name: concept.name || 'Unnamed Concept',
        description: concept.description || 'No description available'
    }));
    
    // Create some sample links between nodes
    // In a real application, you would fetch actual relationships from the backend
    const links = [];
    for (let i = 0; i < nodes.length; i++) {
        // Connect to a random node
        const target = Math.floor(Math.random() * nodes.length);
        if (i !== target) {
            links.push({
                source: i,
                target: target
            });
        }
    }
    
    // Create simulation
    const simulation = d3.forceSimulation(nodes)
        .force('link', d3.forceLink(links).id(d => d.id).distance(100))
        .force('charge', d3.forceManyBody().strength(-300))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('collision', d3.forceCollide().radius(50));
    
    // Create links
    const link = g.append('g')
        .attr('class', 'links')
        .selectAll('line')
        .data(links)
        .enter()
        .append('line')
        .attr('stroke', '#999')
        .attr('stroke-opacity', 0.6)
        .attr('stroke-width', 1);
    
    // Create nodes
    const node = g.append('g')
        .attr('class', 'nodes')
        .selectAll('g')
        .data(nodes)
        .enter()
        .append('g')
        .call(d3.drag()
            .on('start', dragstarted)
            .on('drag', dragged)
            .on('end', dragended));
    
    // Add circles to nodes
    node.append('circle')
        .attr('r', 20)
        .attr('fill', d => stringToColor(d.name))
        .append('title')
        .text(d => d.description);
    
    // Add text to nodes
    node.append('text')
        .text(d => d.name)
        .attr('text-anchor', 'middle')
        .attr('dy', 30)
        .attr('font-size', 10);
    
    // Update positions on tick
    simulation.on('tick', () => {
        link
            .attr('x1', d => d.source.x)
            .attr('y1', d => d.source.y)
            .attr('x2', d => d.target.x)
            .attr('y2', d => d.target.y);
            
        node.attr('transform', d => 'translate(' + d.x + ',' + d.y + ')');
    });
    
    // Drag functions
    function dragstarted(event, d) {
        if (!event.active) simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
    }
    
    function dragged(event, d) {
        d.fx = event.x;
        d.fy = event.y;
    }
    
    function dragended(event, d) {
        if (!event.active) simulation.alphaTarget(0);
        d.fx = null;
        d.fy = null;
    }
    
    // Helper function to generate colors from strings
    function stringToColor(str) {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            hash = str.charCodeAt(i) + ((hash << 5) - hash);
        }
        let color = '#';
        for (let i = 0; i < 3; i++) {
            const value = (hash >> (i * 8)) & 0xFF;
            color += ('00' + value.toString(16)).substr(-2);
        }
        return color;
    }
}`

	return os.WriteFile("public/js/script.js", []byte(content), 0644)
} 