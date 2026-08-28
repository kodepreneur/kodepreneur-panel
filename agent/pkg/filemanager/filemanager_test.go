package filemanager

import (
	"archive/tar"
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
	_ = os.RemoveAll(baseDir)
	_ = os.MkdirAll(baseDir, 0755)
	defer os.RemoveAll(baseDir)

	// 1. Path traversal should fail
	_, err := mgr.Browse(baseDir, "../../etc", true)
	if err == nil || !strings.Contains(err.Error(), "access denied") {
		t.Errorf("Expected path traversal error, got %v", err)
	}

	// 2. Null byte should fail
	_, err = mgr.Browse(baseDir, "test\x00file", true)
	if err == nil || !strings.Contains(err.Error(), "null byte") {
		t.Errorf("Expected null byte error, got %v", err)
	}

	// 3. Create & Write file
	if err := mgr.CreateFile(baseDir, "test.txt"); err != nil {
		t.Fatalf("CreateFile failed: %v", err)
	}
	if err := mgr.WriteFile(baseDir, "test.txt", "Hello Kodepreneur!"); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// 4. Read file
	content, err := mgr.ReadFile(baseDir, "test.txt", 1024)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if content != "Hello Kodepreneur!" {
		t.Errorf("Expected 'Hello Kodepreneur!', got '%s'", content)
	}

	// 5. Create Directory
	if err := mgr.CreateDirectory(baseDir, "subdir"); err != nil {
		t.Fatalf("CreateDirectory failed: %v", err)
	}

	// 6. Copy File
	if err := mgr.Copy(baseDir, "test.txt", "subdir/copied.txt"); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// 7. Move File
	if err := mgr.Move(baseDir, "subdir/copied.txt", "subdir/moved.txt"); err != nil {
		t.Fatalf("Move failed: %v", err)
	}

	// 8. Rename File
	if err := mgr.Rename(baseDir, "subdir/moved.txt", "subdir/renamed.txt"); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	// 9. Browse directory
	entries, err := mgr.Browse(baseDir, "", true)
	if err != nil {
		t.Fatalf("Browse failed: %v", err)
	}
	if len(entries) < 2 {
		t.Errorf("Expected at least 2 entries, got %d", len(entries))
	}

	// 10. Stat
	details, err := mgr.Stat(baseDir, "test.txt")
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if details.Name != "test.txt" || details.SizeBytes != int64(len("Hello Kodepreneur!")) {
		t.Errorf("Stat mismatch: %+v", details)
	}

	// 11. Chmod
	if err := mgr.Chmod(baseDir, "test.txt", "0600", false); err != nil {
		t.Fatalf("Chmod failed: %v", err)
	}

	// 12. Search
	results, err := mgr.Search(baseDir, "renamed", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) == 0 {
		t.Errorf("Expected search to find renamed.txt")
	}

	// 13. Delete file
	if err := mgr.DeleteEntry(baseDir, "test.txt"); err != nil {
		t.Fatalf("DeleteEntry failed: %v", err)
	}
}

func TestZipSlipRejection(t *testing.T) {
	mgr := NewManager(true)
	destDir := filepath.Join(os.TempDir(), "kodepreneur", "test_zip_slip")
	_ = os.RemoveAll(destDir)
	_ = os.MkdirAll(destDir, 0755)
	defer os.RemoveAll(destDir)

	// Create malicious zip with ../../evil.txt
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("../../evil.txt")
	if err != nil {
		t.Fatalf("Failed to create malicious zip entry: %v", err)
	}
	_, _ = f.Write([]byte("malicious content"))
	_ = zw.Close()

	err = mgr.ExtractZipBytes(buf.Bytes(), destDir)
	if err == nil || !strings.Contains(err.Error(), "zip slip") {
		t.Fatalf("Expected Zip Slip error, got: %v", err)
	}
}

func TestTarSlipRejection(t *testing.T) {
	mgr := NewManager(true)
	destDir := filepath.Join(os.TempDir(), "kodepreneur", "test_tar_slip")
	_ = os.RemoveAll(destDir)
	_ = os.MkdirAll(destDir, 0755)
	defer os.RemoveAll(destDir)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: "../../../etc/evil.conf",
		Mode: 0644,
		Size: int64(len("pwned")),
	}
	_ = tw.WriteHeader(hdr)
	_, _ = tw.Write([]byte("pwned"))
	_ = tw.Close()

	err := mgr.extractTarReader(&buf, destDir)
	if err == nil || !strings.Contains(err.Error(), "tar slip") {
		t.Fatalf("Expected Tar Slip error, got: %v", err)
	}
}

func TestCompressAndExtract(t *testing.T) {
	mgr := NewManager(true)
	baseDir := filepath.Join(os.TempDir(), "kodepreneur", "test_archive_mgr")
	_ = os.RemoveAll(baseDir)
	_ = os.MkdirAll(filepath.Join(baseDir, "app"), 0755)
	defer os.RemoveAll(baseDir)

	_ = mgr.WriteFile(baseDir, "app/index.php", "<?php echo 'OK';")
	_ = mgr.WriteFile(baseDir, "app/config.json", `{"app": "kp"}`)

	// 1. Compress to Zip
	if err := mgr.Compress(baseDir, []string{"app"}, "backup.zip", "zip"); err != nil {
		t.Fatalf("Compress Zip failed: %v", err)
	}

	// 2. Extract to extracted/
	if err := mgr.Extract(baseDir, "backup.zip", "extracted"); err != nil {
		t.Fatalf("Extract Zip failed: %v", err)
	}

	content, err := mgr.ReadFile(baseDir, "extracted/app/index.php", 1024)
	if err != nil {
		t.Fatalf("ReadFile from extracted zip failed: %v", err)
	}
	if content != "<?php echo 'OK';" {
		t.Errorf("Content mismatch: %s", content)
	}

	// 3. Compress to tar.gz
	if err := mgr.Compress(baseDir, []string{"app"}, "backup.tar.gz", "tar.gz"); err != nil {
		t.Fatalf("Compress tar.gz failed: %v", err)
	}

	// 4. Extract tar.gz
	if err := mgr.Extract(baseDir, "backup.tar.gz", "extracted_tar"); err != nil {
		t.Fatalf("Extract tar.gz failed: %v", err)
	}

	tarContent, err := mgr.ReadFile(baseDir, "extracted_tar/app/config.json", 1024)
	if err != nil {
		t.Fatalf("ReadFile from extracted tar failed: %v", err)
	}
	if tarContent != `{"app": "kp"}` {
		t.Errorf("Content mismatch from tar: %s", tarContent)
	}

	// 5. Stream Archive
	var streamBuf bytes.Buffer
	if err := mgr.StreamArchive(baseDir, []string{"app/index.php"}, &streamBuf); err != nil {
		t.Fatalf("StreamArchive failed: %v", err)
	}
	if streamBuf.Len() == 0 {
		t.Errorf("Expected streamed zip buffer to have data")
	}
}
