# Alfred Password Generator

> **This is the English (reference) version.**
> For the Japanese canonical version, see [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-password-generator/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-password-generator/actions/workflows/ci.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/y.marui)

| Field | Value |
|---|---|
| Target | Alfred 5 Script Filter workflow |
| License | MIT |
| Runtime | Go (build-time only, see `go.mod`), Alfred 5 |

Generate passwords with customizable length and character set.

## Usage

Open Alfred and type `passgen` followed by a space.

### Basic (default)

```
passgen [length] [pattern]
```

Default length is 18, default pattern is `A-Za-z0-9`.

### With punctuation (panc)

```
passgen panc [length] [pattern]
```

Default pattern becomes `A-Za-z0-9!-*` (includes `!@#^&*`).

### Split into groups (split)

```
passgen split [length] [by] [pattern]
```

Default length is 18, default group size is 6 (e.g. `xxxxxx-xxxxxx-xxxxxx`).
`length` must be a multiple of `by`.

### Split with punctuation (panc split)

```
passgen panc split [length] [by] [pattern]
```

Press Enter to copy the selected password to the clipboard.

## Pattern syntax

Characters can be listed directly (e.g. `ABCabc012!@#`) or as ranges (e.g. `A-Za-z0-9`).

Punctuation range: `!-*` expands to `!@#^&*`.

| Pattern | Expands to |
|---|---|
| `A-Z` | Uppercase letters |
| `a-z` | Lowercase letters |
| `0-9` | Digits |
| `!-*` | `!@#^&*` |
| `A-Za-z0-9` | Alphanumeric |
| `A-Za-z0-9!-*` | Alphanumeric + punctuation |

## Configuration

Access settings via Alfred Preferences (`⌘,`).

| Setting | Default | Description |
|---|---|---|
| Clipboard History | OFF | Save passwords to Alfred's clipboard history (not recommended for security) |

## Installation

```bash
make build-workflow   # Build workflow package
# → dist/*.alfredworkflow
```

Double-click `dist/*.alfredworkflow` to install in Alfred.

## Project Structure

```
alfred-password-generator/
├── cmd/
│   └── password-generator-alfred/  # The binary Alfred invokes
├── internal/
│   ├── passgen/         # Password generation logic (core)
│   ├── passgencmd/      # Command dispatch, argument parsing, help
│   └── scriptfilter/    # Alfred Script Filter JSON types
└── workflow/            # Alfred package (info.plist, icon.png)
```

## License

MIT — see [LICENSE](LICENSE)

---

*This is the reference (English) version. The canonical Japanese version is [README-jp.md](README-jp.md). Update both files in the same commit.*
