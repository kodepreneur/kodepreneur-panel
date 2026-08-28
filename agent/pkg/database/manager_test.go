package database

import (
	"testing"
)

func TestDatabaseManagerLifecycle(t *testing.T) {
	mgr := NewManager(true)

	// 1. MySQL Lifecycle
	if err := mgr.CreateDatabase("mysql", "test_app_db", "utf8mb4", "utf8mb4_unicode_ci"); err != nil {
		t.Fatalf("MySQL CreateDatabase failed: %v", err)
	}

	if err := mgr.CreateUser("mysql", "test_user", "localhost", "SecretPass123!"); err != nil {
		t.Fatalf("MySQL CreateUser failed: %v", err)
	}

	if err := mgr.GrantPrivileges("mysql", "test_app_db", "test_user", "localhost", "ALL PRIVILEGES"); err != nil {
		t.Fatalf("MySQL GrantPrivileges failed: %v", err)
	}

	if err := mgr.ChangePassword("mysql", "test_user", "localhost", "NewSecretPass456!"); err != nil {
		t.Fatalf("MySQL ChangePassword failed: %v", err)
	}

	if err := mgr.RevokePrivileges("mysql", "test_app_db", "test_user", "localhost"); err != nil {
		t.Fatalf("MySQL RevokePrivileges failed: %v", err)
	}

	if err := mgr.DropUser("mysql", "test_user", "localhost"); err != nil {
		t.Fatalf("MySQL DropUser failed: %v", err)
	}

	if err := mgr.DropDatabase("mysql", "test_app_db"); err != nil {
		t.Fatalf("MySQL DropDatabase failed: %v", err)
	}

	// 2. PostgreSQL Lifecycle
	if err := mgr.CreateDatabase("postgres", "pg_app_db", "UTF8", ""); err != nil {
		t.Fatalf("Postgres CreateDatabase failed: %v", err)
	}

	if err := mgr.CreateUser("postgres", "pg_user", "", "PgSecretPass123!"); err != nil {
		t.Fatalf("Postgres CreateUser failed: %v", err)
	}

	if err := mgr.GrantPrivileges("postgres", "pg_app_db", "pg_user", "", "ALL PRIVILEGES"); err != nil {
		t.Fatalf("Postgres GrantPrivileges failed: %v", err)
	}

	if err := mgr.ChangePassword("postgres", "pg_user", "", "NewPgSecretPass456!"); err != nil {
		t.Fatalf("Postgres ChangePassword failed: %v", err)
	}

	if err := mgr.RevokePrivileges("postgres", "pg_app_db", "pg_user", ""); err != nil {
		t.Fatalf("Postgres RevokePrivileges failed: %v", err)
	}

	if err := mgr.DropUser("postgres", "pg_user", ""); err != nil {
		t.Fatalf("Postgres DropUser failed: %v", err)
	}

	if err := mgr.DropDatabase("postgres", "pg_app_db"); err != nil {
		t.Fatalf("Postgres DropDatabase failed: %v", err)
	}
}
