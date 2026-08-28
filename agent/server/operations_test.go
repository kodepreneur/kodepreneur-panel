package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kodepreneur/agent/config"
)

func TestOperationsEndpoints(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Environment.IsDev = true
	router := NewRouter(cfg)
	handler := router.Handler()

	// 1. POST /api/v1/deployments/execute
	deployPayload := []byte(`{
		"system_user": "kp_test",
		"working_dir": "/var/www/test.com",
		"branch": "main",
		"commands": ["git pull", "composer install"]
	}`)
	req := httptest.NewRequest("POST", "/api/v1/deployments/execute", bytes.NewReader(deployPayload))
	signRequest(req, deployPayload, cfg.Security.SecretKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for deploy execute, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. POST /api/v1/cron/sync
	cronPayload := []byte(`{
		"system_user": "kp_test",
		"jobs": [{"schedule": "* * * * *", "command": "php artisan schedule:run", "is_active": true}]
	}`)
	req = httptest.NewRequest("POST", "/api/v1/cron/sync", bytes.NewReader(cronPayload))
	signRequest(req, cronPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for cron sync, got %d: %s", rec.Code, rec.Body.String())
	}

	// 3. POST /api/v1/processes
	procPayload := []byte(`{
		"id": "app-worker",
		"name": "App Worker",
		"command": "php artisan queue:work",
		"system_user": "kp_test",
		"working_dir": "/var/www/test.com"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/processes", bytes.NewReader(procPayload))
	signRequest(req, procPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for process creation, got %d: %s", rec.Code, rec.Body.String())
	}

	// 4. POST /api/v1/firewall/rules
	fwPayload := []byte(`{
		"port": 8080,
		"protocol": "tcp",
		"action": "allow"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/firewall/rules", bytes.NewReader(fwPayload))
	signRequest(req, fwPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for firewall rule, got %d: %s", rec.Code, rec.Body.String())
	}

	// 5. POST /api/v1/files/write
	fileWritePayload := []byte(`{
		"base_path": "/tmp/kodepreneur/test_files",
		"relative_path": "info.txt",
		"content": "Kodepreneur Agent v1.0"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/files/write", bytes.NewReader(fileWritePayload))
	signRequest(req, fileWritePayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for file write, got %d: %s", rec.Code, rec.Body.String())
	}

	// 6. POST /api/v1/files/read
	fileReadPayload := []byte(`{
		"base_path": "/tmp/kodepreneur/test_files",
		"relative_path": "info.txt"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/files/read", bytes.NewReader(fileReadPayload))
	signRequest(req, fileReadPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for file read, got %d: %s", rec.Code, rec.Body.String())
	}
}
