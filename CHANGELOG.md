# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
