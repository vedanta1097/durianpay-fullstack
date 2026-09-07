# DurianPay Payment Dashboard - Implementation Plan

## 1. Objective

Build a small, polished internal dashboard for monitoring incoming payments. The submission must be easy to run on the reviewer machine, demonstrate clean frontend state and component design, and provide a reliable Go API backed by SQLite.

This plan is based on `docs/assignment.pdf`, the current Go starter, and the root `openapi.yaml`. The PDF is the product-requirement source of truth. `openapi.yaml` will become the executable API contract.

## 2. Assignment Requirements

### Required behavior

- A login form accepts email and password.
- A successful login returns a JWT and one of the supplied roles: `cs` or `operation`.
- The frontend stores the token and role in client state and redirects to a protected dashboard.
- Unauthenticated users cannot open the dashboard.
- `GET /dashboard/v1/payments` requires authentication and returns all payments.
- `GET /dashboard/v1/payments?status=<status>` returns only `completed`, `processing`, or `failed` payments.
- The dashboard shows Payment ID, Merchant Name, Date, Amount, and Status.
- The dashboard shows total payments and a status breakdown.
- Both backend and frontend include unit tests for critical behavior.
- Frontend API types/integration are generated from `openapi.yaml`.
- The repository contains `backend/` and `frontend/`, a root `README.md`, a root `Makefile`, seed instructions, exact build/run/test commands, and clear frontend API configuration.

### Reviewer constraints

- Go 1.21 or newer
- Node.js 20 or newer
- Docker and Docker Compose
- Make
- macOS

All checked-in code and documented commands must work at the minimum declared versions, not only on the author's newer local toolchain. The current starter needs adjustment because `backend/go.mod` declares Go 1.25.5 and uses a Go `tool` directive that is unavailable in Go 1.21.

### Explicit interpretations

- The assignment says the role determines access rights but defines no role-specific action. Both roles may view the same dashboard. The role is retained and displayed, but no speculative RBAC behavior will be invented.
- The assignment's summary example says `Success`, while the API status is `completed`. The UI will use `Completed` consistently and will also show `Processing`.
- The OpenAPI starter already defines optional `id` and `sort` query parameters. The backend will implement and test them, but the first frontend scope only needs the required status filter.
- Pagination and date filtering are not explicit requirements. The PDF says the endpoint "returns the full list of payments," and the supplied OpenAPI has no cursor, page, limit, or date parameters. The phrases "data-rich environments" and "production readiness" show that scale awareness matters, but do not by themselves require changing the contract. The core implementation will therefore remain unpaginated. A cursor/date-range extension is a deliberate stretch decision discussed below, not a hidden requirement.

## 3. Technology Decisions

| Area | Choice | Reason and compatibility |
| --- | --- | --- |
| Backend | Go 1.21, `net/http`, Chi | Keeps the provided architecture and supports the reviewer's minimum Go version. |
| API contract | OpenAPI 3.0.3 and `oapi-codegen` v2.4.1 | v2.4.1 targets Go 1.21. Generated server code is committed. Do not use `latest`. |
| Persistence | SQLite through the provided `github.com/mattn/go-sqlite3` driver | Minimizes changes to the interviewer-provided boilerplate. The Docker build will include its CGO compiler requirements; native macOS setup will document Xcode Command Line Tools. |
| Authentication | HMAC-SHA256 JWT and bcrypt | Matches the starter and assignment without introducing an external identity service. |
| Frontend | React 18 + TypeScript + Vite 5 | Vite 5 supports Node 18/20, so Node 20 is a safe baseline. |
| Client state | Zustand | Small API and enough for the required auth/session state without additional architecture. |
| Styling | Tailwind CSS 3.4 | Compatible with Node 20 and suited to a custom responsive UI. |
| Routing | React Router 6 | Provides explicit login/dashboard routes and a small protected-route boundary. |
| OpenAPI client | `openapi-typescript` + `openapi-fetch` | Generates frontend types from the root contract and keeps request/response usage typed. Generated types are committed. |
| Backend tests | Go `testing` + `httptest` | Standard-library-first and easy to run in a clean environment. |
| Frontend tests | Vitest + React Testing Library + jsdom | Exercises components and state using the Vite toolchain. |
| API mocking | Vitest `vi.mock` and `vi.stubGlobal` | Keeps the small test suite dependency-free: feature tests mock the typed API client and API-client tests mock `fetch`, with no real backend request. |
| Date formatting | Native `Intl.DateTimeFormat` initially | Covers the required display without another dependency. Day.js is allowed only if later date-range behavior makes parsing or manipulation materially clearer. |
| Reproducibility | npm lockfile, pinned Go/codegen versions, Make, Docker Compose | Avoids unbounded `latest` installs and provides native and containerized paths. |

Do not add a component framework, ORM, server-state library, Redis, or a monorepo tool. They do not solve a requirement in this assignment.

## 4. Target Repository Shape

```text
.
├── AGENTS.md
├── Makefile
├── README.md
├── compose.yaml
├── openapi.yaml
├── backend/
│   ├── Dockerfile
│   ├── Makefile
│   ├── main.go                     # keep the provided executable entry point
│   ├── internal/
│   │   ├── api/
│   │   ├── config/
│   │   ├── entity/                 # add payment.go beside existing user/error entities
│   │   ├── module/
│   │   │   ├── auth/               # preserve handler/repository/usecase structure
│   │   │   └── payment/            # mirror the existing auth module structure
│   │   ├── openapigen/
│   │   ├── service/http/
│   │   └── transport/
│   └── script/gen-secret/
├── frontend/
│   ├── Dockerfile
│   ├── package.json
│   ├── package-lock.json
│   ├── src/
│   │   ├── app/                    # app root, router, and cross-feature composition
│   │   ├── features/
│   │   │   ├── auth/
│   │   │   │   ├── api/
│   │   │   │   ├── components/
│   │   │   │   ├── pages/
│   │   │   │   └── store/
│   │   │   └── payments/
│   │   │       ├── api/
│   │   │       ├── components/
│   │   │       └── pages/
│   │   ├── shared/                 # only genuinely cross-feature code
│   │   │   ├── api/                # typed client and generated OpenAPI types
│   │   │   ├── components/
│   │   │   └── lib/
│   │   ├── test/
│   │   └── main.tsx
│   └── README.md
└── docs/
    ├── assignment.pdf
    └── implementation-plan.md
```

The backend structure intentionally follows the starter. Do not introduce `cmd/server`, move `main.go`, create a separate database package, or rename existing auth packages. This application has one small executable, so `backend/main.go` is already sufficient. Extend the existing `entity` and `module` conventions with only the payment files needed for the assignment.

The frontend is feature-first: each feature owns its page, API calls, components, and feature state. `app/` is limited to application composition and routing; `shared/` is limited to code used by more than one feature. There is no global `pages/` or `stores/` directory because those would split related feature code by technical type.

## 5. API and Data Contract

### OpenAPI cleanup

Update `openapi.yaml` before implementing payment code, then regenerate both consumers.

- Give the API and operations meaningful names and `operationId` values.
- Mark required response fields as required instead of generating unnecessary pointers.
- Define `UserRole` as `cs | operation`.
- Define `PaymentStatus` as `completed | processing | failed`.
- Define amount as an integer (`int64`) in IDR's smallest used unit; never use floating point for money.
- Define `created_at` as RFC 3339 `date-time` in UTC.
- Make the error response match the implementation: stable string code plus human-readable message.
- Document success, invalid request, unauthenticated, and internal-error responses.
- Keep the existing endpoint paths. Do not add a review endpoint unless the assignment or contract is deliberately expanded later.

### Proposed response shapes

```json
{
  "email": "cs@test.com",
  "role": "cs",
  "token": "<jwt>"
}
```

```json
{
  "payments": [
    {
      "id": "PAY-0001",
      "merchant": "Kopi Nusantara",
      "status": "completed",
      "amount": 125000,
      "created_at": "2026-03-10T08:15:00Z"
    }
  ]
}
```

### Validation and status behavior

- Login requires a syntactically valid email and a non-empty password.
- Invalid JSON or query values return `400` with the documented error shape.
- Missing, malformed, expired, or incorrectly signed JWTs return `401`.
- Unknown credentials return a generic `401 invalid credentials`; do not reveal whether an email exists.
- Unsupported status or sort fields return `400` rather than being silently ignored.
- Successful JSON responses set `Content-Type: application/json`.
- Database/internal details are logged server-side and are never returned to the client.

## 6. Backend Design

### Database and seed

Extend the existing `initDB` flow in `backend/main.go` with `payments` setup rather than moving current code merely for architectural symmetry. Create `users` and `payments` tables using idempotent `CREATE TABLE IF NOT EXISTS` statements. Add an index on payment status and an index suitable for the default newest-first ordering.

Seed deterministic demo data only when the relevant table is empty:

- `cs@test.com` / `password`, role `cs`
- `operation@test.com` / `password`, role `operation`
- Approximately 30 payments across all three statuses, multiple merchants, varied dates, and varied amounts

Run seed inserts in a transaction. Store only bcrypt password hashes. Keep the SQLite file at a configurable path such as `DB_PATH=./data/dashboard.db`, and ignore the generated database file.

Because the workload is read-heavy and has no payment mutations, keep SQLite handling direct: foreign keys enabled, request-context queries, and a conservative connection limit. Do not add a cache, migration framework, database wrapper, or background synchronization.

### Request flow through the provided Go structure

`backend/main.go` constructs the repositories, use cases, handlers, `APIHandler`, and HTTP server once at startup. `backend/internal/openapigen/openapi.gen.go` is generated routing/type plumbing; it should not contain handwritten business logic.

Login execution:

```text
POST /dashboard/v1/auth/login
  -> internal/service/http/server.go
     Chi router + OpenAPI request validation
  -> internal/openapigen/openapi.gen.go
     generated route dispatch
  -> internal/api/api_handler.go
     PostDashboardV1AuthLogin
  -> internal/module/auth/handler/auth.go
     decode HTTP JSON and map response/error
  -> internal/module/auth/usecase/auth.go
     verify credentials, create JWT
  -> internal/module/auth/repository/user.go
     parameterized SELECT
  -> SQLite users table
  <- repository <- use case <- handler <- JSON response
```

Payment-list execution after implementation:

```text
GET /dashboard/v1/payments?status=completed
  -> internal/service/http/server.go
     Chi router + OpenAPI validation + JWT authentication hook
  -> internal/openapigen/openapi.gen.go
     generated route dispatch and typed query params
  -> internal/api/api_handler.go
     GetDashboardV1Payments delegates to payment handler
  -> internal/module/payment/handler/payment.go
     map HTTP/auth/query data to use-case input
  -> internal/module/payment/usecase/payment.go
     validate supported filters and sort
  -> internal/module/payment/repository/payment.go
     parameterized SELECT
  -> SQLite payments table
  <- repository <- use case <- handler <- JSON response
```

The auth middleware validates signature, algorithm, expiry, subject, and role, then attaches a small typed identity to the request context. It must not accept a token merely because the `Authorization` header has the right shape.

### Payments query

- Default order: newest first, with ID as a deterministic tie-breaker.
- Optional exact status filter.
- Optional exact ID filter from the supplied contract.
- Optional allow-listed sort values from the supplied contract; never concatenate arbitrary input into SQL.
- Return an empty array, not `null`, when no records match.

### Pagination and date-range decision

The real-world concern is valid: an unbounded payment feed eventually becomes expensive in database work, payload size, browser memory, and rendering time. A production version would normally combine a bounded date range with cursor pagination. Cursor pagination is preferable to a large page index because it avoids increasingly expensive offsets and is more stable while new payments are arriving.

For this take-home, pagination is deferred by default because adding it correctly is not only four buttons in the UI. It requires changing OpenAPI request/response types, defining stable cursor ordering, deciding summary semantics across all matching records, adding indexes, and expanding backend/frontend tests. Those changes conflict with the supplied "full list" contract and the instruction to keep boilerplate changes small.

If the core assignment is complete early, treat the following as one optional stretch slice and update OpenAPI first:

- `from` and `to` timestamps, with a documented inclusive/exclusive boundary.
- `limit` with a small default and hard maximum.
- An opaque `cursor` based on `(created_at, id)` and `next_cursor` in the response.
- A summary that clearly represents all records matching the filters, not only the current page.
- A UI date range and next/previous navigation matching the reference's restrained controls.

Do not implement offset/page-number pagination simply to make the UI look complete. Do not begin this stretch slice until the required implementation, tests, documentation, and clean reviewer setup all pass.

## 7. Frontend Design

### Routes and session

- `/login`: public login page; an already-authenticated user is redirected to `/dashboard`.
- `/dashboard`: protected dashboard; unauthenticated users are redirected to `/login`.
- Unknown paths redirect based on authentication state.
- A logout action clears the entire session and returns to login.
- Any API `401` clears stale auth state and redirects to login.

Use one Zustand auth store containing only `token`, `role`, `email`, `isAuthenticated` as a derived selector, and `setSession`/`clearSession` actions. Persist the minimal session in `localStorage` so refreshes remain authenticated. This storage follows the assignment's bearer-token contract; the README should acknowledge that an HttpOnly cookie would be preferred for a production system designed against an XSS threat model.

Do not duplicate the token in component state, API-client state, and the store. The API client reads it from the auth store when constructing a request.

### Payment state

Fetch the full seeded payment list once per dashboard visit. Keep the fetched array, loading state, and request error together in the dashboard feature. Keep the selected status filter in one local state value and derive the visible rows and summary with selectors/memoized calculations.

The backend status-filter endpoint is still implemented and tested as part of the API contract. Client-side filtering is intentional here because the contract returns a full, small list and the same canonical array can drive both the stable overall summary and visible rows without a second request or duplicated server state.

### Component boundaries

- `AppShell`: compact internal-tool header, signed-in identity/role, and logout.
- `ProtectedRoute`: authentication gate only.
- `LoginForm`: fields, accessible validation, submit state, and server error.
- `DashboardPage`: owns fetch lifecycle and composes the page.
- `PaymentSummary`: total plus completed/processing/failed counts.
- `PaymentStatusFilter`: all and the three contract statuses.
- `PaymentTable`: semantic table and row rendering.
- `StatusBadge`: the single visual mapping for status label/color.
- `LoadingState`, `EmptyState`, and `ErrorState`: explicit request states, with retry on errors.

Do not create a generic component abstraction until at least two concrete consumers need the same behavior.

### UI/UX direction

- Use the user-provided payout dashboard screenshot as visual direction only, not as a source of product requirements. Do not copy unrelated payout fields/actions such as batch name, sync, export, receipt, bank, or expandable rows.
- Aim for a calm, high-signal operations dashboard rather than a marketing page: white canvas, generous but deliberate spacing, thin neutral dividers, light gray table header, bold key values, and compact controls.
- Use `#7a2cdd` as the primary interactive color with white text for primary buttons. Use a very light purple treatment for the active filter and pale semantic backgrounds for status badges.
- Avoid gradients, oversized marketing headlines, excessive rounded cards, decorative charts, glass effects, fake metrics, and unnecessary icon buttons. The payment table should remain the visual center.
- Use clear visual hierarchy, restrained color, tabular numerals for amounts, and distinct but accessible status treatments.
- Format dates with `Intl.DateTimeFormat` and amounts with `Intl.NumberFormat` for `id-ID` / `IDR`.
- Day.js may be added later if date-range parsing/manipulation becomes non-trivial; it is not needed for initial formatting.
- Never communicate status with color alone; always render the label.
- Provide visible focus states, explicit form labels, useful validation text, and an `aria-live` region for async login errors.
- Disable only the action currently submitting; preserve entered email when login fails.
- On narrow screens, keep summary cards readable and place the semantic table in a clearly signposted horizontal scroll region. Do not hide required columns.
- Include loading, empty, filtered-empty, offline/network-error, invalid-session, and retry states.

## 8. Testing Strategy

### Backend

Use table-driven unit tests with mocks at architectural boundaries. Tests must never contact an external service, start the real application server, or read/write the developer's persistent `dashboard.db`.

- Auth: valid login, unknown email, wrong password, token expiry/signature validation.
- Payment use case/repository: full list, each status, invalid status, ID filter, allow-listed sorting, empty result.
- HTTP: malformed login body, expected JSON/status codes, missing/invalid bearer token, protected happy path.
- Handler tests use `httptest` in process and mocked use cases. `httptest` is not a real network API.
- Use-case tests mock repository interfaces; they do not open SQLite.
- Repository SQL tests may use an isolated in-memory/temporary SQLite database because mocking SQL alone cannot prove SQLite schema/query behavior. This database is created per test and destroyed immediately; it must never be the application database. If repository behavior is fully covered elsewhere, keep this layer small.
- Seed integration test: run initialization twice against isolated temporary SQLite and verify no duplicates.
- OpenAPI: generated handler compiles and representative responses conform to the contract.

### Frontend

Use Vitest-native mocks only. Frontend page/component tests mock the typed API client module; API-client tests use `vi.stubGlobal` to mock `fetch`. Frontend tests must never contact the real Go server or any external API.

- Auth store: set, persist/rehydrate, and clear session.
- Login page: validation, pending state, failed credentials, successful redirect.
- Protected route: unauthenticated redirect and authenticated render.
- Dashboard: loading, successful table/summary, filter behavior, empty/error/retry, and `401` logout.
- Formatting/status components: IDR, date, and accessible status labels.
- Typed API client: request method, URL/query, JSON body, bearer header, successful response, non-OK response, malformed network response, and `401` handling through mocked `fetch`.

### Required verification commands

The completed project should expose these stable root commands:

```bash
make setup
make generate-check
make format-check
make lint
make test
make build
docker compose up --build
```

`generate-check` must fail if committed generated API code differs from `openapi.yaml`. Final verification should also follow the root README from a clean checkout with no existing database or dependency directories.

## 9. Development Phases

### Phase 0 - Normalize the starter and contract

1. Change the backend module baseline to Go 1.21.
2. Remove dependencies and directives that require a newer Go toolchain; pin `oapi-codegen` v2.4.1.
3. Keep the provided SQLite driver, pin a Go 1.21-compatible version, and make its CGO prerequisites reproducible in Docker and explicit for native macOS setup.
4. Correct and complete `openapi.yaml`, then regenerate and commit backend API code.
5. Add root Make targets and ignore rules for generated runtime files.

Exit criteria: a minimal backend build succeeds with Go 1.21, generated code matches the OpenAPI file, and no command installs `latest`.

### Phase 1 - Complete the backend vertical slice

1. Extend the existing database initialization/seed logic in place; extract only if the code becomes genuinely difficult to understand or test.
2. Add the payment entity and `internal/module/payment/{handler,repository,usecase}` alongside the existing auth pattern.
3. Implement JWT authentication middleware and correct error-to-HTTP mapping.
4. Add JSON headers, request-context database calls, configuration, and graceful shutdown fixes.
5. Add backend tests for auth, filtering, authorization, seed idempotency, and errors.

Exit criteria: login followed by an authenticated payment request works; all backend tests pass; invalid credentials and tokens return contract-compliant `401` responses.

### Phase 2 - Scaffold the typed frontend

1. Create `frontend/` with React 18, TypeScript, Vite 5, Tailwind 3.4, and an npm lockfile.
2. Generate and commit frontend API types from the root OpenAPI file.
3. Build the typed API wrapper and the minimal Zustand auth store.
4. Add routing and protected-route behavior.

Exit criteria: Node 20 can install and build the frontend, API calls are typed from OpenAPI, and route/auth-store tests pass.

### Phase 3 - Build login and dashboard UX

1. Implement the polished login flow and errors.
2. Implement dashboard shell, summary, status filter, semantic table, and status badges.
3. Add loading, empty, error, retry, expired-session, and responsive states.
4. Add Vitest-mocked frontend component/integration tests and accessibility assertions without a running backend.

Exit criteria: both seeded users can log in; refresh preserves the session; the dashboard correctly summarizes, filters, formats, and renders seed data on desktop and mobile widths.

### Phase 4 - Integration and reviewer workflow

1. Add Dockerfiles and `compose.yaml` with persistent SQLite storage and explicit health checks.
2. Finish root/backend/frontend Make targets and environment examples.
3. Write the root README with prerequisites, setup, seed credentials, API configuration, exact native/container commands, testing strategy, and OpenAPI instructions.
4. Add a simple Swagger UI or clearly document how to inspect `openapi.yaml`; do not make API docs a runtime dependency for core behavior.

Exit criteria: a clean checkout succeeds by following only the README, both native and Docker paths work, and service readiness is deterministic.

### Phase 5 - Final quality pass

1. Run format, lint, unit/integration tests, generated-code drift check, and production builds.
2. Exercise the complete login-to-dashboard path and failure states manually.
3. Check keyboard navigation, focus visibility, contrast, narrow-screen table behavior, and console/network errors.
4. Recreate the database from seed and repeat the documented clean setup.
5. Remove placeholder text, stale docs, unused dependencies, debug output, and unrequested features.

Exit criteria: every deliverable and evaluation criterion from the assignment has evidence in code, tests, or README, and the reviewer can start the project without undocumented local state.

## 10. Definition of Done

- All assignment requirements are represented in a traceable implementation or test.
- Root `openapi.yaml` is accurate, generated backend/frontend code is committed, and drift checks pass.
- Backend is compatible with Go 1.21 and frontend is compatible with Node 20.
- Login and payment endpoints validate input, enforce JWT authentication, and return consistent errors.
- SQLite schema and deterministic seed are idempotent; no database or secret is committed.
- UI is polished, responsive, keyboard usable, and handles loading/error/empty/session-expiry states.
- Backend and frontend tests, lint, format checks, and production builds pass.
- README commands have been verified from a clean checkout.
- Docker Compose starts a usable end-to-end application.

## 11. Intentionally Out of Scope

- Payment mutation/review workflows
- Pagination, infinite scrolling, or date filtering in the required core scope; cursor/date range remains an explicitly gated stretch slice
- User management, password reset, or role administration
- Different dashboards for `cs` and `operation` without a stated access rule
- Redis, background jobs, websockets, analytics, caching, or audit logs
- A reusable design system beyond components actually used by this dashboard

If a later requirement adds one of these behaviors, update `openapi.yaml`, this plan, tests, and README before introducing the associated complexity.
