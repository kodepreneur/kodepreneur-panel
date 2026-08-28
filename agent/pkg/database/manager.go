package database

import (
	"fmt"
	"strings"
)

type Manager struct {
	mysql    *MysqlManager
	postgres *PostgresManager
	isDev    bool
}

func NewManager(isDev bool) *Manager {
	return &Manager{
		mysql:    NewMysqlManager(isDev),
		postgres: NewPostgresManager(isDev),
		isDev:    isDev,
	}
}

func (m *Manager) CreateDatabase(engine, name, charset, collation string) error {
	switch strings.ToLower(engine) {
	case "mysql", "mariadb":
		return m.mysql.CreateDatabase(name, charset, collation)
	case "postgres", "postgresql":
		return m.postgres.CreateDatabase(name, charset, "")
	default:
		return fmt.Errorf("unsupported database engine: %s", engine)
	}
}

func (m *Manager) DropDatabase(engine, name string) error {
	switch strings.ToLower(engine) {
	case "mysql", "mariadb":
		return m.mysql.DropDatabase(name)
	case "postgres", "postgresql":
		return m.postgres.DropDatabase(name)
	default:
		return fmt.Errorf("unsupported database engine: %s", engine)
	}
}

func (m *Manager) CreateUser(engine, username, host, password string) error {
	switch strings.ToLower(engine) {
	case "mysql", "mariadb":
		return m.mysql.CreateUser(username, host, password)
	case "postgres", "postgresql":
		return m.postgres.CreateUser(username, password)
	default:
		return fmt.Errorf("unsupported database engine: %s", engine)
	}
}

func (m *Manager) DropUser(engine, username, host string) error {
	switch strings.ToLower(engine) {
	case "mysql", "mariadb":
		return m.mysql.DropUser(username, host)
	case "postgres", "postgresql":
		return m.postgres.DropUser(username)
	default:
		return fmt.Errorf("unsupported database engine: %s", engine)
	}
}

func (m *Manager) GrantPrivileges(engine, database, username, host, privileges string) error {
	switch strings.ToLower(engine) {
	case "mysql", "mariadb":
		return m.mysql.GrantPrivileges(database, username, host, privileges)
	case "postgres", "postgresql":
		return m.postgres.GrantPrivileges(database, username, privileges)
	default:
		return fmt.Errorf("unsupported database engine: %s", engine)
	}
}

func (m *Manager) RevokePrivileges(engine, database, username, host string) error {
	switch strings.ToLower(engine) {
	case "mysql", "mariadb":
		return m.mysql.RevokePrivileges(database, username, host)
	case "postgres", "postgresql":
		return m.postgres.RevokePrivileges(database, username)
	default:
		return fmt.Errorf("unsupported database engine: %s", engine)
	}
}

func (m *Manager) ChangePassword(engine, username, host, newPassword string) error {
	switch strings.ToLower(engine) {
	case "mysql", "mariadb":
		return m.mysql.ChangePassword(username, host, newPassword)
	case "postgres", "postgresql":
		return m.postgres.ChangePassword(username, newPassword)
	default:
		return fmt.Errorf("unsupported database engine: %s", engine)
	}
}
