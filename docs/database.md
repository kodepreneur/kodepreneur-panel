# Database Architecture & Schema Specification

## 1. Overview

**Kodepreneur Panel** uses a local **SQLite** database (`/var/lib/kodepreneur/panel.sqlite`) for the internal Control Plane state.

### 1.1 Why SQLite for the Control Plane?
1. **Zero External Dependency**: The panel can boot, recover, and function even if the managed database daemons (MySQL, PostgreSQL) are stopped, broken, or being upgraded.
2. **Atomic Backups**: The entire panel state can be snapshotted or restored with a single file copy or the `VACUUM INTO` command.
3. **High Performance**: In Write-Ahead Logging (WAL) mode with memory caching, SQLite provides sub-millisecond query responses with minimal CPU and RAM overhead.
4. **Self-Contained**: Eliminates the chicken-and-egg bootstrap dilemma during server initialization.

---

## 2. SQLite Engine Optimization PRAGMAs

Kodepreneur Panel configures the following SQLite connection settings in Laravel:

```sql
PRAGMA journal_mode = WAL;          -- Concurrent reads while writing
PRAGMA synchronous = NORMAL;        -- Safe durability with peak write speed
PRAGMA foreign_keys = ON;           -- Enforce relational referential integrity
PRAGMA busy_timeout = 5000;         -- 5-second lock queue before timeout
PRAGMA cache_size = -2000;          -- 2MB page memory cache
PRAGMA temp_store = MEMORY;         -- In-memory temporary tables and indexes
```

---

## 3. Entity-Relationship Model

```text
┌─────────────────┐       ┌─────────────────┐
│      users      │──────<│  activity_logs  │
└────────┬────────┘       └─────────────────┘
         │
         │ 1:1
         ▼
┌─────────────────┐
│      roles      │
└─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│    websites     │──────<│     domains     │
└────────┬────────┘       └─────────────────┘
         │
         ├───< ┌─────────────────┐
         │     │   deployments   │
         │     └─────────────────┘
         │
         ├───< ┌─────────────────┐
         │     │    processes    │
         │     └─────────────────┘
         │
         ├───< ┌─────────────────┐
         │     │    cron_jobs    │
         │     └─────────────────┘
         │
         └───< ┌─────────────────┐
               │ssl_certificates │
               └─────────────────┘

┌─────────────────┐       ┌─────────────────┐
│    databases    │──────<│ database_access │>──────┌─────────────────┐
└─────────────────┘       └─────────────────┘       │ database_users  │
                                                    └─────────────────┘
```

---

## 4. Complete Schema & Field Definitions

### 4.1 `users`
Stores control plane administrators and operators.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key (ULID / UUIDv7) |
| `name` | `VARCHAR(255)` | ❌ | User's display name |
| `email` | `VARCHAR(255)` | ❌ | Unique email address |
| `password` | `VARCHAR(255)` | ❌ | Bcrypt / Argon2id password hash |
| `role_id` | `VARCHAR(36)` | ❌ | FK referencing `roles.id` |
| `two_factor_secret` | `TEXT` | ✅ | Encrypted TOTP 2FA secret |
| `two_factor_confirmed_at`| `DATETIME` | ✅ | TOTP activation timestamp |
| `remember_token` | `VARCHAR(100)` | ✅ | Session remember token |
| `last_login_at` | `DATETIME` | ✅ | Timestamp of most recent authentication |
| `last_login_ip` | `VARCHAR(45)` | ✅ | IP address of most recent authentication |
| `created_at` | `DATETIME` | ❌ | Record creation timestamp |
| `updated_at` | `DATETIME` | ❌ | Record modification timestamp |

---

### 4.2 `roles`
Defines access control tiers and assigned permissions.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `name` | `VARCHAR(100)` | ❌ | Role name (e.g. `Super Admin`, `Operator`) |
| `slug` | `VARCHAR(100)` | ❌ | Unique slug identifier (`super-admin`, `operator`) |
| `permissions` | `JSON` | ❌ | Array of assigned permission keys |
| `created_at` | `DATETIME` | ❌ | Record creation timestamp |
| `updated_at` | `DATETIME` | ❌ | Record modification timestamp |

---

### 4.3 `websites`
Core entity tracking managed virtual hosts and runtime settings.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `domain` | `VARCHAR(255)` | ❌ | Unique primary domain name (e.g. `example.com`) |
| `aliases` | `JSON` | ✅ | Array of secondary domains (e.g. `["www.example.com"]`) |
| `php_version` | `VARCHAR(10)` | ❌ | PHP runtime version (`8.3`, `8.4`, or `none`) |
| `document_root` | `VARCHAR(255)` | ❌ | Absolute path on server (e.g. `/var/www/example.com/public`) |
| `system_user` | `VARCHAR(64)` | ❌ | Isolated Linux system user owning the site files |
| `ssl_enabled` | `BOOLEAN` | ❌ | Default: `0`. Flag indicating active HTTPS |
| `force_https` | `BOOLEAN` | ❌ | Default: `0`. Automatic 301 redirect to HTTPS |
| `status` | `VARCHAR(32)` | ❌ | `active`, `suspended`, `provisioning`, `error` |
| `custom_nginx_config` | `TEXT`| ✅ | Optional custom directives injected into server block |
| `created_at` | `DATETIME` | ❌ | Creation timestamp |
| `updated_at` | `DATETIME` | ❌ | Modification timestamp |

---

### 4.4 `databases`
Tracks managed MySQL and PostgreSQL database instances.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `engine` | `VARCHAR(32)` | ❌ | `mysql` or `postgresql` |
| `name` | `VARCHAR(64)` | ❌ | Database name on the host daemon |
| `character_set`| `VARCHAR(32)` | ❌ | Default: `utf8mb4` (MySQL) or `UTF8` (Postgres) |
| `collation` | `VARCHAR(64)` | ❌ | Default: `utf8mb4_unicode_ci` |
| `size_bytes` | `BIGINT` | ✅ | Cached database size on disk |
| `created_at` | `DATETIME` | ❌ | Creation timestamp |
| `updated_at` | `DATETIME` | ❌ | Modification timestamp |

---

### 4.5 `database_users`
Tracks managed database credentials and access hosts.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `engine` | `VARCHAR(32)` | ❌ | `mysql` or `postgresql` |
| `username` | `VARCHAR(64)` | ❌ | Username on target database engine |
| `host` | `VARCHAR(64)` | ❌ | Allowed host pattern (default `localhost` or `127.0.0.1`) |
| `created_at` | `DATETIME` | ❌ | Creation timestamp |
| `updated_at` | `DATETIME` | ❌ | Modification timestamp |

---

### 4.6 `database_access`
Pivot table defining user grants to specific databases.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `database_id` | `VARCHAR(36)` | ❌ | FK referencing `databases.id` (Cascade on delete) |
| `database_user_id`| `VARCHAR(36)` | ❌ | FK referencing `database_users.id` (Cascade on delete) |
| `permissions` | `VARCHAR(32)` | ❌ | `all`, `read_write`, `read_only` |
| `created_at` | `DATETIME` | ❌ | Creation timestamp |

---

### 4.7 `deployments`
Records Git build runs, logs, execution time, and status.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `website_id` | `VARCHAR(36)` | ❌ | FK referencing `websites.id` |
| `commit_hash` | `VARCHAR(40)` | ✅ | Git commit SHA |
| `commit_message` | `TEXT` | ✅ | Git commit summary message |
| `branch` | `VARCHAR(100)` | ❌ | Target Git branch (e.g. `main`) |
| `status` | `VARCHAR(32)` | ❌ | `pending`, `running`, `success`, `failed`, `cancelled` |
| `trigger_source` | `VARCHAR(32)` | ❌ | `manual`, `webhook`, `rollback` |
| `log_output` | `MEDIUMTEXT` | ✅ | Full stdout / stderr compilation output |
| `duration_seconds` | `INTEGER` | ✅ | Total execution time in seconds |
| `initiated_by_user_id` | `VARCHAR(36)` | ✅ | FK referencing `users.id` |
| `created_at` | `DATETIME` | ❌ | Timestamp of trigger |
| `updated_at` | `DATETIME` | ❌ | Timestamp of completion |

---

### 4.8 `ssl_certificates`
Tracks Let's Encrypt and custom SSL certificates.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `website_id` | `VARCHAR(36)` | ❌ | FK referencing `websites.id` |
| `domain` | `VARCHAR(255)` | ❌ | Certificate subject domain |
| `issuer` | `VARCHAR(100)` | ❌ | `letsencrypt` or `custom` |
| `cert_path` | `VARCHAR(255)` | ❌ | Filesystem path to `fullchain.pem` |
| `key_path` | `VARCHAR(255)` | ❌ | Filesystem path to `privkey.pem` |
| `valid_from` | `DATETIME` | ❌ | Certificate validity start |
| `valid_until` | `DATETIME` | ❌ | Expiration timestamp |
| `auto_renew` | `BOOLEAN` | ❌ | Default: `1`. Automated cron renewal |
| `status` | `VARCHAR(32)` | ❌ | `valid`, `expiring_soon`, `expired`, `revoked` |
| `created_at` | `DATETIME` | ❌ | Creation timestamp |
| `updated_at` | `DATETIME` | ❌ | Modification timestamp |

---

### 4.9 `firewall_rules`
Tracks managed UFW port rules.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `port` | `INTEGER` | ❌ | Target TCP/UDP port number (1 - 65535) |
| `protocol` | `VARCHAR(10)` | ❌ | `tcp`, `udp`, `both` |
| `action` | `VARCHAR(10)` | ❌ | `allow`, `deny`, `reject` |
| `source_ip` | `VARCHAR(45)` | ✅ | Source IP mask (default `0.0.0.0/0` for any) |
| `description` | `VARCHAR(255)` | ✅ | User note (e.g. `Custom SSH`, `Adminer`) |
| `created_at` | `DATETIME` | ❌ | Rule creation timestamp |

---

### 4.10 `activity_logs`
Immutable audit ledger for security tracking.

| Column | Type | Nullable | Constraints & Description |
| :--- | :--- | :---: | :--- |
| `id` | `VARCHAR(36)` | ❌ | Primary Key |
| `user_id` | `VARCHAR(36)` | ✅ | FK referencing `users.id` (Null if system job) |
| `user_email` | `VARCHAR(255)` | ✅ | Snapshot of user email at execution time |
| `ip_address` | `VARCHAR(45)` | ❌ | Request client IP |
| `user_agent` | `TEXT` | ✅ | Browser / API client user-agent |
| `action` | `VARCHAR(100)` | ❌ | Action identifier (e.g. `website.create`) |
| `resource_type`| `VARCHAR(50)` | ❌ | Target entity type (`website`, `database`, `user`) |
| `resource_id` | `VARCHAR(36)` | ✅ | Identifier of the target entity |
| `status` | `VARCHAR(20)` | ❌ | `success`, `failure` |
| `payload_summary`| `JSON` | ✅ | Redacted metadata summary |
| `created_at` | `DATETIME` | ❌ | Timestamp of occurrence |
