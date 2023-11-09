package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/jainnidhi703/ol-processor/pkg/models"
)

// Cache is a global map used to store graphs indexed by their DAG names.
var Cache = make(map[string]*graph.Graph[string, models.GraphData])

// vertexAttributes define common styling for graph vertices.
var vertexAttributes = map[string]string{
	"colorscheme": "blues3",
	"style":       "filled",
	"color":       "2",
	"fillcolor":   "1",
	"shape":       "rectangle",
}

// vertexAttributesDeleted define styling for graph vertices representing deleted elements.
var vertexAttributesDeleted = map[string]string{
	"colorscheme": "reds3",
	"style":       "filled",
	"color":       "2",
	"fillcolor":   "1",
	"shape":       "rectangle",
}

// edgeAttributes define common styling for graph edges.
var edgeAttributes = map[string]string{}

// ProcessEvent takes an event and updates the corresponding graph in the cache.
func ProcessEvent(event models.Event) {
	dagName := getDagName(event)
	g := graph.New(graphDataHash, graph.Directed())
	log.Printf("Processing event for DAG: %s", dagName)
	if val, ok := Cache[dagName]; ok {
		g = *val
		log.Printf("Found existing graph for DAG: %s", dagName)
	} else {
		log.Printf("Creating new graph for DAG: %s", dagName)
	}
	g = buildGraph(event, g)
	Cache[dagName] = &g
	log.Println("Event processed and graph updated in cache.")
}

// RenderGraphToPNG takes a graph and a file path, and generates a PNG image of the graph.
func RenderGraphToPNG(g *graph.Graph[string, models.GraphData], filePath string) {
	file, _ := os.Create(filePath + ".gv")
	_ = draw.DOT(*g, file)
	log.Printf("Rendering graph to PNG: %s.png", filePath)
	_ = convertDotToPng(filePath+".gv", filePath+".png")
}

// convertDotToPng takes a .gv file path and a .png file path, and uses the 'dot' command
// to convert the .gv file to a .png file.
func convertDotToPng(dotFilePath, pngFilePath string) error {
	cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", pngFilePath)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error converting DOT to PNG: %v", err)
		return fmt.Errorf("error converting DOT to PNG: %w", err)
	}
	log.Printf("DOT file %s converted to PNG file %s", dotFilePath, pngFilePath)
	return nil
}

// getDagName extracts the DAG name from a given event by splitting the job name.
func getDagName(event models.Event) string {
	return strings.Split(event.Job.Name, ".")[0]
}

// graphDataHash creates a unique hash for a given GraphData object.
func graphDataHash(data models.GraphData) string {
	return data.Name + " | " + data.Type
}

// buildGraph takes an event and a graph, and incorporates the event data into the graph.
func buildGraph(event models.Event, g graph.Graph[string, models.GraphData]) graph.Graph[string, models.GraphData] {
	log.Println("Building graph...")
	vertexAttr := graph.VertexAttributes(vertexAttributes)
	edgeAtrr := graph.EdgeAttributes(edgeAttributes)
	job := event.Job
	jobGraphData := models.GraphData{Type: "job", Info: job.Facets.SQL.Query, Name: job.Name}
	_ = g.AddVertex(jobGraphData, vertexAttr)

	// Check if the job's SQL query is a DELETE operation
	isSQLDelete := checkSQLDropTable(jobGraphData.Info)

	// Process inputs and outputs, updating the graph with vertices and edges.
	// If the SQL query contains a DROP TABLE operation, it modifies the graph accordingly.
	for _, input := range event.Inputs {
		inputGraphData := models.GraphData{Type: "datasource", Info: input.Facets.DataSource.Name, Name: input.Name}
		_ = g.AddVertex(inputGraphData, vertexAttr)
		_ = g.AddEdge(graphDataHash(inputGraphData), graphDataHash(jobGraphData), edgeAtrr)
	}

	for _, output := range event.Outputs {
		outputGraphData := models.GraphData{Type: "datasource", Info: output.Facets.DataSource.Name, Name: output.Name}
		if isSQLDelete {
			_, err := g.Vertex(graphDataHash(outputGraphData))
			// Vertex doesnt exist, so create a new vertex with deleted status
			if err == nil {
				// Get all edges, remove them and re add
				var filteredEdges []graph.Edge[string]
				edges, _ := g.Edges()
				for _, edge := range edges {
					if edge.Target == graphDataHash(outputGraphData) {
						filteredEdges = append(filteredEdges, edge)
						_ = g.RemoveEdge(edge.Source, graphDataHash(outputGraphData))
					}
				}
				_ = g.RemoveVertex(graphDataHash(outputGraphData))
				// Adding new colored Vertex
				_ = g.AddVertex(outputGraphData, graph.VertexAttributes(vertexAttributesDeleted))
				for _, edge := range filteredEdges {
					_ = g.AddEdge(edge.Source, graphDataHash(outputGraphData))
				}
			} else {
				_ = g.AddVertex(outputGraphData, graph.VertexAttributes(vertexAttributesDeleted))
			}
		} else {
			_ = g.AddVertex(outputGraphData, vertexAttr)
		}
		_ = g.AddEdge(graphDataHash(jobGraphData), graphDataHash(outputGraphData), edgeAtrr)
	}
	log.Println("Graph built successfully.")
	return g
}

// checkSQLDropTable checks if a given SQL query string contains a "DROP TABLE" statement.
func checkSQLDropTable(query string) bool {
	lowercaseQuery := strings.ToLower(query)

	// Check if the query contains "drop table"
	return strings.Index(lowercaseQuery, "drop table") >= 0
}
