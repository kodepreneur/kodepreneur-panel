package ssl

import (
	"os"
	"testing"
)

func TestDevSslIssuance(t *testing.T) {
	mgr := NewManager(true)

	result, err := mgr.IssueCertificate("test.local", []string{"www.test.local"}, "admin@test.local", "/tmp/kodepreneur/test/public")
	if err != nil {
		t.Fatalf("IssueCertificate dev failed: %v", err)
	}

	if result.Domain != "test.local" {
		t.Errorf("Expected domain test.local, got %s", result.Domain)
	}

	if _, err := os.Stat(result.CertPath); err != nil {
		t.Errorf("Certificate file does not exist: %v", err)
	}

	if _, err := os.Stat(result.KeyPath); err != nil {
		t.Errorf("Private key file does not exist: %v", err)
	}

	// Verify certificate can be parsed
	inspect, err := mgr.InspectCertificate(result.CertPath)
	if err != nil {
		t.Fatalf("InspectCertificate failed: %v", err)
	}

	if inspect.ValidUntil.Before(inspect.ValidFrom) {
		t.Errorf("Invalid validity period on certificate")
	}
}
