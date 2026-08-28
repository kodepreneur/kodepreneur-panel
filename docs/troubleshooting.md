# Maintenance, Updating & Debugging Guide 🛠️

This guide outlines standard operating procedures for troubleshooting bugs, inspecting logs, and updating a live, production-installed **Kodepreneur Panel** instance.

---

## 1. Automated Updating

### Option A: In-Panel Web UI (1-Click)
Navigate to **Settings / Updates** (`/updates`) in the Kodepreneur Panel web interface:
1. Click **Check for Updates** to fetch the latest GitHub releases.
2. Click **Update Panel Now** / **Upgrade to Latest**.
3. The UI streams the real-time build and migration logs. The Go Agent daemon executes the update in `--daemon-mode` with an extended 600s execution window, returns the full log to the browser, and then performs an asynchronous, graceful service restart without interrupting HTTP responses.

### Option B: One-Line Remote Updater (SSH)
Run directly from any root SSH terminal on your server:

```bash
curl -fsSL https://raw.githubusercontent.com/kodepreneur/kodepreneur-panel/main/installer/update.sh | sudo bash
```

### Option C: Local Repository Updater
If you have cloned the repository locally on your server:

```bash
cd /path/to/kodepreneur-panel
git pull origin main
sudo ./installer/update.sh
```

This automatically:
1. Updates Laravel production packages (`composer install --no-dev`)
2. Compiles frontend assets (`npm run build`)
3. Runs database migrations (`php artisan migrate --force`)
4. Refreshes and optimizes route, config, and view caches
5. Recompiles and restarts the root Go Agent daemon (`kodepreneur-agent`)
6. Reloads Nginx with zero downtime

---

## 2. Live Debugging & Log Locations

### 🔍 1. Go Agent Daemon Debugging
The Go agent daemon handles all root-privileged operations (Nginx vhosts, PHP-FPM pools, MySQL/PostgreSQL, Certbot, Git deployments).

* **Check Service Status**:
  ```bash
  sudo systemctl status kodepreneur-agent
  ```
* **Real-time Live Logs (Journal)**:
  ```bash
  sudo journalctl -u kodepreneur-agent -f
  ```
* **Test Local Health Endpoint**:
  ```bash
  curl -s http://127.0.0.1:8443/health | jq .
  ```
* **Configuration File**:
  `/etc/kodepreneur/agent.yaml`

---

### 🔍 2. Laravel Control Plane Debugging
The web interface runs as unprivileged `www-data` under `/var/www/kodepreneur-panel`.

* **Laravel Application Error Logs**:
  ```bash
  tail -f /var/www/kodepreneur-panel/storage/logs/laravel.log
  ```
* **Verify Artisan & Database Connection**:
  ```bash
  cd /var/www/kodepreneur-panel
  sudo -u www-data php artisan about
  ```
* **Clear Stale Caches (Fixes most UI/Routing glitches)**:
  ```bash
  cd /var/www/kodepreneur-panel
  sudo -u www-data php artisan optimize:clear
  sudo -u www-data php artisan config:cache
  sudo -u www-data php artisan route:cache
  sudo -u www-data php artisan view:cache
  ```
* **Inspect Database State**:
  The panel database is located at `/var/www/kodepreneur-panel/database/database.sqlite`.

---

### 🔍 3. Nginx Web Server Debugging
* **Check Nginx Configuration Syntax**:
  ```bash
  sudo nginx -t
  ```
* **Nginx Service Status**:
  ```bash
  sudo systemctl status nginx
  ```
* **Global Error Log**:
  ```bash
  sudo tail -f /var/log/nginx/error.log
  ```
* **Per-Website Error Log**:
  ```bash
  sudo tail -f /var/log/nginx/{domain}_error.log
  ```

---

### 🔍 4. PHP-FPM Runtimes Debugging
* **Check PHP-FPM Services**:
  ```bash
  sudo systemctl status php8.3-fpm
  sudo systemctl status php8.4-fpm
  ```
* **PHP-FPM Main Error Logs**:
  ```bash
  sudo tail -f /var/log/php8.3-fpm.log
  sudo tail -f /var/log/php8.4-fpm.log
  ```
* **Per-User Pool Configurations**:
  `/etc/php/8.3/fpm/pool.d/` and `/etc/php/8.4/fpm/pool.d/`

---

### 🔍 5. Database Engines Debugging
* **MySQL / MariaDB Status**:
  ```bash
  sudo systemctl status mariadb
  sudo mariadb-admin status
  ```
* **PostgreSQL Status**:
  ```bash
  sudo systemctl status postgresql
  sudo -u postgres psql -c "\l"
  ```

---

## 3. Common Troubleshooting Scenarios

### Scenario A: "502 Bad Gateway" on a Website
1. Check if PHP-FPM is running: `sudo systemctl status php8.3-fpm`
2. Check if the domain's dedicated socket exists: `ls -l /run/php/`
3. Verify permissions on the webroot: `ls -la /var/www/{domain}` (must be owned by the site's system user `kp_*`).

### Scenario B: "Agent Client Error: Signature Mismatch / Unauthorized"
1. Verify the secret key in `/etc/kodepreneur/agent.yaml` matches `AGENT_SECRET_KEY` in `/var/www/kodepreneur-panel/.env`.
2. Check server time synchronization (`timedatectl` or `sudo ntpdate pool.ntp.org`) — HMAC rejects requests with >300s clock drift.

### Scenario C: "Nginx -t Syntax Error during Deployment"
The Agent automatically rolls back syntax errors before reloading Nginx. Check `journalctl -u kodepreneur-agent -f` to see the exact invalid line reported by `nginx -t`.
