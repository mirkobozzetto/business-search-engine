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
