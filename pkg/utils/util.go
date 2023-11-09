package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/jainnidhi703/ol-processor/pkg/models"
)

var Cache = make(map[string]*graph.Graph[string, models.GraphData])

var vertexAttributes = map[string]string{
	"colorscheme": "blues3",
	"style":       "filled",
	"color":       "2",
	"fillcolor":   "1",
	"shape":       "rectangle",
}

var vertexAttributesDeleted = map[string]string{
	"colorscheme": "reds3",
	"style":       "filled",
	"color":       "2",
	"fillcolor":   "1",
	"shape":       "rectangle",
}

var edgeAttributes = map[string]string{}

func ProcessEvent(event models.Event) {
	dagName := getDagName(event)
	g := graph.New(graphDataHash, graph.Directed())
	if val, ok := Cache[dagName]; ok {
		g = *val
	}
	g = buildGraph(event, g)
	Cache[dagName] = &g
}

func RenderGraphToPNG(g *graph.Graph[string, models.GraphData], filePath string) {
	file, _ := os.Create(filePath + ".gv")
	_ = draw.DOT(*g, file)
	_ = convertDotToPng(filePath+".gv", filePath+".png")
}

func convertDotToPng(dotFilePath, pngFilePath string) error {
	cmd := exec.Command("dot", "-Tpng", dotFilePath, "-o", pngFilePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error converting DOT to PNG: %w", err)
	}
	fmt.Printf("DOT file %s converted to PNG file %s\n", dotFilePath, pngFilePath)
	return nil
}

func getDagName(event models.Event) string {
	return strings.Split(event.Job.Name, ".")[0]
}

func graphDataHash(data models.GraphData) string {
	return data.Name + " | " + data.Type
}

func buildGraph(event models.Event, g graph.Graph[string, models.GraphData]) graph.Graph[string, models.GraphData] {
	vertexAttr := graph.VertexAttributes(vertexAttributes)
	edgeAtrr := graph.EdgeAttributes(edgeAttributes)
	job := event.Job
	jobGraphData := models.GraphData{Type: "job", Info: job.Facets.SQL.Query, Name: job.Name}
	_ = g.AddVertex(jobGraphData, vertexAttr)

	// Check if the query contains sql drop syntax
	isSQLDelete := checkSQLDropTable(jobGraphData.Info)

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
	return g
}

func checkSQLDropTable(query string) bool {
	lowercaseQuery := strings.ToLower(query)

	// Check if the query contains "drop table"
	return strings.Index(lowercaseQuery, "drop table") >= 0
}
