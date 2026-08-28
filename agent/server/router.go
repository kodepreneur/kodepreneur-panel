package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/kodepreneur/agent/config"
	"github.com/kodepreneur/agent/pkg/cron"
	"github.com/kodepreneur/agent/pkg/database"
	"github.com/kodepreneur/agent/pkg/filemanager"
	"github.com/kodepreneur/agent/pkg/firewall"
	"github.com/kodepreneur/agent/pkg/git"
	"github.com/kodepreneur/agent/pkg/nginx"
	"github.com/kodepreneur/agent/pkg/phpfpm"
	"github.com/kodepreneur/agent/pkg/process"
	"github.com/kodepreneur/agent/pkg/ssl"
	"github.com/kodepreneur/agent/pkg/system"
	"github.com/kodepreneur/agent/server/middleware"
)

type Router struct {
	cfg               *config.Config
	mux               *http.ServeMux
	nginxManager      *nginx.Manager
	logManager        *nginx.LogManager
	phpManager        *phpfpm.Manager
	sslManager        *ssl.Manager
	dbManager         *database.Manager
	gitRunner         *git.Runner
	cronManager       *cron.Manager
	processSupervisor *process.Supervisor
	firewallManager   *firewall.Manager
	fileManager       *filemanager.Manager
	updater           *system.UpdateRunner
}

func NewRouter(cfg *config.Config) *Router {
	r := &Router{
		cfg:               cfg,
		mux:               http.NewServeMux(),
		nginxManager:      nginx.NewManager(cfg.Environment.IsDev),
		logManager:        nginx.NewLogManager(cfg.Environment.IsDev),
		phpManager:        phpfpm.NewManager(cfg.Environment.IsDev),
		sslManager:        ssl.NewManager(cfg.Environment.IsDev),
		dbManager:         database.NewManager(cfg.Environment.IsDev),
		gitRunner:         git.NewRunner(cfg.Environment.IsDev),
		cronManager:       cron.NewManager(cfg.Environment.IsDev),
		processSupervisor: process.NewSupervisor(cfg.Environment.IsDev),
		firewallManager:   firewall.NewManager(cfg.Environment.IsDev),
		fileManager:       filemanager.NewManager(cfg.Environment.IsDev),
		updater:           system.NewUpdateRunner(cfg.Environment.IsDev),
	}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
	// Health check
	healthHandler := http.HandlerFunc(r.handleHealth)
	r.mux.Handle("/health", healthHandler)
	r.mux.Handle("/api/v1/health", healthHandler)

	// System Information & Telemetry & Updates
	r.mux.Handle("/api/v1/system/info", http.HandlerFunc(r.handleSystemInfo))
	r.mux.Handle("/api/v1/system/metrics", http.HandlerFunc(r.handleSystemMetrics))
	r.mux.Handle("/api/v1/system/update", http.HandlerFunc(r.handleSystemUpdate))

	// Service Management
	r.mux.Handle("/api/v1/services/", http.HandlerFunc(r.handleServices))

	// Website Management
	r.mux.Handle("/api/v1/websites", http.HandlerFunc(r.handleWebsites))
	r.mux.Handle("/api/v1/websites/", http.HandlerFunc(r.handleWebsiteSubroutes))

	// SSL Management
	r.mux.Handle("/api/v1/ssl/issue", http.HandlerFunc(r.handleSslIssue))

	// Database Management
	r.mux.Handle("/api/v1/databases", http.HandlerFunc(r.handleDatabases))
	r.mux.Handle("/api/v1/databases/", http.HandlerFunc(r.handleDatabaseSubroutes))

	// Git Deployments (Phase 5)
	r.mux.Handle("/api/v1/deployments/execute", http.HandlerFunc(r.handleDeploymentExecute))

	// Cron Jobs (Phase 6)
	r.mux.Handle("/api/v1/cron/sync", http.HandlerFunc(r.handleCronSync))

	// Process Supervisor (Phase 6)
	r.mux.Handle("/api/v1/processes", http.HandlerFunc(r.handleProcesses))
	r.mux.Handle("/api/v1/processes/", http.HandlerFunc(r.handleProcessSubroutes))

	// Firewall (Phase 6)
	r.mux.Handle("/api/v1/firewall/rules", http.HandlerFunc(r.handleFirewallRules))

	// File Manager (Phase 6)
	r.mux.Handle("/api/v1/files/browse", http.HandlerFunc(r.handleFileBrowse))
	r.mux.Handle("/api/v1/files/read", http.HandlerFunc(r.handleFileRead))
	r.mux.Handle("/api/v1/files/write", http.HandlerFunc(r.handleFileWrite))
	r.mux.Handle("/api/v1/files/create", http.HandlerFunc(r.handleFileCreate))
	r.mux.Handle("/api/v1/files/mkdir", http.HandlerFunc(r.handleFileMkdir))
	r.mux.Handle("/api/v1/files/delete", http.HandlerFunc(r.handleFileDelete))
	r.mux.Handle("/api/v1/files/rename", http.HandlerFunc(r.handleFileRename))
	r.mux.Handle("/api/v1/files/copy", http.HandlerFunc(r.handleFileCopy))
	r.mux.Handle("/api/v1/files/move", http.HandlerFunc(r.handleFileMove))
	r.mux.Handle("/api/v1/files/chmod", http.HandlerFunc(r.handleFileChmod))
	r.mux.Handle("/api/v1/files/chown", http.HandlerFunc(r.handleFileChown))
	r.mux.Handle("/api/v1/files/stat", http.HandlerFunc(r.handleFileStat))
	r.mux.Handle("/api/v1/files/compress", http.HandlerFunc(r.handleFileCompress))
	r.mux.Handle("/api/v1/files/extract", http.HandlerFunc(r.handleFileExtract))
	r.mux.Handle("/api/v1/files/search", http.HandlerFunc(r.handleFileSearch))
	r.mux.Handle("/api/v1/files/disk", http.HandlerFunc(r.handleFileDiskUsage))
	r.mux.Handle("/api/v1/files/download", http.HandlerFunc(r.handleFileDownload))
}

func (r *Router) Handler() http.Handler {
	var handler http.Handler = r.mux
	handler = middleware.AuthMiddleware(r.cfg, handler)
	handler = middleware.LoggerMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)
	return handler
}

func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{
		"status":  "healthy",
		"version": "1.0.0",
		"dev":     r.cfg.Environment.IsDev,
	})
}

func (r *Router) handleSystemInfo(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	info := system.GetSystemInfo()
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    info,
	})
}

func (r *Router) handleSystemMetrics(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	metrics := system.GetSystemMetrics()
	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    metrics,
	})
}

func (r *Router) handleSystemUpdate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var updateReq system.UpdateRequest
	if req.Body != nil {
		_ = json.NewDecoder(req.Body).Decode(&updateReq)
	}

	result, err := r.updater.Execute(updateReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": result.Success,
		"data":    result,
	})
}

func (r *Router) handleServices(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/services/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 2 {
		respondError(w, http.StatusBadRequest, "INVALID_PATH", "Format must be /api/v1/services/{service}/{action}")
		return
	}

	serviceName := parts[0]
	action := parts[1]

	if req.Method == http.MethodGet && action == "status" {
		status, err := system.GetServiceStatus(serviceName)
		if err != nil {
			respondError(w, http.StatusBadRequest, "SERVICE_ERROR", err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    status,
		})
		return
	}

	if req.Method == http.MethodPost {
		out, err := system.ManageService(serviceName, action)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "EXEC_FAILED", err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"service": serviceName,
				"action":  action,
				"output":  out,
			},
		})
		return
	}

	respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
}

// handleWebsites handles POST /api/v1/websites
func (r *Router) handleWebsites(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		Domain           string   `json:"domain"`
		Aliases          []string `json:"aliases"`
		PhpVersion       string   `json:"php_version"`
		DocumentRoot     string   `json:"document_root"`
		SystemUser       string   `json:"system_user"`
		SslEnabled       bool     `json:"ssl_enabled"`
		ForceHttps       bool     `json:"force_https"`
		DeploymentSource string   `json:"deployment_source"` // "empty", "zip", "git"
		GitRepository    string   `json:"git_repository"`
		GitBranch        string   `json:"git_branch"`
		ProjectType      string   `json:"project_type"` // "laravel", "generic_php", "static", "auto"
		ZipBase64        string   `json:"zip_base64"`
		ZipPath          string   `json:"zip_path"`
		LaravelSetup     *struct {
			Enabled        bool              `json:"enabled"`
			SetupEnv       bool              `json:"setup_env"`
			EnvVars        map[string]string `json:"env_vars"`
			RunComposer    bool              `json:"run_composer"`
			RunKeyGenerate bool              `json:"run_key_generate"`
			RunMigrations  bool              `json:"run_migrations"`
			RunSeeders     bool              `json:"run_seeders"`
			RunNpmBuild    bool              `json:"run_npm_build"`
			RunOptimize    bool              `json:"run_optimize"`
		} `json:"laravel_setup"`
	}

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON request payload")
		return
	}

	if payload.Domain == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Domain is required")
		return
	}

	if payload.SystemUser == "" {
		payload.SystemUser = fmt.Sprintf("kp_%s", strings.ReplaceAll(strings.Split(payload.Domain, ".")[0], "-", "_"))
	}
	if payload.PhpVersion == "" {
		payload.PhpVersion = "8.3"
	}
	if payload.DeploymentSource == "" {
		payload.DeploymentSource = "empty"
	}

	baseDir := fmt.Sprintf("/var/www/%s", payload.Domain)
	if payload.DocumentRoot != "" && strings.HasSuffix(payload.DocumentRoot, "/public") {
		baseDir = filepath.Dir(payload.DocumentRoot)
	}
	realBaseDir := baseDir
	if r.cfg.Environment.IsDev && strings.HasPrefix(baseDir, "/var/www") {
		realBaseDir = filepath.Join(os.TempDir(), "kodepreneur", "www", strings.TrimPrefix(baseDir, "/var/www"))
		_ = os.MkdirAll(realBaseDir, 0755)
	}

	if err := r.phpManager.ProvisionUser(payload.SystemUser, payload.Domain); err != nil {
		respondError(w, http.StatusInternalServerError, "USER_PROVISION_FAILED", err.Error())
		return
	}

	// 1. Handle Deployment Sources
	switch payload.DeploymentSource {
	case "git":
		if payload.GitRepository != "" {
			if err := r.gitRunner.CloneRepo(payload.GitRepository, payload.GitBranch, realBaseDir, payload.SystemUser); err != nil {
				respondError(w, http.StatusInternalServerError, "GIT_CLONE_FAILED", err.Error())
				return
			}
		}
	case "zip":
		if payload.ZipBase64 != "" {
			zipBytes, err := base64.StdEncoding.DecodeString(payload.ZipBase64)
			if err != nil {
				respondError(w, http.StatusBadRequest, "INVALID_ZIP_PAYLOAD", "Failed to decode base64 zip payload")
				return
			}
			if err := r.fileManager.ExtractZipBytes(zipBytes, realBaseDir); err != nil {
				respondError(w, http.StatusInternalServerError, "ZIP_EXTRACT_FAILED", err.Error())
				return
			}
		} else if payload.ZipPath != "" {
			if err := r.fileManager.ExtractZipFile(payload.ZipPath, realBaseDir); err != nil {
				respondError(w, http.StatusInternalServerError, "ZIP_EXTRACT_FAILED", err.Error())
				return
			}
		}
	case "empty":
		_ = os.MkdirAll(realBaseDir, 0755)
	}

	// 2. Auto-detect Laravel & direct root route to /public
	isLaravel := payload.ProjectType == "laravel"
	if !isLaravel && (payload.ProjectType == "" || payload.ProjectType == "auto") {
		artisanFile := filepath.Join(realBaseDir, "artisan")
		publicDir := filepath.Join(realBaseDir, "public")
		if _, err := os.Stat(artisanFile); err == nil {
			isLaravel = true
		} else if _, err := os.Stat(publicDir); err == nil {
			isLaravel = true
		}
	}

	if isLaravel {
		if payload.DocumentRoot == "" || payload.DocumentRoot == baseDir {
			payload.DocumentRoot = filepath.Join(baseDir, "public")
		}
	} else if payload.DocumentRoot == "" || payload.DocumentRoot == filepath.Join(baseDir, "public") {
		publicDir := filepath.Join(realBaseDir, "public")
		if _, err := os.Stat(publicDir); err == nil {
			payload.DocumentRoot = filepath.Join(baseDir, "public")
		} else {
			payload.DocumentRoot = baseDir
		}
	}

	// 3. Prepare Webroot directory and welcome index if empty
	if err := r.phpManager.PrepareWebroot(payload.DocumentRoot, realBaseDir, payload.SystemUser); err != nil {
		respondError(w, http.StatusInternalServerError, "WEBROOT_CREATION_FAILED", err.Error())
		return
	}

	// 4. Automated Laravel Post-Setup (Env & Build Commands)
	var setupResult *git.DeploymentResult
	if isLaravel && payload.LaravelSetup != nil && payload.LaravelSetup.Enabled {
		if payload.LaravelSetup.SetupEnv && len(payload.LaravelSetup.EnvVars) > 0 {
			_ = configureLaravelEnv(realBaseDir, payload.LaravelSetup.EnvVars, payload.SystemUser)
		}

		var setupCommands []string
		if payload.LaravelSetup.RunComposer {
			setupCommands = append(setupCommands, "if [ -f composer.json ]; then composer install --no-dev --prefer-dist --optimize-autoloader --no-interaction; fi")
		}
		if payload.LaravelSetup.RunKeyGenerate {
			setupCommands = append(setupCommands, "if [ -f artisan ]; then php artisan key:generate --force; fi")
		}
		if payload.LaravelSetup.RunMigrations {
			setupCommands = append(setupCommands, "if [ -f artisan ]; then php artisan migrate --force; fi")
		}
		if payload.LaravelSetup.RunSeeders {
			setupCommands = append(setupCommands, "if [ -f artisan ]; then php artisan db:seed --force; fi")
		}
		if payload.LaravelSetup.RunNpmBuild {
			setupCommands = append(setupCommands, "if [ -f package.json ]; then npm install --silent 2>/dev/null || npm install; npm run build; fi")
		}
		if payload.LaravelSetup.RunOptimize {
			setupCommands = append(setupCommands, "if [ -f artisan ]; then php artisan optimize:clear; php artisan config:cache; php artisan route:cache; php artisan view:cache; fi")
		}
		setupCommands = append(setupCommands, "mkdir -p storage/framework/{sessions,views,cache} bootstrap/cache")
		setupCommands = append(setupCommands, "chmod -R 775 storage bootstrap/cache 2>/dev/null || true")

		if len(setupCommands) > 0 {
			res, err := r.gitRunner.Execute(git.DeploymentRequest{
				SystemUser: payload.SystemUser,
				WorkingDir: realBaseDir,
				Branch:     payload.GitBranch,
				Commands:   setupCommands,
				TimeoutSec: 600,
			})
			if err == nil {
				setupResult = res
			}
		}
	}

	if payload.PhpVersion != "none" {
		poolCfg := phpfpm.PoolConfig{
			Domain:       payload.Domain,
			SystemUser:   payload.SystemUser,
			PhpVersion:   payload.PhpVersion,
			DocumentRoot: payload.DocumentRoot,
		}
		if _, err := r.phpManager.CreatePool(poolCfg); err != nil {
			respondError(w, http.StatusInternalServerError, "PHP_POOL_FAILED", err.Error())
			return
		}
	}

	vhostCfg := nginx.VhostConfig{
		Domain:       payload.Domain,
		Aliases:      payload.Aliases,
		DocumentRoot: payload.DocumentRoot,
		PhpVersion:   payload.PhpVersion,
		SystemUser:   payload.SystemUser,
		SslEnabled:   payload.SslEnabled,
		ForceHttps:   payload.ForceHttps,
	}

	vhostPath, err := r.nginxManager.DeployVhost(vhostCfg)
	if err != nil {
		if payload.PhpVersion != "none" {
			_ = r.phpManager.DeletePool(payload.PhpVersion, payload.SystemUser)
		}
		respondError(w, http.StatusUnprocessableEntity, "NGINX_DEPLOY_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"data": map[string]any{
			"domain":            payload.Domain,
			"vhost_path":        vhostPath,
			"system_user":       payload.SystemUser,
			"document_root":     payload.DocumentRoot,
			"php_version":       payload.PhpVersion,
			"deployment_source": payload.DeploymentSource,
			"is_laravel":        isLaravel,
			"setup_result":      setupResult,
		},
	})
}

// configureLaravelEnv creates or updates .env file with given key-value mappings.
func configureLaravelEnv(baseDir string, envVars map[string]string, systemUser string) error {
	envPath := filepath.Join(baseDir, ".env")
	examplePath := filepath.Join(baseDir, ".env.example")

	var envContent string
	if data, err := os.ReadFile(envPath); err == nil {
		envContent = string(data)
	} else if data, err := os.ReadFile(examplePath); err == nil {
		envContent = string(data)
	} else {
		envContent = "APP_NAME=Laravel\nAPP_ENV=production\nAPP_KEY=\nAPP_DEBUG=false\nAPP_URL=http://localhost\n\nLOG_CHANNEL=stack\nLOG_DEPRECATIONS_CHANNEL=null\nLOG_LEVEL=debug\n\nDB_CONNECTION=mysql\nDB_HOST=127.0.0.1\nDB_PORT=3306\nDB_DATABASE=laravel\nDB_USERNAME=root\nDB_PASSWORD=\n\nBROADCAST_DRIVER=log\nCACHE_DRIVER=file\nFILESYSTEM_DISK=local\nQUEUE_CONNECTION=sync\nSESSION_DRIVER=file\nSESSION_LIFETIME=120\n"
	}

	lines := strings.Split(envContent, "\n")
	foundKeys := make(map[string]bool)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") || !strings.Contains(trimmed, "=") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		key := strings.TrimSpace(parts[0])
		if val, ok := envVars[key]; ok {
			if strings.Contains(val, " ") && !strings.HasPrefix(val, "\"") {
				val = fmt.Sprintf("\"%s\"", val)
			}
			lines[i] = fmt.Sprintf("%s=%s", key, val)
			foundKeys[key] = true
		}
	}

	for k, v := range envVars {
		if !foundKeys[k] {
			if strings.Contains(v, " ") && !strings.HasPrefix(v, "\"") {
				v = fmt.Sprintf("\"%s\"", v)
			}
			lines = append(lines, fmt.Sprintf("%s=%s", k, v))
		}
	}

	finalContent := strings.Join(lines, "\n")
	if err := os.WriteFile(envPath, []byte(finalContent), 0640); err != nil {
		return err
	}

	if runtime.GOOS == "linux" {
		_ = exec.Command("chown", fmt.Sprintf("%s:www-data", systemUser), envPath).Run()
		_ = exec.Command("chmod", "0640", envPath).Run()
	}

	return nil
}

// handleWebsiteSubroutes routes /api/v1/websites/{domain} and /api/v1/websites/{domain}/*
func (r *Router) handleWebsiteSubroutes(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/websites/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		respondError(w, http.StatusBadRequest, "INVALID_PATH", "Domain must be specified")
		return
	}

	domain := parts[0]

	// DELETE /api/v1/websites/{domain}
	if len(parts) == 1 && req.Method == http.MethodDelete {
		var payload struct {
			PhpVersion string `json:"php_version"`
			SystemUser string `json:"system_user"`
		}
		_ = json.NewDecoder(req.Body).Decode(&payload)

		_ = r.nginxManager.RemoveVhost(domain)
		if payload.PhpVersion != "" && payload.SystemUser != "" {
			_ = r.phpManager.DeletePool(payload.PhpVersion, payload.SystemUser)
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": fmt.Sprintf("Website %s removed successfully", domain),
		})
		return
	}

	// PUT /api/v1/websites/{domain}/php
	if len(parts) == 2 && parts[1] == "php" && req.Method == http.MethodPut {
		var payload struct {
			NewPhpVersion string   `json:"new_php_version"`
			OldPhpVersion string   `json:"old_php_version"`
			SystemUser    string   `json:"system_user"`
			DocumentRoot  string   `json:"document_root"`
			Aliases       []string `json:"aliases"`
			SslEnabled    bool     `json:"ssl_enabled"`
			ForceHttps    bool     `json:"force_https"`
			CertPath      string   `json:"cert_path"`
			KeyPath       string   `json:"key_path"`
		}

		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
			return
		}

		poolCfg := phpfpm.PoolConfig{
			Domain:       domain,
			SystemUser:   payload.SystemUser,
			PhpVersion:   payload.NewPhpVersion,
			DocumentRoot: payload.DocumentRoot,
		}
		if err := r.phpManager.SwitchVersion(poolCfg, payload.OldPhpVersion); err != nil {
			respondError(w, http.StatusInternalServerError, "PHP_SWITCH_FAILED", err.Error())
			return
		}

		vhostCfg := nginx.VhostConfig{
			Domain:       domain,
			Aliases:      payload.Aliases,
			DocumentRoot: payload.DocumentRoot,
			PhpVersion:   payload.NewPhpVersion,
			SystemUser:   payload.SystemUser,
			SslEnabled:   payload.SslEnabled,
			ForceHttps:   payload.ForceHttps,
			CertPath:     payload.CertPath,
			KeyPath:      payload.KeyPath,
		}
		if _, err := r.nginxManager.DeployVhost(vhostCfg); err != nil {
			respondError(w, http.StatusUnprocessableEntity, "NGINX_UPDATE_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"domain":      domain,
				"php_version": payload.NewPhpVersion,
			},
		})
		return
	}

	// GET /api/v1/websites/{domain}/logs/{type}
	if len(parts) == 3 && parts[1] == "logs" && req.Method == http.MethodGet {
		logType := parts[2]
		if logType != "access" && logType != "error" {
			respondError(w, http.StatusBadRequest, "INVALID_LOG_TYPE", "Log type must be 'access' or 'error'")
			return
		}

		lines := 100
		if qLines := req.URL.Query().Get("lines"); qLines != "" {
			if l, err := strconv.Atoi(qLines); err == nil && l > 0 {
				lines = l
			}
		}

		logEntries, err := r.logManager.ReadLogLines(domain, logType, lines)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "LOG_READ_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"domain": domain,
				"type":   logType,
				"lines":  logEntries,
			},
		})
		return
	}

	respondError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
}

// handleSslIssue handles POST /api/v1/ssl/issue
func (r *Router) handleSslIssue(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		Domain       string   `json:"domain"`
		Aliases      []string `json:"aliases"`
		Email        string   `json:"email"`
		DocumentRoot string   `json:"document_root"`
		PhpVersion   string   `json:"php_version"`
		SystemUser   string   `json:"system_user"`
		ForceHttps   bool     `json:"force_https"`
	}

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if payload.Domain == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Domain is required")
		return
	}

	certResult, err := r.sslManager.IssueCertificate(payload.Domain, payload.Aliases, payload.Email, payload.DocumentRoot)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "SSL_ISSUANCE_FAILED", err.Error())
		return
	}

	vhostCfg := nginx.VhostConfig{
		Domain:       payload.Domain,
		Aliases:      payload.Aliases,
		DocumentRoot: payload.DocumentRoot,
		PhpVersion:   payload.PhpVersion,
		SystemUser:   payload.SystemUser,
		SslEnabled:   true,
		ForceHttps:   payload.ForceHttps,
		CertPath:     certResult.CertPath,
		KeyPath:      certResult.KeyPath,
	}

	if _, err := r.nginxManager.DeployVhost(vhostCfg); err != nil {
		respondError(w, http.StatusUnprocessableEntity, "NGINX_SSL_DEPLOY_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    certResult,
	})
}

// handleDatabases handles POST /api/v1/databases
func (r *Router) handleDatabases(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		Engine    string `json:"engine"`
		Name      string `json:"name"`
		Charset   string `json:"charset"`
		Collation string `json:"collation"`
	}

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if payload.Name == "" || payload.Engine == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Database engine and name are required")
		return
	}

	if err := r.dbManager.CreateDatabase(payload.Engine, payload.Name, payload.Charset, payload.Collation); err != nil {
		respondError(w, http.StatusInternalServerError, "DB_CREATE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"data": map[string]any{
			"engine": payload.Engine,
			"name":   payload.Name,
		},
	})
}

// handleDatabaseSubroutes routes /api/v1/databases/*
func (r *Router) handleDatabaseSubroutes(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/databases/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// POST /api/v1/databases/users
	if len(parts) == 1 && parts[0] == "users" && req.Method == http.MethodPost {
		var payload struct {
			Engine   string `json:"engine"`
			Username string `json:"username"`
			Host     string `json:"host"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
			return
		}
		if payload.Username == "" || payload.Password == "" || payload.Engine == "" {
			respondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Engine, username, and password are required")
			return
		}

		if err := r.dbManager.CreateUser(payload.Engine, payload.Username, payload.Host, payload.Password); err != nil {
			respondError(w, http.StatusInternalServerError, "USER_CREATE_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, map[string]any{
			"success": true,
			"data": map[string]any{
				"engine":   payload.Engine,
				"username": payload.Username,
				"host":     payload.Host,
			},
		})
		return
	}

	// POST /api/v1/databases/grants
	if len(parts) == 1 && parts[0] == "grants" && req.Method == http.MethodPost {
		var payload struct {
			Engine      string `json:"engine"`
			Database    string `json:"database"`
			Username    string `json:"username"`
			Host        string `json:"host"`
			Permissions string `json:"permissions"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
			return
		}

		if err := r.dbManager.GrantPrivileges(payload.Engine, payload.Database, payload.Username, payload.Host, payload.Permissions); err != nil {
			respondError(w, http.StatusInternalServerError, "GRANT_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": fmt.Sprintf("Granted permissions on %s to %s", payload.Database, payload.Username),
		})
		return
	}

	// POST /api/v1/databases/users/password
	if len(parts) == 2 && parts[0] == "users" && parts[1] == "password" && req.Method == http.MethodPost {
		var payload struct {
			Engine      string `json:"engine"`
			Username    string `json:"username"`
			Host        string `json:"host"`
			NewPassword string `json:"new_password"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
			return
		}

		if err := r.dbManager.ChangePassword(payload.Engine, payload.Username, payload.Host, payload.NewPassword); err != nil {
			respondError(w, http.StatusInternalServerError, "PASSWORD_RESET_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": fmt.Sprintf("Password reset for %s", payload.Username),
		})
		return
	}

	// DELETE /api/v1/databases/users/{engine}/{username}
	if len(parts) == 3 && parts[0] == "users" && req.Method == http.MethodDelete {
		engine := parts[1]
		username := parts[2]
		host := req.URL.Query().Get("host")

		if err := r.dbManager.DropUser(engine, username, host); err != nil {
			respondError(w, http.StatusInternalServerError, "USER_DROP_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": fmt.Sprintf("User %s dropped successfully", username),
		})
		return
	}

	// DELETE /api/v1/databases/{engine}/{name}
	if len(parts) == 2 && req.Method == http.MethodDelete {
		engine := parts[0]
		name := parts[1]

		if err := r.dbManager.DropDatabase(engine, name); err != nil {
			respondError(w, http.StatusInternalServerError, "DB_DROP_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": fmt.Sprintf("Database %s dropped successfully", name),
		})
		return
	}

	respondError(w, http.StatusNotFound, "NOT_FOUND", "Database endpoint not found")
}

// handleDeploymentExecute handles POST /api/v1/deployments/execute
func (r *Router) handleDeploymentExecute(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload git.DeploymentRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	res, err := r.gitRunner.Execute(payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "DEPLOY_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    res,
	})
}

// handleCronSync handles POST /api/v1/cron/sync
func (r *Router) handleCronSync(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		SystemUser string         `json:"system_user"`
		Jobs       []cron.JobItem `json:"jobs"`
	}

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.cronManager.SyncUserCrontab(payload.SystemUser, payload.Jobs); err != nil {
		respondError(w, http.StatusInternalServerError, "CRON_SYNC_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": fmt.Sprintf("Crontab for %s synchronized successfully", payload.SystemUser),
	})
}

// handleProcesses handles POST /api/v1/processes
func (r *Router) handleProcesses(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload process.ProcessConfig
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	unitPath, err := r.processSupervisor.CreateProcess(payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "PROCESS_CREATE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"success":   true,
		"unit_path": unitPath,
	})
}

// handleProcessSubroutes handles /api/v1/processes/{id} and /api/v1/processes/{id}/{action}
func (r *Router) handleProcessSubroutes(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/processes/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) == 0 || parts[0] == "" {
		respondError(w, http.StatusBadRequest, "INVALID_PATH", "Process ID is required")
		return
	}

	id := parts[0]

	// DELETE /api/v1/processes/{id}
	if len(parts) == 1 && req.Method == http.MethodDelete {
		if err := r.processSupervisor.DeleteProcess(id); err != nil {
			respondError(w, http.StatusInternalServerError, "PROCESS_DELETE_FAILED", err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": fmt.Sprintf("Process %s deleted", id),
		})
		return
	}

	// POST /api/v1/processes/{id}/{action}
	if len(parts) == 2 && req.Method == http.MethodPost {
		action := parts[1]
		out, err := r.processSupervisor.ManageProcess(id, action)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "PROCESS_ACTION_FAILED", err.Error())
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"output":  out,
		})
		return
	}

	respondError(w, http.StatusNotFound, "NOT_FOUND", "Endpoint not found")
}

// handleFirewallRules handles POST and DELETE /api/v1/firewall/rules
func (r *Router) handleFirewallRules(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var payload firewall.Rule
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
			return
		}

		if err := r.firewallManager.AddRule(payload); err != nil {
			respondError(w, http.StatusInternalServerError, "FIREWALL_RULE_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, map[string]any{
			"success": true,
			"message": "Firewall rule created",
		})
		return
	}

	if req.Method == http.MethodDelete {
		var payload struct {
			Port     int    `json:"port"`
			Protocol string `json:"protocol"`
			Action   string `json:"action"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
			return
		}

		if err := r.firewallManager.DeleteRule(payload.Port, payload.Protocol, payload.Action); err != nil {
			respondError(w, http.StatusInternalServerError, "FIREWALL_DELETE_FAILED", err.Error())
			return
		}

		respondJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "Firewall rule deleted",
		})
		return
	}

	respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
}

// File Manager Endpoints
func (r *Router) handleFileBrowse(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
		ShowHidden   bool   `json:"show_hidden"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	entries, err := r.fileManager.Browse(payload.BasePath, payload.RelativePath, payload.ShowHidden)
	if err != nil {
		respondError(w, http.StatusBadRequest, "BROWSE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    entries,
	})
}

func (r *Router) handleFileRead(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
		MaxBytes     int64  `json:"max_bytes"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	content, err := r.fileManager.ReadFile(payload.BasePath, payload.RelativePath, payload.MaxBytes)
	if err != nil {
		respondError(w, http.StatusBadRequest, "READ_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"content": content,
	})
}

func (r *Router) handleFileWrite(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
		Content      string `json:"content"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.WriteFile(payload.BasePath, payload.RelativePath, payload.Content); err != nil {
		respondError(w, http.StatusBadRequest, "WRITE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "File written successfully",
	})
}

func (r *Router) handleFileCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.CreateFile(payload.BasePath, payload.RelativePath); err != nil {
		respondError(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "File created successfully",
	})
}

func (r *Router) handleFileMkdir(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.CreateDirectory(payload.BasePath, payload.RelativePath); err != nil {
		respondError(w, http.StatusBadRequest, "MKDIR_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Directory created successfully",
	})
}

func (r *Router) handleFileDelete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.DeleteEntry(payload.BasePath, payload.RelativePath); err != nil {
		respondError(w, http.StatusBadRequest, "DELETE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Entry deleted successfully",
	})
}

func (r *Router) handleFileRename(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath string `json:"base_path"`
		OldPath  string `json:"old_path"`
		NewPath  string `json:"new_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.Rename(payload.BasePath, payload.OldPath, payload.NewPath); err != nil {
		respondError(w, http.StatusBadRequest, "RENAME_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Renamed successfully",
	})
}

func (r *Router) handleFileCopy(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath string `json:"base_path"`
		SrcPath  string `json:"src_path"`
		DestPath string `json:"dest_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.Copy(payload.BasePath, payload.SrcPath, payload.DestPath); err != nil {
		respondError(w, http.StatusBadRequest, "COPY_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Copied successfully",
	})
}

func (r *Router) handleFileMove(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath string `json:"base_path"`
		SrcPath  string `json:"src_path"`
		DestPath string `json:"dest_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.Move(payload.BasePath, payload.SrcPath, payload.DestPath); err != nil {
		respondError(w, http.StatusBadRequest, "MOVE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Moved successfully",
	})
}

func (r *Router) handleFileChmod(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
		Mode         string `json:"mode"`
		Recursive    bool   `json:"recursive"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.Chmod(payload.BasePath, payload.RelativePath, payload.Mode, payload.Recursive); err != nil {
		respondError(w, http.StatusBadRequest, "CHMOD_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Permissions updated successfully",
	})
}

func (r *Router) handleFileChown(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
		UID          int    `json:"uid"`
		GID          int    `json:"gid"`
		Recursive    bool   `json:"recursive"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.Chown(payload.BasePath, payload.RelativePath, payload.UID, payload.GID, payload.Recursive); err != nil {
		respondError(w, http.StatusBadRequest, "CHOWN_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Ownership updated successfully",
	})
}

func (r *Router) handleFileStat(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath     string `json:"base_path"`
		RelativePath string `json:"relative_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	details, err := r.fileManager.Stat(payload.BasePath, payload.RelativePath)
	if err != nil {
		respondError(w, http.StatusBadRequest, "STAT_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    details,
	})
}

func (r *Router) handleFileCompress(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath string   `json:"base_path"`
		Sources  []string `json:"sources"`
		DestPath string   `json:"dest_path"`
		Format   string   `json:"format"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.Compress(payload.BasePath, payload.Sources, payload.DestPath, payload.Format); err != nil {
		respondError(w, http.StatusBadRequest, "COMPRESS_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Archive created successfully",
	})
}

func (r *Router) handleFileExtract(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath    string `json:"base_path"`
		ArchivePath string `json:"archive_path"`
		DestPath    string `json:"dest_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if err := r.fileManager.Extract(payload.BasePath, payload.ArchivePath, payload.DestPath); err != nil {
		respondError(w, http.StatusBadRequest, "EXTRACT_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Archive extracted successfully",
	})
}

func (r *Router) handleFileSearch(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath   string `json:"base_path"`
		Query      string `json:"query"`
		MaxResults int    `json:"max_results"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	results, err := r.fileManager.Search(payload.BasePath, payload.Query, payload.MaxResults)
	if err != nil {
		respondError(w, http.StatusBadRequest, "SEARCH_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    results,
	})
}

func (r *Router) handleFileDiskUsage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath string `json:"base_path"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	usage, err := r.fileManager.GetDiskUsage(payload.BasePath)
	if err != nil {
		respondError(w, http.StatusBadRequest, "DISK_USAGE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    usage,
	})
}

func (r *Router) handleFileDownload(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var payload struct {
		BasePath string   `json:"base_path"`
		Paths    []string `json:"paths"`
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	if len(payload.Paths) == 0 {
		respondError(w, http.StatusBadRequest, "EMPTY_PATHS", "No paths provided for download")
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"download.zip\"")
	if err := r.fileManager.StreamArchive(payload.BasePath, payload.Paths, w); err != nil {
		// Log error if header was already written
		return
	}
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
