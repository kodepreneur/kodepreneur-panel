# Installation Architecture & Bootstrap Guide

## 1. Overview

Kodepreneur Panel is installed on a fresh Ubuntu 24.04 LTS VPS using a modular, idempotent, and fault-tolerant installer. 

```bash
curl -fsSL https://install.kodepreneur.com | sudo bash
```

The installer avoids monolithic, fragile scripts by utilizing an organized modular structure with clear logging, pre-flight validation, failure handling, and idempotent execution.

---

## 2. Pre-Flight System Requirements

| Parameter | Minimum Requirement | Recommended |
| :--- | :--- | :--- |
| **Operating System** | Ubuntu 24.04 LTS (x86_64 / arm64) | Ubuntu 24.04 LTS (x86_64) |
| **Privilege** | `root` or `sudo` access | `root` |
| **vCPU** | 1 Core | 2 Cores+ |
| **RAM** | 1 GB (with swap enabled) | 2 GB+ |
| **Storage** | 10 GB Free Disk Space | 25 GB+ SSD/NVMe |
| **Networking** | Dedicated Public IPv4 / IPv6, Outbound Internet Access | Static Public IP |
| **Required Inbound Ports** | 80 (HTTP), 443 (HTTPS), 8443 (Panel Web UI) | 22, 80, 443, 8443 |

---

## 3. Modular Installer Design

The installer codebase is organized into structured libraries and modular lifecycle steps:

```text
installer/
├── install.sh                  # Main entrypoint & orchestration loop
├── config.env.example          # Default environment configurations
├── lib/
│   ├── common.sh               # Utility helpers, string parsing, spinners
│   ├── logging.sh              # Colored console output & /var/log/kodepreneur-install.log
│   ├── os.sh                   # Ubuntu 24.04 detection & architecture checks
│   ├── packages.sh             # APT wrapper with retry & lock management
│   └── validation.sh           # RAM, disk, network, and port pre-flight checks
└── modules/
    ├── 01_system.sh            # Hostname, timezone, swapfile, core build tools
    ├── 02_packages.sh          # Common repositories, Ondřej Surý PHP PPA, curl, git
    ├── 03_php.sh               # PHP 8.3 & PHP 8.4 runtime, FPM, and required CLI extensions
    ├── 04_nginx.sh             # Nginx mainline/stable installation and base snippets
    ├── 05_composer.sh          # Composer binary installation & global path setup
    ├── 06_agent.sh             # Go Agent binary installation & systemd daemon registration
    ├── 07_panel.sh             # Laravel codebase deployment, SQLite init, keygen, migrations
    ├── 08_firewall.sh          # UFW firewall baseline configuration (22, 80, 443, 8443)
    └── 09_finalize.sh          # Admin user creation, token generation, post-install summary
```

---

## 4. End-to-End Installation Lifecycle (24 Steps)

1. **Pre-execution Check**: Verify execution environment has `bash` and effective UID `0`.
2. **OS Verification**: Validate `lsb_release` strictly returns `Ubuntu` and `VERSION_ID="24.04"`.
3. **Architecture Validation**: Check `uname -m` matches `x86_64` or `aarch64`.
4. **Hardware Sizing Check**: Check available RAM (≥ 1GB) and available disk space (≥ 10GB on `/`).
5. **Network Connectivity**: Perform ping/curl handshake against Ubuntu mirrors and GitHub.
6. **Port Conflict Check**: Validate that ports `80`, `443`, and `8443` are not occupied by conflicting third-party daemons (e.g. Apache2, Lighttpd).
7. **Package Metadata Refresh**: Run `apt-get update -y` with exponential backoff if lock contention occurs.
8. **Base Tooling**: Install `curl`, `wget`, `git`, `unzip`, `tar`, `software-properties-common`, `ca-certificates`, `gnupg`, `ufw`, `fail2ban`.
9. **Automated Swap Provisioning**: If RAM < 2GB and no swap exists, safely provision a 2GB swapfile at `/swapfile`.
10. **PHP Repository Setup**: Add trusted `ppa:ondrej/php` for official multi-PHP support.
11. **PHP Runtime Installation**: Install `php8.3-cli`, `php8.3-fpm`, `php8.3-sqlite3`, `php8.3-mbstring`, `php8.3-xml`, `php8.3-curl`, `php8.3-zip`, `php8.3-bcmath`, `php8.3-intl`.
12. **Nginx Web Server Setup**: Install `nginx`, disable default site, and create `/etc/nginx/kodepreneur/` template directory.
13. **Composer Installation**: Download and verify cryptographic signature of Composer 2.x, install to `/usr/local/bin/composer`.
14. **System User Provisioning**: Create dedicated unprivileged user `kodepreneur` with home `/var/www/kodepreneur-panel`.
15. **Agent Daemon Installation**:
    * Copy compiled binary to `/usr/local/bin/kodepreneur-agent`.
    * Generate secure agent authentication secret stored in `/etc/kodepreneur/agent.yaml`.
    * Register and start `kodepreneur-agent.service` via `systemctl enable --now`.
16. **Panel Codebase Deployment**: Deploy Laravel application to `/var/www/kodepreneur-panel`.
17. **Dependency Resolution**: Execute `composer install --no-dev --optimize-autoloader` as user `kodepreneur`.
18. **Environment Configuration**: Generate `.env` from template, set `APP_KEY`, configure SQLite database at `/var/lib/kodepreneur/panel.sqlite`.
19. **Database Initialization**: Run `php artisan migrate --force` and seed core system settings.
20. **Admin Provisioning**: Create the primary administrative account with a cryptographically secure random password or interactive prompt.
21. **Systemd Worker Setup**: Register and enable `kodepreneur-worker.service` for Laravel queue jobs.
22. **Nginx VHost Configuration**: Create Nginx reverse proxy configuration for the panel on port `8443` (or self-signed HTTPS template).
23. **Firewall Enforcement**: Configure UFW to allow `22/tcp`, `80/tcp`, `443/tcp`, and `8443/tcp`, then enable UFW safely.
24. **Verification & Banner Display**: Probe `http://127.0.0.1:8443/health` to confirm panel readiness and print access credentials to the terminal.

---

## 5. Systemd Service Definitions

### 5.1 Kodepreneur Agent Service (`/etc/systemd/system/kodepreneur-agent.service`)
```ini
[Unit]
Description=Kodepreneur Privileged Server Agent
After=network.target
Wants=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=/usr/local/bin/kodepreneur-agent --config /etc/kodepreneur/agent.yaml
Restart=always
RestartSec=5s
LimitNOFILE=65535
StandardOutput=append:/var/log/kodepreneur/agent.log
StandardError=append:/var/log/kodepreneur/agent.log

# Security Sandboxing
ProtectHome=read-only
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

### 5.2 Kodepreneur Queue Worker Service (`/etc/systemd/system/kodepreneur-worker.service`)
```ini
[Unit]
Description=Kodepreneur Panel Queue Worker
After=network.target kodepreneur-agent.service

[Service]
Type=simple
User=kodepreneur
Group=www-data
WorkingDirectory=/var/www/kodepreneur-panel
ExecStart=/usr/bin/php8.3 artisan queue:work --sleep=3 --tries=3 --max-time=3600
Restart=always
RestartSec=3s
StandardOutput=append:/var/log/kodepreneur/worker.log
StandardError=append:/var/log/kodepreneur/worker.log

[Install]
WantedBy=multi-user.target
```

---

## 6. Idempotency & Re-running Rules

The installer is engineered to be safely re-run on an existing server without data loss:
* **Detection**: Checks if `/var/lib/kodepreneur/panel.sqlite` exists before database operations.
* **Migration**: Executes migrations safely using standard schema comparisons; never drops tables during install/update.
* **Secret Preservation**: Reuses existing `APP_KEY` and Agent token secrets from `/etc/kodepreneur/` instead of overwriting them.
* **Service Health**: Validates and restarts existing services gracefully (`systemctl restart`).
