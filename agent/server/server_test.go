package server_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/kodepreneur/agent/config"
	"github.com/kodepreneur/agent/server"
)

func createAuthenticatedRequest(method, path string, body []byte, secret string) *http.Request {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}
	req := httptest.NewRequest(method, path, bodyReader)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)

	hasher := sha256.New()
	hasher.Write(body)
	bodyHash := hex.EncodeToString(hasher.Sum(nil))

	msg := method + "\n" + path + "\n" + timestamp + "\n" + nonce + "\n" + bodyHash
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	sig := hex.EncodeToString(mac.Sum(nil))

	req.Header.Set("X-Kodepreneur-Timestamp", timestamp)
	req.Header.Set("X-Kodepreneur-Nonce", nonce)
	req.Header.Set("X-Kodepreneur-Signature", sig)
	return req
}

func TestServer_HealthAndMetricsEndpoints(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Security.SecretKey = "test-agent-secret"
	router := server.NewRouter(cfg)
	handler := router.Handler()

	// 1. Health check
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected /health 200 OK, got %d", rr.Code)
	}

	var healthRes map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &healthRes); err != nil || healthRes["status"] != "healthy" {
		t.Fatalf("unexpected health response: %s", rr.Body.String())
	}

	// 2. Authenticated System Info
	authReq := createAuthenticatedRequest("GET", "/api/v1/system/info", nil, cfg.Security.SecretKey)
	infoRR := httptest.NewRecorder()
	handler.ServeHTTP(infoRR, authReq)

	if infoRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /api/v1/system/info, got %d: %s", infoRR.Code, infoRR.Body.String())
	}

	// 3. Authenticated Metrics
	metricsReq := createAuthenticatedRequest("GET", "/api/v1/system/metrics", nil, cfg.Security.SecretKey)
	metricsRR := httptest.NewRecorder()
	handler.ServeHTTP(metricsRR, metricsReq)

	if metricsRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /api/v1/system/metrics, got %d: %s", metricsRR.Code, metricsRR.Body.String())
	}

	var metricsRes struct {
		Success bool `json:"success"`
		Data    struct {
			CPU struct {
				Cores int `json:"cores"`
			} `json:"cpu"`
			Memory struct {
				TotalBytes uint64 `json:"total_bytes"`
			} `json:"memory"`
		} `json:"data"`
	}
	if err := json.Unmarshal(metricsRR.Body.Bytes(), &metricsRes); err != nil || !metricsRes.Success {
		t.Fatalf("failed parsing metrics response: %s", metricsRR.Body.String())
	}

	if metricsRes.Data.CPU.Cores == 0 {
		t.Fatal("expected CPU cores > 0")
	}

	// 4. Authenticated System Update
	updateReqPayload := []byte(`{"repository":"https://github.com/kodepreneur/kodepreneur-panel.git","branch":"main"}`)
	updateReq := createAuthenticatedRequest("POST", "/api/v1/system/update", updateReqPayload, cfg.Security.SecretKey)
	updateRR := httptest.NewRecorder()
	handler.ServeHTTP(updateRR, updateReq)

	if updateRR.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /api/v1/system/update, got %d: %s", updateRR.Code, updateRR.Body.String())
	}

	var updateRes struct {
		Success bool `json:"success"`
		Data    struct {
			Success    bool   `json:"success"`
			CommitHash string `json:"commit_hash"`
			LogOutput  string `json:"log_output"`
		} `json:"data"`
	}
	if err := json.Unmarshal(updateRR.Body.Bytes(), &updateRes); err != nil || !updateRes.Success || !updateRes.Data.Success {
		t.Fatalf("failed parsing update response: %s", updateRR.Body.String())
	}

	if updateRes.Data.CommitHash == "" {
		t.Fatal("expected commit hash in update response")
	}
}
