# Config Package Guide

## Overview

PhotoPrism’s runtime configuration is managed by this package. Fields are defined in [`options.go`](options.go) and then initialized with values from command-line flags, environment variables, and optional YAML files (`storage/config/*.yml`).

## Sources and Precedence

PhotoPrism loads configuration in the following order:

1. **Built-in defaults** defined in this package.
2. **`defaults.yml`** — optional system defaults (typically `/etc/photoprism/defaults.yml`). See [Global Defaults](https://docs.photoprism.app/getting-started/config-files/defaults/) if you package PhotoPrism for other environments and need to override the compiled defaults.
3. **Environment variables** prefixed with `PHOTOPRISM_…` and specified in [flags.go](flags.go) along with the CLI flags. This is the primary override mechanism in container environments.
4. **`options.yml`** — user-level configuration stored under `storage/config/options.yml` (or another directory controlled by `PHOTOPRISM_CONFIG_PATH`). Values here override both defaults and environment variables, see [Config Files](https://docs.photoprism.app/getting-started/config-files/).
5. **CLI flags** (for example `photoprism --cache-path=/tmp/cache`). Flags always win when a conflict exists.

The `PHOTOPRISM_CONFIG_PATH` variable controls where PhotoPrism looks for YAML files (defaults to `storage/config`).

> Any change to configuration (flags, env vars, YAML files) requires a restart. The Go process reads options during startup and does not watch for changes.

## CLI Reference

- `photoprism help` (or `photoprism --help`) lists all subcommands and global flags.
- `photoprism show config` renders every active option along with its current value. Pass `--json`, `--md`, `--tsv`, or `--csv` to change the output format.
- `photoprism show config-options` prints the description and default value for each option. Use this when updating [flags.go](flags.go).
- `photoprism show config-yaml` displays the configuration keys and their expected types in the [same structure that the YAML files use](https://docs.photoprism.app/getting-started/config-files/). It is a read-only helper meant to guide you when editing files under `storage/config`.
- Additional `show` subcommands document search filters, metadata tags, and supported thumbnail sizes; see [`internal/commands/show.go`](../commands/show.go) for the complete list.
