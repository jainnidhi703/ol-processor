package models

// Facet represents a set of attributes that describe additional characteristics of data.
// It can include SQL facets like queries and data source facets like names.
type Facet struct {
	SQL        *struct{ Query string } `json:"sql,omitempty"`
	DataSource *struct{ Name string }  `json:"dataSource,omitempty"`
}

// Input represents an input entity for a job with its specific facets and a name.
type Input struct {
	Facets *Facet `json:"facets,omitempty"`
	Name   string `json:"name"`
}

// Job represents a job entity with its specific facets and a name.
type Job struct {
	Facets *Facet `json:"facets,omitempty"`
	Name   string `json:"name"`
}

// Output represents an output entity from a job with its specific facets and a name.
type Output struct {
	Facets *Facet `json:"facets,omitempty"`
	Name   string `json:"name"`
}

// Run represents a specific execution or run of a job, identified by a unique run ID.
type Run struct {
	RunID string `json:"runId"`
}

// Event encapsulates the entire data event, including inputs, the job itself, outputs, and the run identifier.
type Event struct {
	Inputs  []Input  `json:"inputs"`
	Job     Job      `json:"job"`
	Outputs []Output `json:"outputs"`
	Run     Run      `json:"run"`
}

// GraphData is a generic type that encapsulates information about an entity in a graph.
// It is used to represent nodes or edges within the graph with a type, informational string, and name.
type GraphData struct {
	Type string `json:"type"`
	Info string `json:"info"`
	Name string `json:"name"`
}
