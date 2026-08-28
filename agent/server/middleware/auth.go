package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/kodepreneur/agent/config"
)

type NonceCache struct {
	mu     sync.Mutex
	nonces map[string]time.Time
}

var (
	globalNonceCache = &NonceCache{
		nonces: make(map[string]time.Time),
	}
)

func (nc *NonceCache) AddAndCheck(nonce string, ttl time.Duration) bool {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	now := time.Now()

	// Periodic purge
	for k, exp := range nc.nonces {
		if now.After(exp) {
			delete(nc.nonces, k)
		}
	}

	if _, exists := nc.nonces[nonce]; exists {
		return false // Nonce reused
	}

	nc.nonces[nonce] = now.Add(ttl)
	return true
}

func AuthMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Health check is always unauthenticated
		if r.URL.Path == "/api/v1/health" || r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		timestampStr := r.Header.Get("X-Kodepreneur-Timestamp")
		nonce := r.Header.Get("X-Kodepreneur-Nonce")
		signature := r.Header.Get("X-Kodepreneur-Signature")

		if timestampStr == "" || nonce == "" || signature == "" {
			respondJSONError(w, http.StatusUnauthorized, "MISSING_AUTH_HEADERS", "Required authentication headers missing")
			return
		}

		// 1. Verify Timestamp Freshness
		reqTimeSec, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			respondJSONError(w, http.StatusUnauthorized, "INVALID_TIMESTAMP", "Invalid timestamp format")
			return
		}

		nowSec := time.Now().Unix()
		drift := math.Abs(float64(nowSec - reqTimeSec))
		if drift > float64(cfg.Security.AllowedDriftSec) {
			respondJSONError(w, http.StatusUnauthorized, "EXPIRED_REQUEST", "Request timestamp is outside the allowed time window")
			return
		}

		// 2. Verify Nonce Uniqueness
		if !globalNonceCache.AddAndCheck(nonce, 2*time.Duration(cfg.Security.AllowedDriftSec)*time.Second) {
			respondJSONError(w, http.StatusUnauthorized, "REPLAY_ATTACK_DETECTED", "Duplicate nonce detected")
			return
		}

		// 3. Read Body & Compute SHA-256
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				respondJSONError(w, http.StatusBadRequest, "BODY_READ_ERROR", "Failed to read request body")
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		bodyHasher := sha256.New()
		bodyHasher.Write(bodyBytes)
		bodyHash := hex.EncodeToString(bodyHasher.Sum(nil))

		// 4. Compute Expected HMAC
		message := r.Method + "\n" + r.URL.Path + "\n" + timestampStr + "\n" + nonce + "\n" + bodyHash
		mac := hmac.New(sha256.New, []byte(cfg.Security.SecretKey))
		mac.Write([]byte(message))
		expectedSig := hex.EncodeToString(mac.Sum(nil))

		// 5. Constant Time Comparison
		if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
			respondJSONError(w, http.StatusUnauthorized, "INVALID_SIGNATURE", "HMAC signature verification failed")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
