package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ol-processor/pkg/handlers"
)

func main() {
	router := gin.Default()

	router.POST("/api/v1/lineage", handlers.PostLineage)
	router.GET("/api/get/graph/:dag", handlers.GetLineageGraph)

	router.Run("localhost:3000")
}
