# File Map

> File-level dependency map for the alfred-workflow-template.
> Add to this as you explore the codebase during development.

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
                 ├─ app.commands.search.handle(args)
                 │    └─ app.services.example_service.ExampleService.search(args)
                 │         └─ alfred.cache.Cache.get/set
                 │         └─ app.clients.api_client.ApiClient.search(args)
                 ├─ app.commands.open_cmd.handle(args)
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
| `commands/search.py` | `alfred.response`, `alfred.logger`, `app.services.example_service` | Default command |
| `commands/open_cmd.py` | `alfred.response`, `alfred.logger` | Named shortcut opener |
| `commands/config_cmd.py` | `alfred.response`, `alfred.config`, `alfred.logger` | Config viewer/reset |
| `commands/help_cmd.py` | `alfred.response` | Help display |
| `services/example_service.py` | `alfred.cache`, `alfred.logger`, `app.clients.api_client` | Replace with real service |
| `clients/api_client.py` | stdlib, `alfred.logger` | Replace with real API client |

### Tests (`tests/`)

| File | Tests |
|---|---|
| `test_alfred.py` | Alfred SDK modules (response, router, cache, config, safe_run) |
| `test_commands.py` | Command handlers (search, open, config, help) |
| `test_services.py` | Service layer (ExampleService) |
| `conftest.py` | pytest fixtures — sets Alfred env vars to tmp dirs |

## Key Files for Customization

When building a new workflow from this template, replace or update these files:

| File | What to change |
|---|---|
| `workflow/info.plist` | `bundleid`, keyword, UIDs, category, description |
| `src/app/clients/api_client.py` | Implement real API calls |
| `src/app/services/example_service.py` | Implement real business logic |
| `src/app/commands/search.py` | Adjust result formatting |
| `src/app/commands/open_cmd.py` | Update `_SHORTCUTS` dict |
| `src/app/core.py` | Register new commands |
| `pyproject.toml` | Workflow name and version |
| `workflow/icon.png` | Workflow icon |
