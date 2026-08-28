# Changelog

All notable changes, architectural decisions, and releases for **Kodepreneur Panel** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0-mvp] - Phase 7: Hardening, Modular Installer & v1.0.0-MVP Release (2026-08-28)

### Added
* **Modular Production Installer (`installer/`)**:
  * Master CLI entrypoint (`installer/install.sh`) supporting interactive and non-interactive (`--unattended`, `--email`, `--password`, `--port`) execution.
  * `00_system.sh`: Base packages, UFW firewall defaults, and Fail2ban security.
  * `10_php.sh`: Ondřej Surý PHP PPA, PHP 8.3 LTS & PHP 8.4 runtimes, Composer, and Node.js LTS.
  * `20_nginx.sh`: Nginx web server installation, security headers, and Certbot.
  * `30_databases.sh`: MariaDB (MySQL) and PostgreSQL server setup.
  * `40_agent.sh`: Go Agent daemon compilation, secure HMAC key generation (`/etc/kodepreneur/agent.yaml`), and `kodepreneur-agent.service` systemd registration.
  * `50_panel.sh`: Laravel Control Plane deployment, SQLite storage initialization, initial admin user provisioning, and Nginx panel vhost.
  * `99_summary.sh`: Final installation banner with credentials and management tips.
* **Systemd & Nginx Service Templates**:
  * `installer/systemd/kodepreneur-agent.service`: Root-privileged Go Agent service unit.
  * `installer/nginx/kodepreneur-panel.conf`: Nginx virtual host block for control plane.
* **Complete Documentation (`README.md`)**:
  * Architecture diagrams, 1-line installation command, security model, and developer quickstart.

---

## [1.0.0-alpha.5] - Phase 5 & 6: Git Deployments & Operations Subsystem (2026-08-28)

### Added
* **Git Deployment Runner (`agent/pkg/git`)**:
  * Unprivileged `su - <user>` runner executing git pull, composer install, artisan migrations, and npm build hooks.
  * Real-time deployment log aggregation, duration calculation, and rollback commit checkout.
  * REST endpoint `POST /api/v1/deployments/execute`.
* **Operations Subsystem (`agent/pkg/cron`, `pkg/process`, `pkg/firewall`, `pkg/filemanager`)**:
  * Linux crontab synchronization and 5-field syntax validation (`POST /api/v1/cron/sync`).
  * Systemd application unit generator, daemon-reload, and supervisor for background queue workers (`POST /api/v1/processes`, `POST /api/v1/processes/{id}/{action}`, `DELETE /api/v1/processes/{id}`).
  * UFW firewall rule additions and port deletions (`POST /api/v1/firewall/rules`, `DELETE /api/v1/firewall/rules`).
  * Path-jailed sandboxed File Manager for browsing, reading, saving, and deleting webroot files (`POST /api/v1/files/browse`, `POST /api/v1/files/read`, `POST /api/v1/files/write`, `POST /api/v1/files/delete`).
* **Laravel Control Plane Subsystem**:
  * Added operations methods to `AgentClientInterface`, `HttpAgentClient`, and `MockAgentClient`.
  * `DeploymentController`: Deployment trigger, live log viewing, duration calculation, rollback execution.
  * `CronJobController`: Scheduled tasks CRUD, active toggle, and crontab sync.
  * `ProcessController`: Managed worker daemons CRUD and start/stop/restart actions.
  * `FirewallController`: UFW port rules CRUD.
  * `FileManagerController`: Webroot explorer, code editor modal, and file deletion.
  * Vue 3 UI pages: `Deployments/Show.vue`, `Cron/Index.vue`, `Processes/Index.vue`, `Firewall/Index.vue`, `Files/Index.vue`.
  * Navigation links updated in `AppLayout.vue`.
  * Automated feature tests in `OperationsTest.php` (31 tests, 138 assertions passing).

---

## [1.0.0-alpha.4] - Phase 4: Database Management Engine (2026-08-28)

### Added
* **Go Agent Database Engine (`agent/pkg/database`)**:
  * Unified database manager supporting both **MySQL / MariaDB** and **PostgreSQL**.
  * Database creation with custom charset (`utf8mb4` / `UTF8`) and collation (`utf8mb4_unicode_ci` / `en_US.UTF-8`).
  * Database user creation with host restrictions (`localhost` / `%`) and secure password generation.
  * Privilege grant assignments (`ALL PRIVILEGES` / `READ ONLY`) and revocations.
  * User password reset and secure database dropping.
  * Cross-platform dev mode simulation.
  * REST endpoints (`POST /databases`, `DELETE /databases/{engine}/{name}`, `POST /databases/users`, `DELETE /databases/users/{engine}/{username}`, `POST /databases/grants`, `POST /databases/users/password`).
* **Laravel Control Plane Subsystem**:
  * Added database methods to `AgentClientInterface`, `HttpAgentClient`, and `MockAgentClient`.
  * Enhanced `DatabaseController` with database CRUD, user management, grant assignments, and password resets.
  * Registered subroutes in `routes/web.php`.
  * Enhanced Vue 3 UI (`Databases/Index.vue`) with multi-tab layout (**Databases**, **Database Users**, **Assign User Grants**), random password generator, and modal workflows.
  * Added automated feature tests in `DatabaseTest.php` (25 tests, 104 assertions passing).

---

## [1.0.0-alpha.3] - Phase 3: Website & PHP Runtime Management (2026-08-28)

### Added
* **Nginx Virtual Host Engine (`agent/pkg/nginx`)**:
  * Dynamic Nginx vhost template generator with HTTP, HTTPS/SSL, FastCGI PHP pass, ACME challenge webroot, and security headers.
  * Atomic configuration deployment with `nginx -t` validation harness and automatic rollback on syntax errors.
  * Real-time Nginx `access.log` and `error.log` tail reader with line filters.
* **PHP-FPM Pool Isolation (`agent/pkg/phpfpm`)**:
  * Dedicated PHP-FPM pool generation for PHP 8.3 & PHP 8.4 (`/etc/php/{version}/fpm/pool.d/{user}.conf`).
  * Isolated Linux system user provisioning with hardened permissions (`0750`).
  * Live runtime version switcher with zero-downtime pool recreation and Nginx socket rebinding.
* **Certbot & SSL Automation (`agent/pkg/ssl`)**:
  * Automated Let's Encrypt ACME HTTP-01 challenge issuance via Certbot.
  * Development mode fallback with self-signed certificate generator.
  * Certificate expiration date inspector and automatic HTTPS redirect enforcement.
* **Laravel Control Plane Subsystem**:
  * Updated `AgentClientInterface`, `HttpAgentClient`, and `MockAgentClient` with website lifecycle, PHP switching, SSL issuance, and log streaming methods.
  * Enhanced `WebsiteController` with SSL issuance (`POST /websites/{website}/ssl`), PHP switching (`PUT /websites/{website}/php`), and log retrieval (`GET /websites/{website}/logs/{type}`).
  * Multi-tab management UI in Vue 3 (`Overview`, `PHP Runtime`, `SSL / HTTPS`, `Nginx Logs`, `Deployments`).
  * Added automated feature tests in `WebsiteTest.php` (18 tests, 71 assertions passing).

---

## [1.0.0-alpha.2] - Upgrade to Laravel 13 (2026-08-28)

### Changed
* Upgraded **Laravel Framework** core to `v13.29.0`.
* Upgraded **Inertia.js Laravel Adapter** to `v3.3.1`.
* Upgraded **Laravel Tinker** to `v3.0.2`.
* Upgraded **PHPUnit** to `v12.5.34` and Symfony components to `v8.1.5`.
* Verified 100% test suite compatibility (14 tests, 55 assertions passing with zero deprecations).
* Updated architecture documentation, ADRs, and roadmap to reflect Laravel 13.

## [1.0.0-alpha.1] - Phase 1: Project Foundation & Agent Skeleton (2026-08-28)

### Added
* **Go Agent Skeleton (`agent/`)**:
  * Compiled Go daemon with command-line flags (`--config`, `--dev`, `--port`, `--socket`).
  * HMAC-SHA256 request authentication middleware with timestamp freshness window and nonce de-duplication cache.
  * System metrics collector inspecting `/proc` and `/sys` (CPU, RAM, Disk, Load average) with cross-platform dev fallbacks.
  * System info and systemd service management abstractions.
  * Automated Go test suite (`auth_test.go`, `server_test.go`).
* **Laravel Control Plane (`panel/`)**:
  * Clean Laravel 11 application with Inertia.js v2, Vue 3, TypeScript, and Tailwind CSS.
  * Full relational SQLite schema migrations with optimized PRAGMAs (WAL mode, foreign keys).
  * RBAC foundation (`roles`, `users` relationships, permissions checking middleware).
  * Seeders for `Super Administrator`, `Server Operator`, `Developer`, and system settings.
  * Typed Agent Client interface with `HttpAgentClient` (HMAC signed requests) and `MockAgentClient` for local development.
  * Vue 3 UI pages:
    * `Auth/Login.vue` (Sleek dark theme authentication)
    * `Dashboard/Index.vue` (Live gauges, system specs, service status, recent audit activity)
    * `Websites/Index.vue`, `Websites/Create.vue`, `Websites/Show.vue` (Virtual host management)
    * `Databases/Index.vue` (MySQL / PostgreSQL database management)
    * `Deployments/Index.vue`, `Deployments/Show.vue` (Git build pipeline and terminal log viewer)
    * `Services/Index.vue` (Interactive systemd service manager)
    * `Activity/Index.vue` (Audit log ledger)
    * `Settings/Index.vue` (Server panel configuration)
  * Comprehensive Laravel feature tests covering Auth, Dashboard, Websites, and Agent client.
