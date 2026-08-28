# ADR-003: SQLite for Control Plane Internal State Storage

## Status
**Accepted**

## Date
2026-08-28

## Context
A server control panel needs to persist records of websites, SSL certificates, users, deployments, and logs. Traditional panels often require a running MySQL or PostgreSQL instance just to serve the panel itself.

This creates critical problems:
1. **Bootstrapping Dependency**: The panel cannot start if MySQL fails to initialize or experiences a crash.
2. **Resource Overhead**: Running an additional MySQL instance on a small 1GB RAM VPS consumes significant baseline memory.
3. **Upgrade Fragility**: Upgrading or maintaining user database services risks breaking the control panel itself.

## Decision
We choose **SQLite** (`/var/lib/kodepreneur/panel.sqlite`) as the dedicated database engine for the Kodepreneur Control Plane.

We enforce standard production optimizations:
* Write-Ahead Logging (`PRAGMA journal_mode = WAL;`)
* Memory temp stores (`PRAGMA temp_store = MEMORY;`)
* Enforced foreign key checks (`PRAGMA foreign_keys = ON;`)
* Pragmatic lock timeout (`PRAGMA busy_timeout = 5000;`)

Managed user workloads (websites/applications) can independently install and use PostgreSQL or MySQL/MariaDB without any coupling to the panel's internal state storage.

## Consequences

### Positive
* **High Availability**: Panel is 100% resilient to crashes or restarts of MySQL/PostgreSQL daemons.
* **Instant Backups**: Panel state can be snapshotted simply by copying the single `.sqlite` file.
* **Ultra-low Footprint**: Negligible RAM usage (< 5MB) for panel data persistence.
* **Zero Configuration**: No database passwords or user grants required during bootstrap.

### Negative / Tradeoffs
* Single-writer concurrency model: SQLite serializes write operations. For a single-node VPS panel handling administrative requests, WAL mode handles hundreds of concurrent reads and serialized writes with sub-millisecond latency, making this limitation irrelevant for our scope.
