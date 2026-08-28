# Development & Contributing Guide

## 1. Development Principles

All contributors and engineers working on **Kodepreneur Panel** must adhere to the core principles:

1. **Documentation First**: Every feature, schema change, or architectural decision must be documented and reviewed before code is implemented.
2. **Minimize Dependencies**: Do not introduce third-party libraries unless strictly justified and documented.
3. **Keep It Boring & Simple**: Write clean, idiomatic code adhering to standard PSR-12/PER (PHP), standard Go conventions (`gofmt`, `golangci-lint`), and TypeScript strict typing.
4. **Security by Default**: Never execute unsanitized shell commands, never bypass path jail boundaries, and always validate input at both the Laravel and Go Agent layers.

---

## 2. Project Directory Structure

```text
kodepreneur_panel/
├── .github/                    # CI/CD Workflows (linting, tests, build)
├── docs/                       # Complete documentation & ADRs
│   ├── decisions/              # Architecture Decision Records
│   ├── architecture.md
│   ├── installation.md
│   └── ...
├── installer/                  # Modular VPS bootstrap installer
│   ├── install.sh
│   ├── lib/
│   └── modules/
├── agent/                      # Go Privileged Daemon Codebase
│   ├── cmd/agent/main.go
│   ├── config/
│   ├── server/
│   ├── pkg/
│   ├── go.mod
│   └── go.sum
└── panel/                      # Laravel 13 + Vue 3 Control Plane
    ├── app/
    │   ├── Http/Controllers/
    │   ├── Models/
    │   ├── Services/Agent/     # Typed Go Agent Client
    │   └── Jobs/
    ├── config/
    ├── database/
    │   ├── migrations/
    │   └── seeders/
    ├── resources/
    │   ├── js/
    │   │   ├── Components/
    │   │   ├── Layouts/
    │   │   ├── Pages/
    │   │   ├── app.ts
    │   │   └── types/
    │   └── css/
    ├── routes/
    │   ├── web.php
    │   └── api.php
    ├── composer.json
    ├── package.json
    ├── vite.config.ts
    └── tailwind.config.js
```

---

## 3. Local Environment Setup

### 3.1 Prerequisites
* PHP 8.3+ with `sqlite3`, `mbstring`, `xml`, `curl`, `zip` extensions.
* Composer 2.x.
* Node.js 20+ and NPM / PNPM.
* Go 1.22+.

### 3.2 Setting up the Laravel Control Plane (`panel/`)
```bash
cd panel
cp .env.example .env
composer install
npm install
php artisan key:generate
touch database/panel.sqlite
php artisan migrate --seed
npm run build # or npm run dev for HMR
php artisan serve --port=8000
```

### 3.3 Building & Running the Go Agent (`agent/`)
```bash
cd agent
go mod tidy
go build -o kodepreneur-agent cmd/agent/main.go

# Run in mock / dev mode (uses mock system managers when not running on root Linux)
./kodepreneur-agent --config agent.dev.yaml --dev
```

---

## 4. Coding Standards & Tooling

### 4.1 PHP / Laravel
* Adhere to **PSR-12** and **PER Coding Style**.
* Run static analysis:
  ```bash
  ./vendor/bin/pint
  ./vendor/bin/phpstan analyse --level=8
  ```
* Run automated tests:
  ```bash
  php artisan test
  ```

### 4.2 Go (Agent)
* Adhere to standard Go formatting and idioms.
* Run linter & tests:
  ```bash
  gofmt -s -w .
  golangci-lint run
  go test -v -race ./...
  ```

### 4.3 Frontend (Vue 3 + TypeScript)
* Use Vue 3 Composition API with `<script setup lang="ts">`.
* Strict TypeScript: No `any` types unless explicitly justified.
* Type checking:
  ```bash
  npm run typecheck
  ```

---

## 5. Workflow: Adding a New Feature

Follow this mandatory 8-step sequence:
1. **Analyze Requirement**: Define user story, security boundaries, and failure modes.
2. **Update Documentation**: Update relevant docs in `/docs` and create an ADR if architectural.
3. **Define Implementation Plan**: Create an implementation plan detailing files and test cases.
4. **Agent Layer (if applicable)**: Implement typed Go handler, validator, and tests.
5. **Panel Layer**: Implement Laravel models, migrations, Agent DTO client, controllers, and tests.
6. **Frontend Layer**: Implement Vue 3 Inertia page, components, and responsive layout.
7. **End-to-End Test**: Verify complete flow on target Ubuntu 24.04 test VM.
8. **Update Changelog**: Record the change in `docs/changelog.md`.
