# ADR-004: Authentication Architecture & Agent Communication Security

## Status
**Accepted**

## Date
2026-08-28

## Context
The system has two distinct authentication boundaries:
1. **User / Administrator to Web Control Plane**: Human operators interacting with the browser UI or developer tools invoking the public Control Plane API.
2. **Control Plane to Privileged Agent**: Internal communication between unprivileged Laravel PHP processes and the privileged Go Agent daemon.

Both boundaries require strict authentication, anti-tampering, and replay protection.

## Decision

### 1. User & Web Authentication
* **Web UI Session**: Secure HTTP-only, `SameSite=Lax/Strict` cookies managed via Laravel session middleware, protected by Inertia CSRF tokens.
* **REST API Authentication**: Standard Laravel Sanctum tokens with hashed tokens in the database and explicit permission abilities.
* **Two-Factor Authentication**: RFC 6238 Time-based One-Time Password (TOTP) supported for all administrative accounts.
* **Password Hashing**: Argon2id or Bcrypt with high work factor.

### 2. Control Plane to Agent IPC Authentication
* **Channel**: Restrict communication to a local Unix Domain Socket (`/run/kodepreneur/agent.sock`) with filesystem permissions `0660` owned by `root:kodepreneur`, or localhost HTTPS.
* **Protocol**: **HMAC-SHA256 Signed Payloads**.
  * A 256-bit cryptographically secure random key is generated during installation and stored in `/etc/kodepreneur/agent.yaml` (`0600 root:root`) and the Laravel `.env` (`0600 kodepreneur:kodepreneur`).
  * Every request calculates a signature over `Method + URI + Timestamp + Nonce + BodyHash`.
  * The Go Agent rejects requests with timestamp skew > 30s or reused nonces to prevent replay attacks.

## Consequences

### Positive
* Prevents unauthorized local users from querying the privileged agent daemon directly.
* Cryptographically verifies payload integrity before the Go Agent parses and executes actions.
* Standardized, battle-tested web authentication patterns for operators.

### Negative / Tradeoffs
* Clock drift between host clock and container (if containerized) could cause request rejection if timestamp skew exceeds 30 seconds (mitigated by system NTP synchronization).
