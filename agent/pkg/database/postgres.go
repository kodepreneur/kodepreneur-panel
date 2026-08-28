package database

import (
	"fmt"
	"os/exec"
	"runtime"
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

	// Connect as postgres system user via peer socket
	cmd := exec.Command("sudo", "-u", "postgres", "psql", "-c", query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("postgres error: %s: %s", err.Error(), strings.TrimSpace(string(out)))
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

	// PostgreSQL cannot execute CREATE DATABASE in a transaction block
	query := fmt.Sprintf("CREATE DATABASE \"%s\" ENCODING '%s' %s;", name, encoding, ownerClause)
	return p.execQuery(query)
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

	query := fmt.Sprintf("GRANT %s ON DATABASE \"%s\" TO \"%s\";", privileges, database, username)
	return p.execQuery(query)
}

func (p *PostgresManager) RevokePrivileges(database, username string) error {
	query := fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE \"%s\" FROM \"%s\";", database, username)
	return p.execQuery(query)
}

func (p *PostgresManager) ChangePassword(username, newPassword string) error {
	query := fmt.Sprintf("ALTER ROLE \"%s\" WITH PASSWORD '%s';", username, newPassword)
	return p.execQuery(query)
}
