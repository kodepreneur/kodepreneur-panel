# Kodepreneur Panel 🚀

[![Laravel 13](https://img.shields.io/badge/Laravel-13.x-FF2D20?style=for-the-badge&logo=laravel)](https://laravel.com)
[![Vue 3](https://img.shields.io/badge/Vue-3.x-4FC08D?style=for-the-badge&logo=vue.js)](https://vuejs.org)
[![Go Agent](https://img.shields.io/badge/Go_Agent-1.24-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Ubuntu 24.04 LTS](https://img.shields.io/badge/Ubuntu-24.04_LTS-E95420?style=for-the-badge&logo=ubuntu)](https://ubuntu.com)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)

**Kodepreneur Panel** is an ultra-fast, lightweight, and modern cloud server control plane engineered for modern web applications. It replaces bloated, slow legacy panels with an unprivileged **Laravel 13 + Inertia.js (Vue 3)** control plane and a root-privileged **Go Agent daemon** communicating over HMAC-SHA256 authenticated Unix/HTTP APIs.

---

## 🏛️ Architecture Overview

```
                          ┌─────────────────────────────┐
                          │   Browser / Admin User      │
                          └──────────────┬──────────────┘
                                         │ HTTPS / Web UI
                                         ▼
                     ┌───────────────────────────────────────┐
                     │   Laravel 13 Control Plane (Panel)   │
                     │  - Vue 3 + Inertia.js SPA             │
                     │  - SQLite Config & Audit Storage     │
                     │  - Unprivileged (runs as www-data)    │
                     └───────────────────┬───────────────────┘
                                         │
                                         │ HMAC-SHA256 Signed API
                                         │ (127.0.0.1:8443 / Unix Socket)
                                         ▼
                     ┌───────────────────────────────────────┐
                     │      Go Agent Daemon (Core)           │
                     │  - Root Privileges (Systemd)          │
                     │  - Real-time CPU, RAM, Disk Telemetry │
                     │  - Atomic Nginx vhost Validator       │
                     │  - PHP-FPM Pool & Isolation Manager   │
                     │  - MySQL & PostgreSQL Engine          │
                     │  - Certbot Automated Let's Encrypt    │
                     │  - Unprivileged Git Deploy Runner     │
                     │  - Crontab, Systemd, UFW, File Jail   │
                     └───────────────────┬───────────────────┘
                                         │
         ┌───────────────────────────────┼───────────────────────────────┐
         ▼                               ▼                               ▼
  ┌──────────────┐                ┌──────────────┐                ┌──────────────┐
  │ Nginx Engine │                │ PHP 8.3/8.4  │                │ DB (MySQL/PG)│
  │ /etc/nginx/  │                │ Pool Sockets │                │ Databases/Users
  └──────────────┘                └──────────────┘                └──────────────┘
```

---

## ⚡ Quick 1-Line Installation & Updates

### Fresh Installation
Run the automated installer on a clean **Ubuntu 24.04 LTS** or **Ubuntu 22.04 LTS** server:

```bash
curl -fsSL https://raw.githubusercontent.com/kodepreneur/kodepreneur-panel/main/installer/install.sh | sudo bash
```

### Upgrade / Update Existing Installation
To upgrade a live server to the latest version directly from GitHub:

```bash
curl -fsSL https://raw.githubusercontent.com/kodepreneur/kodepreneur-panel/main/installer/update.sh | sudo bash
```

Or update via local repository:
```bash
cd /path/to/kodepreneur-panel
git pull origin main
sudo ./installer/update.sh
```

---

## ✨ Features Matrix

### 🌐 Website & Domain Management
- **Atomic Nginx Virtual Hosts**: Automatic configuration generation with pre-reload `nginx -t` validation and auto-rollback on syntax errors.
- **PHP Version Switching**: Zero-downtime runtime switching between **PHP 8.3 LTS** and **PHP 8.4 Latest** per site.
- **Dedicated Pool Isolation**: Each domain runs as an isolated Linux system user (`0750` permissions) with a dedicated FastCGI Unix socket (`/run/php/php{version}-fpm-{user}.sock`).
- **Free Automated SSL**: 1-click Let's Encrypt SSL certificate issuance via Certbot HTTP-01 challenges with auto-renewal and force-HTTPS redirection.
- **Real-Time Logs**: Stream live Nginx `access.log` and `error.log` directly into the web UI.

### 🗄️ Database Engine (MySQL & PostgreSQL)
- **Multi-Engine Support**: Manage both **MySQL / MariaDB** and **PostgreSQL** from a single unified interface.
- **User & Privileges Management**: Create database users with host restrictions (`localhost` or `%`), assign database grants (`ALL PRIVILEGES` or `READ ONLY`), and reset passwords.

### 🚀 Git Deployments & CI/CD
- **Unprivileged Build Runner**: All deployments execute unprivileged as the domain's Linux user (`su - <user>`) with a 600s execution timeout.
- **Automated Build Hooks**: Runs `git pull`, `composer install --no-dev`, `php artisan migrate --force`, `npm run build`, and route/config caching.
- **Live Terminal Log Viewer**: Watch real-time terminal output with ANSI colors, duration tracking, and status badges.
- **1-Click Rollback**: Instantly revert any website to a previous deployment commit.

### ⚙️ Operations & System Tools
- **Cron Job Manager**: Schedule periodic server tasks with a visual 5-field cron builder and automatic crontab synchronization.
- **Managed Background Processes**: Provision and supervise persistent background workers (e.g. `php artisan queue:work`) via managed Systemd application unit files (`/etc/systemd/system/kp-proc-{id}.service`).
- **Path-Jailed File Manager**: Browse, upload, delete, and edit files in the webroot (`/var/www/{domain}`) with a dark code editor modal. Strict path canonicalization prevents directory traversal.
- **UFW Firewall Security**: Open and close TCP/UDP ports, restrict access by IP address, and inspect active firewall rules.
- **System Telemetry**: Real-time CPU, RAM, Disk, System Load, and uptime gauges with 3-second live polling.

---

## 🔒 Security Architecture

1. **Least-Privilege Isolation**: The Laravel web panel runs as unprivileged `www-data` and has **zero** direct root or shell execution privileges.
2. **HMAC-SHA256 Signature Verification**: Every request between the Laravel Panel and the Go Agent requires:
   - `X-Kodepreneur-Timestamp` (rejects requests older than 300s)
   - `X-Kodepreneur-Nonce` (cached to prevent replay attacks)
   - `X-Kodepreneur-Signature` (`hash_hmac('sha256', ...)`)
3. **No Database Password Storage**: Database passwords are never stored in the panel's SQLite database; they are sent directly to the agent during provisioning.

---

## 🛠️ Local Development & Testing

### Running the Go Agent locally (Dev Mode)
```bash
cd agent
go test -v ./...
go run cmd/agent/main.go --dev
```

### Running the Laravel Control Plane
```bash
cd panel
composer install
npm install
npm run build
php artisan test
php artisan serve
```

---

## 📄 License
Kodepreneur Panel is open-source software licensed under the [MIT license](LICENSE).
