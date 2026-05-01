# Specification

> Functional specification, behavior definition, and data flow for the
> alfred-workflow-template. Adapt this document when building a workflow from
> this template.

## Overview

This workflow is an Alfred 5 Script Filter that accepts a keyword + query,
dispatches to a command handler, and returns a JSON result list to Alfred.

## Commands

### `search` (default)

**Trigger:** `wf <query>` or `wf search <query>`

**Behavior:**
1. If `query` is empty or whitespace → display "Type to search" placeholder item (valid: false).
2. Call `ExampleService.search(query)` — checks cache first, falls back to `ApiClient.search`.
3. If no results → display "No results" item (valid: false).
4. Otherwise → display one result item per record.

**Result item fields:**

| Field | Source | Notes |
|---|---|---|
| `title` | `result["title"]` | Primary display text |
| `subtitle` | `result["subtitle"]` | Secondary display text |
| `arg` | `result["url"]` or `result["title"]` | Passed to Alfred action on Enter |
| `uid` | `result["id"]` | Used by Alfred for learned ordering |

**Cache:** Results are cached by query string with TTL from `cache_ttl` env var (default 300 s).

---

### `open`

**Trigger:** `wf open <name>`

**Behavior:**
1. Filter `_SHORTCUTS` dict by whether `name` appears in the shortcut key.
2. If no matches → display "No shortcut" error item (valid: false).
3. Otherwise → display matching shortcuts; pressing Enter opens the URL in the default browser.

**Shortcuts (template defaults — replace with project-specific URLs):**

| Name | URL |
|---|---|
| `repo` | `https://github.com/yourname/your-workflow` |
| `docs` | `https://github.com/yourname/your-workflow/tree/main/docs` |
| `issues` | `https://github.com/yourname/your-workflow/issues` |

---

### `config`

**Trigger:** `wf config` / `wf config reset`

**Behavior:**
- `wf config` → list all keys in the persistent config store, plus a "Reset" action item.
- `wf config reset` → call `Config.reset()`, display confirmation item.

**Config storage:** `alfred_workflow_data` directory (set by Alfred at runtime).

---

### `help`

**Trigger:** `wf help`

**Behavior:** Display all registered commands with descriptions and autocomplete strings (valid: false for all items).

---

## Data Flow

```
Alfred input (keyword + query string)
  │
  ▼
workflow/scripts/entry.py         reads sys.argv[1]
  │
  ▼
alfred.safe_run.safe_run(main)    catches any uncaught exception → error item
  │
  ▼
app.core.run(query)               passes query to router
  │
  ▼
alfred.router.Router.dispatch     splits "search foo" → command="search", args="foo"
  │
  ▼
Command handler (e.g. search.handle("foo"))
  │
  ├─ [cache hit]  alfred.cache.Cache.get(key) → returns cached list
  │
  └─ [cache miss] app.clients.api_client.ApiClient.search(query)
                       → HTTP GET → parse response → return list[dict]
                  alfred.cache.Cache.set(key, result)
  │
  ▼
alfred.response.output(items)     prints JSON to stdout → Alfred renders result list
```

## Error Handling

- Any uncaught exception in `main()` is caught by `safe_run`, which emits a
  single error result item containing the exception message.
- Callers never crash Alfred silently (empty output = Alfred spinner forever).

## Configuration Variables

Managed via Alfred Configuration Builder (see `docs/configuration-builder.md`).

| Variable | Type | Default | Effect |
|---|---|---|---|
| `use_uv` | checkbox | `1` (on) | Use `uv run python` when uv is available |
| `log_level` | popupbutton | `WARNING` | Controls log verbosity |
| `cache_ttl` | textfield | `300` | Cache lifetime in seconds |
| `api_base_url` | textfield | `https://api.example.com/v1` | API endpoint |
| `api_timeout` | textfield | `5` | HTTP request timeout in seconds |

## Constraints

- Script Filter response time target: **< 100 ms**
- Cache miss (network call) must complete within `api_timeout` seconds.
- All output must go through `alfred.response.output()` — never `print()` directly.
- `entry.py` contains no business logic; it only sets `sys.path` and calls `safe_run(main)`.
