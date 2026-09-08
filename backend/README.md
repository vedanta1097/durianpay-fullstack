# Backend

Go API for the internal payment dashboard. It uses SQLite, creates the schema on startup, and safely reuses the same deterministic seed records on subsequent runs.

## Run locally

The easiest way to run the complete application is from the repository root:

```bash
make setup
make run
```

To run only the backend from this directory:

```bash
make setup
make run
```

`make setup` creates `.env` from `env.sample` only when `.env` does not already exist.

The API listens on `http://localhost:8080` by default. The root OpenAPI contract is at `../openapi.yaml`.

For the complete Docker Compose workflow, use the root [README](../README.md).

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
