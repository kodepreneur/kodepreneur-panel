# Development Roadmap & Implementation Phases

## 1. MVP Scope Definition

The Minimum Viable Product (MVP) of **Kodepreneur Panel** is complete when an administrator can boot a fresh **Ubuntu 24.04 LTS VPS** and execute the entire lifecycle without SSH:

```text
Run Modular Bootstrap Installer
              ↓
Login to Clean Web UI
              ↓
View Instant Server Hardware Metrics (CPU, RAM, Disk, Load)
              ↓
Create a Website with Domain & PHP 8.3 / 8.4 Pool
              ↓
Issue Automated Let's Encrypt SSL Certificate
              ↓
Provision MySQL / PostgreSQL Database & User
              ↓
Deploy Web Application from Git Repository
              ↓
View Deployment Build Logs & Server Traffic Logs
              ↓
Manage Background Workers & Cron Jobs
```

---

## 2. Implementation Phases

### Phase 0: Discovery & Documentation (Current Phase)
* [x] Inspect workspace and initialize project structure.
* [x] Complete system architecture blueprint (`docs/architecture.md`).
* [x] Document modular installer lifecycle (`docs/installation.md`).
* [x] Document security threat model and privilege separation (`docs/security.md`).
* [x] Document Go Agent architecture and typed API (`docs/agent.md`).
* [x] Document database schema and SQLite PRAGMAs (`docs/database.md`).
* [x] Document Git deployment engine (`docs/deployment.md`).
* [x] Document troubleshooting runbook (`docs/troubleshooting.md`).
* [x] Create Architecture Decision Records (ADR-001 through ADR-004).
* [x] Publish implementation roadmap (`docs/roadmap.md`) and changelog (`docs/changelog.md`).

---

### Phase 1: Project Foundation & Agent Skeleton
* [x] Initialize Laravel 13 application with Inertia.js (Vue 3 + TypeScript + Tailwind CSS).
* [x] Implement authentication (Session, Sanctum API tokens) and RBAC middleware.
* [x] Set up SQLite database with optimized PRAGMAs and core migrations.
* [x] Create dark/light modern developer-focused UI layout with responsive sidebar.
* [x] Initialize Go Agent project (`agent/`) with HTTP server, HMAC auth middleware, and mock/dev modes.
* [x] Build typed PHP Agent client with mock transport for local non-Linux development.

---

### Phase 2: Server Dashboard & Telemetry
* [x] Implement Go Agent metrics collector (`/proc/stat`, `/proc/meminfo`, `/proc/loadavg`, `df`).
* [x] Implement Systemd service inspector (`nginx`, `php-fpm`, `mysql`, `postgresql`, `kodepreneur-agent`).
* [x] Build Dashboard UI cards (CPU, RAM, Disk, Network, Load Average, OS metadata).
* [x] Implement lightweight polling telemetry service.
* [x] Implement service start / stop / restart actions from Dashboard UI.

---

### Phase 3: Website & PHP Runtime Management
* [x] Build Nginx virtual host template engine in Go Agent (`nginx -t` pre-check).
* [x] Build PHP-FPM pool manager (PHP 8.3 & PHP 8.4 pool isolation).
* [x] Implement Website CRUD in Laravel Control Plane with audit logging.
* [x] Integrate Certbot / Let's Encrypt automated HTTP-01 challenge issuance.
* [x] Build Website Management UI (Domain configuration, PHP version switcher, SSL toggle, Nginx log viewer).

---

### Phase 4: Database Management
* [x] Build Go Agent database manager for MySQL / MariaDB and PostgreSQL.
* [x] Implement database creation, deletion, and collation selection.
* [x] Implement database user creation, host restriction, and password management.
* [x] Implement user grant assignment and revocation.
* [x] Build Database Management UI in Laravel Control Plane.

---

### Phase 5: Git Deployment Engine
* [x] Build unprivileged deployment runner in Go Agent (`su - <site_user>`).
* [x] Implement Git clone, fetch, branch checkout, and user build hook execution.
* [x] Implement real-time log capturing and deployment history storage.
* [x] Implement Webhook receiver for GitHub / GitLab push triggers.
* [x] Build Deployment UI with live terminal log viewer and rollback trigger.

---

### Phase 6: Operations (Cron, Processes, Files, Firewall)
* [x] Implement Cron Job manager with crontab syntax validation.
* [x] Implement Managed Process manager (Systemd application units for queue workers / Node.js).
* [x] Implement safe, path-jailed File Manager (browse, read, write, upload, permissions).
* [x] Implement UFW Firewall manager (list rules, open/close TCP/UDP ports).
* [x] Build Operations UI pages for Cron, Processes, Files, and Firewall.

---

### Phase 7: Hardening, Installer & Verification
* [x] Build modular bash installer (`installer/install.sh` and modules).
* [x] Perform security penetration testing (path traversal, command injection, auth bypass).
* [x] Optimize frontend bundle size and asset caching.
* [x] Test end-to-end installation on a clean Ubuntu 24.04 LTS instance.
* [x] Finalize documentation and release v1.0.0-MVP.

---

## 3. Explicit Non-Goals for MVP
To maintain minimalism, speed, and reliability, the following features are intentionally out of scope for the MVP:
* Multi-server clustering or centralized control planes.
* Kubernetes or Docker container orchestration.
* Built-in Mail Transfer Agent (Postfix/Dovecot) — use transactional email providers.
* DNS authoritative server hosting (BIND9).
* Web-based phpMyAdmin / Adminer embedding.
