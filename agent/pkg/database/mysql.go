package database

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

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

	// Connect as root using defaults file or socket
	cmd := exec.Command("mysql", "--defaults-file=/etc/mysql/debian.cnf", "-e", query)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback without defaults-file for standard socket auth
		fallbackCmd := exec.Command("mysql", "-u", "root", "-e", query)
		fallbackOut, fallbackErr := fallbackCmd.CombinedOutput()
		if fallbackErr != nil {
			return fmt.Errorf("mysql error: %s: %s (fallback: %s)", err.Error(), strings.TrimSpace(string(out)), strings.TrimSpace(string(fallbackOut)))
		}
	}
	return nil
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

	// Create user with caching_sha2_password / mysql_native_password
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
