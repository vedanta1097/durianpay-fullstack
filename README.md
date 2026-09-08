# Fullstack App

A full-stack internal dashboard for Customer Support and Operation users to monitor incoming payments. The backend uses Go and SQLite; the frontend uses React, TypeScript, Zustand, and Tailwind CSS.

## Tool versions

The project supports the reviewer environment defined in the assignment:

```text
Go 1.21+
Node.js 20+
Docker and Docker Compose
Make
macOS
```

---

## Install requirements

For native development, install the Xcode Command Line Tools on macOS. The SQLite driver uses CGO and needs the C compiler (`clang`) provided by these tools:

```bash
xcode-select --install
```

Then install the Go and npm dependencies from the repository root:

```bash
make setup
```

The setup command creates `backend/.env` from `backend/env.sample` only when it does not exist.

---

## Run the full app locally

```bash
make run
```

This starts both the Go backend at `http://localhost:8080` and the Vite frontend at `http://localhost:5173` for local development.

Open [http://localhost:5173](http://localhost:5173). Press `Ctrl+C` to stop both services.

---

## Run the backend locally

```bash
make run-backend
```

The API runs at [http://localhost:8080](http://localhost:8080).

## Run the backend production build

```bash
make -C backend build
cd backend
./bin/mygolangapp
```

---

## Run the frontend locally

Start the backend first, then run:

```bash
make run-frontend
```

The frontend runs at [http://localhost:5173](http://localhost:5173). No frontend `.env` is required because Vite proxies `/dashboard/v1` to the backend.

To use another API origin, copy `frontend/.env.example` to `frontend/.env` and set `VITE_API_BASE_URL`.

## Run the frontend production build

Docker Compose builds the frontend and serves it through Nginx together with the Go API:

```bash
make docker-up
```

Open [http://localhost:8080](http://localhost:8080). Press `Ctrl+C` or run `make docker-down` to stop it. Docker Desktop and Make are the only host requirements for this workflow.

---

## OpenAPI documentation

After the backend is running, open Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

The source specification is also available at:

```text
http://localhost:8080/openapi.yaml
```

The hand-authored [`openapi.yaml`](openapi.yaml) at the repository root remains the API source of truth. Backend and frontend types are generated from it; no `swag init` annotations are used.

---

## Login to the frontend

Use [http://localhost:5173](http://localhost:5173) for native development or [http://localhost:8080](http://localhost:8080) for Docker.

| Role             | Email                | Password   |
| ---------------- | -------------------- | ---------- |
| Customer Support | `cs@test.com`        | `password` |
| Operation        | `operation@test.com` | `password` |

The backend creates the SQLite schema and deterministic payment seed data at startup.

---

## Build and test

Build both applications:

```bash
make build
```

Run all backend and frontend tests:

```bash
make test
```

Run the complete pre-submission check:

```bash
make generate-check
make format-check
make lint
make test
make build
```

## Evidence

See the [testing video](https://drive.google.com/file/d/1UxIimV77z-y5UEvmxJkpcMsTqzIeTG-D/view?usp=sharing).

---

See the [backend README](backend/README.md) and [frontend README](frontend/README.md) for architecture, testing details, and trade-offs.
