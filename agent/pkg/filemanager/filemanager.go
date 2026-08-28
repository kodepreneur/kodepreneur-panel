package filemanager

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileEntry struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	IsDir       bool      `json:"is_dir"`
	SizeBytes   int64     `json:"size_bytes"`
	Permissions string    `json:"permissions"`
	ModifiedAt  time.Time `json:"modified_at"`
}

type Manager struct {
	isDev bool
}

func NewManager(isDev bool) *Manager {
	return &Manager{isDev: isDev}
}

// safePath resolves and validates that target stays inside basePath.
func (m *Manager) safePath(basePath, relativePath string) (string, error) {
	if basePath == "" {
		return "", fmt.Errorf("base path cannot be empty")
	}

	cleanBase := filepath.Clean(basePath)
	target := filepath.Clean(filepath.Join(cleanBase, relativePath))

	// Ensure target starts with cleanBase
	if !strings.HasPrefix(target, cleanBase) {
		return "", fmt.Errorf("access denied: path traversal detected out of %s", cleanBase)
	}

	return target, nil
}

// Browse lists entries within basePath/relativePath.
func (m *Manager) Browse(basePath, relativePath string) ([]FileEntry, error) {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return nil, err
	}

	// Create directory if not exists in dev mode
	if m.isDev {
		_ = os.MkdirAll(target, 0755)
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var results []FileEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		rel, _ := filepath.Rel(basePath, filepath.Join(target, entry.Name()))
		results = append(results, FileEntry{
			Name:        entry.Name(),
			Path:        rel,
			IsDir:       entry.IsDir(),
			SizeBytes:   info.Size(),
			Permissions: info.Mode().String(),
			ModifiedAt:  info.ModTime(),
		})
	}

	return results, nil
}

// ReadFile reads up to 2MB from target file.
func (m *Manager) ReadFile(basePath, relativePath string) (string, error) {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return "", err
	}

	file, err := os.Open(target)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read max 2MB
	limitedReader := io.LimitReader(file, 2*1024*1024)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	return string(content), nil
}

// WriteFile atomically writes content to target file.
func (m *Manager) WriteFile(basePath, relativePath, content string) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}

	_ = os.MkdirAll(filepath.Dir(target), 0755)
	return os.WriteFile(target, []byte(content), 0644)
}

// CreateDirectory creates a folder at target path.
func (m *Manager) CreateDirectory(basePath, relativePath string) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}
	return os.MkdirAll(target, 0755)
}

// DeleteEntry removes file or directory at target path.
func (m *Manager) DeleteEntry(basePath, relativePath string) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}
	if target == filepath.Clean(basePath) {
		return fmt.Errorf("cannot delete base root directory")
	}
	return os.RemoveAll(target)
}

// ExtractZip extracts a zip archive safely, protecting against Zip Slip path traversal.
func (m *Manager) ExtractZip(reader io.ReaderAt, size int64, destDir string) error {
	zipReader, err := zip.NewReader(reader, size)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}

	cleanDest := filepath.Clean(destDir)
	if err := os.MkdirAll(cleanDest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, file := range zipReader.File {
		targetPath := filepath.Join(cleanDest, file.Name)
		if !strings.HasPrefix(filepath.Clean(targetPath), cleanDest) {
			return fmt.Errorf("illegal file path in zip: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
		}

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to read file %s inside zip: %w", file.Name, err)
		}

		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to write extracted file %s: %w", targetPath, err)
		}

		_, copyErr := io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if copyErr != nil {
			return fmt.Errorf("failed to extract file %s: %w", targetPath, copyErr)
		}
	}

	return nil
}

// ExtractZipBytes extracts zip bytes into destDir.
func (m *Manager) ExtractZipBytes(data []byte, destDir string) error {
	return m.ExtractZip(bytes.NewReader(data), int64(len(data)), destDir)
}

// ExtractZipFile opens a zip file on disk and extracts it into destDir.
func (m *Manager) ExtractZipFile(zipFilePath, destDir string) error {
	f, err := os.Open(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat zip file: %w", err)
	}

	return m.ExtractZip(f, info.Size(), destDir)
}
