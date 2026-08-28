package filemanager

import (
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
