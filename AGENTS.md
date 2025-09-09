# Repository Guidelines

## Purpose

This file tells automated coding agents (and humans) where to find the single sources of truth for building, testing, and contributing to PhotoPrism. Keep it short to prevent drift.

## Sources of Truth

- Makefile targets: https://github.com/photoprism/photoprism/blob/develop/Makefile
- Developer Guide (setup): https://docs.photoprism.app/developer-guide/setup/
- Developer Guide (tests): https://docs.photoprism.app/developer-guide/tests/
- CONTRIBUTING: https://github.com/photoprism/photoprism/blob/develop/CONTRIBUTING.md
- SECURITY: https://github.com/photoprism/photoprism/blob/develop/SECURITY.md
- REST API docs: https://docs.photoprism.dev/

## Build & run (local)

- Run `make help` to see all targets.
- Starting the Development Environment with Docker/Compose:
  1. Run `make docker-build` once to build a local image based on the `Dockerfile`.
  2. Then, run `docker compose up -d` to start the Development Environment, if it is not already running.
  3. Open a terminal session with `make terminal` once the services have been started.
- From within the Development Environment:
  - Install deps: `make dep`
  - Build frontend/backend: `make build-js` and `make build-go`
  - Start PhotoPrism server on port 2342: `./photoprism start`

## Tests

- From within the Development Environment:
  - Full unit test suite: `make test` (runs backend and frontend tests)
  - Test frontend/backend: `make test-js` and `make test-go`
- Frontend unit test scripts are defined in `frontend/package.json`
- Acceptance tests: use the `acceptance-*` targets in the `Makefile`

## Code Style & Lint

- Go: `go fmt` / `goimports` (see `Makefile` targets such as `fmt` / `fmt-go`)
- JS/Vue: use the lint/format scripts in `frontend/package.json` (ESLint + Prettier)

## Safety & Data

- Never commit secrets. Use environment variables or a local `.env`.
- Do not run destructive commands against production data. Prefer ephemeral volumes and test fixtures when running acceptance tests.

If anything in this file conflicts with the `Makefile` or the Developer Guide, the Makefile and docs win.
