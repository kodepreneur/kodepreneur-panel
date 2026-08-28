# Architecture Overview

## 1. Executive Summary

**Kodepreneur Panel** is engineered as a two-tier, privilege-separated server control panel designed for single-node Linux VPS administration. It eliminates the security risks of traditional monolithic web panels that run as `root` by cleanly separating the **Control Plane** (Web UI, API, Business Logic, RBAC) from the **Execution Plane** (Privileged Go Daemon).

---

## 2. High-Level Architecture Diagram

```text
                        ┌────────────────────────┐
                        │      Web Browser       │
                        └───────────┬────────────┘
                                    │ HTTPS (Port 8443 or Domain)
                                    ▼
       ┌────────────────────────────────────────────────────────┐
       │             Kodepreneur Control Plane                  │
       │  (Unprivileged User: `kodepreneur`, Group: `www-data`)  │
       │                                                        │
       │   ┌────────────────────────────────────────────────┐   │
       │   │        Inertia.js + Vue 3 + Tailwind CSS       │   │
       │   └───────────────────────┬────────────────────────┘   │
       │                           │ HTTP / JSON                │
       │   ┌───────────────────────▼────────────────────────┐   │
       │   │         Laravel 13 Backend (PHP 8.3+)          │   │
       │   │    - Authentication, RBAC, Rate Limiting       │   │
       │   │    - Business Validation & Desired State       │   │
       │   │    - Internal SQLite Database (`panel.sqlite`) │   │
       │   └───────────────────────┬────────────────────────┘   │
       └───────────────────────────┼────────────────────────────┘
                                   │
                                   │ Internal HTTPS / Unix Socket
                                   │ Mutual TLS / Signed Token Auth
                                   ▼
       ┌────────────────────────────────────────────────────────┐
       │             Kodepreneur Execution Agent                │
       │              (Privileged User: `root`)                 │
       │                                                        │
       │   ┌────────────────────────────────────────────────┐   │
       │   │           Go Daemon (`kodepreneur-agent`)      │   │
       │   │    - Typed RPC/REST Handlers                   │   │
       │   │    - Input Sanitization & Path Verification    │   │
       │   │    - Atomic Config Writers & Linters           │   │
       │   │    - Process, Service & Metric Collectors      │   │
       │   └───────────────────────┬────────────────────────┘   │
       └───────────────────────────┼────────────────────────────┘
                                   │
                                   ▼
 ┌─────────────────────────────────────────────────────────────────────────┐
 │                            Host Linux Subsystems                        │
 │                                                                         │
 │  ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌──────────────┐ │
 │  │     Nginx     │ │    PHP-FPM    │ │  PostgreSQL/  │ │   Systemd &  │ │
 │  │ Configuration │ │  Pools (8.3/  │ │     MySQL     │ │   Cgroups    │ │
 │  │   & Vhosts    │ │     8.4)      │ │   Databases   │ │   Services   │ │
 │  └───────────────┘ └───────────────┘ └───────────────┘ └──────────────┘ │
 │  ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌──────────────┐ │
 │  │   Certbot /   │ │  UFW Firewall │ │ Linux Users & │ │  Crontab &   │ │
 │  │ Let's Encrypt │ │     Rules     │ │  Permissions  │ │  Schedulers  │ │
 │  └───────────────┘ └───────────────┘ └───────────────┘ └──────────────┘ │
 └─────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Tier Responsibilities & Boundaries

### 3.1 Tier 1: Laravel Control Plane
* **Runtime**: PHP 8.3+ FPM, served by Nginx or internal high-performance worker.
* **Privilege Level**: Unprivileged system user (`kodepreneur`).
* **Database**: Local SQLite (`/var/lib/kodepreneur/panel.sqlite`).
* **Responsibilities**:
  * User Authentication (Session, API Tokens, MFA/2FA).
  * Role-Based Access Control (Admin, Operator, Read-Only).
  * Web Interface (Single-Page Application powered by Inertia.js + Vue 3).
  * Storing Desired State (Sites, Databases, Domains, SSL Metadata).
  * Dispatching Asynchronous Jobs (Deployments, Backups, Health Checks).
  * Comprehensive Audit Logging (Who did what, from where, and result).
* **Strict Constraints**:
  * **Zero Root Access**: Laravel cannot execute `sudo` or raw shell strings.
  * **No Direct Config File Writing**: Laravel cannot directly modify `/etc/nginx/sites-available` or `/etc/systemd/system`.
  * All privileged actions must be delegated to the Go Agent via strongly typed DTOs.

### 3.2 Tier 2: Kodepreneur Go Agent
* **Runtime**: Single compiled static Go binary (`/usr/local/bin/kodepreneur-agent`).
* **Privilege Level**: `root` (managed via `systemd`).
* **Communication Channel**:
  * Local Unix Domain Socket (`/run/kodepreneur/agent.sock`) with strict file permissions (`0660`, owner `root:kodepreneur`), OR
  * Localhost HTTPS (`127.0.0.1:8765`) protected with mTLS / Bearer HMAC Token.
* **Responsibilities**:
  * Inspecting & validating system health and resource metrics (CPU, Memory, Disk, Network, Load).
  * Generating and testing Nginx virtual host configurations (`nginx -t` validation before reload).
  * Managing PHP-FPM pools (custom user isolation, socket creation, memory limits).
  * Managing database users, databases, and permissions for PostgreSQL & MySQL.
  * Automating Let's Encrypt certificate issuance via ACME/Certbot with safe renewals.
  * Managing UFW firewall rules (whitelisting/blocking ports safely).
  * Managing managed user processes (Node.js, Python, Laravel queue workers) via systemd units.
  * Controlled, sandboxed file operations confined to `/var/www` web roots.

---

## 4. State Management: Declarative vs. Imperative

Kodepreneur Panel operates on a **Declarative Desired State** model with **Idempotent Reconciliation**:

1. **User Action**: An administrator adds a new domain `example.com` with PHP 8.4.
2. **Control Plane**: Laravel validates the request, records the record in `panel.sqlite` with status `provisioning`, and dispatches an API request to the Agent:
   ```json
   POST /api/v1/websites
   {
     "id": "site_01hx...",
     "domain": "example.com",
     "php_version": "8.4",
     "document_root": "/var/www/example.com/public",
     "system_user": "kodepreneur_site_1"
   }
   ```
3. **Agent Execution**:
   * Creates system user and isolated web root directory `/var/www/example.com`.
   * Generates PHP-FPM pool config `/etc/php/8.4/fpm/pool.d/example.com.conf`.
   * Generates Nginx configuration in `/etc/nginx/sites-available/example.com.conf`.
   * Performs atomic syntax check: `nginx -t`.
   * Enables symlink in `/etc/nginx/sites-enabled/`.
   * Reloads services gracefully (`systemctl reload php8.4-fpm`, `systemctl reload nginx`).
4. **State Confirmation**: Agent returns `200 OK` with execution report. Laravel updates database state to `active`.

---

## 5. Security & Isolation Matrix

| Subsystem | Managed Mechanism | Isolation Model |
| :--- | :--- | :--- |
| **Websites** | Nginx + PHP-FPM Pools | Dedicated unprivileged Linux user per site or shared secure pool with `open_basedir`. |
| **Databases** | MySQL / MariaDB / PostgreSQL | Distinct database user with grants restricted to the target database. |
| **Deployments** | Git + Shell scripts | Run strictly under the site's unprivileged user (never as root). |
| **Processes** | Systemd Units (`kp-app-<id>.service`) | Sandboxed with `ProtectSystem=full`, `PrivateTmp=yes`, and CPU/Memory limits. |
| **Files** | File Manager Backend | Strict path normalization and jail verification restricted to `/var/www/<site>`. |

---

## 6. Directory Layout on the Target VPS

```text
/
├── etc/
│   ├── kodepreneur/
│   │   ├── agent.yaml           # Agent configuration & secrets
│   │   └── panel.env            # Laravel environment configuration
│   ├── nginx/
│   │   ├── kodepreneur/         # Global snippet configurations (SSL, security, PHP)
│   │   ├── sites-available/     # Managed site vhosts
│   │   └── sites-enabled/       # Symlinks to active vhosts
│   └── systemd/system/
│       ├── kodepreneur-agent.service
│       ├── kodepreneur-panel.service
│       └── kodepreneur-worker.service
├── var/
│   ├── lib/kodepreneur/
│   │   ├── panel.sqlite         # Panel SQLite database
│   │   └── back-ups/            # Internal state snapshots
│   ├── log/kodepreneur/
│   │   ├── agent.log            # Agent execution logs
│   │   ├── panel.log            # Laravel application logs
│   │   └── deployments/         # Git deployment build logs
│   └── www/
│       ├── kodepreneur-panel/   # Panel UI & API codebase
│       └── <domain.com>/        # Managed user websites
└── usr/local/bin/
    └── kodepreneur-agent        # Compiled Go Agent binary
```
