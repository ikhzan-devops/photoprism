# Repository Guidelines

## Purpose

This file tells automated coding agents (and humans) where to find the single sources of truth for building, testing, and contributing to PhotoPrism. Keep it short to prevent drift.

## Sources of Truth

- Makefile targets: https://github.com/photoprism/photoprism/blob/develop/Makefile
- Developer Guide (setup): https://docs.photoprism.app/developer-guide/setup/
- Developer Guide (tests): https://docs.photoprism.app/developer-guide/tests/
- CONTRIBUTING: https://github.com/photoprism/photoprism/blob/develop/CONTRIBUTING.md
- SECURITY: https://github.com/photoprism/photoprism/blob/develop/SECURITY.md
- REST API docs:
  - https://docs.photoprism.dev/ (Swagger)
  - https://docs.photoprism.app/developer-guide/api/

## Build & run (local)

- Run `make help` to see common targets (or open the `Makefile`).
- Start the Development Environment with Docker/Compose:
  - Run `make docker-build` once to build a local image based on the `Dockerfile`.
  - Then, run `docker compose up` to start the Development Environment, if it is not already running (add `-d` to start in the background)
  - Execute a single command: `docker compose exec photoprism <command>`, e.g. `docker compose exec photoprism ./photoprism help`
    - Why `./photoprism`? It runs the locally built binary in the project directory (as used in the setup guide).
    - Run as non-root to avoid root-owned files on bind mounts: `docker compose exec -u "$(id -u):$(id -g)" photoprism <command>`
    - Durable alternative: set the service user or `PHOTOPRISM_UID` in `compose.yaml`; if you hit issues, run `make fix-permissions`.
  - Open a terminal session: `make terminal`
- From within the Development Environment:
  - Install deps: `make dep`
  - Build frontend/backend: `make build-js` and `make build-go`
  - Start PhotoPrism server on port 2342: `./photoprism start` (open http://localhost:2342/)
- Stop the Development Environment with Docker/Compose: `docker compose --profile=all down --remove-orphans` (`--profile=all` ensures all services are stopped; `make down` does the same)

## Tests

- From within the Development Environment:
  - Full unit test suite: `make test` (runs backend and frontend tests)
  - Test frontend/backend: `make test-js` and `make test-go`
- Frontend unit tests are driven by Vitest; see scripts in `frontend/package.json`
  - Vitest watch/coverage: `make vitest-watch` and `make vitest-coverage`
- Acceptance tests: use the `acceptance-*` targets in the `Makefile`

## Code Style & Lint

- Go: `go fmt` / `goimports` (see `Makefile` targets such as `fmt` / `fmt-go`)
- JS/Vue: use the lint/format scripts in `frontend/package.json` (ESLint + Prettier)

## Safety & Data

- Never commit secrets. Use environment variables or a local `.env`.
- Do not run destructive commands against production data. Prefer ephemeral volumes and test fixtures when running acceptance tests.

If anything in this file conflicts with the `Makefile` or the Developer Guide, the `Makefile` and docs win.
