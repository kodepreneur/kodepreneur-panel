package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kodepreneur/agent/config"
)

func signRequest(req *http.Request, body []byte, secret string) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := fmt.Sprintf("nonce-%d-%s", time.Now().UnixNano(), req.URL.Path)

	bodyHash := sha256.Sum256(body)
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	message := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", req.Method, req.URL.Path, timestamp, nonce, bodyHashHex)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	req.Header.Set("X-Kodepreneur-Timestamp", timestamp)
	req.Header.Set("X-Kodepreneur-Nonce", nonce)
	req.Header.Set("X-Kodepreneur-Signature", signature)
	req.Header.Set("Content-Type", "application/json")
}

func TestWebsiteEndpoints(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Environment.IsDev = true
	router := NewRouter(cfg)
	handler := router.Handler()

	// 1. POST /api/v1/websites
	createPayload := []byte(`{
		"domain": "testsite.org",
		"php_version": "8.3",
		"document_root": "/tmp/kodepreneur/testsite/public",
		"system_user": "kp_testsite"
	}`)

	req := httptest.NewRequest("POST", "/api/v1/websites", bytes.NewReader(createPayload))
	signRequest(req, createPayload, cfg.Security.SecretKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for website creation, got %d: %s", rec.Code, rec.Body.String())
	}

	// 2. PUT /api/v1/websites/testsite.org/php
	switchPayload := []byte(`{
		"new_php_version": "8.4",
		"old_php_version": "8.3",
		"system_user": "kp_testsite",
		"document_root": "/tmp/kodepreneur/testsite/public"
	}`)

	req = httptest.NewRequest("PUT", "/api/v1/websites/testsite.org/php", bytes.NewReader(switchPayload))
	signRequest(req, switchPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for PHP switch, got %d: %s", rec.Code, rec.Body.String())
	}

	// 3. POST /api/v1/ssl/issue
	sslPayload := []byte(`{
		"domain": "testsite.org",
		"document_root": "/tmp/kodepreneur/testsite/public",
		"php_version": "8.4",
		"system_user": "kp_testsite"
	}`)

	req = httptest.NewRequest("POST", "/api/v1/ssl/issue", bytes.NewReader(sslPayload))
	signRequest(req, sslPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for SSL issuance, got %d: %s", rec.Code, rec.Body.String())
	}

	// 4. GET /api/v1/websites/testsite.org/logs/access
	req = httptest.NewRequest("GET", "/api/v1/websites/testsite.org/logs/access?lines=10", nil)
	signRequest(req, nil, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for log read, got %d: %s", rec.Code, rec.Body.String())
	}

	var logResp struct {
		Success bool `json:"success"`
		Data    struct {
			Lines []string `json:"lines"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&logResp); err != nil {
		t.Fatalf("Failed to decode log response: %v", err)
	}
	if len(logResp.Data.Lines) == 0 {
		t.Errorf("Expected log entries in response")
	}

	// 5. DELETE /api/v1/websites/testsite.org
	deletePayload := []byte(`{"php_version": "8.4", "system_user": "kp_testsite"}`)
	req = httptest.NewRequest("DELETE", "/api/v1/websites/testsite.org", bytes.NewReader(deletePayload))
	signRequest(req, deletePayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for website deletion, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestWebsiteDeploymentSources(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Environment.IsDev = true
	router := NewRouter(cfg)
	handler := router.Handler()

	// 1. Create with Git deployment source
	gitPayload := []byte(`{
		"domain": "git-site.com",
		"php_version": "8.3",
		"deployment_source": "git",
		"git_repository": "https://github.com/example/laravel-demo.git",
		"git_branch": "main",
		"project_type": "laravel"
	}`)
	req := httptest.NewRequest("POST", "/api/v1/websites", bytes.NewReader(gitPayload))
	signRequest(req, gitPayload, cfg.Security.SecretKey)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for git website creation, got %d: %s", rec.Code, rec.Body.String())
	}

	var gitResp struct {
		Success bool `json:"success"`
		Data    struct {
			DocumentRoot string `json:"document_root"`
			IsLaravel    bool   `json:"is_laravel"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&gitResp); err != nil {
		t.Fatalf("Failed to decode git response: %v", err)
	}
	if !gitResp.Data.IsLaravel {
		t.Errorf("Expected is_laravel to be true")
	}
	if gitResp.Data.DocumentRoot != "/var/www/git-site.com/public" {
		t.Errorf("Expected document_root to be /var/www/git-site.com/public, got %s", gitResp.Data.DocumentRoot)
	}

	// 2. Create with generic PHP Zip deployment source (no public dir)
	zipPayload := []byte(`{
		"domain": "zip-site.com",
		"php_version": "8.3",
		"deployment_source": "zip",
		"project_type": "generic_php"
	}`)
	req = httptest.NewRequest("POST", "/api/v1/websites", bytes.NewReader(zipPayload))
	signRequest(req, zipPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for zip website creation, got %d: %s", rec.Code, rec.Body.String())
	}

	var zipResp struct {
		Success bool `json:"success"`
		Data    struct {
			DocumentRoot string `json:"document_root"`
			IsLaravel    bool   `json:"is_laravel"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&zipResp); err != nil {
		t.Fatalf("Failed to decode zip response: %v", err)
	}
	if zipResp.Data.DocumentRoot != "/var/www/zip-site.com" {
		t.Errorf("Expected document_root to be /var/www/zip-site.com, got %s", zipResp.Data.DocumentRoot)
	}

	// 3. Create Laravel site with automated post-setup and .env
	laravelPayload := []byte(`{
		"domain": "laravel-app.com",
		"php_version": "8.4",
		"deployment_source": "empty",
		"project_type": "laravel",
		"laravel_setup": {
			"enabled": true,
			"setup_env": true,
			"env_vars": {
				"APP_NAME": "MyShop",
				"DB_CONNECTION": "mysql",
				"DB_DATABASE": "db_myshop",
				"DB_USERNAME": "u_myshop",
				"DB_PASSWORD": "secretpassword123"
			},
			"run_composer": true,
			"run_key_generate": true,
			"run_migrations": true,
			"run_npm_build": true,
			"run_optimize": true
		}
	}`)
	req = httptest.NewRequest("POST", "/api/v1/websites", bytes.NewReader(laravelPayload))
	signRequest(req, laravelPayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for Laravel automated setup website, got %d: %s", rec.Code, rec.Body.String())
	}

	var laravelResp struct {
		Success bool `json:"success"`
		Data    struct {
			DocumentRoot string `json:"document_root"`
			IsLaravel    bool   `json:"is_laravel"`
			SetupResult  *struct {
				Success   bool   `json:"success"`
				LogOutput string `json:"log_output"`
			} `json:"setup_result"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&laravelResp); err != nil {
		t.Fatalf("Failed to decode laravel response: %v", err)
	}
	if !laravelResp.Data.IsLaravel {
		t.Errorf("Expected is_laravel to be true")
	}
	if laravelResp.Data.SetupResult == nil || !laravelResp.Data.SetupResult.Success {
		t.Errorf("Expected setup_result to be successful")
	}

	// 4. Create Laravel site with string/number booleans (e.g. "1", "0", "true", "false", 1, 0)
	flexiblePayload := []byte(`{
		"domain": "flexible-site.com",
		"php_version": "8.3",
		"deployment_source": "empty",
		"project_type": "laravel",
		"ssl_enabled": "0",
		"force_https": "false",
		"laravel_setup": {
			"enabled": "1",
			"setup_env": "true",
			"env_vars": {
				"APP_NAME": "FlexApp"
			},
			"run_composer": 1,
			"run_key_generate": "1",
			"run_migrations": "1",
			"run_seeders": "0",
			"run_npm_build": "true",
			"run_optimize": 1
		}
	}`)
	req = httptest.NewRequest("POST", "/api/v1/websites", bytes.NewReader(flexiblePayload))
	signRequest(req, flexiblePayload, cfg.Security.SecretKey)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for flexible booleans payload, got %d: %s", rec.Code, rec.Body.String())
	}
}
