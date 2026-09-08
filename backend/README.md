# Backend

Go HTTP API for authentication and incoming-payment monitoring. It uses SQLite for persistence, JWT bearer authentication, Chi for routing, and the root OpenAPI contract for generated server types and request validation.

## Requirements

- Go 1.21+
- Make
- CGO and a C compiler for `github.com/mattn/go-sqlite3`

On macOS, the compiler is provided by Xcode Command Line Tools:

```bash
xcode-select --install
```

## Setup and run

From the repository root:

```bash
make setup
make run-backend
```

Or from `backend/`:

```bash
make setup
make run
```

`make setup` downloads Go modules and creates `.env` from `env.sample` when it does not exist. Existing `.env` files are not overwritten.

The API listens at [http://localhost:8080](http://localhost:8080).

## Build

From the repository root:

```bash
make -C backend build
```

Run the compiled binary from `backend/` so its relative configuration paths resolve correctly:

```bash
cd backend
./bin/mygolangapp
```

The Docker production-like workflow is documented in the [root README](../README.md).

## Test and quality checks

From `backend/`:

```bash
make test
make format-check
make lint
make openapi-check
make build
```

Use `make format` to apply `gofmt`. `make openapi-gen` regenerates `internal/openapigen/openapi.gen.go` from the root `openapi.yaml` using the pinned generator version.

## Configuration

Defaults are defined in `env.sample`:

| Variable | Default | Purpose |
| --- | --- | --- |
| `HTTP_ADDR` | `:8080` | HTTP listen address |
| `OPENAPIYAML_LOCATION` | `../openapi.yaml` | OpenAPI file served by the documentation route |
| `DB_PATH` | `dashboard.db` | SQLite database path |
| `JWT_SECRET` | `your-very-secret` | HMAC key used to sign JWTs |
| `JWT_EXPIRED` | `24h` | JWT lifetime |

The committed values are development defaults. Do not reuse the example JWT secret in a real deployment.

## API and documentation

| Method | Path | Authentication | Purpose |
| --- | --- | --- | --- |
| `GET` | `/healthz` | Public | Container readiness check |
| `POST` | `/dashboard/v1/auth/login` | Public | Authenticate with email and password |
| `GET` | `/dashboard/v1/payments` | Bearer JWT | List and filter payments |
| `GET` | `/swagger/index.html` | Public | Interactive Swagger UI |
| `GET` | `/openapi.yaml` | Public | Source OpenAPI document |

The payment endpoint supports `status`, `id`, and allow-listed `sort` query parameters. Its response contains the matching payments and a global status summary.

Swagger UI is only a view over the root OpenAPI document. This project does not use `swag init` or handler annotations to create a second specification.

## Architecture

The backend keeps the starter entry point and feature-oriented package layout:

```text
main.go
  -> initializes SQLite schema and deterministic seed data
  -> wires repositories, use cases, handlers, and HTTP server

internal/
  api/                 generated-interface adapter
  config/              environment configuration
  entity/              domain data and application errors
  module/
    auth/
      handler/          HTTP request/response mapping
      usecase/          credential verification and JWT logic
      repository/       user queries
    payment/
      handler/          payment HTTP mapping
      usecase/          filter and sort validation
      repository/       payment list and summary SQL
  openapigen/           generated OpenAPI types and Chi server interface
  service/http/         router, OpenAPI validation, auth middleware, docs
  transport/            shared JSON error responses
```

### Request flow

```text
HTTP request
  -> Chi router
  -> OpenAPI request validation and JWT authentication
  -> internal/api adapter
  -> feature handler
  -> feature use case
  -> repository
  -> SQLite
  -> generated response type -> JSON
```

Handlers own HTTP concerns, use cases own application validation, and repositories own SQL. `context.Context` is passed from each request to database calls.

### Authentication

```text
POST /auth/login
  -> find user by email
  -> compare bcrypt password hash
  -> sign HS256 JWT with subject, role, issued-at, and expiry claims
```

Protected requests validate the bearer token signature, signing algorithm, expiry, subject, and supported role before reaching the handler.

### Persistence and seed

SQLite schema and seed initialization remain in `main.go` to preserve the supplied boilerplate and keep the small application easy to follow. Initialization is idempotent, and multi-row seed writes run in a transaction.

The seed contains two hashed-password users and 30 deterministic payments. Native runs use `backend/dashboard.db`; Docker uses `/data/dashboard.db` in the persistent `dashboard-data` volume. Those databases are separate, and tests never use either runtime database.

## Testing strategy

- Server tests use `httptest` and stub handlers/token verification.
- Use-case tests mock repository interfaces.
- Repository tests use isolated SQLite databases to verify real SQL behavior.
- `main_test.go` exercises schema, idempotent seed behavior, and the HTTP vertical slice with an isolated database.
- Tests do not start the real application server or access external services.

## Trade-offs

### SQLite instead of a database server

SQLite makes setup deterministic and requires no external service. It is appropriate for this read-only assignment, but a high-write or horizontally scaled system would normally use a server database such as PostgreSQL.

### Schema and seed in `main.go`

Keeping initialization in the supplied entry point avoids introducing migration infrastructure for two tables. The trade-off is that schema history and rollback are not managed; a larger application should use versioned migrations.

### Stateless JWT authentication

JWTs avoid server-side session storage and keep protected requests simple. They cannot be revoked individually in the current design, so production logout/revocation requirements would need token rotation, shorter lifetimes, or a session store.

### OpenAPI-first generation

A single contract keeps backend validation and frontend types aligned. Contract changes require regeneration, and generated code adds repository size, but reviewers can build the project without installing a generator globally.

### Global summary as a second query

The list query applies filters while the summary query always describes the complete dataset, matching the dashboard behavior. This costs one additional query and does not provide a transactional snapshot across both reads; that is acceptable while payments are read-only in this assignment.

### No pagination

The assignment asks for the full payment list, so the required slice remains small and direct. A real high-volume dashboard should add cursor pagination and optionally a date range through an OpenAPI-first change.
