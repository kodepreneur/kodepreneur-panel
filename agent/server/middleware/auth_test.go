package middleware_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/kodepreneur/agent/config"
	"github.com/kodepreneur/agent/server/middleware"
)

func TestAuthMiddleware_ValidHMAC(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Security.SecretKey = "test-secret"
	cfg.Security.AllowedDriftSec = 30

	nextCalled := false
	handler := middleware.AuthMiddleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := "test-nonce-1"
	body := []byte(`{"key":"value"}`)

	hasher := sha256.New()
	hasher.Write(body)
	bodyHash := hex.EncodeToString(hasher.Sum(nil))

	msg := "POST\n/api/v1/test\n" + timestamp + "\n" + nonce + "\n" + bodyHash
	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write([]byte(msg))
	sig := hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest("POST", "/api/v1/test", bytes.NewBuffer(body))
	req.Header.Set("X-Kodepreneur-Timestamp", timestamp)
	req.Header.Set("X-Kodepreneur-Nonce", nonce)
	req.Header.Set("X-Kodepreneur-Signature", sig)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rr.Code, rr.Body.String())
	}
	if !nextCalled {
		t.Fatal("expected downstream handler to be called")
	}
}

func TestAuthMiddleware_RejectsInvalidSignature(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Security.SecretKey = "test-secret"

	handler := middleware.AuthMiddleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/system/metrics", nil)
	req.Header.Set("X-Kodepreneur-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Set("X-Kodepreneur-Nonce", "nonce-bad")
	req.Header.Set("X-Kodepreneur-Signature", "bad-signature")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rr.Code)
	}
}

func TestAuthMiddleware_HealthBypassesAuth(t *testing.T) {
	cfg := config.DefaultConfig()
	handler := middleware.AuthMiddleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /health without headers, got %d", rr.Code)
	}
}
