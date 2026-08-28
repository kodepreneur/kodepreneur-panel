package filemanager

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"mime"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type FileEntry struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	IsDir       bool      `json:"is_dir"`
	SizeBytes   int64     `json:"size_bytes"`
	Permissions string    `json:"permissions"`
	ModeOctal   string    `json:"mode_octal"`
	Owner       string    `json:"owner"`
	Group       string    `json:"group"`
	ModifiedAt  time.Time `json:"modified_at"`
	ItemCount   int       `json:"item_count"`
	MimeType    string    `json:"mime_type"`
	Extension   string    `json:"extension"`
}

type FileDetails struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	IsDir       bool      `json:"is_dir"`
	SizeBytes   int64     `json:"size_bytes"`
	Permissions string    `json:"permissions"`
	ModeOctal   string    `json:"mode_octal"`
	Owner       string    `json:"owner"`
	Group       string    `json:"group"`
	UID         int       `json:"uid"`
	GID         int       `json:"gid"`
	ModifiedAt  time.Time `json:"modified_at"`
	CreatedAt   time.Time `json:"created_at"`
	MimeType    string    `json:"mime_type"`
	Extension   string    `json:"extension"`
	ItemCount   int       `json:"item_count"`
}

type DiskUsageInfo struct {
	Path         string  `json:"path"`
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	PathSize     int64   `json:"path_size"`
}

type Manager struct {
	isDev bool
}

func NewManager(isDev bool) *Manager {
	return &Manager{isDev: isDev}
}

// SafePath resolves and validates that target stays inside basePath.
func (m *Manager) SafePath(basePath, relativePath string) (string, error) {
	return m.safePath(basePath, relativePath)
}

// safePath resolves and validates that target stays inside basePath.
func (m *Manager) safePath(basePath, relativePath string) (string, error) {
	if basePath == "" {
		return "", fmt.Errorf("base path cannot be empty")
	}

	// Reject null-bytes
	if strings.Contains(basePath, "\x00") || strings.Contains(relativePath, "\x00") {
		return "", fmt.Errorf("access denied: null byte detected in path")
	}

	cleanBase := filepath.Clean(basePath)
	target := filepath.Clean(filepath.Join(cleanBase, relativePath))

	// Ensure target stays within cleanBase boundary
	if target != cleanBase && !strings.HasPrefix(target, cleanBase+string(filepath.Separator)) {
		return "", fmt.Errorf("access denied: path traversal detected out of %s", cleanBase)
	}

	// Check if target is a symlink pointing outside cleanBase
	if fi, err := os.Lstat(target); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		resolved, err := filepath.EvalSymlinks(target)
		if err == nil {
			cleanResolved := filepath.Clean(resolved)
			if cleanResolved != cleanBase && !strings.HasPrefix(cleanResolved, cleanBase+string(filepath.Separator)) {
				return "", fmt.Errorf("access denied: symlink points outside %s", cleanBase)
			}
		}
	}

	return target, nil
}

func (m *Manager) getTargetOwnership(basePath string) (int, int) {
	if m.isDev || runtime.GOOS != "linux" {
		return 0, 0
	}

	cleanBase := filepath.Clean(basePath)

	// 1. Check if basePath has a specific non-root owner
	if fi, err := os.Stat(cleanBase); err == nil {
		if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
			uid := int(stat.Uid)
			gid := int(stat.Gid)
			if uid != 0 || gid != 0 {
				return uid, gid
			}
		}
	}

	// 2. If basePath is inside /var/www/<domain>, look up domain system user kp_<domain>
	if strings.HasPrefix(cleanBase, "/var/www/") {
		parts := strings.Split(strings.TrimPrefix(cleanBase, "/var/www/"), "/")
		if len(parts) > 0 && parts[0] != "" {
			domainSlug := strings.ReplaceAll(strings.Split(parts[0], ".")[0], "-", "_")
			systemUser := fmt.Sprintf("kp_%s", domainSlug)
			if u, err := user.Lookup(systemUser); err == nil {
				if uid, err := strconv.Atoi(u.Uid); err == nil {
					gid := uid
					if g, err := user.LookupGroup("www-data"); err == nil {
						if gidInt, err := strconv.Atoi(g.Gid); err == nil {
							gid = gidInt
						}
					}
					return uid, gid
				}
			}
		}
	}

	// 3. Fallback to www-data if available
	if u, err := user.Lookup("www-data"); err == nil {
		if uid, err := strconv.Atoi(u.Uid); err == nil {
			gid := uid
			if g, err := user.LookupGroup("www-data"); err == nil {
				if gidInt, err := strconv.Atoi(g.Gid); err == nil {
					gid = gidInt
				}
			}
			return uid, gid
		}
	}

	return 0, 0
}

func (m *Manager) applyOwnership(basePath, targetPath string, recursive bool) {
	if m.isDev || runtime.GOOS != "linux" {
		return
	}
	uid, gid := m.getTargetOwnership(basePath)
	if uid == 0 && gid == 0 {
		return
	}
	if recursive {
		_ = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
			if err == nil {
				_ = os.Chown(path, uid, gid)
			}
			return nil
		})
	} else {
		_ = os.Chown(targetPath, uid, gid)
	}
}

func getOwnerGroup(info os.FileInfo) (string, string, int, int) {
	uid := 0
	gid := 0
	ownerName := "www-data"
	groupName := "www-data"

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = int(stat.Uid)
		gid = int(stat.Gid)

		if u, err := user.LookupId(strconv.Itoa(uid)); err == nil {
			ownerName = u.Username
		} else {
			ownerName = strconv.Itoa(uid)
		}

		if g, err := user.LookupGroupId(strconv.Itoa(gid)); err == nil {
			groupName = g.Name
		} else {
			groupName = strconv.Itoa(gid)
		}
	}

	return ownerName, groupName, uid, gid
}

// Browse lists entries within basePath/relativePath.
func (m *Manager) Browse(basePath, relativePath string, showHidden bool) ([]FileEntry, error) {
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
		name := entry.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		rel, _ := filepath.Rel(basePath, filepath.Join(target, name))
		owner, group, _, _ := getOwnerGroup(info)

		ext := strings.ToLower(filepath.Ext(name))
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			if entry.IsDir() {
				mimeType = "inode/directory"
			} else {
				mimeType = "application/octet-stream"
			}
		}

		itemCount := 0
		if entry.IsDir() {
			if subEntries, err := os.ReadDir(filepath.Join(target, name)); err == nil {
				itemCount = len(subEntries)
			}
		}

		modeOctal := fmt.Sprintf("%04o", info.Mode().Perm())

		results = append(results, FileEntry{
			Name:        name,
			Path:        rel,
			IsDir:       entry.IsDir(),
			SizeBytes:   info.Size(),
			Permissions: info.Mode().String(),
			ModeOctal:   modeOctal,
			Owner:       owner,
			Group:       group,
			ModifiedAt:  info.ModTime(),
			ItemCount:   itemCount,
			MimeType:    mimeType,
			Extension:   strings.TrimPrefix(ext, "."),
		})
	}

	return results, nil
}

// ReadFile reads up to maxBytes from target file.
func (m *Manager) ReadFile(basePath, relativePath string, maxBytes int64) (string, error) {
	data, err := m.ReadFileBytes(basePath, relativePath, maxBytes)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadFileBytes reads binary data from target file.
func (m *Manager) ReadFileBytes(basePath, relativePath string, maxBytes int64) ([]byte, error) {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(target)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if maxBytes <= 0 {
		maxBytes = 50 * 1024 * 1024 // 50MB default
	}

	limitedReader := io.LimitReader(file, maxBytes)
	return io.ReadAll(limitedReader)
}

// WriteFile atomically writes string content to target file.
func (m *Manager) WriteFile(basePath, relativePath, content string) error {
	return m.WriteFileBytes(basePath, relativePath, []byte(content))
}

// WriteFileBytes atomically writes binary content to target file.
func (m *Manager) WriteFileBytes(basePath, relativePath string, data []byte) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}

	_ = os.MkdirAll(filepath.Dir(target), 0755)

	var perm os.FileMode = 0644
	if fi, err := os.Stat(target); err == nil {
		perm = fi.Mode().Perm()
	}

	if err := os.WriteFile(target, data, perm); err != nil {
		return err
	}

	m.applyOwnership(basePath, target, false)
	return nil
}

// CreateFile creates an empty file.
func (m *Manager) CreateFile(basePath, relativePath string) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("file already exists: %s", relativePath)
	}

	_ = os.MkdirAll(filepath.Dir(target), 0755)
	f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	_ = f.Close()

	m.applyOwnership(basePath, target, false)
	return nil
}

// CreateDirectory creates a folder at target path.
func (m *Manager) CreateDirectory(basePath, relativePath string) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	m.applyOwnership(basePath, target, true)
	return nil
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

// Rename renames a file or directory.
func (m *Manager) Rename(basePath, oldRelPath, newRelPath string) error {
	oldTarget, err := m.safePath(basePath, oldRelPath)
	if err != nil {
		return err
	}
	newTarget, err := m.safePath(basePath, newRelPath)
	if err != nil {
		return err
	}

	if oldTarget == filepath.Clean(basePath) || newTarget == filepath.Clean(basePath) {
		return fmt.Errorf("cannot rename root directory")
	}

	return os.Rename(oldTarget, newTarget)
}

// Copy copies a file or directory recursively.
func (m *Manager) Copy(basePath, srcRelPath, destRelPath string) error {
	srcTarget, err := m.safePath(basePath, srcRelPath)
	if err != nil {
		return err
	}
	destTarget, err := m.safePath(basePath, destRelPath)
	if err != nil {
		return err
	}

	srcInfo, err := os.Stat(srcTarget)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if srcInfo.IsDir() {
		return m.copyDir(srcTarget, destTarget)
	}
	return m.copyFile(srcTarget, destTarget, srcInfo.Mode().Perm())
}

func (m *Manager) copyFile(src, dest string, perm os.FileMode) error {
	_ = os.MkdirAll(filepath.Dir(dest), 0755)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func (m *Manager) copyDir(src, dest string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, srcInfo.Mode().Perm()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := m.copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if err := m.copyFile(srcPath, destPath, info.Mode().Perm()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Move moves a file or directory.
func (m *Manager) Move(basePath, srcRelPath, destRelPath string) error {
	srcTarget, err := m.safePath(basePath, srcRelPath)
	if err != nil {
		return err
	}
	destTarget, err := m.safePath(basePath, destRelPath)
	if err != nil {
		return err
	}

	if srcTarget == filepath.Clean(basePath) {
		return fmt.Errorf("cannot move base root directory")
	}

	_ = os.MkdirAll(filepath.Dir(destTarget), 0755)

	// Try atomic rename first
	if err := os.Rename(srcTarget, destTarget); err == nil {
		return nil
	}

	// Fallback to copy + remove
	if err := m.Copy(basePath, srcRelPath, destRelPath); err != nil {
		return err
	}
	return os.RemoveAll(srcTarget)
}

// Chmod modifies permissions for file/directory.
func (m *Manager) Chmod(basePath, relativePath, modeStr string, recursive bool) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}

	// Parse octal mode e.g. "0755", "755", "0644"
	modeStr = strings.TrimPrefix(modeStr, "0")
	modeVal, err := strconv.ParseUint(modeStr, 8, 32)
	if err != nil {
		return fmt.Errorf("invalid mode format: %s", modeStr)
	}
	mode := os.FileMode(modeVal)

	if !recursive {
		return os.Chmod(target, mode)
	}

	return filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, mode)
	})
}

// Chown modifies owner and group for file/directory.
func (m *Manager) Chown(basePath, relativePath string, uid, gid int, recursive bool) error {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return err
	}

	if !recursive {
		return os.Chown(target, uid, gid)
	}

	return filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chown(path, uid, gid)
	})
}

// Stat retrieves detailed metadata.
func (m *Manager) Stat(basePath, relativePath string) (*FileDetails, error) {
	target, err := m.safePath(basePath, relativePath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	owner, group, uid, gid := getOwnerGroup(info)
	ext := strings.ToLower(filepath.Ext(info.Name()))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		if info.IsDir() {
			mimeType = "inode/directory"
		} else {
			mimeType = "application/octet-stream"
		}
	}

	itemCount := 0
	if info.IsDir() {
		if entries, err := os.ReadDir(target); err == nil {
			itemCount = len(entries)
		}
	}

	rel, _ := filepath.Rel(basePath, target)
	if rel == "." {
		rel = ""
	}

	return &FileDetails{
		Name:        info.Name(),
		Path:        rel,
		IsDir:       info.IsDir(),
		SizeBytes:   info.Size(),
		Permissions: info.Mode().String(),
		ModeOctal:   fmt.Sprintf("%04o", info.Mode().Perm()),
		Owner:       owner,
		Group:       group,
		UID:         uid,
		GID:         gid,
		ModifiedAt:  info.ModTime(),
		CreatedAt:   info.ModTime(),
		MimeType:    mimeType,
		Extension:   strings.TrimPrefix(ext, "."),
		ItemCount:   itemCount,
	}, nil
}

// Search searches for files matching a query inside basePath.
func (m *Manager) Search(basePath, query string, maxResults int) ([]FileEntry, error) {
	cleanBase := filepath.Clean(basePath)
	if maxResults <= 0 {
		maxResults = 100
	}

	queryLower := strings.ToLower(query)
	var results []FileEntry

	err := filepath.Walk(cleanBase, func(path string, info os.FileInfo, err error) error {
		if err != nil || len(results) >= maxResults {
			return nil
		}

		if path == cleanBase {
			return nil
		}

		if strings.Contains(strings.ToLower(info.Name()), queryLower) {
			rel, _ := filepath.Rel(cleanBase, path)
			owner, group, _, _ := getOwnerGroup(info)
			ext := strings.ToLower(filepath.Ext(info.Name()))
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				if info.IsDir() {
					mimeType = "inode/directory"
				} else {
					mimeType = "application/octet-stream"
				}
			}

			results = append(results, FileEntry{
				Name:        info.Name(),
				Path:        rel,
				IsDir:       info.IsDir(),
				SizeBytes:   info.Size(),
				Permissions: info.Mode().String(),
				ModeOctal:   fmt.Sprintf("%04o", info.Mode().Perm()),
				Owner:       owner,
				Group:       group,
				ModifiedAt:  info.ModTime(),
				MimeType:    mimeType,
				Extension:   strings.TrimPrefix(ext, "."),
			})
		}
		return nil
	})

	return results, err
}

// GetDiskUsage retrieves filesystem partition disk info for basePath.
func (m *Manager) GetDiskUsage(basePath string) (DiskUsageInfo, error) {
	cleanBase := filepath.Clean(basePath)
	var stat syscall.Statfs_t
	if err := syscall.Statfs(cleanBase, &stat); err != nil {
		return DiskUsageInfo{Path: cleanBase}, err
	}

	// Calculate partition bytes
	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bfree * uint64(stat.Bsize)
	usedBytes := totalBytes - freeBytes
	var percent float64 = 0
	if totalBytes > 0 {
		percent = float64(usedBytes) / float64(totalBytes) * 100.0
	}

	// Calculate immediate path size (shallow sum of direct children)
	var pathSize int64 = 0
	if entries, err := os.ReadDir(cleanBase); err == nil {
		for _, entry := range entries {
			if info, err := entry.Info(); err == nil {
				pathSize += info.Size()
			}
		}
	}

	return DiskUsageInfo{
		Path:         cleanBase,
		TotalBytes:   totalBytes,
		UsedBytes:    usedBytes,
		FreeBytes:    freeBytes,
		UsagePercent: percent,
		PathSize:     pathSize,
	}, nil
}

// Compress creates an archive (.zip or .tar.gz) from sources.
func (m *Manager) Compress(basePath string, sources []string, destRelPath, format string) error {
	destTarget, err := m.safePath(basePath, destRelPath)
	if err != nil {
		return err
	}

	_ = os.MkdirAll(filepath.Dir(destTarget), 0755)

	outFile, err := os.Create(destTarget)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %w", err)
	}
	defer outFile.Close()

	if strings.ToLower(format) == "tar.gz" || strings.ToLower(format) == "tgz" || strings.HasSuffix(destRelPath, ".tar.gz") {
		return m.compressTarGz(basePath, sources, outFile)
	}
	return m.compressZip(basePath, sources, outFile)
}

func (m *Manager) compressZip(basePath string, sources []string, writer io.Writer) error {
	zw := zip.NewWriter(writer)
	defer zw.Close()

	for _, src := range sources {
		srcTarget, err := m.safePath(basePath, src)
		if err != nil {
			continue
		}

		info, err := os.Stat(srcTarget)
		if err != nil {
			continue
		}

		if info.IsDir() {
			baseDir := filepath.Dir(srcTarget)
			err = filepath.Walk(srcTarget, func(path string, fileInfo os.FileInfo, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}

				rel, err := filepath.Rel(baseDir, path)
				if err != nil {
					return err
				}

				header, err := zip.FileInfoHeader(fileInfo)
				if err != nil {
					return err
				}
				header.Name = filepath.ToSlash(rel)
				if fileInfo.IsDir() {
					header.Name += "/"
				} else {
					header.Method = zip.Deflate
				}

				w, err := zw.CreateHeader(header)
				if err != nil {
					return err
				}

				if !fileInfo.IsDir() {
					file, err := os.Open(path)
					if err != nil {
						return err
					}
					defer file.Close()
					if _, err := io.Copy(w, file); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = filepath.Base(srcTarget)
			header.Method = zip.Deflate

			w, err := zw.CreateHeader(header)
			if err != nil {
				return err
			}

			file, err := os.Open(srcTarget)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(w, file); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Manager) compressTarGz(basePath string, sources []string, writer io.Writer) error {
	gw := gzip.NewWriter(writer)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, src := range sources {
		srcTarget, err := m.safePath(basePath, src)
		if err != nil {
			continue
		}

		info, err := os.Stat(srcTarget)
		if err != nil {
			continue
		}

		if info.IsDir() {
			baseDir := filepath.Dir(srcTarget)
			err = filepath.Walk(srcTarget, func(path string, fileInfo os.FileInfo, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}

				rel, err := filepath.Rel(baseDir, path)
				if err != nil {
					return err
				}

				header, err := tar.FileInfoHeader(fileInfo, "")
				if err != nil {
					return err
				}
				header.Name = filepath.ToSlash(rel)

				if err := tw.WriteHeader(header); err != nil {
					return err
				}

				if !fileInfo.IsDir() {
					file, err := os.Open(path)
					if err != nil {
						return err
					}
					defer file.Close()
					if _, err := io.Copy(tw, file); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			header.Name = filepath.Base(srcTarget)

			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			file, err := os.Open(srcTarget)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(tw, file); err != nil {
				return err
			}
		}
	}
	return nil
}

// Extract extracts archives (.zip, .tar.gz, .tgz, .tar) with Zip Slip / Tar Slip protection.
func (m *Manager) Extract(basePath, archiveRelPath, destRelPath string) error {
	archiveTarget, err := m.safePath(basePath, archiveRelPath)
	if err != nil {
		return err
	}
	destTarget, err := m.safePath(basePath, destRelPath)
	if err != nil {
		return err
	}

	cleanDest := filepath.Clean(destTarget)
	if err := os.MkdirAll(cleanDest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	lower := strings.ToLower(archiveRelPath)
	if strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz") {
		return m.extractTarGz(archiveTarget, cleanDest)
	} else if strings.HasSuffix(lower, ".tar") {
		return m.extractTar(archiveTarget, cleanDest)
	} else {
		return m.ExtractZipFile(archiveTarget, cleanDest)
	}
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
		// Clean and evaluate destination path
		targetPath := filepath.Clean(filepath.Join(cleanDest, file.Name))
		if targetPath != cleanDest && !strings.HasPrefix(targetPath, cleanDest+string(filepath.Separator)) {
			return fmt.Errorf("zip slip detected: illegal file path %s", file.Name)
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

	m.applyOwnership(destDir, cleanDest, true)
	return nil
}

func (m *Manager) extractTarGz(archivePath, destDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gr.Close()

	return m.extractTarReader(gr, destDir)
}

func (m *Manager) extractTar(archivePath, destDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return m.extractTarReader(f, destDir)
}

func (m *Manager) extractTarReader(r io.Reader, destDir string) error {
	tr := tar.NewReader(r)
	cleanDest := filepath.Clean(destDir)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Clean(filepath.Join(cleanDest, header.Name))
		if targetPath != cleanDest && !strings.HasPrefix(targetPath, cleanDest+string(filepath.Separator)) {
			return fmt.Errorf("tar slip detected: illegal file path %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return err
			}
			outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
	m.applyOwnership(destDir, cleanDest, true)
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

// StreamArchive creates an in-memory or streamed zip of given paths directly to writer.
func (m *Manager) StreamArchive(basePath string, paths []string, writer io.Writer) error {
	return m.compressZip(basePath, paths, writer)
}
