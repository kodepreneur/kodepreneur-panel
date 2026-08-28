# ADR-001: Two-Tier Architecture (Control Plane vs. Privileged Agent)

## Status
**Accepted**

## Date
2026-08-28

## Context
Traditional VPS management panels (e.g. early cPanel scripts, aaPanel, Webmin) often run the web application tier with full `root` or unrestricted `sudo` privileges. Any vulnerability in the web framework, third-party PHP package, or file upload endpoint immediately results in complete root-level compromise of the underlying host server.

Furthermore, monolithic architectures couple the presentation layer directly with system command execution, making testing, auditing, and maintenance hazardous.

## Decision
We adopt a strict **Two-Tier Architecture**:
1. **Control Plane (Tier 1)**: Built with **Laravel 13 + Vue 3 (Inertia.js)**, running as an unprivileged system user (`kodepreneur`). It handles authentication, UI, business validation, and state persistence in SQLite. It has **no root access** and **cannot execute shell commands**.
2. **Execution Plane (Tier 2)**: Built as a compiled **Go Daemon (`kodepreneur-agent`)**, running as `root` via systemd. It receives strongly typed DTO requests from the control plane over a local socket/TLS and executes discrete, validated Linux system operations.

## Consequences

### Positive
* **Defense in Depth**: A security compromise in Laravel does not automatically grant host root execution.
* **Separation of Concerns**: UI and business logic evolve independently of low-level Linux system interfaces.
* **Auditability**: All privileged operations pass through a single, auditable Go API.
* **Reliability**: Web server crashes or PHP worker restarts do not take down the privileged management daemon.

### Negative / Tradeoffs
* Minor overhead of IPC (Inter-Process Communication) between Laravel and Go Agent.
* Requires coordinating updates between two codebases (Laravel panel and Go binary) during major releases.
