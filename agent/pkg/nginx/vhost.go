package nginx

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type VhostConfig struct {
	Domain       string   `json:"domain"`
	Aliases      []string `json:"aliases"`
	DocumentRoot string   `json:"document_root"`
	PhpVersion   string   `json:"php_version"` // "8.3", "8.4", or "none"
	SystemUser   string   `json:"system_user"`
	SslEnabled   bool     `json:"ssl_enabled"`
	ForceHttps   bool     `json:"force_https"`
	CertPath     string   `json:"cert_path"`
	KeyPath      string   `json:"key_path"`
	AccessLog    string   `json:"access_log"`
	ErrorLog     string   `json:"error_log"`
}

const vhostTemplate = `# Managed by Kodepreneur Control Plane - DO NOT EDIT MANUALLY
{{- if and .SslEnabled .ForceHttps }}
server {
    listen 80;
    listen [::]:80;
    server_name {{ .Domain }}{{ range .Aliases }} {{ . }}{{ end }};

    # ACME-challenge webroot for Let's Encrypt
    location ^~ /.well-known/acme-challenge/ {
        root {{ .DocumentRoot }};
        default_type "text/plain";
        allow all;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}
{{- end }}

server {
{{- if .SslEnabled }}
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    {{- if not .ForceHttps }}
    listen 80;
    listen [::]:80;
    {{- end }}
{{- else }}
    listen 80;
    listen [::]:80;
{{- end }}

    server_name {{ .Domain }}{{ range .Aliases }} {{ . }}{{ end }};
    root {{ .DocumentRoot }};
    index index.php index.html index.htm;

    charset utf-8;

{{- if .SslEnabled }}
    ssl_certificate {{ .CertPath }};
    ssl_certificate_key {{ .KeyPath }};
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
{{- end }}

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # Logging
    access_log {{ .AccessLog }};
    error_log {{ .ErrorLog }} warn;

    # ACME-challenge
    location ^~ /.well-known/acme-challenge/ {
        root {{ .DocumentRoot }};
        default_type "text/plain";
        allow all;
    }

    # Favicon and Robots
    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    # Main application routing
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

{{- if ne .PhpVersion "none" }}
    # PHP-FPM FastCGI Handler
    location ~ \.php$ {
        fastcgi_pass unix:/run/php/php{{ .PhpVersion }}-fpm-{{ .SystemUser }}.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
        fastcgi_hide_header X-Powered-By;
        fastcgi_read_timeout 300;
        fastcgi_buffer_size 128k;
        fastcgi_buffers 4 256k;
        fastcgi_busy_buffers_size 256k;
    }
{{- end }}

    # Deny access to hidden files (.env, .git, etc)
    location ~ /\.(?!well-known).* {
        deny all;
        access_log off;
        log_not_found off;
    }
}
`

// GenerateVhost creates the Nginx configuration string based on VhostConfig.
func GenerateVhost(cfg VhostConfig) (string, error) {
	if cfg.Domain == "" {
		return "", fmt.Errorf("domain cannot be empty")
	}
	if cfg.DocumentRoot == "" {
		cfg.DocumentRoot = fmt.Sprintf("/var/www/%s/public", cfg.Domain)
	}
	if cfg.AccessLog == "" {
		cfg.AccessLog = fmt.Sprintf("/var/log/nginx/%s.access.log", cfg.Domain)
	}
	if cfg.ErrorLog == "" {
		cfg.ErrorLog = fmt.Sprintf("/var/log/nginx/%s.error.log", cfg.Domain)
	}
	if cfg.SystemUser == "" {
		cfg.SystemUser = fmt.Sprintf("kp_%s", strings.ReplaceAll(strings.Split(cfg.Domain, ".")[0], "-", "_"))
	}
	if cfg.PhpVersion == "" {
		cfg.PhpVersion = "8.3"
	}

	tmpl, err := template.New("nginx_vhost").Parse(vhostTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return "", fmt.Errorf("failed to render vhost template: %w", err)
	}

	return buf.String(), nil
}
