# Frontend

React single-page application for login and incoming-payment monitoring. It uses TypeScript, React Router, Zustand, Tailwind CSS, and an `openapi-fetch` client typed from the root OpenAPI contract.

## Requirements

- Node.js 20+
- npm
- The backend running at `http://localhost:8080` for an end-to-end local session

## Setup and run

From the repository root:

```bash
make setup
make run-frontend
```

Or from `frontend/`:

```bash
npm install
npm run dev
```

Open [http://localhost:5173](http://localhost:5173). Start the backend separately with `make run-backend`, or use `make run` from the root to start both services.

## API configuration

No `.env` file is needed for the default workflow. Vite proxies `/dashboard/v1` to `http://localhost:8080`, while the Docker Nginx server proxies the same path to the backend container.

To target another API origin:

```bash
cp .env.example .env
```

Then set:

```dotenv
VITE_API_BASE_URL=https://api.example.com
```

`VITE_API_BASE_URL` is read at build/start time. The target server must allow the frontend origin when it is not using the same-origin proxy.

## Build and preview

```bash
npm run build
npm run preview
```

The build runs TypeScript checking before Vite and writes static assets to `dist/`. Vite preview inherits the local API proxy. For the production-like Nginx workflow, run `make docker-up` from the repository root.

## Test and quality checks

```bash
npm test
npm run typecheck
npm run generate:api:check
npm run build
```

Use `npm run test:watch` during development. Run `npm run generate:api` after changing the root `openapi.yaml`.

## Architecture

The source is organized by feature so each product area owns its API, page, components, and state:

```text
src/
  app/
    router.tsx              route definitions and auth guard
    AppShell.tsx            authenticated layout and logout
  features/
    auth/
      api/                  typed login request
      pages/                login page and form state
      store/                persisted Zustand auth session
    payments/
      api/                  typed payment request
      components/           summary, filters, table, status/state UI
      pages/                dashboard orchestration and request state
  shared/
    api/                    OpenAPI types, client, and API error
    lib/                    IDR and browser-timezone formatting
  test/                     shared Vitest setup
```

### Application flow

```text
main.tsx
  -> React Router
  -> ProtectedRoute checks Zustand token
  -> AppShell
  -> LoginPage or DashboardPage
  -> feature API function
  -> shared typed openapi-fetch client
  -> Vite/Nginx proxy
  -> Go API
```

### State ownership

- Zustand stores the single persisted auth session: token, email, and role.
- Login form values, submission state, and validation errors stay inside `LoginPage`.
- Payment data, global summary, selected filter, and request status stay inside `DashboardPage`.
- Selecting a payment status sends a new backend request; the frontend does not filter a previously downloaded list.
- A `401` clears the persisted session and redirects to login.

### OpenAPI client

`src/shared/api/openapi.ts` is generated from the root `openapi.yaml`. `openapi-fetch` uses those generated paths and schemas, while a single request middleware reads the latest bearer token from Zustand.

The typed client is the only network boundary. Feature API modules translate unsuccessful responses into `ApiError`, leaving pages responsible for user-facing behavior.

### UI behavior

- The protected dashboard models loading, success, empty, filtered-empty, error/retry, and expired-session states.
- Required payment columns remain available on narrow screens through a keyboard-focusable horizontal table region.
- Amounts use `Intl.NumberFormat` with IDR currency.
- Dates use `Intl.DateTimeFormat`, so the browser's local timezone controls the displayed time.
- Semantic labels, focus rings, `aria-live`, and `aria-pressed` support keyboard and assistive-technology use.

## Testing strategy

- Vitest runs all tests in a JSDOM environment.
- Page/component tests mock feature API modules and assert visible behavior rather than implementation details.
- API-client tests stub global `fetch` and never contact the Go server or external APIs.
- Auth-store tests use isolated browser storage.
- Formatting tests cover IDR and browser-timezone date output.

## Trade-offs

### Zustand only for authentication

The application has one small piece of shared persisted state, so Zustand keeps it explicit and lightweight. Payment request state remains page-local; adding a second server-state library would duplicate ownership for the current scope.

### `openapi-fetch` without React Query

The client provides contract-safe requests with little abstraction. It does not provide caching, background refresh, or request deduplication. Those features are unnecessary for one dashboard request flow but would be worth adding if several screens reused payment data.

### JWT persistence in local storage

Persisting the session survives browser refresh and matches the assignment's client-state requirement. Local storage is accessible to JavaScript, so a production system with stronger security requirements should prefer secure, HTTP-only cookies and corresponding CSRF protection.

### Server-side filtering

Each filter selection performs a backend request, avoiding client-side filtering of a potentially large dataset. The trade-off is additional network activity when switching filters; caching can be introduced if real usage demonstrates a need.

### Feature-based folders

Feature ownership keeps auth and payments independently maintainable and provides a clear place for future feature-specific code. It creates slightly more directory nesting than a flat structure, which is acceptable for preserving boundaries as the dashboard grows.

### Native controls and Tailwind

Avoiding a component library keeps the bundle and visual language focused. The trade-off is that accessibility states and component behavior must be implemented and tested directly.
