# Security Architecture & Threat Model

## 1. Security Philosophy

Because **Kodepreneur Panel** manages critical system infrastructure on Linux servers, security is treated as an uncompromising first-class constraint. 

The core security tenet is:
> **Never trust web layer inputs, never grant root privileges to the web application, and never expose generic shell execution endpoints.**

---

## 2. Threat Model & Mitigation Matrix

| Threat Category | Potential Vector | Mitigation Strategy |
| :--- | :--- | :--- |
| **Privilege Escalation** | Remote Code Execution (RCE) in Laravel web app attempting to gain `root` access. | Laravel runs as unprivileged user `kodepreneur`. No `sudoers` privileges are granted to `kodepreneur`. Privileged actions require typed requests to the Go agent. |
| **Command Injection** | Unsanitized domain, database name, or username passed to bash/system shell. | **Zero shell interpolation**. The Go Agent uses Go's `os/exec` with explicit argument slices (`exec.Command("nginx", "-t")`), completely avoiding shell string evaluation. |
| **Path Traversal** | Malicious paths (e.g. `../../etc/shadow`) in File Manager or Git deploy paths. | Agent strictly canonicalizes all paths via `filepath.Abs` and `filepath.EvalSymlinks`, asserting that target paths are strictly within authorized directory boundaries (e.g. `/var/www/<domain>`). |
| **Unauthorized Agent Access** | Attacker on the VPS attempts to call the Go Agent directly. | Agent listens on a Unix Domain Socket (`/run/kodepreneur/agent.sock`) with `0660` permissions restricted to `root:kodepreneur`, or localhost with HMAC-SHA256 request signing and timestamp anti-replay checks. |
| **Brute Force & Credential Stuffing** | Automated login attacks against the web control panel. | Laravel rate limiters, progressive IP delays, account lockout after 5 consecutive failures, and Fail2ban jail integration. |
| **Cross-Site Request Forgery (CSRF)** | Malicious third-party site triggering actions on an active session. | Inertia.js + Laravel CSRF tokens validated on all state-changing `POST`, `PUT`, `PATCH`, and `DELETE` requests with `SameSite=Lax/Strict` cookies. |
| **Information Leakage in Logs** | Database passwords or secret keys dumped into audit or system logs. | Strict log sanitization filter in Laravel and Go Agent that automatically masks keys named `password`, `secret`, `private_key`, `token`, and `api_key`. |

---

## 3. Privilege Separation & Boundary Enforcement

```text
[ Web Request ]
      │
      ▼
┌───────────────────────────────────────────────────────────┐
│ Laravel Web Application (User: `kodepreneur`)             │
│                                                           │
│  1. Authenticates session / JWT / API Token               │
│  2. Checks RBAC permissions (`can('manage-websites')`)    │
│  3. Validates input schema strictly (FormRequest)         │
│  4. Formulates typed DTO: CreateWebsiteRequest            │
│  5. Signs payload with HMAC-SHA256 Secret                 │
│  6. Dispatches to Agent over local socket/TLS             │
└─────────────────────────────┬─────────────────────────────┘
                              │ Signed Typed Request
                              ▼
┌───────────────────────────────────────────────────────────┐
│ Go Privileged Agent (User: `root`)                        │
│                                                           │
│  1. Verifies HMAC signature, timestamp (< 30s drift), auth│
│  2. Validates DTO types, regex patterns, domain syntax    │
│  3. Verifies path containment (`filepath.Clean`)          │
│  4. Executes explicit internal handlers (NO generic exec) │
│  5. Writes configuration files atomically                 │
│  6. Verifies syntax (`nginx -t`) before service reload    │
│  7. Returns typed response: 200 OK / 422 Unprocessable    │
└───────────────────────────────────────────────────────────┘
```

---

## 4. Agent Communication Protocol & Authentication

The communication channel between the Laravel Control Plane and the Go Agent is secured using **Shared Secret HMAC-SHA256 Authentication** over a restricted Unix Domain Socket or Localhost HTTPS.

### 4.1 Request Headers
Every request sent from Laravel to the Go Agent must contain:
```http
X-Kodepreneur-Timestamp: 1772154000
X-Kodepreneur-Nonce: 8f3c2a1e-4b9d-4e7a-9c1a-2e3f4a5b6c7d
X-Kodepreneur-Signature: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Content-Type: application/json
```

### 4.2 Signature Calculation
```text
Signature = HMAC_SHA256(
    Key: AGENT_SHARED_SECRET,
    Data: Method + "\n" + RequestURI + "\n" + Timestamp + "\n" + Nonce + "\n" + SHA256(RequestBody)
)
```

### 4.3 Agent Verification Rules
1. **Timestamp Freshness**: Rejects any request where `|currentTime - Timestamp| > 30 seconds` to prevent replay attacks.
2. **Nonce De-duplication**: Stores recent nonces in a memory cache for 60 seconds; rejects duplicate nonces.
3. **Constant-Time Comparison**: Compares the calculated signature using `crypto/subtle.ConstantTimeCompare` to mitigate timing attacks.

---

## 5. Path Sanitization & Jail Rules

For all operations touching the filesystem (File Manager, Git Deployments, SSL certs, Nginx logs):

```go
func ValidateSafePath(baseDir, requestedPath string) (string, error) {
    cleanBase := filepath.Clean(baseDir)
    targetPath := filepath.Clean(filepath.Join(cleanBase, requestedPath))
    
    // Resolve any symlinks to find the real physical path
    realTarget, err := filepath.EvalSymlinks(targetPath)
    if err != nil && !os.IsNotExist(err) {
        return "", errors.New("failed to resolve path")
    }
    
    // Ensure real physical path begins with clean base directory
    if !strings.HasPrefix(realTarget, cleanBase+string(filepath.Separator)) && realTarget != cleanBase {
        return "", errors.New("security violation: path traversal detected")
    }
    
    return targetPath, nil
}
```

---

## 6. Role-Based Access Control (RBAC)

Kodepreneur Panel implements a granular RBAC matrix stored in `panel.sqlite`:

| Permission Key | Super Admin | Server Operator | Developer / Deployer |
| :--- | :---: | :---: | :---: |
| `system.metrics.view` | ✅ | ✅ | ✅ |
| `websites.view` | ✅ | ✅ | ✅ |
| `websites.create_edit` | ✅ | ✅ | ❌ |
| `websites.delete` | ✅ | ❌ | ❌ |
| `deployments.trigger` | ✅ | ✅ | ✅ |
| `databases.manage` | ✅ | ✅ | ❌ |
| `services.restart` | ✅ | ✅ | ❌ |
| `firewall.manage` | ✅ | ❌ | ❌ |
| `settings.manage` | ✅ | ❌ | ❌ |
| `audit.view` | ✅ | ✅ | ❌ |

---

## 7. Audit Logging Specification

Every state-changing administrative action is permanently recorded in the immutable `activity_logs` table:

```json
{
  "id": "act_01hx...",
  "user_id": "usr_01hx...",
  "user_email": "admin@kodepreneur.com",
  "ip_address": "203.0.113.42",
  "user_agent": "Mozilla/5.0 ...",
  "action": "website.create",
  "resource_type": "website",
  "resource_id": "site_01hx...",
  "status": "success",
  "payload_summary": {
    "domain": "api.kodepreneur.com",
    "php_version": "8.4"
  },
  "created_at": "2026-08-28T05:10:00Z"
}
```
Sensitive fields (passwords, tokens, database secrets) are redacted before serialization.
