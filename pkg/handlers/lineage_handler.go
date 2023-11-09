package handlers

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/jainnidhi703/ol-processor/pkg/models"
	"github.com/jainnidhi703/ol-processor/pkg/utils"
)

var cache = utils.Cache

func PostLineage(c *gin.Context) {
	var event models.Event

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to read request body"})
		return
	}

	err = json.Unmarshal(data, &event)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to parse JSON"})
		return
	}

	utils.ProcessEvent(event)

	c.JSON(200, gin.H{"message": "JSON data processed successfully"})
}

func GetLineageGraph(c *gin.Context) {
	dagName := c.Param("dag")
	if val, ok := cache[dagName]; ok {
		filePath := "./" + dagName
		utils.RenderGraphToPNG(val, filePath)
		c.File(filePath + ".png")
	} else {
		c.JSON(400, gin.H{"error": "Invalid Dag Name"})
	}
}
