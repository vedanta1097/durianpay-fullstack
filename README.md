# DurianPay Payment Dashboard

An internal dashboard for Customer Support and Operation users to monitor incoming payments.

- Backend: Go 1.21+, SQLite, OpenAPI
- Frontend: React, TypeScript, Zustand, Tailwind CSS

## Quick start

### Option 1 - Native with Make

Requirements: Go 1.21+, Node.js 20+, npm, Make, and a C compiler for SQLite. On macOS, install the compiler with:

```bash
xcode-select --install
```

From the repository root:

```bash
# 1. Download Go and npm dependencies
make setup

# 2. Start the backend and frontend
make run
```

Open [http://localhost:5173](http://localhost:5173). Press `Ctrl+C` to stop both services.

`make setup` creates `backend/.env` from `backend/env.sample` only when it does not already exist. Edit that file only when you need to override the backend address, database path, or JWT settings.

No frontend `.env` is required for local development because Vite proxies `/dashboard/v1` to the backend. To use another API origin, copy `frontend/.env.example` to `frontend/.env` and set `VITE_API_BASE_URL`.

To run each service separately:

```bash
# Terminal 1 - Go API at http://localhost:8080
make run-backend

# Terminal 2 - Vite app at http://localhost:5173
make run-frontend
```

### Option 2 - Docker Compose

Requirements: Docker Desktop with Docker Compose, and Make. Go and Node.js are not required on the host.

```bash
make docker-up
```

Open [http://localhost:8080](http://localhost:8080). Press `Ctrl+C` or run `make docker-down` to stop the application.

## Login and seed data

| Role             | Email                | Password   |
| ---------------- | -------------------- | ---------- |
| Customer Support | `cs@test.com`        | `password` |
| Operation        | `operation@test.com` | `password` |

The backend creates the SQLite schema and deterministic seed data at startup.

- Native database: `backend/dashboard.db`
- Docker database: `dashboard-data` named volume

The databases are separate. `make docker-down` preserves the Docker volume. To reset it, run `docker compose down --volumes`.

## Build and test

Build both applications:

```bash
make build
```

For a production-like end-to-end run, use:

```bash
make docker-up
```

Run all backend and frontend tests:

```bash
make test
```

To run the complete quality check:

```bash
make generate-check
make format-check
make lint
make test
make build
```

Testing strategy:

- Backend handlers and use cases use mocks or `httptest`.
- SQLite integration tests use isolated in-memory or temporary databases.
- Frontend component and API tests mock the API client or `fetch` with Vitest.
- Tests never contact a real API or the runtime database.

## API documentation

The OpenAPI v3 contract is [`openapi.yaml`](openapi.yaml). Open it with an OpenAPI or Swagger preview extension.

| Method | Path                       | Authentication |
| ------ | -------------------------- | -------------- |
| `POST` | `/dashboard/v1/auth/login` | Public         |
| `GET`  | `/dashboard/v1/payments`   | Bearer token   |

The payments endpoint supports the `status` and `sort` query parameters. Its response contains the filtered `payments` and a global status `summary`.

Backend and frontend types are generated from `openapi.yaml` and committed. After changing the contract, run:

```bash
make generate
```

## Technical notes

- The frontend uses generated OpenAPI types through a typed API client.
- Zustand stores only the persisted login session.
- Payment filtering runs in the backend; dashboard totals come from the API summary.
- More backend details are available in [`backend/README.md`](backend/README.md).
