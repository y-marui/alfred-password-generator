# Specification

> Functional specification, behavior definition, and data flow for alfred-password-generator.

## Overview

This workflow is an Alfred 5 Script Filter that accepts a keyword + query, dispatches to a
command handler, and returns a JSON result list to Alfred. Selecting a result copies the
generated password to the clipboard.

## Commands

### `passgen` (default)

**Trigger:** `passgen [length] [pattern]`

**Behavior:**
1. Parse `length` (int, default 18) and `pattern` (default `A-Za-z0-9`) from the query.
2. Generate `_NUM_SUGGESTIONS` (5) passwords using `passgen_service.generate(pattern, length)`.
3. Return each password as a valid result item (arg = password, skipknowledge = true).

**Result item fields:**

| Field | Value |
|---|---|
| `title` | Generated password |
| `subtitle` | `{length} chars, pattern: {pattern}` |
| `arg` | Password (copied to clipboard on Enter) |

---

### `panc`

**Trigger:** `passgen panc [length] [pattern]`

**Behavior:** Same as basic, but default pattern is `A-Za-z0-9!-*` (includes punctuation `!@#^&*`).

If the first token after `panc` is `split`, delegates to split mode with punctuation default pattern.

---

### `split`

**Trigger:** `passgen split [length] [by] [pattern]`

**Behavior:**
1. Parse `length` (default 18), `by` (default 6), `pattern` (default `A-Za-z0-9`).
2. Generate `_NUM_SUGGESTIONS` (5) passwords using `passgen_service.generate_split(pattern, length, by)`.
3. `length` must be a multiple of `by`; otherwise returns an error item.

---

### `config`

**Trigger:** `passgen config` / `passgen config reset`

**Behavior:**
- `passgen config` → list all keys in the persistent config store, plus a "Reset" action item.
- `passgen config reset` → call `Config.reset()`, display confirmation item.

**Config storage:** `alfred_workflow_data` directory (set by Alfred at runtime).

---

### `help`

**Trigger:** `passgen help`

**Behavior:** Display all registered commands with descriptions and autocomplete strings (valid: false).

---

## Pattern Syntax

Characters can be listed directly (e.g. `ABCabc012`) or as ranges.

| Range | Expands to |
|---|---|
| `A-Z` | Uppercase letters A through Z |
| `a-z` | Lowercase letters a through z |
| `0-9` | Digits 0 through 9 |
| `!-*` | Punctuation `!@#^&*` |

Ranges within different character classes cannot be mixed (e.g. `A-z` is invalid).

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
alfred.router.Router.dispatch     splits "panc split 18 6" → command="panc", args="split 18 6"
  │
  ▼
Command handler (e.g. passgen_cmd.handle_panc("split 18 6"))
  │
  └─ app.services.passgen_service.generate_split(pattern, length, by)
  │
  ▼
alfred.response.output(items)     prints JSON to stdout → Alfred renders result list
  │
  ▼
User selects item → arg (password) passed to Conditional node
  │
  ├─ {var:history} = "1" → Clipboard (transient=false, saved to history)
  └─ {var:history} = ""  → Clipboard (transient=true, not saved)
```

## Error Handling

- Any uncaught exception in `main()` is caught by `safe_run`, which emits a single error
  result item containing the exception message.
- Invalid pattern or `length % by != 0` → `ValueError` caught in `_output()` → error item shown.

## Configuration Variables

Managed via Alfred Configuration Builder (see `docs/configuration-builder.md`).

| Variable | Type | Default | Effect |
|---|---|---|---|
| `use_uv` | checkbox | `1` (on) | Use `uv run python` when uv is available |
| `history` | checkbox | `""` (off) | Save password to Alfred clipboard history |
| `log_level` | popupbutton | `WARNING` | Controls log verbosity |

## Constraints

- Script Filter response time target: **< 100 ms** (pure Python, no I/O)
- All output must go through `alfred.response.output()` — never `print()` directly.
- `entry.py` contains no business logic; it only sets `sys.path` and calls `safe_run(main)`.
