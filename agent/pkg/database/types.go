package database

// TableInfo represents metadata about a database table or view.
type TableInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Engine      string `json:"engine"`
	Rows        int64  `json:"rows"`
	DataLength  int64  `json:"data_length"`
	IndexLength int64  `json:"index_length"`
	TotalSize   int64  `json:"total_size"`
	Collation   string `json:"collation"`
	Comment     string `json:"comment"`
}

// ColumnInfo represents the schema definition of a single table column.
type ColumnInfo struct {
	Name       string  `json:"name"`
	Position   int     `json:"position"`
	Type       string  `json:"type"`
	DataType   string  `json:"data_type"`
	IsNullable bool    `json:"is_nullable"`
	Key        string  `json:"key"`
	Default    *string `json:"default"`
	Extra      string  `json:"extra"`
	Comment    string  `json:"comment"`
}

// IndexInfo represents index metadata for a table.
type IndexInfo struct {
	Name       string `json:"name"`
	Column     string `json:"column"`
	NonUnique  bool   `json:"non_unique"`
	IsPrimary  bool   `json:"is_primary"`
	SeqInIndex int    `json:"seq_in_index"`
	Type       string `json:"type"`
}

// ForeignKeyInfo represents a foreign key relationship.
type ForeignKeyInfo struct {
	ConstraintName   string `json:"constraint_name"`
	Column           string `json:"column"`
	ReferencedTable  string `json:"referenced_table"`
	ReferencedColumn string `json:"referenced_column"`
}

// TableStructure contains the complete schema metadata for a table.
type TableStructure struct {
	TableName       string           `json:"table_name"`
	Columns         []ColumnInfo     `json:"columns"`
	Indexes         []IndexInfo      `json:"indexes"`
	ForeignKeys     []ForeignKeyInfo `json:"foreign_keys"`
	CreateStatement string           `json:"create_statement"`
}

// TableDataResult represents paginated rows and column headers from a table.
type TableDataResult struct {
	TableName  string                   `json:"table_name"`
	Columns    []string                 `json:"columns"`
	Rows       []map[string]interface{} `json:"rows"`
	TotalRows  int64                    `json:"total_rows"`
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	TotalPages int                      `json:"total_pages"`
}
