# Git Deployment Engine Architecture

## 1. Overview

**Kodepreneur Panel** includes a built-in Git deployment engine that allows automated and manual application deployment directly from Git repositories (GitHub, GitLab, Bitbucket, or self-hosted Git).

---

## 2. Core Security & Isolation Principles

1. **Strict Unprivileged Execution**: Deployment commands (e.g. `composer install`, `npm run build`, `php artisan migrate`) **never run as `root`**. They are executed strictly under the dedicated site user (e.g. `kp_example`).
2. **Environment Scrubbing**: Sensitive system environment variables (`AGENT_SECRET`, `PANEL_DB_PASSWORD`) are scrubbed before invoking deploy hooks.
3. **Directory Confinement**: Deployment scripts are restricted to the site's root directory (`/var/www/<domain>`).
4. **Execution Timeouts**: Every deployment has an enforced timeout (default: 600 seconds) to prevent hung processes.

---

## 3. Deployment Workflow

```text
[ Developer Push / Manual Trigger ]
                │
                ▼
┌──────────────────────────────────────────────┐
│ Laravel Control Plane                        │
│ 1. Validate Webhook Signature / User RBAC    │
│ 2. Create `deployments` record (status: pending)
│ 3. Dispatch Typed Job to Go Agent            │
└───────────────────────┬──────────────────────┘
                        │
                        ▼
┌──────────────────────────────────────────────┐
│ Go Agent (`kodepreneur-agent`)               │
│ 1. Switch UID/GID to Site User (`kp_example`)│
│ 2. Execute Git Clone / Fetch / Pull          │
│ 3. Execute User Build Script Lines           │
│ 4. Capture Stderr & Stdout to Memory Buffer │
│ 5. Reload PHP-FPM / Worker if requested      │
│ 6. Return Structured Output & Exit Code      │
└───────────────────────┬──────────────────────┘
                        │
                        ▼
┌──────────────────────────────────────────────┐
│ Laravel Control Plane                        │
│ 1. Update `deployments` (status: success/fail)
│ 2. Store full execution logs in database     │
│ 3. Notify UI via Inertia / Polling           │
└──────────────────────────────────────────────┘
```

---

## 4. Webhook Payload Validation

For automated deployments, each website receives a unique Webhook URL and Secret:
```text
POST https://panel.example.com:8443/api/v1/webhooks/deploy/{site_id}
X-Hub-Signature-256: sha256=...
```

The Laravel webhook receiver verifies the signature before dispatching the deployment job.

---

## 5. Typical Deployment Script Templates

### 5.1 Laravel Application
```bash
git pull origin main
composer install --no-dev --optimize-autoloader
php artisan migrate --force
php artisan optimize:clear
php artisan config:cache
php artisan route:cache
php artisan view:cache
```

### 5.2 Node.js Application
```bash
git pull origin main
npm ci --production
npm run build
# Systemd application restart is triggered via agent hook
```

---

## 6. Rollback Mechanism

When a deployment fails or causes an application error:
1. The administrator selects a previous successful deployment from the **Deployment History** list.
2. The agent executes a checkout of the previous commit SHA:
   ```bash
   git checkout <PREVIOUS_COMMIT_SHA>
   ```
3. Re-runs the dependency install and service reload.
4. Records the action as a new deployment with `trigger_source = rollback`.
