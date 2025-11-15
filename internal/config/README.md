# Config Package Guide

## Overview

PhotoPrism’s runtime configuration is managed by this package. Defaults are defined in [`options.go`](options.go) and then merged with values from command-line flags, environment variables, and optional YAML files (`storage/config/*.yml`).

## Sources and Precedence

PhotoPrism loads configuration in the following order:

1. **Built-in defaults** defined in [`internal/config/options.go`](https://github.com/photoprism/photoprism/blob/develop/internal/config/options.go).
2. **`defaults.yml`** — optional system defaults (typically `/etc/photoprism/defaults.yml`). See [Global Defaults](https://docs.photoprism.app/getting-started/config-files/defaults/) if you package PhotoPrism for other environments and need to override the compiled defaults.
3. **Environment variables** prefixed with `PHOTOPRISM_…`. The list mirrors [Config Options](https://docs.photoprism.app/getting-started/config-options/). In Docker/Compose this is the primary override mechanism.
4. **`options.yml`** — user-level configuration stored under `storage/config/options.yml` (or another directory controlled by `PHOTOPRISM_CONFIG_PATH`). Values here override both defaults and environment variables, see [Config Files](https://docs.photoprism.app/getting-started/config-files/).
5. **Command-line flags** (for example `photoprism --cache-path=/tmp/cache`). Flags always win when a conflict exists.

The `PHOTOPRISM_CONFIG_PATH` variable controls where PhotoPrism looks for YAML files (defaults to `storage/config`).

> Any change to configuration (flags, env vars, YAML files) requires a restart. The Go process reads options during startup and does not watch for changes.

## CLI Reference

- `photoprism help` (or `photoprism --help`) lists all subcommands and global flags.
- `photoprism show config` renders every active option along with its current value. Pass `--json`, `--md`, `--tsv`, or `--csv` to change the output format.
- `photoprism show config-options` prints the description and default value for each option. Use this when updating [Config Options](https://docs.photoprism.app/getting-started/config-options/).
- `photoprism show config-yaml` displays the configuration keys and their expected types in the [same structure that the YAML files use](https://docs.photoprism.app/getting-started/config-files/). It is a read-only helper meant to guide you when editing files under `storage/config`.
- Additional `show` subcommands document search filters, metadata tags, and supported thumbnail sizes; see [`internal/commands/show.go`](../commands/show.go) for the complete list.
