# Backend

Go API for the internal payment dashboard. It uses SQLite, creates the schema on startup, and safely reuses the same deterministic seed records on subsequent runs.

## Run locally

```bash
cp env.sample .env
make setup
make run
```

The API listens on `http://localhost:8080` by default. The root OpenAPI contract is at `../openapi.yaml`.

Seed users:

- `cs@test.com` / `password`
- `operation@test.com` / `password`

Useful commands:

```bash
make openapi-gen
make openapi-check
make format-check
make lint
make test
make build
make gen-secret
```

`openapi-gen` uses the repository-pinned generator version; it never installs a global or `latest` version.
