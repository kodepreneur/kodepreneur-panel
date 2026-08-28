package phpfpm

import (
	"os"
	"strings"
	"testing"
)

func TestPoolLifecycle(t *testing.T) {
	mgr := NewManager(true)

	cfg := PoolConfig{
		Domain:       "test-domain.dev",
		SystemUser:   "kp_test_dev",
		PhpVersion:   "8.3",
		DocumentRoot: "/tmp/kodepreneur/test-domain/public",
	}

	path, err := mgr.CreatePool(cfg)
	if err != nil {
		t.Fatalf("CreatePool failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read created pool file %s: %v", path, err)
	}

	content := string(data)
	if !strings.Contains(content, "[kp_test_dev]") {
		t.Errorf("Expected pool user header [kp_test_dev], got:\n%s", content)
	}
	if !strings.Contains(content, "listen = /run/php/php8.3-fpm-kp_test_dev.sock") {
		t.Errorf("Expected socket path in pool config, got:\n%s", content)
	}

	// Test Switch Version
	cfg.PhpVersion = "8.4"
	if err := mgr.SwitchVersion(cfg, "8.3"); err != nil {
		t.Fatalf("SwitchVersion failed: %v", err)
	}

	// Delete Pool
	if err := mgr.DeletePool("8.4", cfg.SystemUser); err != nil {
		t.Fatalf("DeletePool failed: %v", err)
	}
}
