package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kodepreneur/agent/config"
)

func TestDatabaseEndpoints(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Environment.IsDev = true
	router := NewRouter(cfg)
	handler := router.Handler()

	// 1. POST /api/v1/databases (Create Database)
	createDbPayload := []byte(`{
		"engine": "mysql",
		"name": "ecommerce_prod",
		"charset": "utf8mb4",
		"collation": "utf8mb4_unicode_ci"
	}`)
	req := httptest.NewRequest("POST", "/api/v1/databases", bytes.NewReader(createDbPayload))
	signRequest(req, createDbPayload, cfg.Security.SecretKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for database creation, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. POST /api/v1/databases/users (Create User)
	createUserPayload := []byte(`{
		"engine": "mysql",
		"username": "ecom_user",
		"host": "localhost",
		"password": "SuperSecretPassword123!"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/databases/users", bytes.NewReader(createUserPayload))
	signRequest(req, createUserPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for database user creation, got %d: %s", rec.Code, rec.Body.String())
	}

	// 3. POST /api/v1/databases/grants (Grant Privileges)
	grantPayload := []byte(`{
		"engine": "mysql",
		"database": "ecommerce_prod",
		"username": "ecom_user",
		"host": "localhost",
		"permissions": "ALL"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/databases/grants", bytes.NewReader(grantPayload))
	signRequest(req, grantPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for grant privileges, got %d: %s", rec.Code, rec.Body.String())
	}

	// 4. POST /api/v1/databases/users/password (Reset Password)
	resetPassPayload := []byte(`{
		"engine": "mysql",
		"username": "ecom_user",
		"host": "localhost",
		"new_password": "NewSuperSecretPassword456!"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/databases/users/password", bytes.NewReader(resetPassPayload))
	signRequest(req, resetPassPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for password reset, got %d: %s", rec.Code, rec.Body.String())
	}

	// 5. GET /api/v1/databases/mysql/ecommerce_prod/tables (List Tables)
	req = httptest.NewRequest("GET", "/api/v1/databases/mysql/ecommerce_prod/tables", nil)
	signRequest(req, nil, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for list tables, got %d: %s", rec.Code, rec.Body.String())
	}

	// 6. GET /api/v1/databases/mysql/ecommerce_prod/tables/orders/structure (Table Structure)
	req = httptest.NewRequest("GET", "/api/v1/databases/mysql/ecommerce_prod/tables/orders/structure", nil)
	signRequest(req, nil, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for table structure, got %d: %s", rec.Code, rec.Body.String())
	}

	// 7. GET /api/v1/databases/mysql/ecommerce_prod/tables/orders/data (Table Data)
	req = httptest.NewRequest("GET", "/api/v1/databases/mysql/ecommerce_prod/tables/orders/data?page=1&per_page=10", nil)
	signRequest(req, nil, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for table data, got %d: %s", rec.Code, rec.Body.String())
	}

	// 8. GET /api/v1/databases/postgres/app_pg/tables (Postgres List Tables)
	req = httptest.NewRequest("GET", "/api/v1/databases/postgres/app_pg/tables", nil)
	signRequest(req, nil, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for postgres list tables, got %d: %s", rec.Code, rec.Body.String())
	}

	// 9. DELETE /api/v1/databases/users/mysql/ecom_user
	req = httptest.NewRequest("DELETE", "/api/v1/databases/users/mysql/ecom_user?host=localhost", nil)
	signRequest(req, nil, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for user delete, got %d: %s", rec.Code, rec.Body.String())
	}

	// 10. DELETE /api/v1/databases/mysql/ecommerce_prod
	req = httptest.NewRequest("DELETE", "/api/v1/databases/mysql/ecommerce_prod", nil)
	signRequest(req, nil, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for database delete, got %d: %s", rec.Code, rec.Body.String())
	}
}
