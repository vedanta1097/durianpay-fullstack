# DurianPay Payment Dashboard

A small full-stack internal dashboard for authenticating support/operations users and reviewing incoming payments. The backend is a Go HTTP API backed by SQLite; the frontend is a React, TypeScript, Zustand, and Tailwind CSS single-page application.

## Prerequisites

- Go 1.21 or newer
- Node.js 20 or newer and npm
- Make
- A C compiler for `github.com/mattn/go-sqlite3`
  - On macOS: install Xcode Command Line Tools with `xcode-select --install`

## Quick start

From the repository root, install both backend and frontend dependencies:

```bash
make setup
```

The setup command also creates `backend/.env` from `backend/env.sample` when it does not already exist. Existing environment configuration is never overwritten.

Run both development servers with one command:

```bash
make run
```

Open [http://localhost:5173](http://localhost:5173). The Vite development server proxies `/dashboard/v1` API requests to the backend at `http://localhost:8080`.

Stop both servers with `Ctrl+C`.

### Seed accounts

| Role | Email | Password |
| --- | --- | --- |
| Customer Support | `cs@test.com` | `password` |
| Operation | `operation@test.com` | `password` |

The SQLite schema and deterministic demo data are initialized idempotently when the backend starts. The default local database is `backend/dashboard.db`.

## Run one service at a time

Use two terminal windows when you want to inspect each service independently.

Backend:

```bash
make run-backend
```

Frontend:

```bash
make run-frontend
```

The equivalent service-local commands are `make -C backend run` and `npm --prefix frontend run dev`.

## Configuration

Backend defaults live in `backend/env.sample`:

| Variable | Default | Purpose |
| --- | --- | --- |
| `HTTP_ADDR` | `:8080` | Backend listen address |
| `OPENAPIYAML_LOCATION` | `../openapi.yaml` | OpenAPI contract path from the backend directory |
| `DB_PATH` | `dashboard.db` | SQLite database path from the backend directory |
| `JWT_SECRET` | `your-very-secret` | Local JWT signing secret; replace outside local development |
| `JWT_EXPIRED` | `24h` | JWT lifetime |

The frontend normally needs no local environment file because its development proxy targets port 8080. To use another API origin, copy `frontend/.env.example` to `frontend/.env` and set `VITE_API_BASE_URL`.

## API

The hand-authored API contract is [`openapi.yaml`](openapi.yaml). Generated backend and frontend types are committed, so code generation is not required to run the application.

| Method | Path | Authentication |
| --- | --- | --- |
| `POST` | `/dashboard/v1/auth/login` | Public |
| `GET` | `/dashboard/v1/payments` | Bearer token |

The payments endpoint supports the `status` and `sort` query parameters defined in the OpenAPI contract. Its response contains the matching `payments` array plus a global `summary` with total, completed, processing, and failed counts. Selecting a dashboard status makes a new server request; the summary stays global while the table is filtered.

## Development commands

Run these commands from the repository root:

| Command | Purpose |
| --- | --- |
| `make help` | List available commands |
| `make setup` | Install backend and frontend dependencies |
| `make run` | Run backend and frontend together |
| `make run-backend` | Run only the Go API |
| `make run-frontend` | Run only the Vite frontend |
| `make generate` | Regenerate backend and frontend OpenAPI code |
| `make generate-check` | Check generated OpenAPI code for drift |
| `make format` | Format Go code |
| `make format-check` | Check Go formatting |
| `make lint` | Run Go vet and TypeScript type checking |
| `make test` | Run all backend and frontend tests |
| `make build` | Build both applications |

## Build and preview

Build both applications:

```bash
make build
```

Run the compiled backend:

```bash
cd backend
./bin/mygolangapp
```

Preview the built frontend locally in another terminal:

```bash
npm --prefix frontend run preview
```

`vite preview` is intended for reviewing the production build locally. A real deployment should serve `frontend/dist` through a production static-file server and route `/dashboard` to the Go API.

## Testing strategy

```bash
make test
```

- Frontend page/component tests mock their typed API modules with Vitest.
- Frontend API-client tests stub `fetch`; they never contact a real server.
- Backend handler and use-case tests use mocks or `httptest`.
- SQLite integration tests use isolated in-memory or temporary test databases and never touch `backend/dashboard.db`.

Before submitting or reviewing a change, run:

```bash
make generate-check
make format-check
make lint
make test
make build
```
