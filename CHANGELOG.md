# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Breaking (implementation):** Reimplemented the workflow in Go
  (`cmd/password-generator-alfred` + `internal/passgen`, `internal/passgencmd`,
  `internal/scriptfilter`), replacing the Python `src/alfred`/`src/app` implementation.
  The `passgen` keyword, bundle ID, and the behavior of `passgen`/`panc`/`split`/`help`
  are unchanged; results are byte-for-byte equivalent JSON.
- Alfred now invokes a compiled universal (amd64+arm64) binary directly instead of
  `python3`/`uv run python`; the `Use uv` Config Builder toggle is removed.
- Build/test tooling moved from `uv`/`ruff`/`mypy`/`pytest` to `go build`/`gofmt`/`go vet`/`go test`.

### Removed

- The `config` subcommand and its persistent settings store (`alfred_workflow_data`).
  It was never written to by any `passgen` command and was undocumented; the workflow's
  only real setting (clipboard history) remains available via Alfred's Config Builder
  (`⌘,`).
- The `Log Level` Config Builder toggle and file logging — errors already surface as a
  visible Script Filter error item.

## [1.0.0] - 2026-05-02

### Added

- Password generation with customizable length and character set
- Four modes: basic, panc (with punctuation), split, panc split
- Alfred SDK: `response`, `cache`, `config`, `logger`, `router`, `safe_run`
- Shows 5 password suggestions per query in Alfred
- Clipboard history toggle via Config Builder
- Vendor packaging via `scripts/vendor.sh`
- Build pipeline via `scripts/build.sh`
- GitHub Actions CI (lint, test, build)
- GitHub Actions Release (tag → `.alfredworkflow` → GitHub Release)
- Full pytest test suite

[Unreleased]: https://github.com/y-marui/alfred-password-generator/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/y-marui/alfred-password-generator/releases/tag/v1.0.0
