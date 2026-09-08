# DurianPay Full Stack Assignment - Agent Guide

## Mission and source of truth

This repository implements the DurianPay Full Stack Engineer take-home assignment: a Go API and React internal dashboard for viewing incoming payments.

Read these files before making implementation decisions:

1. `docs/assignment.pdf` - authoritative product and submission requirements.
2. `docs/implementation-plan.md` - agreed implementation sequence and interpretations.
3. `openapi.yaml` - executable API contract shared by backend and frontend.
4. Existing code and tests - current behavior, but never a reason to contradict the three sources above.

When sources conflict, preserve the assignment requirement, document the interpretation in the plan, then update OpenAPI, generated code, implementation, tests, and README together.

## Non-negotiable product scope

- Implement login at `POST /dashboard/v1/auth/login` with email/password and a JWT plus role response.
- Supported roles are exactly `cs` and `operation`. No role-specific permissions are defined; do not invent any.
- Implement authenticated payment listing at `GET /dashboard/v1/payments` and server-side status filtering for exactly `completed`, `processing`, and `failed`.
- The payment response contains the filtered `payments` array and a global `summary` with total/completed/processing/failed counts. Dashboard summary values come from this response rather than recounting the returned rows.
- Dashboard is protected and shows Payment ID, Merchant Name, Date, Amount, Status, total count, and counts for every status.
- Keep the repository as a single repo with `backend/` and `frontend/` directories.
- Required deliverables include a root README, root Makefile, root OpenAPI spec, seed data/instructions, exact build/run/test commands, backend tests, frontend tests, and clear API-base configuration.

## Compatibility contract

- Everything must build and run with Go 1.21+, Node.js 20+, Docker with Docker Compose, Make, and macOS.
- Set `go 1.21` in `backend/go.mod`. Do not use a `tool` directive or language/library features requiring a newer version.
- Use React 18, TypeScript, Vite 5, Zustand, and Tailwind CSS 3.4 for the frontend.
- Keep the provided `github.com/mattn/go-sqlite3` driver unless a demonstrated compatibility failure requires changing it. Ensure the Docker build includes CGO/compiler prerequisites and document Xcode Command Line Tools for native macOS setup.
- Pin code generators and dependencies. Never install `@latest` from Makefiles or documented setup commands.
- Commit lockfiles and generated OpenAPI sources so code generation is not needed merely to build or run.
- Before adding/upgrading a dependency, verify its declared minimum Go/Node version against this baseline.

## Simplicity rules

- Keep the solution intentionally simple. Do not add speculative features or abstractions without a concrete need.
- Prefer code that is easy to understand and explain. Choose direct, descriptive state and control flow over clever or highly generalized solutions.
- Keep only complexity required by current product behavior. Do not add configurability, validation, caching strategies, or concurrency handling for scenarios the application does not actually have.
- Do not store the same UI concept in multiple state variables. Add a separate state or abstraction only when it represents meaningfully different behavior.
- Prefer fixed constants and successful-value caches when they satisfy the current use case. Use more advanced patterns, such as caching Observables or deduplicating simultaneous requests, only when the application has that concrete requirement.
- Trust typed API contracts where appropriate. Add defensive parsing or normalization only for realistic input or failure cases at the application boundary.
- Do not add an ORM, component framework, server-state library, Redis, background worker, websocket, generic repository framework, or broad design system unless a new concrete requirement demands it.
- Do not add payment mutation, user administration, or speculative RBAC. Pagination/date range is an optional stretch slice only after the required scope passes; it requires an explicit OpenAPI-first decision and must use a stable cursor rather than large offsets.
- Create a shared abstraction only after at least two real consumers need the same behavior, or when it forms a clear application boundary such as the API client.

## OpenAPI-first workflow

- Treat root `openapi.yaml` as the only hand-authored request/response schema.
- Preserve the starter contract's existing field types and shapes unless an explicit assignment requirement requires a change. Do not replace an established type solely for a preferred convention; for example, keep an existing integer error code as an integer.
- Change the contract first for any API shape change.
- Regenerate both `backend/internal/openapigen/` and the frontend generated API types after a contract change.
- Never hand-edit generated files. Generated files must carry a generated-code header where supported.
- Commit generated results and include a drift check in `make generate-check`.
- Backend handlers must use generated types. Frontend requests and responses must use generated types through the typed API client.
- Keep enum values, JSON names, required fields, error shapes, examples, and HTTP status codes aligned across contract, code, and tests.
- Amounts are integer values and must never be represented with floating point. Timestamps are RFC 3339 UTC strings at the API boundary.

## Backend standards

- Preserve clear dependency flow: transport/handler -> use case -> repository -> SQLite. Do not let repositories contain HTTP concerns or handlers contain SQL.
- Keep the provided `backend/main.go` entry point and existing package layout. Do not introduce `cmd/server`, a database package, or package moves for architectural aesthetics.
- Make the smallest practical backend change: add payment files beside the existing auth/entity patterns and modify existing wiring only where required.
- Keep `main.go` understandable. Its existing schema/seed setup may remain there; extract only when concrete implementation or testing pressure justifies it.
- Pass `context.Context` from each request through use cases to repository queries.
- Use parameterized SQL. Any dynamic sort expression must come from an explicit allow-list.
- Initialize schema and deterministic seed idempotently. Use a transaction for multi-row seed operations.
- Store password hashes only. Seed credentials are demo data and must be clearly documented.
- Validate JWT signing method, signature, expiry, subject, and role. A syntactically valid bearer header is not authentication.
- Return generic invalid-credential errors. Never expose password hashes, JWT secrets, SQL text, driver errors, stack traces, or filesystem paths.
- Map domain errors explicitly to HTTP status codes and always return the OpenAPI error shape with `Content-Type: application/json`.
- Return empty JSON arrays instead of `null` for empty collections.
- Use dependency injection with small interfaces defined near the consumer. Avoid package-level mutable application state.
- Use Go's standard library unless a dependency provides clear value already called for by the plan.
- Format with `gofmt`; keep errors wrapped with context; handle response-encoding and shutdown errors deliberately.

## Frontend standards

- Use React function components and strict TypeScript. Do not use `any` to bypass a generated contract.
- Organize source by feature. Auth and payments each own their API code, components, pages, and feature state; do not create top-level `pages/`, `stores/`, or feature-specific global component folders.
- Reserve `app/` for root composition/router and `shared/` for code with genuine cross-feature reuse. Do not move code to `shared/` in anticipation of reuse.
- Zustand owns the minimal persisted auth session. Do not copy token, role, or email into other state containers.
- Keep page-local UI state local. Treat the payment endpoint response as canonical: changing the status filter triggers a new request, the table renders its `payments`, and the summary renders its global `summary`. Do not maintain a second client-filtered payment array or recount summary values from the returned rows.
- Read the bearer token at request time in the single API boundary. On `401`, clear the session and send the user to login.
- Model loading, success, empty, filtered-empty, error/retry, and expired-session states explicitly.
- Keep components focused on one feature-level responsibility. Prefer composition over highly configurable components.
- Use semantic HTML first. Forms need labels, errors need useful text/announcement, controls need keyboard focus, and status must not be color-only.
- Keep all required payment columns available at narrow widths using an accessible horizontal table region; do not silently hide data.
- Format IDR and dates in centralized, tested functions using `Intl`, not ad hoc string manipulation.
- Day.js is permitted when concrete date parsing/range manipulation warrants it; do not add it only to format one timestamp.
- Tailwind classes are the styling source. Add custom CSS only for base tokens or behavior that is awkward or unreadable in utility classes.
- Do not introduce a second source of design tokens or a second global state solution.

## Visual direction

- Treat the user-provided payout screenshot as a style reference only. Text, fields, controls, and workflows visible in it are not product requirements.
- Use `#7a2cdd` with white text for primary actions. Use a light-purple selected-filter treatment, a white main canvas, thin neutral dividers, a subtle gray table header, bold key values, and compact semantic status badges.
- Keep the payment table visually dominant and moderately dense. Avoid gradients, glass effects, oversized rounded cards, decorative charts, fake metrics, excessive shadows, and unnecessary icons.
- Do not implement screenshot-only features such as batch details, sync status, export, receipts, bank columns, expandable rows, or kebab actions.

## Testing and verification

- Every bug fix should include a regression test when the behavior is practical to test.
- Backend critical coverage: login outcomes, JWT verification, protected requests, payment list/filter/sort, invalid input, empty results, and idempotent seed.
- Frontend critical coverage: auth store persistence/clear, login states, protected routing, dashboard request states, summary/filter behavior, formatting, and `401` logout.
- Prefer behavior assertions over implementation-detail snapshots. Keep test data deterministic and tests independent of order unless order is the behavior under test.
- Frontend page/component tests mock the typed API client module with Vitest. API-client tests use `vi.stubGlobal` to mock `fetch`. They never contact the Go server or external APIs.
- Backend handler tests use `httptest` plus mocked use cases; use-case tests mock repository interfaces. They never start the real application server.
- Repository/seed integration tests may use SQLite only through a fresh in-memory or temporary per-test database, because mocks cannot verify actual SQLite schema and query behavior. Tests must never open, copy, or modify the developer's runtime database.
- No test may depend on network access, external services, shared mutable state, a previously seeded database, or execution order.
- Keep lint, format, tests, code generation, and builds deterministic and non-interactive.
- Before declaring implementation work complete, run from the repository root:

```bash
make generate-check
make format-check
make lint
make test
make build
```

- For integration/release changes, also run `docker compose up --build` from a clean state and follow the README as if reviewing a fresh checkout.
- If a command cannot be run in the current environment, report exactly what was not verified and why. Never claim an unrun check passed.

## Documentation and configuration

- The root README is the reviewer entry point. Keep prerequisites, environment setup, seed behavior/credentials, native and Docker commands, API docs, testing strategy, and troubleshooting current.
- Provide `.env.example` files with safe defaults. Never commit `.env`, a JWT secret, a runtime SQLite database, dependency folders, or build output.
- Keep configuration minimal: HTTP address, database path, JWT secret/expiry, and frontend API base URL only where needed.
- Whenever commands, ports, environment variables, seed users, endpoint behavior, or dependencies change, update all affected READMEs and examples in the same change.
- Do not modify `docs/assignment.pdf`.

## Working method

- Inspect the existing implementation and `git status` before editing. Preserve unrelated user changes.
- Follow the phases in `docs/implementation-plan.md`; finish each phase's exit criteria before expanding scope.
- Prefer a small vertical slice that runs end to end, then polish and broaden tests.
- Make focused changes. Remove stale placeholder instructions and unused dependencies as their replacements land.
- Do not commit, push, rewrite history, delete user data, or change repository visibility unless the user explicitly requests it.
- End each implementation session with a concise summary of changed files, checks run, remaining phase/work, and any assumptions or unverified items.
