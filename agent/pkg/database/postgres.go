package database

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type PostgresManager struct {
	isDev bool
}

func NewPostgresManager(isDev bool) *PostgresManager {
	return &PostgresManager{
		isDev: isDev || runtime.GOOS != "linux",
	}
}

func (p *PostgresManager) execQuery(query string) error {
	if p.isDev {
		return nil
	}

	cmd := exec.Command("sudo", "-u", "postgres", "psql", "-c", query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("postgres error: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}
	return nil
}

func (p *PostgresManager) execQueryOutput(database, query string) (string, error) {
	if p.isDev {
		return "", nil
	}

	args := []string{"-u", "postgres", "psql"}
	if database != "" {
		args = append(args, "-d", database)
	}
	args = append(args, "-t", "-A", "-F", "\t", "-c", query)

	cmd := exec.Command("sudo", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("postgres error: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func (p *PostgresManager) execDatabaseQuery(database, query string) error {
	if p.isDev {
		return nil
	}

	args := []string{"-u", "postgres", "psql"}
	if database != "" {
		args = append(args, "-d", database)
	}
	args = append(args, "-c", query)

	cmd := exec.Command("sudo", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("postgres error on database '%s': %s: %s", database, err.Error(), strings.TrimSpace(string(out)))
	}
	return nil
}

func (p *PostgresManager) CreateDatabase(name, encoding, owner string) error {
	if encoding == "" {
		encoding = "UTF8"
	}

	ownerClause := ""
	if owner != "" {
		ownerClause = fmt.Sprintf("OWNER \"%s\"", owner)
	}

	query := fmt.Sprintf("CREATE DATABASE \"%s\" ENCODING '%s' %s;", name, encoding, ownerClause)
	if err := p.execQuery(query); err != nil {
		return err
	}

	if owner != "" {
		schemaQuery := fmt.Sprintf(
			"GRANT ALL ON SCHEMA public TO \"%s\"; "+
				"GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO \"%s\"; "+
				"GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO \"%s\"; "+
				"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO \"%s\"; "+
				"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO \"%s\";",
			owner, owner, owner, owner, owner,
		)
		_ = p.execDatabaseQuery(name, schemaQuery)
	}
	return nil
}

func (p *PostgresManager) DropDatabase(name string) error {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\";", name)
	return p.execQuery(query)
}

func (p *PostgresManager) CreateUser(username, password string) error {
	query := fmt.Sprintf("CREATE ROLE \"%s\" WITH LOGIN PASSWORD '%s';", username, password)
	return p.execQuery(query)
}

func (p *PostgresManager) DropUser(username string) error {
	query := fmt.Sprintf("DROP ROLE IF EXISTS \"%s\";", username)
	return p.execQuery(query)
}

func (p *PostgresManager) GrantPrivileges(database, username, privileges string) error {
	if privileges == "" || strings.ToLower(privileges) == "all" {
		privileges = "ALL PRIVILEGES"
	}

	// 1. Grant database-level privileges
	query := fmt.Sprintf("GRANT %s ON DATABASE \"%s\" TO \"%s\";", privileges, database, username)
	if err := p.execQuery(query); err != nil {
		return err
	}

	// 2. Grant schema-level and default privileges inside the target database.
	// In PostgreSQL 15+, public schema permissions are revoked from PUBLIC by default.
	// Granting ALL on SCHEMA public and default object privileges allows migrations to create tables and sequences seamlessly.
	if database != "" {
		schemaQuery := fmt.Sprintf(
			"GRANT ALL ON SCHEMA public TO \"%s\"; "+
				"GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO \"%s\"; "+
				"GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO \"%s\"; "+
				"GRANT ALL PRIVILEGES ON ALL ROUTINES IN SCHEMA public TO \"%s\"; "+
				"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO \"%s\"; "+
				"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO \"%s\"; "+
				"ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON ROUTINES TO \"%s\";",
			username, username, username, username, username, username, username,
		)
		if err := p.execDatabaseQuery(database, schemaQuery); err != nil {
			return fmt.Errorf("failed to grant schema privileges on %s to %s: %w", database, username, err)
		}
	}

	return nil
}

func (p *PostgresManager) RevokePrivileges(database, username string) error {
	query := fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE \"%s\" FROM \"%s\";", database, username)
	if err := p.execQuery(query); err != nil {
		return err
	}

	if database != "" {
		schemaQuery := fmt.Sprintf(
			"REVOKE ALL ON SCHEMA public FROM \"%s\"; "+
				"REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM \"%s\"; "+
				"REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM \"%s\"; "+
				"REVOKE ALL PRIVILEGES ON ALL ROUTINES IN SCHEMA public FROM \"%s\";",
			username, username, username, username,
		)
		_ = p.execDatabaseQuery(database, schemaQuery)
	}

	return nil
}

func (p *PostgresManager) ChangePassword(username, newPassword string) error {
	query := fmt.Sprintf("ALTER ROLE \"%s\" WITH PASSWORD '%s';", username, newPassword)
	return p.execQuery(query)
}

// ListTables returns metadata for all tables in a PostgreSQL database.
func (p *PostgresManager) ListTables(database string) ([]TableInfo, error) {
	if !validIdentifierRegex.MatchString(database) {
		return nil, fmt.Errorf("invalid database name: %s", database)
	}

	if p.isDev {
		return p.mockTables(), nil
	}

	query := `
SELECT 
    t.table_name,
    t.table_type,
    'PostgreSQL' AS engine,
    COALESCE(c.reltuples::bigint, 0) AS estimated_rows,
    COALESCE(pg_relation_size(quote_ident(t.table_schema) || '.' || quote_ident(t.table_name)), 0) AS data_bytes,
    COALESCE(pg_indexes_size(quote_ident(t.table_schema) || '.' || quote_ident(t.table_name)), 0) AS index_bytes,
    COALESCE(pg_total_relation_size(quote_ident(t.table_schema) || '.' || quote_ident(t.table_name)), 0) AS total_bytes,
    'UTF8' AS collation,
    COALESCE(obj_description(c.oid, 'pg_class'), '') AS comment
FROM information_schema.tables t
LEFT JOIN pg_namespace n ON n.nspname = t.table_schema
LEFT JOIN pg_class c ON c.relname = t.table_name AND c.relnamespace = n.oid
WHERE t.table_schema = 'public'
ORDER BY t.table_name ASC;`

	out, err := p.execQueryOutput(database, query)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var tables []TableInfo
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		cols := strings.Split(line, "\t")
		if len(cols) < 9 {
			continue
		}
		rows, _ := strconv.ParseInt(cols[3], 10, 64)
		dataLen, _ := strconv.ParseInt(cols[4], 10, 64)
		idxLen, _ := strconv.ParseInt(cols[5], 10, 64)
		totalSize, _ := strconv.ParseInt(cols[6], 10, 64)

		tables = append(tables, TableInfo{
			Name:        cols[0],
			Type:        cols[1],
			Engine:      cols[2],
			Rows:        rows,
			DataLength:  dataLen,
			IndexLength: idxLen,
			TotalSize:   totalSize,
			Collation:   cols[7],
			Comment:     cols[8],
		})
	}
	return tables, nil
}

// GetTableStructure returns columns, indexes, foreign keys, and DDL for a table in PostgreSQL.
func (p *PostgresManager) GetTableStructure(database, table string) (*TableStructure, error) {
	if !validIdentifierRegex.MatchString(database) || !validIdentifierRegex.MatchString(table) {
		return nil, fmt.Errorf("invalid database or table identifier")
	}

	if p.isDev {
		return p.mockTableStructure(table), nil
	}

	// 1. Columns
	colQuery := fmt.Sprintf(`
SELECT 
    c.column_name,
    c.ordinal_position,
    CASE 
        WHEN c.character_maximum_length IS NOT NULL THEN c.data_type || '(' || c.character_maximum_length || ')'
        WHEN c.numeric_precision IS NOT NULL AND c.numeric_scale IS NOT NULL THEN c.data_type || '(' || c.numeric_precision || ',' || c.numeric_scale || ')'
        ELSE c.data_type
    END AS formatted_type,
    c.data_type,
    CASE WHEN c.is_nullable = 'YES' THEN true ELSE false END AS is_nullable,
    COALESCE((
        SELECT tc.constraint_type
        FROM information_schema.table_constraints tc
        JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
        WHERE tc.table_schema = 'public' AND tc.table_name = c.table_name AND kcu.column_name = c.column_name
        LIMIT 1
    ), '') AS constraint_type,
    COALESCE(c.column_default, 'NULL_MARKER') AS column_default,
    COALESCE(col_description(pc.oid, c.ordinal_position), '') AS column_comment
FROM information_schema.columns c
LEFT JOIN pg_namespace pn ON pn.nspname = c.table_schema
LEFT JOIN pg_class pc ON pc.relname = c.table_name AND pc.relnamespace = pn.oid
WHERE c.table_schema = 'public' AND c.table_name = '%s'
ORDER BY c.ordinal_position ASC;`, table)

	colOut, err := p.execQueryOutput(database, colQuery)
	if err != nil {
		return nil, err
	}

	colLines := strings.Split(strings.TrimSpace(colOut), "\n")
	var columns []ColumnInfo
	for _, line := range colLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 8 {
			continue
		}
		pos, _ := strconv.Atoi(parts[1])
		isNullable := parts[4] == "t" || parts[4] == "true"
		keyStr := ""
		if parts[5] == "PRIMARY KEY" {
			keyStr = "PRI"
		} else if parts[5] == "UNIQUE" {
			keyStr = "UNI"
		} else if parts[5] == "FOREIGN KEY" {
			keyStr = "MUL"
		}

		var defVal *string
		if parts[6] != "NULL_MARKER" {
			val := parts[6]
			defVal = &val
		}

		columns = append(columns, ColumnInfo{
			Name:       parts[0],
			Position:   pos,
			Type:       parts[2],
			DataType:   parts[3],
			IsNullable: isNullable,
			Key:        keyStr,
			Default:    defVal,
			Extra:      "",
			Comment:    parts[7],
		})
	}

	// 2. Indexes
	idxQuery := fmt.Sprintf(`
SELECT 
    i.relname AS index_name,
    a.attname AS column_name,
    NOT idx.indisunique AS non_unique,
    idx.indisprimary AS is_primary,
    1 AS seq_in_index,
    am.amname AS index_type
FROM pg_index idx
JOIN pg_class c ON c.oid = idx.indrelid
JOIN pg_class i ON i.oid = idx.indexrelid
JOIN pg_am am ON am.oid = i.relam
JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = ANY(idx.indkey)
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'public' AND c.relname = '%s'
ORDER BY i.relname;`, table)

	idxOut, _ := p.execQueryOutput(database, idxQuery)
	idxLines := strings.Split(strings.TrimSpace(idxOut), "\n")
	var indexes []IndexInfo
	for _, line := range idxLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 6 {
			continue
		}
		nonUnique := parts[2] == "t" || parts[2] == "true"
		isPrimary := parts[3] == "t" || parts[3] == "true"
		seq, _ := strconv.Atoi(parts[4])

		indexes = append(indexes, IndexInfo{
			Name:       parts[0],
			Column:     parts[1],
			NonUnique:  nonUnique,
			IsPrimary:  isPrimary,
			SeqInIndex: seq,
			Type:       parts[5],
		})
	}

	// 3. Foreign Keys
	fkQuery := fmt.Sprintf(`
SELECT 
    tc.constraint_name,
    kcu.column_name,
    ccu.table_name AS referenced_table,
    ccu.column_name AS referenced_column
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name AND ccu.table_schema = tc.table_schema
WHERE tc.table_schema = 'public' AND tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = '%s';`, table)

	fkOut, _ := p.execQueryOutput(database, fkQuery)
	fkLines := strings.Split(strings.TrimSpace(fkOut), "\n")
	var fks []ForeignKeyInfo
	for _, line := range fkLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 4 {
			continue
		}
		fks = append(fks, ForeignKeyInfo{
			ConstraintName:   parts[0],
			Column:           parts[1],
			ReferencedTable:  parts[2],
			ReferencedColumn: parts[3],
		})
	}

	// 4. Generate DDL Representation
	ddlBuilder := strings.Builder{}
	ddlBuilder.WriteString(fmt.Sprintf("CREATE TABLE public.%s (\n", table))
	for i, col := range columns {
		nullStr := ""
		if !col.IsNullable {
			nullStr = " NOT NULL"
		}
		defStr := ""
		if col.Default != nil {
			defStr = fmt.Sprintf(" DEFAULT %s", *col.Default)
		}
		comma := ","
		if i == len(columns)-1 && len(indexes) == 0 {
			comma = ""
		}
		ddlBuilder.WriteString(fmt.Sprintf("    %s %s%s%s%s\n", col.Name, col.Type, nullStr, defStr, comma))
	}
	ddlBuilder.WriteString(");")

	return &TableStructure{
		TableName:       table,
		Columns:         columns,
		Indexes:         indexes,
		ForeignKeys:     fks,
		CreateStatement: ddlBuilder.String(),
	}, nil
}

// GetTableData returns paginated rows from a PostgreSQL table.
func (p *PostgresManager) GetTableData(database, table string, page, perPage int, sortField, sortDirection, search, searchColumn string) (*TableDataResult, error) {
	if !validIdentifierRegex.MatchString(database) || !validIdentifierRegex.MatchString(table) {
		return nil, fmt.Errorf("invalid database or table identifier")
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 500 {
		perPage = 50
	}
	offset := (page - 1) * perPage

	if p.isDev {
		return p.mockTableData(table, page, perPage), nil
	}

	// WHERE clause
	whereClause := ""
	if search != "" && searchColumn != "" && validIdentifierRegex.MatchString(searchColumn) {
		escapedSearch := strings.ReplaceAll(search, "'", "''")
		whereClause = fmt.Sprintf("WHERE \"%s\"::text ILIKE '%%%s%%'", searchColumn, escapedSearch)
	}

	// ORDER BY clause
	orderClause := ""
	if sortField != "" && validIdentifierRegex.MatchString(sortField) {
		dir := "ASC"
		if strings.ToUpper(sortDirection) == "DESC" {
			dir = "DESC"
		}
		orderClause = fmt.Sprintf("ORDER BY \"%s\" %s", sortField, dir)
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\" %s;", table, whereClause)
	countOut, err := p.execQueryOutput(database, countQuery)
	if err != nil {
		return nil, err
	}
	totalRows, _ := strconv.ParseInt(strings.TrimSpace(countOut), 10, 64)

	// Column list query
	colQuery := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_schema = 'public' AND table_name = '%s' ORDER BY ordinal_position;", table)
	colOut, _ := p.execQueryOutput(database, colQuery)
	colLines := strings.Split(strings.TrimSpace(colOut), "\n")
	var headers []string
	for _, l := range colLines {
		if strings.TrimSpace(l) != "" {
			headers = append(headers, strings.TrimSpace(l))
		}
	}

	// JSON rows query
	jsonQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json) FROM (SELECT * FROM "%s" %s %s LIMIT %d OFFSET %d) t;`, table, whereClause, orderClause, perPage, offset)
	jsonOut, err := p.execQueryOutput(database, jsonQuery)
	if err != nil {
		return nil, err
	}

	var rows []map[string]interface{}
	_ = json.Unmarshal([]byte(strings.TrimSpace(jsonOut)), &rows)

	totalPages := 0
	if totalRows > 0 {
		totalPages = int((totalRows + int64(perPage) - 1) / int64(perPage))
	}

	return &TableDataResult{
		TableName:  table,
		Columns:    headers,
		Rows:       rows,
		TotalRows:  totalRows,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// Dev mock helpers for PostgreSQL
func (p *PostgresManager) mockTables() []TableInfo {
	return []TableInfo{
		{Name: "accounts", Type: "BASE TABLE", Engine: "PostgreSQL", Rows: 5400, DataLength: 524288, IndexLength: 131072, TotalSize: 655360, Collation: "UTF8", Comment: "Organization billing accounts"},
		{Name: "subscriptions", Type: "BASE TABLE", Engine: "PostgreSQL", Rows: 3200, DataLength: 262144, IndexLength: 65536, TotalSize: 327680, Collation: "UTF8", Comment: "SaaS subscriptions"},
		{Name: "invoices", Type: "BASE TABLE", Engine: "PostgreSQL", Rows: 15400, DataLength: 2097152, IndexLength: 524288, TotalSize: 2621440, Collation: "UTF8", Comment: "Customer invoices"},
		{Name: "audit_events", Type: "BASE TABLE", Engine: "PostgreSQL", Rows: 84000, DataLength: 8388608, IndexLength: 2097152, TotalSize: 10485760, Collation: "UTF8", Comment: "System audit logs"},
	}
}

func (p *PostgresManager) mockTableStructure(table string) *TableStructure {
	defNow := "now()"
	defActive := "'active'::text"

	return &TableStructure{
		TableName: table,
		Columns: []ColumnInfo{
			{Name: "id", Position: 1, Type: "uuid", DataType: "uuid", IsNullable: false, Key: "PRI", Default: nil, Extra: "", Comment: "Primary Key"},
			{Name: "account_id", Position: 2, Type: "uuid", DataType: "uuid", IsNullable: false, Key: "MUL", Default: nil, Extra: "", Comment: "Account FK"},
			{Name: "status", Position: 3, Type: "character varying(32)", DataType: "varchar", IsNullable: false, Key: "", Default: &defActive, Extra: "", Comment: "State"},
			{Name: "metadata", Position: 4, Type: "jsonb", DataType: "jsonb", IsNullable: true, Key: "", Default: nil, Extra: "", Comment: "Arbitrary JSON payload"},
			{Name: "created_at", Position: 5, Type: "timestamp with time zone", DataType: "timestamp", IsNullable: false, Key: "", Default: &defNow, Extra: "", Comment: ""},
		},
		Indexes: []IndexInfo{
			{Name: fmt.Sprintf("%s_pkey", table), Column: "id", NonUnique: false, IsPrimary: true, SeqInIndex: 1, Type: "btree"},
			{Name: fmt.Sprintf("%s_account_id_idx", table), Column: "account_id", NonUnique: true, IsPrimary: false, SeqInIndex: 1, Type: "btree"},
		},
		ForeignKeys: []ForeignKeyInfo{
			{ConstraintName: fmt.Sprintf("%s_account_id_fkey", table), Column: "account_id", ReferencedTable: "accounts", ReferencedColumn: "id"},
		},
		CreateStatement: fmt.Sprintf("CREATE TABLE public.%s (\n    id uuid NOT NULL,\n    account_id uuid NOT NULL,\n    status character varying(32) NOT NULL DEFAULT 'active'::text,\n    metadata jsonb,\n    created_at timestamp with time zone NOT NULL DEFAULT now(),\n    CONSTRAINT %s_pkey PRIMARY KEY (id)\n);", table, table),
	}
}

func (p *PostgresManager) mockTableData(table string, page, perPage int) *TableDataResult {
	cols := []string{"id", "account_id", "status", "metadata", "created_at"}
	rows := []map[string]interface{}{
		{"id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "account_id": "c1f9b3e1-22e4-4e9b-b6d8-111122223333", "status": "active", "metadata": `{"plan": "pro", "seats": 10}`, "created_at": "2026-08-28T08:00:00Z"},
		{"id": "b1ffcd88-8d0a-4fe7-aa5c-5aa8ac270b22", "account_id": "c1f9b3e1-22e4-4e9b-b6d8-111122223333", "status": "active", "metadata": `{"plan": "enterprise", "seats": 50}`, "created_at": "2026-08-28T09:15:00Z"},
		{"id": "c2eedd77-7e09-4fd6-994b-4bb7bd160c33", "account_id": "d2e8a4f2-33f5-4f0c-a7e9-444455556666", "status": "trialing", "metadata": `{"plan": "starter", "seats": 2}`, "created_at": "2026-08-28T10:30:00Z"},
	}

	return &TableDataResult{
		TableName:  table,
		Columns:    cols,
		Rows:       rows,
		TotalRows:  3200,
		Page:       page,
		PerPage:    perPage,
		TotalPages: 64,
	}
}
