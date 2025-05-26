package models

type ColumnInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type TableStructure struct {
	Name    string       `json:"name"`
	Rows    int64        `json:"rows"`
	Columns []ColumnInfo `json:"columns"`
}


type APIResponse struct {
	Success bool        `json:"success"`
	Data    any         `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Count  int `json:"count,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type Table struct {
	Name    string `json:"name"`
	Rows    int64  `json:"rows"`
	Columns int    `json:"columns"`
}

type TableInfo struct {
	Table   string   `json:"table"`
	Rows    int64    `json:"rows"`
	Columns int      `json:"columns"`
	Fields  []string `json:"fields,omitempty"`
}

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

func Success(data any) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func SuccessWithMeta(data any, meta Meta) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Meta:    &meta,
	}
}

func Error(message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   message,
	}
}
