package server

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	r.mux.Handle("/api/v1/files/delete", http.HandlerFunc(r.handleFileDelete))
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
		Domain       string   `json:"domain"`
		Aliases      []string `json:"aliases"`
		PhpVersion   string   `json:"php_version"`
		DocumentRoot string   `json:"document_root"`
		SystemUser   string   `json:"system_user"`
		SslEnabled   bool     `json:"ssl_enabled"`
		ForceHttps   bool     `json:"force_https"`
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
	if payload.DocumentRoot == "" {
		payload.DocumentRoot = fmt.Sprintf("/var/www/%s/public", payload.Domain)
	}
	if payload.PhpVersion == "" {
		payload.PhpVersion = "8.3"
	}

	if err := r.phpManager.ProvisionUser(payload.SystemUser, payload.Domain); err != nil {
		respondError(w, http.StatusInternalServerError, "USER_PROVISION_FAILED", err.Error())
		return
	}

	if err := r.phpManager.PrepareWebroot(payload.DocumentRoot, payload.SystemUser); err != nil {
		respondError(w, http.StatusInternalServerError, "WEBROOT_CREATION_FAILED", err.Error())
		return
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
			"domain":        payload.Domain,
			"vhost_path":    vhostPath,
			"system_user":   payload.SystemUser,
			"document_root": payload.DocumentRoot,
			"php_version":   payload.PhpVersion,
		},
	})
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
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	entries, err := r.fileManager.Browse(payload.BasePath, payload.RelativePath)
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
	}
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_PAYLOAD", "Invalid JSON payload")
		return
	}

	content, err := r.fileManager.ReadFile(payload.BasePath, payload.RelativePath)
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
