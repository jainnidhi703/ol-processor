package models

type Facet struct {
	SQL        *struct{ Query string } `json:"sql,omitempty"`
	DataSource *struct{ Name string }  `json:"dataSource,omitempty"`
}

type Input struct {
	Facets *Facet `json:"facets,omitempty"`
	Name   string `json:"name"`
}

type Job struct {
	Facets *Facet `json:"facets,omitempty"`
	Name   string `json:"name"`
}

type Output struct {
	Facets *Facet `json:"facets,omitempty"`
	Name   string `json:"name"`
}

type Run struct {
	RunID string `json:"runId"`
}

type Event struct {
	Inputs  []Input  `json:"inputs"`
	Job     Job      `json:"job"`
	Outputs []Output `json:"outputs"`
	Run     Run      `json:"run"`
}

type GraphData struct {
	Type string `json:"type"`
	Info string `json:"info"`
	Name string `json:"name"`
}
