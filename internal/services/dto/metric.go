package dto

//go:generate easyjson -all
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

//go:generate easyjson -all
type MetricsCollection struct {
	Metrics []Metrics `json:"metrics"`
}