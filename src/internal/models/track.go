package models

type Track struct {
	Id         uint64
	Source     string
	Producers  []string
	Authors    []string
	Performers []string
	Name       string
	Genre      string
	Embedding  []float64
}
