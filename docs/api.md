# Kodepreneur Control Plane API Reference

## 1. Overview

The **Kodepreneur Control Plane** exposes a standardized JSON REST API designed for headless management, CLI integration, and programmatic operations.

### Base URL
```text
https://your-server-ip:8443/api/v1
```

### Authentication
All requests must include a valid Bearer token generated from the panel:
```http
Authorization: Bearer 1|kp_admintoken_xxxxxxxxxxxxxxxx
Accept: application/json
Content-Type: application/json
```

---

## 2. Standard Response Envelope

### 2.1 Success Response (`200 OK`, `201 Created`)
```json
{
  "success": true,
  "data": {
    "id": "site_01hxabc123",
    "domain": "example.com",
    "status": "active"
  },
  "message": "Resource created successfully"
}
```

### 2.2 Error Response (`400 Bad Request`, `422 Unprocessable`, `500 Server Error`)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "The given data was invalid.",
    "details": {
      "domain": ["The domain has already been taken."]
    }
  }
}
```

---

## 3. Endpoints Matrix

### 3.1 Authentication
* `POST /auth/login` — Authenticate and receive API token or session.
* `POST /auth/logout` — Revoke active token.
* `GET  /auth/me` — Retrieve authenticated user profile and permissions.

### 3.2 System & Metrics
* `GET  /system/info` — Hostname, OS version, kernel, uptime, public IP.
* `GET  /system/metrics` — Real-time CPU, RAM, Disk, Load average metrics.
* `GET  /system/services` — Status of core daemons (`nginx`, `php8.3-fpm`, `php8.4-fpm`, `mysql`, `postgresql`, `kodepreneur-agent`).
* `POST /system/services/{service}/{action}` — Restart, start, stop, or reload a system service.

### 3.3 Website Management
* `GET    /websites` — List all managed websites with pagination & search.
* `POST   /websites` — Provision a new website, document root, and Nginx vhost.
* `GET    /websites/{id}` — Get single website configuration details.
* `PUT    /websites/{id}` — Update website runtime, PHP version, or aliases.
* `DELETE /websites/{id}` — Remove website and optionally delete files.
* `POST   /websites/{id}/ssl` — Request or renew Let's Encrypt SSL certificate.
* `GET    /websites/{id}/logs/{type}` — Stream or fetch `access.log` or `error.log`.

### 3.4 Database Management
* `GET    /databases` — List all MySQL / PostgreSQL databases.
* `POST   /databases` — Create database instance with specified collation.
* `DELETE /databases/{id}` — Drop database.
* `GET    /databases/users` — List database users.
* `POST   /databases/users` — Create database user and assign grants.
* `DELETE /databases/users/{id}` — Delete database user.

### 3.5 Git Deployments
* `GET    /websites/{id}/deployments` — List deployment history.
* `POST   /websites/{id}/deployments/trigger` — Trigger manual Git pull & build.
* `GET    /deployments/{id}/logs` — Retrieve full terminal build log output.
* `POST   /deployments/{id}/rollback` — Rollback to an earlier deployment state.

### 3.6 Process Management
* `GET    /processes` — List managed application processes (workers, Node.js, Python).
* `POST   /processes` — Create new managed systemd service.
* `POST   /processes/{id}/{action}` — Start, stop, or restart process.
* `DELETE /processes/{id}` — Remove managed process unit.

### 3.7 Cron & Scheduled Tasks
* `GET    /cron-jobs` — List scheduled tasks.
* `POST   /cron-jobs` — Register a new cron schedule with command.
* `DELETE /cron-jobs/{id}` — Remove scheduled task.

### 3.8 Firewall Management
* `GET    /firewall/rules` — List active UFW rules and open ports.
* `POST   /firewall/rules` — Add new allowed port / protocol rule.
* `DELETE /firewall/rules/{id}` — Remove firewall rule.

### 3.9 Audit & Logs
* `GET    /audit-logs` — Query immutable audit activity trail.
