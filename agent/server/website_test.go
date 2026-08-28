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
}
