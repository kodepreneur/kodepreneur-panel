package filemanager

import (
	"archive/zip"
	"bytes"
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

func TestExtractZip(t *testing.T) {
	mgr := NewManager(true)
	destDir := filepath.Join(os.TempDir(), "kodepreneur", "test_extract_zip")
	_ = os.RemoveAll(destDir)
	defer os.RemoveAll(destDir)

	// Create test in-memory zip
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// File 1: index.php
	f1, err := zw.Create("public/index.php")
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}
	_, _ = f1.Write([]byte("<?php echo 'Laravel App';"))

	// File 2: artisan
	f2, err := zw.Create("artisan")
	if err != nil {
		t.Fatalf("Failed to create zip file: %v", err)
	}
	_, _ = f2.Write([]byte("#!/usr/bin/env php\n<?php // artisan"))

	if err := zw.Close(); err != nil {
		t.Fatalf("Failed to close zip writer: %v", err)
	}

	// Extract
	if err := mgr.ExtractZipBytes(buf.Bytes(), destDir); err != nil {
		t.Fatalf("ExtractZipBytes failed: %v", err)
	}

	// Verify files exist
	artisanPath := filepath.Join(destDir, "artisan")
	if _, err := os.Stat(artisanPath); os.IsNotExist(err) {
		t.Errorf("Expected artisan file to exist at %s", artisanPath)
	}

	indexPath := filepath.Join(destDir, "public", "index.php")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Errorf("Expected public/index.php to exist at %s", indexPath)
	}
}
