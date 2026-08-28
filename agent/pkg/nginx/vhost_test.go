package nginx

import (
	"strings"
	"testing"
)

func TestGenerateVhost(t *testing.T) {
	cfg := VhostConfig{
		Domain:       "test-app.com",
		Aliases:      []string{"www.test-app.com"},
		DocumentRoot: "/var/www/test-app.com/public",
		PhpVersion:   "8.3",
		SystemUser:   "kp_test_app",
		SslEnabled:   false,
	}

	content, err := GenerateVhost(cfg)
	if err != nil {
		t.Fatalf("GenerateVhost failed: %v", err)
	}

	if !strings.Contains(content, "server_name test-app.com www.test-app.com;") {
		t.Errorf("Expected server_name directive, got:\n%s", content)
	}

	if !strings.Contains(content, "root /var/www/test-app.com/public;") {
		t.Errorf("Expected root directive, got:\n%s", content)
	}

	if !strings.Contains(content, "fastcgi_pass unix:/run/php/php8.3-fpm-kp_test_app.sock;") {
		t.Errorf("Expected PHP-FPM fastcgi_pass socket, got:\n%s", content)
	}
}

func TestGenerateVhostWithSsl(t *testing.T) {
	cfg := VhostConfig{
		Domain:       "secure.example.com",
		DocumentRoot: "/var/www/secure.example.com/public",
		PhpVersion:   "8.4",
		SystemUser:   "kp_secure",
		SslEnabled:   true,
		ForceHttps:   true,
		CertPath:     "/etc/letsencrypt/live/secure.example.com/fullchain.pem",
		KeyPath:      "/etc/letsencrypt/live/secure.example.com/privkey.pem",
	}

	content, err := GenerateVhost(cfg)
	if err != nil {
		t.Fatalf("GenerateVhost with SSL failed: %v", err)
	}

	if !strings.Contains(content, "return 301 https://$host$request_uri;") {
		t.Errorf("Expected HTTP to HTTPS redirect, got:\n%s", content)
	}

	if !strings.Contains(content, "ssl_certificate /etc/letsencrypt/live/secure.example.com/fullchain.pem;") {
		t.Errorf("Expected ssl_certificate directive, got:\n%s", content)
	}

	if !strings.Contains(content, "fastcgi_pass unix:/run/php/php8.4-fpm-kp_secure.sock;") {
		t.Errorf("Expected PHP 8.4 fastcgi_pass socket, got:\n%s", content)
	}
}
