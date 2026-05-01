# File Map

> File-level dependency map for alfred-password-generator.
> Update this as the codebase evolves.

## Entry Points

| File | Role |
|---|---|
| `workflow/scripts/entry.py` | Alfred executes this file — the sole entry point |
| `src/app/core.py` | Wires Router to command handlers |

## Call Flow

```
workflow/scripts/entry.py
  └─ alfred.safe_run.safe_run(main)
       └─ app.core.run(query)
            └─ alfred.router.Router.dispatch(query)
                 ├─ app.commands.passgen_cmd.handle_basic(args)   [default]
                 │    └─ app.services.passgen_service.generate(pattern, length)
                 ├─ app.commands.passgen_cmd.handle_panc(args)
                 │    └─ passgen_service.generate / generate_split
                 ├─ app.commands.passgen_cmd.handle_split(args)
                 │    └─ app.services.passgen_service.generate_split(pattern, length, by)
                 ├─ app.commands.config_cmd.handle(args)
                 │    └─ alfred.config.Config.all/reset
                 └─ app.commands.help_cmd.handle(args)
```

## Module Dependency Table

### Alfred SDK (`src/alfred/`)

| Module | Imports from | Notes |
|---|---|---|
| `response.py` | stdlib only | Emits Script Filter JSON to stdout |
| `router.py` | stdlib only | Parses query string, dispatches to handler |
| `safe_run.py` | `alfred.response` | Wraps `main()` to catch uncaught exceptions |
| `cache.py` | stdlib only | TTL disk cache; reads `alfred_workflow_cache` env var |
| `config.py` | stdlib only | Persistent JSON store; reads `alfred_workflow_data` env var |
| `logger.py` | stdlib only | File logger to `~/Library/Logs/Alfred/Workflow/` |

### Application Layer (`src/app/`)

| Module | Imports from | Notes |
|---|---|---|
| `core.py` | `alfred.router`, `app.commands.*` | Dependency injection point |
| `commands/passgen_cmd.py` | `alfred.response`, `alfred.logger`, `app.services.passgen_service` | Default + panc + split handlers |
| `commands/config_cmd.py` | `alfred.response`, `alfred.config`, `alfred.logger` | Config viewer/reset |
| `commands/help_cmd.py` | `alfred.response` | Help display |
| `services/passgen_service.py` | stdlib only | Core password generation logic |

### Tests (`tests/`)

| File | Tests |
|---|---|
| `test_alfred.py` | Alfred SDK modules (response, router, cache, config, safe_run) |
| `test_commands.py` | Command handlers (passgen, config, help) |
| `test_services.py` | `passgen_service` (generate, generate_split, pattern expansion) |
| `conftest.py` | pytest fixtures — sets Alfred env vars to tmp dirs |
