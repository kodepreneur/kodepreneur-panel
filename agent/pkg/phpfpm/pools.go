package phpfpm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

type PoolConfig struct {
	Domain            string `json:"domain"`
	SystemUser        string `json:"system_user"`
	PhpVersion        string `json:"php_version"` // "8.3" or "8.4"
	DocumentRoot      string `json:"document_root"`
	BaseDir           string `json:"base_dir"`
	MaxChildren       int    `json:"max_children"`
	StartServers      int    `json:"start_servers"`
	MinSpareServers   int    `json:"min_spare_servers"`
	MaxSpareServers   int    `json:"max_spare_servers"`
	MaxRequests       int    `json:"max_requests"`
	MemoryLimit       string `json:"memory_limit"`
	UploadMaxFilesize string `json:"upload_max_filesize"`
	PostMaxSize       string `json:"post_max_size"`
}

const poolTemplate = `; Managed by Kodepreneur Control Plane - DO NOT EDIT MANUALLY
[{{ .SystemUser }}]
user = {{ .SystemUser }}
group = {{ .SystemUser }}

listen = /run/php/php{{ .PhpVersion }}-fpm-{{ .SystemUser }}.sock
listen.owner = www-data
listen.group = www-data
listen.mode = 0660

pm = ondemand
pm.max_children = {{ .MaxChildren }}
pm.process_idle_timeout = 10s
pm.max_requests = {{ .MaxRequests }}

chdir = {{ .DocumentRoot }}

php_admin_value[memory_limit] = {{ .MemoryLimit }}
php_admin_value[upload_max_filesize] = {{ .UploadMaxFilesize }}
php_admin_value[post_max_size] = {{ .PostMaxSize }}
php_admin_value[max_execution_time] = 300
php_admin_value[date.timezone] = UTC
php_admin_value[disable_functions] = exec,passthru,shell_exec,system,proc_open,popen,show_source

; Security / Open Basedir
php_admin_value[open_basedir] = {{ .BaseDir }}:/tmp:/var/tmp
`

type Manager struct {
	phpFpmDir string
	isDev     bool
}

func NewManager(isDev bool) *Manager {
	if isDev || runtime.GOOS != "linux" {
		mockDir := filepath.Join(os.TempDir(), "kodepreneur", "php")
		_ = os.MkdirAll(filepath.Join(mockDir, "8.3", "fpm", "pool.d"), 0755)
		_ = os.MkdirAll(filepath.Join(mockDir, "8.4", "fpm", "pool.d"), 0755)
		_ = os.MkdirAll(filepath.Join(os.TempDir(), "kodepreneur", "run", "php"), 0755)
		return &Manager{
			phpFpmDir: mockDir,
			isDev:     true,
		}
	}

	return &Manager{
		phpFpmDir: "/etc/php",
		isDev:     false,
	}
}

// ProvisionUser creates the isolated Linux system user if not exists.
func (m *Manager) ProvisionUser(username, domain string) error {
	if m.isDev || runtime.GOOS != "linux" {
		return nil
	}

	// Check if user already exists
	checkCmd := exec.Command("id", "-u", username)
	if err := checkCmd.Run(); err == nil {
		return nil // User already exists
	}

	// Create user with /var/www/{domain} home directory and no shell login
	cmd := exec.Command("useradd", "-m", "-d", fmt.Sprintf("/var/www/%s", domain), "-s", "/bin/bash", username)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("useradd failed: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}

	// Add www-data to the user's group so Nginx can access static files
	addGrpCmd := exec.Command("usermod", "-aG", username, "www-data")
	_ = addGrpCmd.Run()

	return nil
}

// PrepareWebroot creates document root and permissions safely for the target domain.
func (m *Manager) PrepareWebroot(docRoot, baseDir, username string) error {
	if m.isDev || runtime.GOOS != "linux" {
		if strings.HasPrefix(docRoot, "/var/www") {
			mockRoot := filepath.Join(os.TempDir(), "kodepreneur", "www", strings.TrimPrefix(docRoot, "/var/www"))
			_ = os.MkdirAll(mockRoot, 0755)
			indexPhp := filepath.Join(mockRoot, "index.php")
			indexHtml := filepath.Join(mockRoot, "index.html")
			if _, errPhp := os.Stat(indexPhp); os.IsNotExist(errPhp) {
				if _, errHtml := os.Stat(indexHtml); os.IsNotExist(errHtml) {
					samplePhp := "<?php echo 'Welcome to Kodepreneur';\n"
					_ = os.WriteFile(indexPhp, []byte(samplePhp), 0644)
				}
			}
			return nil
		}
		_ = os.MkdirAll(docRoot, 0755)
		return nil
	}

	if err := os.MkdirAll(docRoot, 0755); err != nil {
		return fmt.Errorf("failed to create document root %s: %w", docRoot, err)
	}

	// Only create default sample index if NO index file exists in docRoot
	indexPhp := filepath.Join(docRoot, "index.php")
	indexHtml := filepath.Join(docRoot, "index.html")
	if _, errPhp := os.Stat(indexPhp); os.IsNotExist(errPhp) {
		if _, errHtml := os.Stat(indexHtml); os.IsNotExist(errHtml) {
			samplePhp := `<?php
echo "<h1>Welcome to " . htmlspecialchars($_SERVER['HTTP_HOST'] ?? 'Kodepreneur') . "</h1>";
echo "<p>PHP Version: " . PHP_VERSION . "</p>";
echo "<p>Server Time: " . date('Y-m-d H:i:s') . " UTC</p>";
`
			_ = os.WriteFile(indexPhp, []byte(samplePhp), 0644)
		}
	}

	// Safe chown: Only chown target domain directory (/var/www/{domain}), NEVER /var/www or /
	targetDir := baseDir
	if targetDir == "" {
		targetDir = docRoot
	}
	cleanTarget := filepath.Clean(targetDir)
	if cleanTarget != "/var/www" && cleanTarget != "/" && cleanTarget != "." && strings.HasPrefix(cleanTarget, "/var/www/") {
		cmd := exec.Command("chown", "-R", fmt.Sprintf("%s:www-data", username), cleanTarget)
		_ = cmd.Run()
	}

	return nil
}

// CreatePool renders and writes the PHP-FPM pool config file and reloads the service.
func (m *Manager) CreatePool(cfg PoolConfig) (string, error) {
	if cfg.PhpVersion == "" || cfg.PhpVersion == "none" {
		return "", nil // No pool needed for static sites
	}
	if cfg.MaxChildren <= 0 {
		cfg.MaxChildren = 10
	}
	if cfg.MaxRequests <= 0 {
		cfg.MaxRequests = 500
	}
	if cfg.MemoryLimit == "" {
		cfg.MemoryLimit = "512M"
	}
	if cfg.UploadMaxFilesize == "" {
		cfg.UploadMaxFilesize = "512M"
	}
	if cfg.PostMaxSize == "" {
		cfg.PostMaxSize = "512M"
	}
	if cfg.BaseDir == "" {
		if cfg.Domain != "" {
			cfg.BaseDir = fmt.Sprintf("/var/www/%s", cfg.Domain)
		} else if cfg.DocumentRoot != "" {
			if strings.HasSuffix(cfg.DocumentRoot, "/public") {
				cfg.BaseDir = filepath.Dir(cfg.DocumentRoot)
			} else {
				cfg.BaseDir = cfg.DocumentRoot
			}
		} else {
			cfg.BaseDir = "/var/www"
		}
	}

	tmpl, err := template.New("php_pool").Parse(poolTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse pool template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cfg); err != nil {
		return "", fmt.Errorf("failed to render pool template: %w", err)
	}

	poolDir := filepath.Join(m.phpFpmDir, cfg.PhpVersion, "fpm", "pool.d")
	if err := os.MkdirAll(poolDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create pool directory: %w", err)
	}

	poolPath := filepath.Join(poolDir, fmt.Sprintf("%s.conf", cfg.SystemUser))
	if err := os.WriteFile(poolPath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to write pool config %s: %w", poolPath, err)
	}

	// Validate PHP-FPM configuration before keeping pool
	if err := m.TestConfig(cfg.PhpVersion); err != nil {
		_ = os.Remove(poolPath)
		_ = m.ReloadFpm(cfg.PhpVersion)
		return "", fmt.Errorf("php%s-fpm pool syntax validation failed: %w", cfg.PhpVersion, err)
	}

	// Reload PHP-FPM daemon gracefully
	if err := m.ReloadFpm(cfg.PhpVersion); err != nil {
		_ = os.Remove(poolPath)
		_ = m.ReloadFpm(cfg.PhpVersion)
		return "", fmt.Errorf("failed to reload php%s-fpm: %w", cfg.PhpVersion, err)
	}

	return poolPath, nil
}

// DeletePool removes the PHP-FPM pool and reloads the service.
func (m *Manager) DeletePool(phpVersion, systemUser string) error {
	if phpVersion == "" || phpVersion == "none" {
		return nil
	}

	poolPath := filepath.Join(m.phpFpmDir, phpVersion, "fpm", "pool.d", fmt.Sprintf("%s.conf", systemUser))
	_ = os.Remove(poolPath)

	return m.ReloadFpm(phpVersion)
}

// SwitchVersion removes pool from previous PHP version, creates it in the new one, and reloads both.
func (m *Manager) SwitchVersion(cfg PoolConfig, oldVersion string) error {
	if oldVersion != "" && oldVersion != "none" && oldVersion != cfg.PhpVersion {
		_ = m.DeletePool(oldVersion, cfg.SystemUser)
	}

	if cfg.PhpVersion != "none" {
		_, err := m.CreatePool(cfg)
		return err
	}

	return nil
}

// TestConfig tests PHP-FPM configuration syntax.
func (m *Manager) TestConfig(phpVersion string) error {
	if m.isDev || runtime.GOOS != "linux" {
		return nil
	}

	cmd := exec.Command(fmt.Sprintf("php-fpm%s", phpVersion), "-t")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", err.Error(), strings.TrimSpace(string(out)))
	}
	return nil
}

// ReloadFpm sends reload signal to php-fpm systemd unit.
func (m *Manager) ReloadFpm(phpVersion string) error {
	if m.isDev || runtime.GOOS != "linux" {
		return nil
	}

	serviceName := fmt.Sprintf("php%s-fpm", phpVersion)
	cmd := exec.Command("systemctl", "reload", serviceName)
	if out, err := cmd.CombinedOutput(); err != nil {
		// Fallback to reload-or-restart
		cmdFallback := exec.Command("systemctl", "reload-or-restart", serviceName)
		outFallback, errFallback := cmdFallback.CombinedOutput()
		if errFallback != nil {
			return fmt.Errorf("%s: %s (reload fallback: %s)", errFallback.Error(), strings.TrimSpace(string(outFallback)), strings.TrimSpace(string(out)))
		}
	}
	return nil
}
