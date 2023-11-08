package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jainnidhi703/ol-processor/pkg/models"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

var cache = make(map[string]*graph.Graph[string, models.GraphData])

func PostLineage(c *gin.Context) {
	var event models.Event

	// Read the JSON data from the request body
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	// Parse the JSON data into the Event object
	err = json.Unmarshal(data, &event)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse JSON"})
		return
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Print the JSON string
	fmt.Println(string(eventJSON))

	dagName := getDagName(event)
	g := graph.New(graphDataHash, graph.Directed())
	if val, ok := cache[dagName]; ok {
		g = *val
	}
	g = buildGraph(event, g)
	cache[dagName] = &g

	// You can also send a response back to the client if needed
	c.JSON(200, gin.H{"message": "JSON data processed successfully"})
}

func GetLineageGraph(c *gin.Context) {
	dagName := c.Param("dag")
	// dagName := "lineage_combine"
	if val, ok := cache[dagName]; ok {
		filePath := "./" + dagName + ".gv"
		file, _ := os.Create(filePath)
		_ = draw.DOT(*val, file)
		c.JSON(200, gin.H{"message": "Graph built successfully"})
	} else {
		c.JSON(400, gin.H{"error": "Invalid Dag Name"})
	}

}

func getDagName(event models.Event) string {
	return strings.Split(event.Job.Name, ".")[0]
}

func graphDataHash(data models.GraphData) string {
	return data.Name
}

func buildGraph(event models.Event, g graph.Graph[string, models.GraphData]) graph.Graph[string, models.GraphData] {

	job := event.Job
	_ = g.AddVertex(models.GraphData{Type: "job", Info: job.Facets.SQL.Query, Name: job.Name})

	for _, input := range event.Inputs {
		_ = g.AddVertex(models.GraphData{Type: "datasource", Info: input.Facets.DataSource.Name, Name: input.Name})
		_ = g.AddEdge(input.Name, job.Name)
	}

	for _, output := range event.Outputs {
		_ = g.AddVertex(models.GraphData{Type: "datasource", Info: output.Facets.DataSource.Name, Name: output.Name})
		_ = g.AddEdge(job.Name, output.Name)
	}
	return g
}
