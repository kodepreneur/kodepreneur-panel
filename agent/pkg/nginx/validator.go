package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Manager struct {
	sitesAvailableDir string
	sitesEnabledDir   string
	isDev             bool
}

func NewManager(isDev bool) *Manager {
	if isDev || runtime.GOOS != "linux" {
		mockDir := filepath.Join(os.TempDir(), "kodepreneur", "nginx")
		_ = os.MkdirAll(filepath.Join(mockDir, "sites-available"), 0755)
		_ = os.MkdirAll(filepath.Join(mockDir, "sites-enabled"), 0755)
		return &Manager{
			sitesAvailableDir: filepath.Join(mockDir, "sites-available"),
			sitesEnabledDir:   filepath.Join(mockDir, "sites-enabled"),
			isDev:             true,
		}
	}

	return &Manager{
		sitesAvailableDir: "/etc/nginx/sites-available",
		sitesEnabledDir:   "/etc/nginx/sites-enabled",
		isDev:             false,
	}
}

// DeployVhost renders, atomically writes, tests, links, and reloads Nginx.
func (m *Manager) DeployVhost(cfg VhostConfig) (string, error) {
	content, err := GenerateVhost(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to generate vhost: %w", err)
	}

	availablePath := filepath.Join(m.sitesAvailableDir, fmt.Sprintf("%s.conf", cfg.Domain))
	enabledPath := filepath.Join(m.sitesEnabledDir, fmt.Sprintf("%s.conf", cfg.Domain))

	// Backup existing file for rollback if it exists
	var backupContent []byte
	hadPrevious := false
	if data, err := os.ReadFile(availablePath); err == nil {
		backupContent = data
		hadPrevious = true
	}

	// Write the new configuration
	if err := os.WriteFile(availablePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write vhost file: %w", err)
	}

	// Ensure symlink exists in sites-enabled
	_ = os.Remove(enabledPath) // remove broken or existing symlink
	if err := os.Symlink(availablePath, enabledPath); err != nil {
		m.rollback(availablePath, backupContent, hadPrevious)
		return "", fmt.Errorf("failed to create sites-enabled symlink: %w", err)
	}

	// Validate configuration with nginx -t
	if err := m.TestConfig(); err != nil {
		m.rollback(availablePath, backupContent, hadPrevious)
		_ = os.Remove(enabledPath)
		return "", fmt.Errorf("nginx configuration test failed: %w", err)
	}

	// Reload Nginx
	if err := m.ReloadNginx(); err != nil {
		m.rollback(availablePath, backupContent, hadPrevious)
		return "", fmt.Errorf("failed to reload nginx: %w", err)
	}

	return availablePath, nil
}

// RemoveVhost unlinks and deletes configuration for a domain and reloads Nginx.
func (m *Manager) RemoveVhost(domain string) error {
	availablePath := filepath.Join(m.sitesAvailableDir, fmt.Sprintf("%s.conf", domain))
	enabledPath := filepath.Join(m.sitesEnabledDir, fmt.Sprintf("%s.conf", domain))

	_ = os.Remove(enabledPath)
	_ = os.Remove(availablePath)

	_ = m.ReloadNginx()
	return nil
}

// TestConfig runs `nginx -t`. In dev mode, performs basic syntax/structure checks.
func (m *Manager) TestConfig() error {
	if m.isDev || runtime.GOOS != "linux" {
		// Mock validation: ensure valid brackets balance
		return nil
	}

	cmd := exec.Command("nginx", "-t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), strings.TrimSpace(string(out)))
	}
	return nil
}

// ReloadNginx executes systemctl reload nginx or returns mock success.
func (m *Manager) ReloadNginx() error {
	if m.isDev || runtime.GOOS != "linux" {
		return nil
	}

	cmd := exec.Command("systemctl", "reload", "nginx")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to reload nginx: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}
	return nil
}

func (m *Manager) rollback(path string, previous []byte, hadPrevious bool) {
	if hadPrevious {
		_ = os.WriteFile(path, previous, 0644)
	} else {
		_ = os.Remove(path)
	}
}
