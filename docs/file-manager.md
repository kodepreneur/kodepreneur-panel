# Kodepreneur File Manager Specification & Architecture

## 1. Current Implementation Analysis

The initial File Manager implementation provided basic browsing, single-file reading, writing, and deletion. However, routine VPS server and website operations require full-featured file management comparable to cPanel, Plesk, and CloudPanel.

### Limitations of Previous Implementation:
* **Navigation**: Dependent on full Inertia page reloads for folder browsing rather than snappy AJAX/API updates.
* **Operations**: Missing file/folder creation, copying, moving, renaming, permissions (`chmod`), ownership (`chown`), search, and multi-item bulk operations.
* **Archives**: Zip extraction existed in the Agent package, but compression (`.zip`, `.tar.gz`), format selection, and UI extraction were not integrated.
* **Uploads & Downloads**: Missing drag-and-drop chunked/streaming uploads, progress bars, and streaming multi-file/directory zip downloads.
* **Editor & Preview**: Basic single-textarea editor lacked line numbers, syntax highlighting cues, search & replace, word wrap toggle, and binary/media preview capabilities.
* **Security & Symlink Containment**: Basic path prefix checks needed enhanced symlink resolution (`filepath.EvalSymlinks`), null-byte sanitization, Zip Slip & Tar Slip guarantees, and strict website webroot boundary enforcement.

---

## 2. Architecture & System Flow

The File Manager adheres strictly to the Kodepreneur three-tier architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    Browser Client (Vue 3)                   │
│   • Responsive Explorer  • Drag & Drop Upload  • Code Editor│
│   • Context Menus        • Modal Dialogs       • Previews   │
└──────────────────────────────┬──────────────────────────────┘
                               │ AJAX / JSON API & Downloads
                               │ (CSRF + Session / Sanctum)
┌──────────────────────────────▼──────────────────────────────┐
│             Kodepreneur Panel (Laravel 11 / PHP)            │
│   • Website Sandboxing   • Auth & Policy Checks             │
│   • Activity Logging     • Request Validation & Sanitization│
│   • Stream Proxying      • Agent Client Interface           │
└──────────────────────────────┬──────────────────────────────┘
                               │ HMAC SHA-256 Signed HTTP API
                               │ (Internal Unix Socket or Port)
┌──────────────────────────────▼──────────────────────────────┐
│               Kodepreneur Agent (Go Daemon)                 │
│   • Canonical Path Validation  • Filesystem OS Operations   │
│   • Archive Compress/Extract   • Streaming Upload/Download  │
│   • Stat & Disk Usage Metrics  • Chmod / Chown Isolation    │
└──────────────────────────────┬──────────────────────────────┘
                               │ POSIX / Linux System Calls
┌──────────────────────────────▼──────────────────────────────┐
│                    Host Linux Filesystem                    │
│      (/var/www/example.com, /var/www/other.com, etc.)       │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Filesystem Security Model & Sandboxing

Filesystem operations are privileged and require defense-in-depth at both the **Panel (PHP)** and **Agent (Go)** layers.

### 3.1 Base Path Jail Enforcement
* Every operation is bound to a validated `base_path` (e.g. `/var/www/example.com`).
* The Panel ensures users cannot supply arbitrary root paths (such as `/etc`, `/root`, `/proc`, `/sys`) unless explicitly authorized.
* The Agent verifies that `canonical(base_path + relative_path)` strictly resides within `canonical(base_path)`.

### 3.2 Path Traversal & Injection Prevention
* **Null-byte rejection**: Strings containing `\x00` are rejected immediately.
* **Traversal tokens**: `..`, `/../`, `%2e%2e%2f` and double-encoded variants are normalized and verified against jail boundaries.
* **Symlink Escape Protection**: `filepath.EvalSymlinks` is evaluated to ensure symlinks pointing outside the sandbox root are rejected.
* **Root Protection**: Deletion, moving, or renaming the base sandbox root (`/var/www/domain`) itself is strictly prohibited.

### 3.3 Archive Slip Protection (Zip Slip & Tar Slip)
* When uncompressing archives, every internal entry path is normalized.
* If any entry resolves outside the destination folder (e.g. `../../etc/passwd` or `/etc/cron.d/evil`), the entire extraction job aborts with an error before writing.

---

## 4. API Endpoints Matrix

### Panel API Endpoints (`panel/routes/web.php` & `panel/routes/api.php`)

| Method | Endpoint | Description | Auth / Middleware |
|---|---|---|---|
| `GET` | `/files` | Inertia file manager page view | `auth` |
| `POST` | `/files/browse` | AJAX directory listing (JSON) | `auth` |
| `GET` | `/files/read` | Read text file content (JSON) | `auth` |
| `POST` | `/files/write` | Save file content (JSON) | `auth` |
| `POST` | `/files/create-file` | Create new blank file | `auth` |
| `POST` | `/files/create-folder` | Create new directory | `auth` |
| `POST` | `/files/delete` | Delete file(s) or folder(s) | `auth` |
| `POST` | `/files/rename` | Rename a file or directory | `auth` |
| `POST` | `/files/copy` | Copy file(s) or folder(s) | `auth` |
| `POST` | `/files/move` | Move / cut file(s) or folder(s) | `auth` |
| `POST` | `/files/chmod` | Change file/directory permissions | `auth` |
| `POST` | `/files/chown` | Change file/directory ownership | `auth` |
| `POST` | `/files/compress` | Compress items into .zip / .tar.gz | `auth` |
| `POST` | `/files/extract` | Extract archive into destination | `auth` |
| `POST` | `/files/search` | Search files by name | `auth` |
| `POST` | `/files/stat` | Get detailed file/directory metadata | `auth` |
| `POST` | `/files/upload` | Multipart file upload | `auth` |
| `GET` | `/files/download` | Download single file or zip stream | `auth` |
| `GET` | `/files/preview` | Stream file for inline preview (image/pdf) | `auth` |

---

## 5. Agent Operations (Go Daemon)

The Go Agent daemon (`agent/pkg/filemanager/filemanager.go` and `agent/server/router.go`) exposes low-overhead, controlled operations:

1. `Browse(basePath, relativePath string, showHidden bool)`: Scans single directory, returns names, sizes, permissions, mod times, MIME types, and item counts.
2. `ReadFile(basePath, relativePath string, maxBytes int64)`: Reads text files safely up to configurable limits (e.g. 5MB).
3. `WriteFile(basePath, relativePath, content string)`: Atomically writes file content.
4. `CreateFile(basePath, relativePath string)`: Creates empty file if non-existent.
5. `CreateDirectory(basePath, relativePath string)`: Creates directory tree (`0755`).
6. `DeleteEntry(basePath, relativePath string)`: Recursively removes path (except root).
7. `Rename(basePath, oldPath, newPath string)`: Renames/moves item within sandbox.
8. `Copy(basePath, srcPath, destPath string)`: Copies file or recursively copies directory.
9. `Move(basePath, srcPath, destPath string)`: Moves file or directory.
10. `Chmod(basePath, relativePath string, mode os.FileMode, recursive bool)`: Applies POSIX permissions.
11. `Chown(basePath, relativePath string, uid, gid int, recursive bool)`: Applies POSIX user/group.
12. `Stat(basePath, relativePath string)`: Returns detailed POSIX metadata (size, mod/access/change time, inode, uid, gid, permissions).
13. `Compress(basePath string, sources []string, destRelPath, format string)`: Archives into `.zip` or `.tar.gz`.
14. `Extract(basePath, archiveRelPath, destRelPath string)`: Extracts `.zip`, `.tar.gz`, `.tgz`, `.tar` with Zip Slip validation.
15. `Search(basePath, query string, maxResults int)`: Non-blocking scoped file search.
16. `GetDiskUsage(basePath string)`: Returns total/used/free disk space on partition and path size.
17. `DownloadArchive(basePath string, paths []string, writer io.Writer)`: Dynamically generates zip archive stream to HTTP response without temporary disk storage.

---

## 6. Authentication & Authorization

* **Panel Level**: Standard Laravel Session / Sanctum token authentication with role-based permission checks.
* **Agent Level**: HMAC SHA-256 signature validation on all requests using `X-Kodepreneur-Timestamp`, `X-Kodepreneur-Nonce`, and shared secret.
* **Website Isolation**: Users can only access file roots matching their assigned websites or global webroot if superadmin.

---

## 7. Upload & Download Flow

### 7.1 Upload Flow
1. User drops files or selects them via browser UI.
2. Frontend sends files via `FormData` with progress tracking (`XMLHttpRequest` / `axios` `onUploadProgress`).
3. Panel receives multipart stream, validates file size against configured limit, and streams directly to Agent.
4. Agent writes file to destination path safely.

### 7.2 Download Flow
1. **Single File**: Streamed directly with `Content-Disposition: attachment; filename="..."` and accurate `Content-Type`.
2. **Directory / Multi-Select**: Dynamic on-the-fly Zip archive streaming. No orphan zip files left on the disk.

---

## 8. Archive Management

* Supported compression formats: `.zip`, `.tar.gz`, `.tgz`, `.tar`.
* Compression includes directory tree structure preservation.
* Extraction supports auto-detecting format from extension and extracting to selected destination directory.
* Pre-flight checks guarantee no destination collisions or path escapes.

---

## 9. File Editor & Preview

### 9.1 Professional Code Editor
* Clean monospaced editing environment with line numbers.
* Syntax highlighting hints for PHP, HTML, CSS, JavaScript, TypeScript, JSON, YAML, SQL, Shell, Markdown, and Nginx conf.
* Keyboard shortcuts (`Cmd/Ctrl + S` to save, `Esc` to close).
* Unsaved changes indicator and confirmation before navigating away.
* Search & Replace bar inside the editor.
* Word wrap toggle and indentation helpers.

### 9.2 File Preview
* **Images**: Direct rendering (PNG, JPG, JPEG, GIF, SVG, WEBP) with dimensions and file size.
* **PDF**: Embedded PDF viewer object.
* **Text / Log**: Read-only fast preview with line counts.
* **Binary Files**: Fallback property summary card.

---

## 10. Permissions & Ownership

* **Permissions Modal**:
  * 3x3 POSIX permission matrix (Owner: R/W/X, Group: R/W/X, Others: R/W/X).
  * Real-time Octal calculation (e.g. `0755`, `0644`).
  * "Apply recursively to subdirectories and files" toggle.
* **Ownership Modal**:
  * Dropdown/input for user and group (e.g. `www-data`, website user).
  * Restrictable to administrative roles.

---

## 11. Audit Logging

Every destructive or configuration-altering action is persisted to `activity_logs`:
* `file.create`, `file.upload`, `file.download`, `file.delete`, `file.rename`, `file.copy`, `file.move`, `file.compress`, `file.extract`, `file.write`, `file.chmod`, `file.chown`.
* Captured metadata: `user_id`, `user_email`, `ip_address`, `action`, `resource_type: 'file'`, `resource_id: relative_path`, `status: 'success'|'failure'`, `payload_summary`.
* File contents and sensitive secrets are **never** logged.

---

## 12. Error Handling & User Feedback

* Raw internal errors (e.g. POSIX `openat ENOENT`, `permission denied`) are caught and transformed into user-friendly messages:
  * "Unable to read directory: the directory does not exist or permission is denied."
  * "Cannot overwrite existing file without confirmation."
* Detailed technical errors are retained in server / agent logs for administrator inspection.

---

## 13. Performance Considerations

* **Non-blocking current directory listing**: Never recursively scans subdirectories when browsing a folder.
* **Memory Efficiency**: Chunked streaming for large uploads and downloads. Large text files (>5MB) are safeguarded against memory blowouts.
* **AJAX Navigation**: Folder changes load instantaneously via JSON API without reloading layout or destroying UI state.
* **Stat Cache**: Rapid retrieval of file size and modified timestamps using standard POSIX `stat`.

---

## 14. Testing Strategy

* **Agent Unit & Integration Tests (Go)**:
  * Path traversal attacks (`../../etc/passwd`, `/../../etc/shadow`, URL-encoded paths).
  * Zip Slip and Tar Slip attack vectors.
  * CRUD operations (create, read, write, copy, move, rename, delete, chmod).
  * In-memory zip creation and streaming.
* **Panel Feature Tests (PHP / Laravel)**:
  * Authentication & authorization enforcement.
  * JSON API endpoints responses & error conditions.
  * Activity logging verification.
  * Single and bulk operations.
* **Frontend Verification**:
  * TypeScript typecheck (`vue-tsc --noEmit`).
  * Production bundle compilation (`vite build`).
