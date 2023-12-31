package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jainnidhi703/ol-processor/pkg/handlers"
)

func main() {
	router := gin.Default()

	router.POST("/api/v1/lineage", handlers.PostLineage)
	router.GET("/api/get/graph/:dag", handlers.GetLineageGraph)

	router.Run("0.0.0.0:3000")
}
