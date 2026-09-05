# UI Design

Alfred Script Filter workflows present results as a list of items in the Alfred
launcher. This document defines the UI conventions for result items in this
workflow.

## Result Item Structure

Alfred result items are JSON objects with the following fields used in this workflow:

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string | yes | Primary text (large, always visible) — the generated password itself, or a command's usage string in `passgen help` |
| `subtitle` | string | no | Secondary text (small, below title) — the length/pattern description, or a command's description in `passgen help` |
| `arg` | string | no | The generated password; copied to the clipboard on Enter |
| `uid` | string | no | Unique ID for Alfred's learned ordering |
| `valid` | bool | yes | If false, Enter does not trigger an action (`passgen help` rows) |
| `autocomplete` | string | no | Text inserted into Alfred's input on Tab — used by `passgen help` rows |

## Text Guidelines

### No Unicode Emoji in `title` / `subtitle`

- **Prohibited:** `🔍 Search`, `✅ Done`, `📄 Document`
- **Allowed:** ASCII symbols — `>`, `*`, `[x]`, `(!)`, `--`
- **Reason:** Emoji rendering is inconsistent across Alfred versions and macOS
  updates. ASCII symbols are universally stable.

### Empty / Error States

- Bare `passgen` (or a length with no pattern) → an overview: one generated
  password per registered variant (`pin`, `code`, `basic`, `panc`, `split`,
  `panc split`), each `valid: true`, not a placeholder.
- No valid pattern for the requested length → a single error item,
  `title: "Error: No valid patterns for length <n>"`, `valid: false`.
- Generation failure (e.g. an exhausted pattern) → a single error item,
  `title: "Error: <message>"`, `subtitle: "Press ⌘C to copy the error"`,
  `valid: false`.

## Icon

- Workflow icon: `workflow/icon.png` (PNG, any size — Alfred scales it).
- Alfred controls light/dark mode; do not ship separate light/dark icons.
- No per-item icons are used in this workflow.

## Keyboard Shortcuts

These are standard Alfred behaviors — do not override them in the workflow:

| Key | Action |
|---|---|
| ↩ Enter | Copy the generated password (`arg`) to clipboard |
| ⇥ Tab | Insert `autocomplete` text into Alfred's input (`passgen help` rows) |
| ⌘C | Copy `arg` to clipboard (also how to copy an error message) |
| ⌘L | Show `title` in Large Type |

## Layout Conventions by Command

### `passgen` overview (bare keyword, or a length with no pattern)

One item per registered variant:

```
title:    <generated password>
subtitle: <variant name> · <length> chars[ in groups of <n>]
arg:      <generated password>
uid:      passgen-<index>
valid:    true
```

### `passgen [length] [pattern]`, `passgen panc [...]`, `passgen pin [...]`, `passgen code [...]`

One or more suggestions (5 for `panc`/`pin`/`code`, 1 for a fully-specified
`passgen [length] [pattern]`):

```
title:    <generated password>
subtitle: <length> chars, pattern: <pattern>
arg:      <generated password>
uid:      passgen-<index>
valid:    true
```

### `passgen split [...]` / `passgen panc split [...]`

```
title:    <generated password, grouped>
subtitle: <length> chars in groups of <by>, pattern: <pattern>
arg:      <generated password>
uid:      passgen-<index>
valid:    true
```

### `passgen help`

```
title:    passgen <command usage>
subtitle: <command description>
valid:    false
autocomplete: <command trigger string>
```

### Error rows

```
title:    "Error: <message>"
subtitle: "Press ⌘C to copy the error"
valid:    false
```
