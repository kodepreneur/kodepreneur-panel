package nginx

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseLogLine(t *testing.T) {
	line := `192.168.1.100 - - [31/Aug/2026:14:30:15 +0000] "GET /api/v1/users?page=1 HTTP/1.1" 200 4520 "https://example.com" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"`

	entry, parsedTime, err := ParseLogLine(line)
	if err != nil {
		t.Fatalf("unexpected error parsing log line: %v", err)
	}

	if entry.ClientIP != "192.168.1.100" {
		t.Errorf("expected client IP 192.168.1.100, got %s", entry.ClientIP)
	}
	if entry.Method != "GET" {
		t.Errorf("expected method GET, got %s", entry.Method)
	}
	if entry.Path != "/api/v1/users" {
		t.Errorf("expected path /api/v1/users, got %s", entry.Path)
	}
	if entry.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", entry.StatusCode)
	}
	if entry.BytesSent != 4520 {
		t.Errorf("expected bytes sent 4520, got %d", entry.BytesSent)
	}
	if entry.Browser != "Google Chrome" {
		t.Errorf("expected browser Google Chrome, got %s", entry.Browser)
	}
	if parsedTime.Year() != 2026 {
		t.Errorf("expected year 2026, got %d", parsedTime.Year())
	}
}

func TestDetectBrowser(t *testing.T) {
	tests := []struct {
		ua       string
		expected string
	}{
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0.0.0 Safari/537.36", "Google Chrome"},
		{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 Version/17.0 Safari/605.1.15", "Apple Safari"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0", "Mozilla Firefox"},
		{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Edg/120.0.0.0", "Microsoft Edge"},
		{"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)", "Bot / Crawler"},
		{"curl/7.68.0", "API Client (cURL/Postman)"},
		{"PostmanRuntime/7.28.4", "API Client (cURL/Postman)"},
	}

	for _, tt := range tests {
		result := detectBrowser(tt.ua)
		if result != tt.expected {
			t.Errorf("detectBrowser(%q) = %q, expected %q", tt.ua, result, tt.expected)
		}
	}
}

func TestAnalyzeTrafficSampleGeneration(t *testing.T) {
	tm := NewTrafficManager(true)

	periods := []string{"24h", "today", "7d", "30d"}
	for _, p := range periods {
		summary, err := tm.AnalyzeTraffic("test-domain.com", p)
		if err != nil {
			t.Fatalf("unexpected error analyzing traffic for period %s: %v", p, err)
		}

		if summary.Domain != "test-domain.com" {
			t.Errorf("expected domain test-domain.com, got %s", summary.Domain)
		}
		if summary.Period != p {
			t.Errorf("expected period %s, got %s", p, summary.Period)
		}
		if summary.TotalRequests <= 0 {
			t.Errorf("expected positive total requests for period %s", p)
		}
		if summary.TotalBytesSent <= 0 {
			t.Errorf("expected positive total bytes sent for period %s", p)
		}
		if len(summary.TimeSeries) == 0 {
			t.Errorf("expected non-empty time series for period %s", p)
		}
		if len(summary.TopPaths) == 0 {
			t.Errorf("expected non-empty top paths for period %s", p)
		}
		if len(summary.RecentRequests) == 0 {
			t.Errorf("expected non-empty recent requests for period %s", p)
		}
		if summary.SuccessRate <= 0 || summary.SuccessRate > 100 {
			t.Errorf("invalid success rate %f for period %s", summary.SuccessRate, p)
		}
	}
}

func TestAnalyzeTrafficRealFile(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "mysite.com.access.log")

	now := time.Now()
	t1 := now.Add(-2 * time.Hour).Format("02/Jan/2006:15:04:05 -0700")
	t2 := now.Add(-1 * time.Hour).Format("02/Jan/2006:15:04:05 -0700")

	logContent := "192.168.1.10 - - [" + t1 + `] "GET /home HTTP/1.1" 200 1200 "-" "Mozilla/5.0 Chrome/120.0"` + "\n" +
		"192.168.1.20 - - [" + t2 + `] "POST /api/login HTTP/1.1" 200 800 "https://mysite.com" "Mozilla/5.0 Safari/604.1"` + "\n" +
		"192.168.1.30 - - [" + t2 + `] "GET /unknown HTTP/1.1" 404 400 "-" "curl/7.68.0"` + "\n"

	if err := os.WriteFile(logFile, []byte(logContent), 0644); err != nil {
		t.Fatalf("failed to write test log file: %v", err)
	}

	tm := &TrafficManager{
		logDir: tempDir,
		isDev:  false,
	}

	summary, err := tm.AnalyzeTraffic("mysite.com", "24h")
	if err != nil {
		t.Fatalf("unexpected error analyzing real log file: %v", err)
	}

	if summary.TotalRequests != 3 {
		t.Errorf("expected 3 total requests, got %d", summary.TotalRequests)
	}
	if summary.UniqueVisitors != 3 {
		t.Errorf("expected 3 unique visitors, got %d", summary.UniqueVisitors)
	}
	if summary.TotalBytesSent != 2400 {
		t.Errorf("expected 2400 total bytes sent, got %d", summary.TotalBytesSent)
	}
	if summary.StatusCategories["2xx"] != 2 {
		t.Errorf("expected 2 2xx responses, got %d", summary.StatusCategories["2xx"])
	}
	if summary.StatusCategories["4xx"] != 1 {
		t.Errorf("expected 1 4xx response, got %d", summary.StatusCategories["4xx"])
	}
	if len(summary.RecentRequests) != 3 {
		t.Errorf("expected 3 recent requests, got %d", len(summary.RecentRequests))
	}
}
