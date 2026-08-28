# Kodepreneur Agent Architecture & API Specification

## 1. Overview

The **Kodepreneur Agent** (`kodepreneur-agent`) is a lightweight, compiled Go daemon running with `root` privileges. It serves as the single point of execution for privileged server operations, communicating exclusively with the Laravel Control Plane via a strictly typed, authenticated REST API over a Unix Domain Socket or Localhost HTTPS.

### Key Characteristics
* **Language**: Go 1.22+ (statically compiled binary, zero external runtime dependencies).
* **Binary Location**: `/usr/local/bin/kodepreneur-agent`.
* **Config Location**: `/etc/kodepreneur/agent.yaml`.
* **Execution User**: `root`.
* **Zero Shell Interpretation**: All system actions invoke explicit binaries using `exec.Command` with structured arguments.
* **Strict Non-Generic API**: No arbitrary command execution endpoint (`POST /execute`) exists.

---

## 2. Internal Subsystem Modules

```text
kodepreneur-agent/
├── cmd/
│   └── agent/
│       └── main.go                 # Startup, CLI flags, signal handling
├── config/
│   └── config.go                   # YAML configuration loader
├── server/
│   ├── server.go                   # HTTP server (Socket / Localhost)
│   ├── router.go                   # Route registration
│   └── middleware/
│       ├── auth.go                 # HMAC-SHA256 signature verification
│       ├── logger.go               # Structured JSON request logging
│       └── recovery.go             # Panic recovery & error serialization
└── pkg/
    ├── system/
    │   ├── metrics.go              # CPU, RAM, Disk, Net via /proc & /sys
    │   ├── info.go                 # OS, Kernel, Hostname, Uptime
    │   └── service.go              # Systemd unit management
    ├── nginx/
    │   ├── vhost.go                # Template rendering & atomic write
    │   └── validator.go            # `nginx -t` validation harness
    ├── php/
    │   ├── pools.go                # PHP-FPM pool generation
    │   ├── versions.go             # Detection and extension inspect
    │   └── service.go              # Restart / reload FPM daemons
    ├── database/
    │   ├── mysql.go                # MySQL / MariaDB user & db management
    │   └── postgres.go             # PostgreSQL user & db management
    ├── ssl/
    │   ├── certbot.go              # ACME / Let's Encrypt automation
    │   └── renewal.go              # Expiration check & renewal
    ├── firewall/
    │   └── ufw.go                  # UFW status and rule management
    ├── process/
    │   └── supervisor.go           # Managed systemd application units
    ├── git/
    │   └── runner.go               # Unprivileged Git clone/pull runner
    └── filemanager/
        ├── paths.go                # Path canonicalization & jail check
        └── ops.go                  # Read, write, delete, chmod, list
```

---

## 3. Typed Agent API Reference

All requests must include standard HMAC-SHA256 authentication headers (see [Security Architecture](file:///docs/security.md)).

### 3.1 System & Health

#### `GET /api/v1/health`
Checks if the agent daemon is operational.
* **Response `200 OK`**:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime_seconds": 86400
}
```

#### `GET /api/v1/system/info`
Retrieves core OS and host metadata.
* **Response `200 OK`**:
```json
{
  "hostname": "vps-lon1-01",
  "os": "Ubuntu 24.04 LTS",
  "kernel": "6.8.0-31-generic",
  "architecture": "x86_64",
  "public_ip": "203.0.113.10",
  "uptime_seconds": 142850
}
```

#### `GET /api/v1/system/metrics`
Collects instantaneous hardware utilization directly from Linux `/proc` and `/sys`.
* **Response `200 OK`**:
```json
{
  "cpu": {
    "usage_percent": 14.5,
    "cores": 2
  },
  "memory": {
    "total_bytes": 2097152000,
    "used_bytes": 838860800,
    "free_bytes": 1258291200,
    "usage_percent": 40.0
  },
  "disk": {
    "mount": "/",
    "total_bytes": 53687091200,
    "used_bytes": 10737418240,
    "free_bytes": 42949672960,
    "usage_percent": 20.0
  },
  "load_average": {
    "load1": 0.25,
    "load5": 0.18,
    "load15": 0.12
  }
}
```

---

### 3.2 Service Management

#### `POST /api/v1/services/{service}/{action}`
Controls a registered systemd service (`nginx`, `php8.3-fpm`, `php8.4-fpm`, `mysql`, `postgresql`).
* **Parameters**:
  * `service`: Whitelisted service name.
  * `action`: `start`, `stop`, `restart`, `reload`, `status`.
* **Response `200 OK`**:
```json
{
  "service": "nginx",
  "action": "restart",
  "success": true,
  "output": ""
}
```

---

### 3.3 Website Management

#### `POST /api/v1/websites`
Creates an isolated website, system user, web root, PHP-FPM pool, and Nginx virtual host.
* **Request Body**:
```json
{
  "domain": "example.com",
  "aliases": ["www.example.com"],
  "php_version": "8.3",
  "document_root": "/var/www/example.com/public",
  "system_user": "kp_example",
  "ssl_enabled": false
}
```
* **Execution Flow**:
  1. Creates Linux user `kp_example` with `/var/www/example.com` home.
  2. Generates PHP-FPM pool in `/etc/php/8.3/fpm/pool.d/example.com.conf`.
  3. Generates Nginx configuration in `/etc/nginx/sites-available/example.com.conf`.
  4. Runs `nginx -t`. If invalid, rolls back files and returns `422 Unprocessable Entity`.
  5. Links to `/etc/nginx/sites-enabled/example.com.conf` and reloads Nginx.
* **Response `201 Created`**:
```json
{
  "success": true,
  "domain": "example.com",
  "vhost_path": "/etc/nginx/sites-available/example.com.conf"
}
```

#### `DELETE /api/v1/websites/{domain}`
Removes website configurations and safely archives/deletes web root.
* **Request Body**:
```json
{
  "delete_files": false
}
```
* **Response `200 OK`**:
```json
{
  "success": true,
  "message": "Website configurations removed and Nginx reloaded"
}
```

---

### 3.4 SSL / ACME Management

#### `POST /api/v1/ssl/issue`
Requests and installs a Let's Encrypt SSL certificate via ACME / Certbot HTTP-01 challenge.
* **Request Body**:
```json
{
  "domain": "example.com",
  "aliases": ["www.example.com"],
  "email": "admin@example.com",
  "force_https": true
}
```
* **Response `200 OK`**:
```json
{
  "success": true,
  "certificate_path": "/etc/letsencrypt/live/example.com/fullchain.pem",
  "valid_until": "2026-11-26T12:00:00Z"
}
```

---

### 3.5 Database Management

#### `POST /api/v1/databases`
Creates a managed database and grants privileges to a database user.
* **Request Body**:
```json
{
  "engine": "mysql",
  "database_name": "app_production",
  "username": "app_user",
  "password": "SecurePassword123!",
  "host": "localhost"
}
```
* **Response `201 Created`**:
```json
{
  "success": true,
  "engine": "mysql",
  "database": "app_production",
  "username": "app_user"
}
```

---

### 3.6 File Operations (Sandboxed)

#### `POST /api/v1/files/browse`
Lists directory contents strictly within an authorized base path.
* **Request Body**:
```json
{
  "base_path": "/var/www/example.com",
  "relative_path": "public"
}
```
* **Response `200 OK`**:
```json
{
  "path": "/var/www/example.com/public",
  "entries": [
    {
      "name": "index.php",
      "type": "file",
      "size_bytes": 1420,
      "permissions": "0644",
      "modified_at": "2026-08-28T04:00:00Z"
    }
  ]
}
```

---

### 3.7 Git Deployment Execution

#### `POST /api/v1/deployments/execute`
Executes an unprivileged Git pull/clone and build sequence as the site's system user.
* **Request Body**:
```json
{
  "system_user": "kp_example",
  "working_dir": "/var/www/example.com",
  "repository": "https://github.com/org/repo.git",
  "branch": "main",
  "commands": [
    "git pull origin main",
    "composer install --no-dev --optimize-autoloader",
    "php artisan migrate --force",
    "php artisan config:cache"
  ]
}
```
* **Response `200 OK`**:
```json
{
  "success": true,
  "exit_code": 0,
  "log_output": "[git pull] Already up to date.\n[composer] Installing dependencies...\n..."
}
```
*(Commands are executed strictly as `system_user` with environment scrubbing and a 600-second execution timeout).*
