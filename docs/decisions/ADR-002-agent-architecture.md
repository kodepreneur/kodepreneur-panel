# ADR-002: Go-Based Privileged Agent with Zero Generic Shell Execution

## Status
**Accepted**

## Date
2026-08-28

## Context
Many open-source control panels implement server management by exposing a generic `/api/exec` endpoint that accepts arbitrary bash script strings. This pattern is vulnerable to command injection, escaping flaws, and unverified parameters.

Additionally, using interpreted languages (like Python or PHP) for the privileged root daemon introduces heavy runtime dependencies, package conflicts with user applications, and high memory footprints.

## Decision
1. **Implementation Language**: Go (Golang) compiled into a single static binary (`/usr/local/bin/kodepreneur-agent`) with zero external runtime dependencies.
2. **Zero Generic Shell Execution**: We explicitly reject generic `POST /execute` or `sh -c` endpoints. Every supported operation must be mapped to a dedicated typed endpoint (e.g. `POST /api/v1/websites`, `POST /api/v1/php/{version}/restart`, `POST /api/v1/databases`).
3. **Structured Binary Invocation**: When external system binaries must be invoked (`nginx`, `systemctl`, `ufw`, `certbot`), the agent invokes Go's `os/exec` directly passing arguments as slices (e.g. `exec.Command("systemctl", "reload", "nginx")`) rather than invoking a shell.
4. **Pre-flight Validation & Rollback**: Any configuration file (such as Nginx vhosts) is validated syntactically (`nginx -t`) before being committed or symlinked. If validation fails, changes are automatically reverted and a descriptive error is returned.

## Consequences

### Positive
* **Immunity to Shell Injection**: Passing argument arrays directly to syscalls avoids shell string interpolation bugs.
* **Minimal Memory Footprint**: The Go agent consumes < 25MB RAM at idle and requires no separate Python/Ruby runtime.
* **Predictable Failure Handling**: Syntactic checks prevent broken Nginx or PHP-FPM configurations from bringing down all hosted sites.

### Negative / Tradeoffs
* Adding new system capabilities requires updating the Go agent's endpoint and data transfer models rather than simply dispatching an arbitrary bash script.
