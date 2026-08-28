package process

import (
	"os"
	"testing"
)

func TestProcessSupervisorLifecycle(t *testing.T) {
	sup := NewSupervisor(true)

	cfg := ProcessConfig{
		ID:          "worker-1",
		Name:        "Queue Worker",
		Command:     "php artisan queue:work --tries=3",
		SystemUser:  "kp_test",
		WorkingDir:  "/var/www/test.com",
		AutoRestart: true,
	}

	path, err := sup.CreateProcess(cfg)
	if err != nil {
		t.Fatalf("CreateProcess failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("Created unit file does not exist: %v", err)
	}

	if _, err := sup.ManageProcess(cfg.ID, "restart"); err != nil {
		t.Fatalf("ManageProcess failed: %v", err)
	}

	if err := sup.DeleteProcess(cfg.ID); err != nil {
		t.Fatalf("DeleteProcess failed: %v", err)
	}
}
