package filemanager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileManagerSandboxing(t *testing.T) {
	mgr := NewManager(true)
	baseDir := filepath.Join(os.TempDir(), "kodepreneur", "test_jail")
	_ = os.MkdirAll(baseDir, 0755)

	// 1. Path traversal should fail
	_, err := mgr.Browse(baseDir, "../../etc")
	if err == nil || !strings.Contains(err.Error(), "access denied") {
		t.Errorf("Expected path traversal error, got %v", err)
	}

	// 2. Write file
	if err := mgr.WriteFile(baseDir, "test.txt", "Hello Kodepreneur!"); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// 3. Read file
	content, err := mgr.ReadFile(baseDir, "test.txt")
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if content != "Hello Kodepreneur!" {
		t.Errorf("Expected 'Hello Kodepreneur!', got '%s'", content)
	}

	// 4. Browse directory
	entries, err := mgr.Browse(baseDir, "")
	if err != nil {
		t.Fatalf("Browse failed: %v", err)
	}
	if len(entries) == 0 {
		t.Errorf("Expected at least 1 entry, got 0")
	}

	// 5. Delete file
	if err := mgr.DeleteEntry(baseDir, "test.txt"); err != nil {
		t.Fatalf("DeleteEntry failed: %v", err)
	}
}
