package models

type PreviewData struct {
	Table   string           `json:"table"`
	Columns []string         `json:"columns"`
	Data    []map[string]any `json:"data"`
	Meta    Meta             `json:"meta"`
}

type ColumnValue struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type ColumnValues struct {
	Table  string        `json:"table"`
	Column string        `json:"column"`
	Values []ColumnValue `json:"values"`
	Meta   Meta          `json:"meta"`
}

type NaceSearchResult struct {
	Query   string              `json:"query"`
	Results []map[string]any    `json:"results"`
	Meta    Meta                `json:"meta"`
}
