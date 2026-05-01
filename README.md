# Alfred Workflow Template

> **This is the English (reference) version.**
> For the Japanese canonical version, see [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-workflow-template/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-workflow-template/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-workflow-template/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-workflow-template/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/y.marui)

| Field | Value |
|---|---|
| Target | Alfred 5 Script Filter workflow |
| Team size | Individual to small team (1–3 people) |
| Language | English (OSS) |
| License | MIT |
| Runtime | Python 3.9+, Alfred 5 |
| AI tools | Claude Code / GitHub Copilot / Gemini CLI |

> Production-ready template for building Alfred 5 Script Filter workflows.
> Start shipping in 10 minutes.

## Features

- ✅ **Layered architecture** — Alfred boundary isolated from business logic
- ✅ **Lightweight Alfred SDK** — response builder, router, cache, config, logger
- ✅ **Command-based UX** — `wf search`, `wf open`, `wf config`, `wf help`
- ✅ **Full test suite** — pytest, no Alfred required to run tests
- ✅ **CI/CD** — lint, test, build, and release via GitHub Actions
- ✅ **Vendor packaging** — third-party deps bundled in `vendor/`
- ✅ **AI-ready** — `AI_CONTEXT.md` + `CLAUDE.md` for AI assistant context

## Requirements

- Alfred 5 (Powerpack required for Script Filter)
- Python 3.9+
- [pre-commit](https://pre-commit.com/) (for security hooks)

## Quick Start

### Using this template

1. Click **"Use this template"** → **"Create a new repository"** on GitHub
2. Clone your new repository and open it in your AI tool
3. Tell the AI: "Run the initial setup from AI_CONTEXT.md" — it will:
   - Apply GitHub repository settings
   - Rename `README_TEMPLATE.md` → `README.md` and `README_TEMPLATE-jp.md` → `README-jp.md`
   - Replace `{user}`, `{repo}`, `{keyword}` placeholders
   - Randomize the dev-charter cron schedule
4. Customize the workflow (see `DEVELOPING.md`)

### Development (this template)

```bash
git clone https://github.com/yourname/alfred-workflow-template
cd alfred-workflow-template

# Install dev dependencies
make install

# Simulate Alfred locally
make run Q="search foo"
make run Q="help"

# Run tests
make test

# Build workflow package
make build
# → dist/workflow-template-0.1.0.alfredworkflow
```

Double-click `dist/*.alfredworkflow` to install in Alfred.

## Usage

Open Alfred and type `wf` followed by a space.

### Search (default)

```
wf <query>
wf search <query>
```

Type any query to search. Press Enter to open the result.

| Key | Action |
|---|---|
| ↩ Enter | Open result |
| ⌘C | Copy result URL |

### Open

```
wf open <name>
```

Open a named shortcut. Available shortcuts: `repo`, `docs`, `issues`

### Config

```
wf config
wf config reset
```

View current settings or reset all configuration.

### Help

```
wf help
```

Show all available commands.

### Tips

- The workflow remembers your most-used results (Alfred learns from usage).
- Results are cached for 5 minutes to minimize API calls.
- Use `⌘,` in Alfred to access Workflow Preferences.

### Troubleshooting

**No results appear**
- Check Alfred's debugger: open Alfred → `⌘D`
- Check logs: `~/Library/Logs/Alfred/Workflow/<bundle-id>.log`

**Results are stale**
- The cache TTL is 5 minutes. Wait for it to expire, or clear manually: `wf config reset`

## Project Structure

```
alfred-workflow-template/
├── src/
│   ├── alfred/         # Alfred SDK (response, router, cache, config, logger, safe_run)
│   └── app/            # Application layer (commands, services, clients)
├── workflow/           # Alfred package (info.plist, scripts/entry.py, vendor/)
├── tests/              # pytest test suite
├── scripts/            # build.sh, dev.sh, release.sh, vendor.sh
└── docs/               # Architecture and reference documentation
```

## Documentation

| Document | Description |
|---|---|
| [docs/architecture.md](docs/architecture.md) | Module and layer design |
| [docs/file-map.md](docs/file-map.md) | File-level dependency map |
| [docs/specification.md](docs/specification.md) | Feature specification and data flow |
| [docs/ui-design.md](docs/ui-design.md) | Alfred result item UI conventions |
| [docs/configuration-builder.md](docs/configuration-builder.md) | Alfred Configuration Builder reference |

## AI-Assisted Development

This template is configured for AI-assisted development.

| Tool | Role |
|---|---|
| Claude Code | Architecture, large-scale changes, refactoring |
| GitHub Copilot | Bug fixes, small implementation, unit tests |
| Gemini CLI | Documentation management |

See [`AI_CONTEXT.md`](AI_CONTEXT.md) and [`CLAUDE.md`](CLAUDE.md) for session context.

## Customizing This Template

After running the initial setup (see Quick Start above), customize the workflow:

1. Edit `workflow/info.plist`:
   - Replace `bundleid` with your bundle ID (`com.yourname.workflowname`)
   - Replace the `keyword` (`wf`) with your trigger keyword
   - Run `uuidgen` and replace the placeholder UIDs
2. Replace `src/app/clients/api_client.py` with your API
3. Update the workflow name in `pyproject.toml`
4. Update shortcuts in `src/app/commands/open_cmd.py`
5. Add your `workflow/icon.png`

## Release

```bash
# 1. Bump version in pyproject.toml
# 2. Tag and push
git tag v1.2.3
git push --tags
# GitHub Actions builds .alfredworkflow and creates a GitHub Release
```

## License

MIT — see [LICENSE](LICENSE)

---

*This is the reference (English) version. The canonical Japanese version is [README-jp.md](README-jp.md). Update both files in the same commit.*
