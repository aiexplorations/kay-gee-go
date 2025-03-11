from typing import Dict, List, Optional, Any, Union
import os
import subprocess
import json
from fastapi import FastAPI, HTTPException, Query
from fastapi.staticfiles import StaticFiles
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import uvicorn
from neo4j import GraphDatabase
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, 
                    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Create FastAPI app
app = FastAPI(title="Knowledge Graph Visualizer API")

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Neo4j connection settings
NEO4J_URI = os.getenv("NEO4J_URI", "bolt://neo4j:7687")
NEO4J_USER = os.getenv("NEO4J_USER", "neo4j")
NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD", "password")

# Create Neo4j driver
neo4j_driver = GraphDatabase.driver(
    NEO4J_URI, 
    auth=(NEO4J_USER, NEO4J_PASSWORD)
)

# Pydantic models for request/response validation
class BuilderParams(BaseModel):
    seedConcept: str
    maxNodes: int
    timeout: int
    randomRelationships: int
    concurrency: int

class EnricherParams(BaseModel):
    batchSize: int
    interval: int
    maxRelationships: int
    concurrency: int

class RelationshipCreate(BaseModel):
    source: str
    target: str
    type: str

class GraphData(BaseModel):
    nodes: List[Dict[str, Any]]
    links: List[Dict[str, Any]]

class Statistics(BaseModel):
    conceptCount: int
    relationshipCount: int

# Helper function to run shell commands
def run_command(command: List[str]) -> str:
    try:
        logger.info(f"Running command: {' '.join(command)}")
        result = subprocess.run(command, capture_output=True, text=True, check=True)
        return result.stdout
    except subprocess.CalledProcessError as e:
        logger.error(f"Command failed: {e.stderr}")
        raise HTTPException(status_code=500, detail=f"Command failed: {e.stderr}")

# API routes
@app.get("/api/graph", response_model=GraphData)
async def get_graph_data():
    """
    Get the current graph data from Neo4j
    """
    try:
        with neo4j_driver.session() as session:
            # Query nodes
            nodes_result = session.run("""
                MATCH (n:Concept)
                RETURN id(n) AS id, n.name AS name
            """)
            nodes = [{"id": str(record["id"]), "name": record["name"]} for record in nodes_result]
            
            # Query relationships
            links_result = session.run("""
                MATCH (a:Concept)-[r]->(b:Concept)
                RETURN id(a) AS source, id(b) AS target, type(r) AS type
            """)
            links = [{"source": str(record["source"]), "target": str(record["target"]), "type": record["type"]} 
                    for record in links_result]
            
            return {"nodes": nodes, "links": links}
    except Exception as e:
        logger.error(f"Error fetching graph data: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error fetching graph data: {str(e)}")

@app.post("/api/builder/start")
async def start_builder(params: BuilderParams):
    """
    Start the knowledge graph builder
    """
    try:
        # Build command to start the builder
        command = [
            "/bin/sh", 
            "/app/start-builder.sh",
            "--seed", params.seedConcept,
            "--max-nodes", str(params.maxNodes),
            "--timeout", str(params.timeout),
            "--random-relationships", str(params.randomRelationships),
            "--concurrency", str(params.concurrency)
        ]
        
        # Run the command
        output = run_command(command)
        
        return {"status": "success", "message": "Builder started successfully", "output": output}
    except Exception as e:
        logger.error(f"Error starting builder: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error starting builder: {str(e)}")

@app.post("/api/builder/stop")
async def stop_builder():
    """
    Stop the knowledge graph builder
    """
    try:
        # Build command to stop the builder
        command = ["/bin/sh", "/app/stop-builder.sh"]
        
        # Run the command
        output = run_command(command)
        
        return {"status": "success", "message": "Builder stopped successfully", "output": output}
    except Exception as e:
        logger.error(f"Error stopping builder: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error stopping builder: {str(e)}")

@app.post("/api/enricher/start")
async def start_enricher(params: EnricherParams):
    """
    Start the knowledge graph enricher
    """
    try:
        # Build command to start the enricher
        command = [
            "/bin/sh", 
            "/app/start-enricher.sh",
            "--batch-size", str(params.batchSize),
            "--interval", str(params.interval),
            "--max-relationships", str(params.maxRelationships),
            "--concurrency", str(params.concurrency)
        ]
        
        # Run the command
        output = run_command(command)
        
        return {"status": "success", "message": "Enricher started successfully", "output": output}
    except Exception as e:
        logger.error(f"Error starting enricher: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error starting enricher: {str(e)}")

@app.post("/api/enricher/stop")
async def stop_enricher():
    """
    Stop the knowledge graph enricher
    """
    try:
        # Build command to stop the enricher
        command = ["/bin/sh", "/app/stop-enricher.sh"]
        
        # Run the command
        output = run_command(command)
        
        return {"status": "success", "message": "Enricher stopped successfully", "output": output}
    except Exception as e:
        logger.error(f"Error stopping enricher: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error stopping enricher: {str(e)}")

@app.get("/api/concepts/search")
async def search_concepts(q: str = Query(..., min_length=1)):
    """
    Search for concepts in the graph
    """
    try:
        with neo4j_driver.session() as session:
            # Query concepts that match the search term
            result = session.run("""
                MATCH (n:Concept)
                WHERE n.name CONTAINS $query
                RETURN id(n) AS id, n.name AS name
                LIMIT 10
            """, {"query": q})
            
            concepts = [{"id": str(record["id"]), "name": record["name"]} for record in result]
            
            return concepts
    except Exception as e:
        logger.error(f"Error searching concepts: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error searching concepts: {str(e)}")

@app.post("/api/relationships")
async def create_relationship(relationship: RelationshipCreate):
    """
    Create a relationship between two concepts
    """
    try:
        with neo4j_driver.session() as session:
            # Create relationship between concepts
            result = session.run("""
                MATCH (a:Concept), (b:Concept)
                WHERE id(a) = $source AND id(b) = $target
                CREATE (a)-[r:`{}`]->(b)
                RETURN id(a) AS source, id(b) AS target, type(r) AS type
            """.format(relationship.type), {"source": int(relationship.source), "target": int(relationship.target)})
            
            record = result.single()
            if not record:
                raise HTTPException(status_code=404, detail="Concepts not found")
            
            return {
                "source": str(record["source"]), 
                "target": str(record["target"]), 
                "type": record["type"]
            }
    except Exception as e:
        logger.error(f"Error creating relationship: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error creating relationship: {str(e)}")

@app.get("/api/statistics", response_model=Statistics)
async def get_statistics():
    """
    Get statistics about the graph
    """
    try:
        with neo4j_driver.session() as session:
            # Query concept count
            concept_result = session.run("""
                MATCH (n:Concept)
                RETURN count(n) AS count
            """)
            concept_count = concept_result.single()["count"]
            
            # Query relationship count
            relationship_result = session.run("""
                MATCH ()-[r]->()
                RETURN count(r) AS count
            """)
            relationship_count = relationship_result.single()["count"]
            
            return {"conceptCount": concept_count, "relationshipCount": relationship_count}
    except Exception as e:
        logger.error(f"Error fetching statistics: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error fetching statistics: {str(e)}")

# Mount static files
app.mount("/", StaticFiles(directory="/app/public", html=True), name="static")

# Startup and shutdown events
@app.on_event("startup")
async def startup_event():
    logger.info("Starting Knowledge Graph Visualizer API")

@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Shutting down Knowledge Graph Visualizer API")
    neo4j_driver.close()

if __name__ == "__main__":
    uvicorn.run("server:app", host="0.0.0.0", port=8080, reload=True) 