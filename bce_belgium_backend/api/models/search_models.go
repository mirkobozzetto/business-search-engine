package models

type SearchResult struct {
	Table   string   `json:"table"`
	Column  string   `json:"column"`
	Query   string   `json:"query"`
	Results []string `json:"results"`
	Meta    Meta     `json:"meta"`
}

type CountResult struct {
	Table  string `json:"table"`
	Column string `json:"column"`
	Query  string `json:"query"`
	Count  int64  `json:"count"`
}
