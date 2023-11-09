package handlers

import (
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jainnidhi703/ol-processor/pkg/models"
	"github.com/jainnidhi703/ol-processor/pkg/utils"
)

// cache is to store dag -> graph.
var cache = utils.Cache

// PostLineage reads the event from the request body, processes it, and updates the cache with graph.
func PostLineage(c *gin.Context) {
	var event models.Event

	data, err := io.ReadAll(c.Request.Body)

	if err != nil {
		log.Println("Error reading request body:", err)
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	err = json.Unmarshal(data, &event)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		c.JSON(400, gin.H{"error": "Failed to parse JSON"})
		return
	}

	utils.ProcessEvent(event)
	log.Printf("Processed lineage for DAG: %s", event.Job.Name)

	c.JSON(200, gin.H{"message": "JSON data processed successfully"})
}

// GetLineageGraph looks up the graph by the DAG name and renders it to a PNG file.
func GetLineageGraph(c *gin.Context) {
	dagName := c.Param("dag")
	log.Printf("Request received to get lineage graph for DAG: %s", dagName)

	if val, ok := cache[dagName]; ok {
		filePath := "./" + dagName
		utils.RenderGraphToPNG(val, filePath)
		log.Printf("Rendered lineage graph for DAG: %s", dagName)
		c.File(filePath + ".png")
	} else {
		log.Printf("Invalid DAG name: %s", dagName)
		c.JSON(400, gin.H{"error": "Invalid Dag Name"})
	}

}
