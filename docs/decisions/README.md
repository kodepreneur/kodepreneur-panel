# Architecture Decision Records (ADRs)

## Overview
This directory contains the formal records of all architecturally significant decisions made during the design and development of **Kodepreneur Panel**.

An Architecture Decision Record (ADR) captures a single architectural decision along with its context, considered options, decision rationale, and positive/negative consequences.

---

## Index of ADRs

| ADR ID | Title | Status | Date |
| :--- | :--- | :--- | :--- |
| [ADR-001](file:///docs/decisions/ADR-001-panel-architecture.md) | Two-Tier Architecture (Control Plane vs. Privileged Agent) | **Accepted** | 2026-08-28 |
| [ADR-002](file:///docs/decisions/ADR-002-agent-architecture.md) | Go-Based Privileged Agent with Zero Generic Shell Execution | **Accepted** | 2026-08-28 |
| [ADR-003](file:///docs/decisions/ADR-003-database-choice.md) | SQLite for Control Plane Internal State Storage | **Accepted** | 2026-08-28 |
| [ADR-004](file:///docs/decisions/ADR-004-authentication-and-agent-communication.md) | HMAC-SHA256 Authenticated Agent Communication & Web RBAC | **Accepted** | 2026-08-28 |

---

## ADR Template
When proposing a new ADR, follow this structure:
1. **Title**: `ADR-XXX-[short-title].md`
2. **Status**: `Proposed` | `Accepted` | `Deprecated` | `Superseded`
3. **Context**: What is the problem or requirement being addressed?
4. **Decision**: What is the chosen solution and technical approach?
5. **Consequences**: What are the positive tradeoffs and negative implications?
