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
2. Generate `numSuggestions` (5) passwords using `internal/passgen.Generate(pattern, length)`.
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
2. Generate `numSuggestions` (5) passwords using `internal/passgen.GenerateSplit(pattern, length, by)`.
3. `length` must be a multiple of `by`; otherwise returns an error item.

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
cmd/password-generator-alfred/main.go   reads os.Args[1]
  │
  ▼
dispatch(query)                         recovers any panic → error item
  │
  ▼
internal/passgencmd.Dispatch(query)     splits "panc split 18 6" → command="panc", args="split 18 6"
  │
  ▼
Command handler (e.g. handlePanc("split 18 6"))
  │
  └─ internal/passgen.GenerateSplit(pattern, length, by)
  │
  ▼
scriptfilter.Response.Write(os.Stdout)  prints JSON to stdout → Alfred renders result list
  │
  ▼
User selects item → arg (password) passed to Conditional node
  │
  ├─ {var:history} = "1" → Clipboard (transient=false, saved to history)
  └─ {var:history} = ""  → Clipboard (transient=true, not saved)
```

## Error Handling

- Any panic in `dispatch()` is recovered in `cmd/password-generator-alfred/main.go`, which
  emits a single error result item containing the panic message.
- Invalid pattern or `length % by != 0` → error returned from `internal/passgen`, surfaced
  as a single `Error: <message>` result item.

## Configuration Variables

Managed via Alfred Configuration Builder (see `docs/configuration-builder.md`).

| Variable | Type | Default | Effect |
|---|---|---|---|
| `history` | checkbox | `""` (off) | Save password to Alfred clipboard history |

## Constraints

- Script Filter response time target: **< 100 ms** (compiled Go binary, no I/O)
- All output must go through `scriptfilter.Response.Write()` — never `fmt.Print*` directly.
- `cmd/password-generator-alfred` contains no business logic; it only reads `os.Args[1]`,
  recovers panics, and writes the response `internal/passgencmd.Dispatch` returns.
