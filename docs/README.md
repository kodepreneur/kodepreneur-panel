# Kodepreneur Panel Documentation

Welcome to the official technical documentation and architecture blueprints for **Kodepreneur Panel**.

> **Philosophy**: Simple server management without unnecessary complexity.

Kodepreneur Panel is a lightweight, fast, and secure VPS control panel designed to allow developers and system administrators to manage websites, PHP runtimes, databases, services, SSL certificates, Git deployments, cron jobs, and basic server resources without wrestling with SSH.

---

## 📚 Documentation Index

| Document | Description |
| :--- | :--- |
| [Architecture Overview](file:///docs/architecture.md) | High-level system design, control plane vs. agent model, component boundaries. |
| [Installation Guide](file:///docs/installation.md) | Modular bootstrap installer design, system requirements, idempotency, systemd setup. |
| [Development Guide](file:///docs/development.md) | Local development workflow, code standards, branching model, testing protocols. |
| [Security Architecture](file:///docs/security.md) | Threat model, privilege separation, typed execution, authentication, path traversal prevention. |
| [Agent Architecture & API](file:///docs/agent.md) | Go-based privileged daemon design, typed RPC/REST API, system interaction patterns. |
| [Panel API Reference](file:///docs/api.md) | Laravel Control Plane REST API endpoints, schemas, authentication, and responses. |
| [Database Architecture](file:///docs/database.md) | SQLite panel schema, entity-relationship models, normalization, indexing, and migration rules. |
| [Deployment & Git Engine](file:///docs/deployment.md) | Zero-downtime & atomic Git deployments, unprivileged execution, rollbacks, and logs. |
| [Troubleshooting & Runbook](file:///docs/troubleshooting.md) | Diagnosing agent failures, Nginx validation issues, SSL errors, log locations, recovery flows. |
| [Roadmap & Phases](file:///docs/roadmap.md) | Development roadmap from Phase 0 to Phase 7, MVP definition, future scope. |
| [Changelog](file:///docs/changelog.md) | Historical record of architectural shifts, releases, and updates. |
| [Architecture Decision Records (ADRs)](file:///docs/decisions/README.md) | Formal records of major architectural choices and rationales. |

---

## 🏛️ Core Principles

1. **Keep it Simple & Boring**: Predictable, maintainable, and standard Linux tooling (`systemd`, `nginx`, `php-fpm`, `ufw`, `certbot`, `sqlite3`) over bespoke orchestrators.
2. **Strict Privilege Separation**: The Laravel web panel runs as an unprivileged user (`kodepreneur-panel`) and **never** executes root commands or raw shell inputs. Privileged actions are handled exclusively by the Go agent (`kodepreneur-agent`) over an authenticated, strictly typed API.
3. **Minimal Dependencies**: Every dependency must be justified and documented.
4. **Documentation First**: No code, command, migration, or configuration is executed without prior analysis, documentation, and implementation planning.

---

## 🎯 Supported Target Environments

* **Operating System**: Ubuntu 24.04 LTS (Noble Numbat)
* **Architecture**: `x86_64` (amd64), `aarch64` (arm64)
* **Minimum Specifications**: 1 vCPU, 1 GB RAM, 10 GB Disk space
