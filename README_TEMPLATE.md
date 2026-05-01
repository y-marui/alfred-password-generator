# {WorkflowName}

> **This is the English (reference) version.**
> For the Japanese canonical version, see [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/{user}/{repo}/actions/workflows/ci.yml/badge.svg)](https://github.com/{user}/{repo}/actions/workflows/ci.yml)
[![Charter Check](https://github.com/{user}/{repo}/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/{user}/{repo}/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/{user}?style=social)](https://github.com/sponsors/{user})
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/{bmc_username})

{One-line description: what · for whom · how it solves it}

## Requirements

- Alfred 5 (Powerpack required for Script Filter)
- Python 3.11+

## Setup

```bash
# Download the latest release:
# https://github.com/{user}/{repo}/releases/latest
```

Double-click `*.alfredworkflow` to install in Alfred.

Or build from source:

```bash
git clone https://github.com/{user}/{repo}
cd {repo}
make install
make build
# → dist/*.alfredworkflow
```

## Usage

Open Alfred and type `{keyword}` followed by a space.

### Search (default)

```
{keyword} <query>
{keyword} search <query>
```

### Help

```
{keyword} help
```

Show all available commands.

## License

MIT — see [LICENSE](LICENSE)

---

*This is the reference (English) version. The canonical Japanese version is [README-jp.md](README-jp.md). Update both files in the same commit.*
