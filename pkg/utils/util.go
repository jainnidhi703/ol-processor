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
