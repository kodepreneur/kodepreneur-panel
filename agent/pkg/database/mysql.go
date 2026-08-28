package database

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var validIdentifierRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

type MysqlManager struct {
	isDev bool
}

func NewMysqlManager(isDev bool) *MysqlManager {
	return &MysqlManager{
		isDev: isDev || runtime.GOOS != "linux",
	}
}

func (m *MysqlManager) execQuery(query string) error {
	if m.isDev {
		return nil
	}

	cmd := exec.Command("mysql", "--defaults-file=/etc/mysql/debian.cnf", "-e", query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fallbackCmd := exec.Command("mysql", "-u", "root", "-e", query)
		fallbackOut, fallbackErr := fallbackCmd.CombinedOutput()
		if fallbackErr != nil {
			return fmt.Errorf("mysql error: %s: %s (fallback: %s)", err.Error(), strings.TrimSpace(string(out)), strings.TrimSpace(string(fallbackOut)))
		}
	}
	return nil
}

func (m *MysqlManager) execQueryOutput(query string) (string, error) {
	if m.isDev {
		return "", nil
	}

	cmd := exec.Command("mysql", "--defaults-file=/etc/mysql/debian.cnf", "-B", "-e", query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fallbackCmd := exec.Command("mysql", "-u", "root", "-B", "-e", query)
		fallbackOut, fallbackErr := fallbackCmd.CombinedOutput()
		if fallbackErr != nil {
			return "", fmt.Errorf("mysql error: %s: %s (fallback: %s)", err.Error(), strings.TrimSpace(string(out)), strings.TrimSpace(string(fallbackOut)))
		}
		return string(fallbackOut), nil
	}
	return string(out), nil
}

func (m *MysqlManager) CreateDatabase(name, charset, collation string) error {
	if charset == "" {
		charset = "utf8mb4"
	}
	if collation == "" {
		collation = "utf8mb4_unicode_ci"
	}

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET %s COLLATE %s;", name, charset, collation)
	return m.execQuery(query)
}

func (m *MysqlManager) DropDatabase(name string) error {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", name)
	return m.execQuery(query)
}

func (m *MysqlManager) CreateUser(username, host, password string) error {
	if host == "" {
		host = "localhost"
	}

	query := fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%s' IDENTIFIED BY '%s'; FLUSH PRIVILEGES;", username, host, password)
	return m.execQuery(query)
}

func (m *MysqlManager) DropUser(username, host string) error {
	if host == "" {
		host = "localhost"
	}

	query := fmt.Sprintf("DROP USER IF EXISTS '%s'@'%s'; FLUSH PRIVILEGES;", username, host)
	return m.execQuery(query)
}

func (m *MysqlManager) GrantPrivileges(database, username, host, privileges string) error {
	if host == "" {
		host = "localhost"
	}
	if privileges == "" || strings.ToLower(privileges) == "all" {
		privileges = "ALL PRIVILEGES"
	}

	query := fmt.Sprintf("GRANT %s ON `%s`.* TO '%s'@'%s'; FLUSH PRIVILEGES;", privileges, database, username, host)
	return m.execQuery(query)
}

func (m *MysqlManager) RevokePrivileges(database, username, host string) error {
	if host == "" {
		host = "localhost"
	}

	query := fmt.Sprintf("REVOKE ALL PRIVILEGES ON `%s`.* FROM '%s'@'%s'; FLUSH PRIVILEGES;", database, username, host)
	return m.execQuery(query)
}

func (m *MysqlManager) ChangePassword(username, host, newPassword string) error {
	if host == "" {
		host = "localhost"
	}

	query := fmt.Sprintf("ALTER USER '%s'@'%s' IDENTIFIED BY '%s'; FLUSH PRIVILEGES;", username, host, newPassword)
	return m.execQuery(query)
}

// ListTables returns metadata for all tables in a MySQL database.
func (m *MysqlManager) ListTables(database string) ([]TableInfo, error) {
	if !validIdentifierRegex.MatchString(database) {
		return nil, fmt.Errorf("invalid database name: %s", database)
	}

	if m.isDev {
		return m.mockTables(), nil
	}

	query := fmt.Sprintf("SELECT TABLE_NAME, TABLE_TYPE, COALESCE(ENGINE, ''), COALESCE(TABLE_ROWS, 0), COALESCE(DATA_LENGTH, 0), COALESCE(INDEX_LENGTH, 0), COALESCE(TABLE_COLLATION, ''), COALESCE(TABLE_COMMENT, '') FROM information_schema.TABLES WHERE TABLE_SCHEMA = '%s' ORDER BY TABLE_NAME ASC;", database)
	out, err := m.execQueryOutput(query)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) <= 1 {
		return []TableInfo{}, nil
	}

	var tables []TableInfo
	for _, line := range lines[1:] {
		cols := strings.Split(line, "\t")
		if len(cols) < 8 {
			continue
		}
		rows, _ := strconv.ParseInt(cols[3], 10, 64)
		dataLen, _ := strconv.ParseInt(cols[4], 10, 64)
		idxLen, _ := strconv.ParseInt(cols[5], 10, 64)

		tables = append(tables, TableInfo{
			Name:        cols[0],
			Type:        cols[1],
			Engine:      cols[2],
			Rows:        rows,
			DataLength:  dataLen,
			IndexLength: idxLen,
			TotalSize:   dataLen + idxLen,
			Collation:   cols[6],
			Comment:     cols[7],
		})
	}
	return tables, nil
}

// GetTableStructure returns columns, indexes, foreign keys, and DDL for a table.
func (m *MysqlManager) GetTableStructure(database, table string) (*TableStructure, error) {
	if !validIdentifierRegex.MatchString(database) || !validIdentifierRegex.MatchString(table) {
		return nil, fmt.Errorf("invalid database or table identifier")
	}

	if m.isDev {
		return m.mockTableStructure(table), nil
	}

	// 1. Fetch columns
	colQuery := fmt.Sprintf("SELECT COLUMN_NAME, ORDINAL_POSITION, COLUMN_TYPE, DATA_TYPE, IS_NULLABLE, COLUMN_KEY, COALESCE(COLUMN_DEFAULT, 'NULL_MARKER'), EXTRA, COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s' ORDER BY ORDINAL_POSITION ASC;", database, table)
	colOut, err := m.execQueryOutput(colQuery)
	if err != nil {
		return nil, err
	}

	colLines := strings.Split(strings.TrimSpace(colOut), "\n")
	var columns []ColumnInfo
	if len(colLines) > 1 {
		for _, line := range colLines[1:] {
			parts := strings.Split(line, "\t")
			if len(parts) < 9 {
				continue
			}
			pos, _ := strconv.Atoi(parts[1])
			isNullable := parts[4] == "YES"
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
				Key:        parts[5],
				Default:    defVal,
				Extra:      parts[7],
				Comment:    parts[8],
			})
		}
	}

	// 2. Fetch indexes
	idxQuery := fmt.Sprintf("SELECT INDEX_NAME, COLUMN_NAME, NON_UNIQUE, SEQ_IN_INDEX, INDEX_TYPE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s' ORDER BY INDEX_NAME, SEQ_IN_INDEX;", database, table)
	idxOut, _ := m.execQueryOutput(idxQuery)
	idxLines := strings.Split(strings.TrimSpace(idxOut), "\n")
	var indexes []IndexInfo
	if len(idxLines) > 1 {
		for _, line := range idxLines[1:] {
			parts := strings.Split(line, "\t")
			if len(parts) < 5 {
				continue
			}
			nonUnique := parts[2] == "1"
			seq, _ := strconv.Atoi(parts[3])
			isPrimary := parts[0] == "PRIMARY"

			indexes = append(indexes, IndexInfo{
				Name:       parts[0],
				Column:     parts[1],
				NonUnique:  nonUnique,
				IsPrimary:  isPrimary,
				SeqInIndex: seq,
				Type:       parts[4],
			})
		}
	}

	// 3. Fetch foreign keys
	fkQuery := fmt.Sprintf("SELECT k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE k WHERE k.TABLE_SCHEMA = '%s' AND k.TABLE_NAME = '%s' AND k.REFERENCED_TABLE_NAME IS NOT NULL;", database, table)
	fkOut, _ := m.execQueryOutput(fkQuery)
	fkLines := strings.Split(strings.TrimSpace(fkOut), "\n")
	var fks []ForeignKeyInfo
	if len(fkLines) > 1 {
		for _, line := range fkLines[1:] {
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
	}

	// 4. Fetch Create Statement (DDL)
	ddlQuery := fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`;", database, table)
	ddlOut, _ := m.execQueryOutput(ddlQuery)
	ddlLines := strings.Split(strings.TrimSpace(ddlOut), "\n")
	createStatement := ""
	if len(ddlLines) > 1 {
		parts := strings.Split(ddlLines[1], "\t")
		if len(parts) >= 2 {
			createStatement = parts[1]
		}
	}

	return &TableStructure{
		TableName:       table,
		Columns:         columns,
		Indexes:         indexes,
		ForeignKeys:     fks,
		CreateStatement: createStatement,
	}, nil
}

// GetTableData returns paginated rows from a table.
func (m *MysqlManager) GetTableData(database, table string, page, perPage int, sortField, sortDirection, search, searchColumn string) (*TableDataResult, error) {
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

	if m.isDev {
		return m.mockTableData(table, page, perPage), nil
	}

	// WHERE clause
	whereClause := ""
	if search != "" && searchColumn != "" && validIdentifierRegex.MatchString(searchColumn) {
		escapedSearch := strings.ReplaceAll(search, "'", "''")
		whereClause = fmt.Sprintf("WHERE `%s` LIKE '%%%s%%'", searchColumn, escapedSearch)
	}

	// ORDER BY clause
	orderClause := ""
	if sortField != "" && validIdentifierRegex.MatchString(sortField) {
		dir := "ASC"
		if strings.ToUpper(sortDirection) == "DESC" {
			dir = "DESC"
		}
		orderClause = fmt.Sprintf("ORDER BY `%s` %s", sortField, dir)
	}

	// Total count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s` %s;", database, table, whereClause)
	countOut, err := m.execQueryOutput(countQuery)
	if err != nil {
		return nil, err
	}
	countLines := strings.Split(strings.TrimSpace(countOut), "\n")
	var totalRows int64 = 0
	if len(countLines) > 1 {
		totalRows, _ = strconv.ParseInt(countLines[1], 10, 64)
	}

	// Data query
	dataQuery := fmt.Sprintf("SELECT * FROM `%s`.`%s` %s %s LIMIT %d OFFSET %d;", database, table, whereClause, orderClause, perPage, offset)
	dataOut, err := m.execQueryOutput(dataQuery)
	if err != nil {
		return nil, err
	}

	dataLines := strings.Split(strings.TrimSpace(dataOut), "\n")
	var headers []string
	var rows []map[string]interface{}

	if len(dataLines) > 0 && dataLines[0] != "" {
		headers = strings.Split(dataLines[0], "\t")
		if len(dataLines) > 1 {
			for _, line := range dataLines[1:] {
				cols := strings.Split(line, "\t")
				rowMap := make(map[string]interface{})
				for i, h := range headers {
					if i < len(cols) {
						if cols[i] == "NULL" {
							rowMap[h] = nil
						} else {
							rowMap[h] = cols[i]
						}
					} else {
						rowMap[h] = nil
					}
				}
				rows = append(rows, rowMap)
			}
		}
	}

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

// Dev mock helpers for MySQL
func (m *MysqlManager) mockTables() []TableInfo {
	return []TableInfo{
		{Name: "users", Type: "BASE TABLE", Engine: "InnoDB", Rows: 1420, DataLength: 163840, IndexLength: 65536, TotalSize: 229376, Collation: "utf8mb4_unicode_ci", Comment: "System and app users"},
		{Name: "orders", Type: "BASE TABLE", Engine: "InnoDB", Rows: 8930, DataLength: 1048576, IndexLength: 393216, TotalSize: 1441792, Collation: "utf8mb4_unicode_ci", Comment: "Customer orders"},
		{Name: "order_items", Type: "BASE TABLE", Engine: "InnoDB", Rows: 24500, DataLength: 2097152, IndexLength: 786432, TotalSize: 2883584, Collation: "utf8mb4_unicode_ci", Comment: "Items per order"},
		{Name: "products", Type: "BASE TABLE", Engine: "InnoDB", Rows: 340, DataLength: 98304, IndexLength: 32768, TotalSize: 131072, Collation: "utf8mb4_unicode_ci", Comment: "Product catalog"},
		{Name: "settings", Type: "BASE TABLE", Engine: "InnoDB", Rows: 28, DataLength: 16384, IndexLength: 16384, TotalSize: 32768, Collation: "utf8mb4_unicode_ci", Comment: "Configuration settings"},
		{Name: "migrations", Type: "BASE TABLE", Engine: "InnoDB", Rows: 18, DataLength: 16384, IndexLength: 0, TotalSize: 16384, Collation: "utf8mb4_unicode_ci", Comment: "Database migration ledger"},
	}
}

func (m *MysqlManager) mockTableStructure(table string) *TableStructure {
	defZero := "0"
	defOne := "1"
	defCurrentTime := "CURRENT_TIMESTAMP"

	switch table {
	case "orders":
		return &TableStructure{
			TableName: "orders",
			Columns: []ColumnInfo{
				{Name: "id", Position: 1, Type: "bigint unsigned", DataType: "bigint", IsNullable: false, Key: "PRI", Default: nil, Extra: "auto_increment", Comment: "Primary Key"},
				{Name: "user_id", Position: 2, Type: "bigint unsigned", DataType: "bigint", IsNullable: false, Key: "MUL", Default: nil, Extra: "", Comment: "Customer FK"},
				{Name: "order_number", Position: 3, Type: "varchar(64)", DataType: "varchar", IsNullable: false, Key: "UNI", Default: nil, Extra: "", Comment: "Invoice Ref"},
				{Name: "total_amount", Position: 4, Type: "decimal(12,2)", DataType: "decimal", IsNullable: false, Key: "", Default: &defZero, Extra: "", Comment: "Total in USD"},
				{Name: "status", Position: 5, Type: "varchar(32)", DataType: "varchar", IsNullable: false, Key: "", Default: nil, Extra: "", Comment: "Order Status"},
				{Name: "created_at", Position: 6, Type: "timestamp", DataType: "timestamp", IsNullable: true, Key: "", Default: &defCurrentTime, Extra: "", Comment: ""},
				{Name: "updated_at", Position: 7, Type: "timestamp", DataType: "timestamp", IsNullable: true, Key: "", Default: &defCurrentTime, Extra: "on update CURRENT_TIMESTAMP", Comment: ""},
			},
			Indexes: []IndexInfo{
				{Name: "PRIMARY", Column: "id", NonUnique: false, IsPrimary: true, SeqInIndex: 1, Type: "BTREE"},
				{Name: "orders_order_number_unique", Column: "order_number", NonUnique: false, IsPrimary: false, SeqInIndex: 1, Type: "BTREE"},
				{Name: "orders_user_id_index", Column: "user_id", NonUnique: true, IsPrimary: false, SeqInIndex: 1, Type: "BTREE"},
			},
			ForeignKeys: []ForeignKeyInfo{
				{ConstraintName: "orders_user_id_foreign", Column: "user_id", ReferencedTable: "users", ReferencedColumn: "id"},
			},
			CreateStatement: "CREATE TABLE `orders` (\n  `id` bigint unsigned NOT NULL AUTO_INCREMENT,\n  `user_id` bigint unsigned NOT NULL,\n  `order_number` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `total_amount` decimal(12,2) NOT NULL DEFAULT '0.00',\n  `status` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `orders_order_number_unique` (`order_number`),\n  KEY `orders_user_id_index` (`user_id`),\n  CONSTRAINT `orders_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
		}
	default:
		return &TableStructure{
			TableName: table,
			Columns: []ColumnInfo{
				{Name: "id", Position: 1, Type: "bigint unsigned", DataType: "bigint", IsNullable: false, Key: "PRI", Default: nil, Extra: "auto_increment", Comment: "Unique Identifier"},
				{Name: "name", Position: 2, Type: "varchar(255)", DataType: "varchar", IsNullable: false, Key: "", Default: nil, Extra: "", Comment: "Name field"},
				{Name: "email", Position: 3, Type: "varchar(255)", DataType: "varchar", IsNullable: false, Key: "UNI", Default: nil, Extra: "", Comment: "Email Address"},
				{Name: "is_active", Position: 4, Type: "tinyint(1)", DataType: "tinyint", IsNullable: false, Key: "", Default: &defOne, Extra: "", Comment: "Active Flag"},
				{Name: "created_at", Position: 5, Type: "timestamp", DataType: "timestamp", IsNullable: true, Key: "", Default: &defCurrentTime, Extra: "", Comment: ""},
			},
			Indexes: []IndexInfo{
				{Name: "PRIMARY", Column: "id", NonUnique: false, IsPrimary: true, SeqInIndex: 1, Type: "BTREE"},
				{Name: "users_email_unique", Column: "email", NonUnique: false, IsPrimary: false, SeqInIndex: 1, Type: "BTREE"},
			},
			ForeignKeys: []ForeignKeyInfo{},
			CreateStatement: fmt.Sprintf("CREATE TABLE `%s` (\n  `id` bigint unsigned NOT NULL AUTO_INCREMENT,\n  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,\n  `is_active` tinyint(1) NOT NULL DEFAULT '1',\n  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `users_email_unique` (`email`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;", table),
		}
	}
}

func (m *MysqlManager) mockTableData(table string, page, perPage int) *TableDataResult {
	if table == "orders" {
		cols := []string{"id", "user_id", "order_number", "total_amount", "status", "created_at"}
		rows := []map[string]interface{}{
			{"id": "1001", "user_id": "42", "order_number": "ORD-2026-001", "total_amount": "149.99", "status": "completed", "created_at": "2026-08-28 10:15:00"},
			{"id": "1002", "user_id": "19", "order_number": "ORD-2026-002", "total_amount": "89.50", "status": "processing", "created_at": "2026-08-28 11:20:00"},
			{"id": "1003", "user_id": "77", "order_number": "ORD-2026-003", "total_amount": "320.00", "status": "completed", "created_at": "2026-08-28 12:05:00"},
			{"id": "1004", "user_id": "105", "order_number": "ORD-2026-004", "total_amount": "24.95", "status": "pending", "created_at": "2026-08-28 13:45:00"},
			{"id": "1005", "user_id": "31", "order_number": "ORD-2026-005", "total_amount": "512.10", "status": "completed", "created_at": "2026-08-28 14:10:00"},
		}
		return &TableDataResult{
			TableName:  table,
			Columns:    cols,
			Rows:       rows,
			TotalRows:  8930,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 179,
		}
	}

	cols := []string{"id", "name", "email", "is_active", "created_at"}
	rows := []map[string]interface{}{
		{"id": "1", "name": "Admin User", "email": "admin@kodepreneur.com", "is_active": "1", "created_at": "2026-01-01 08:00:00"},
		{"id": "2", "name": "Sarah Connor", "email": "sarah@cyberdyne.io", "is_active": "1", "created_at": "2026-01-15 09:30:00"},
		{"id": "3", "name": "Alex Murphy", "email": "alex.murphy@ocp.corp", "is_active": "1", "created_at": "2026-02-10 14:22:00"},
		{"id": "4", "name": "Thomas Anderson", "email": "neo@matrix.net", "is_active": "0", "created_at": "2026-03-01 19:45:00"},
		{"id": "5", "name": "Ellen Ripley", "email": "ripley@weyland.org", "is_active": "1", "created_at": "2026-04-12 11:15:00"},
	}

	return &TableDataResult{
		TableName:  table,
		Columns:    cols,
		Rows:       rows,
		TotalRows:  1420,
		Page:       page,
		PerPage:    perPage,
		TotalPages: 29,
	}
}
