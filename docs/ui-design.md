# UI Design

Alfred Script Filter workflows present results as a list of items in the Alfred
launcher. This document defines the UI conventions for result items in this
workflow.

## Result Item Structure

Alfred result items are JSON objects with the following fields used in this workflow:

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string | yes | Primary text (large, always visible) |
| `subtitle` | string | no | Secondary text (small, below title) |
| `arg` | string | no | Value passed to Alfred's action on Enter |
| `uid` | string | no | Unique ID for Alfred's learned ordering |
| `valid` | bool | yes | If false, Enter does not trigger an action |
| `autocomplete` | string | no | Text inserted into Alfred's input on Tab |
| `icon` | object | no | Custom icon (`{ "path": "icon.png" }`) |

## Text Guidelines

### No Unicode Emoji in `title` / `subtitle`

- **Prohibited:** `🔍 Search`, `✅ Done`, `📄 Document`
- **Allowed:** ASCII symbols — `>`, `*`, `[x]`, `(!)`, `--`
- **Reason:** Emoji rendering is inconsistent across Alfred versions and macOS
  updates. ASCII symbols are universally stable.

### Capitalization

- `title`: Sentence case for action labels; preserve original casing for data values.
- `subtitle`: Sentence case. Use short phrases, not full sentences.

### Empty / Error States

- Empty query → show a placeholder item with `valid: false` to guide the user.
- No results → show an informative item (e.g., `No results for "foo"`) with `valid: false`.
- Error → `safe_run` automatically shows an error item; do not hide errors silently.

## Icon

- Workflow icon: `workflow/icon.png` (PNG, any size — Alfred scales it).
- Alfred controls light/dark mode; do not ship separate light/dark icons.
- Per-item icons are optional. If omitted, the workflow icon is used.

## Keyboard Shortcuts

These are standard Alfred behaviors — do not override them in the workflow:

| Key | Action |
|---|---|
| ↩ Enter | Run action with `arg` |
| ⌘↩ | Alfred action picker |
| ⇥ Tab | Insert `autocomplete` text into Alfred input |
| ⌘C | Copy `arg` to clipboard |
| ⌘L | Show `title` in Large Type |

## Layout Conventions by Command

### `search` results

```
title:    <result title>
subtitle: <result subtitle or URL>
arg:      <URL or identifier>
uid:      <unique result ID>
valid:    true
```

### `open` shortcut list

```
title:    <shortcut name>
subtitle: <URL>
arg:      <URL>
uid:      open-<name>
valid:    true
autocomplete: open <name>
```

### `config` items

```
title:    <key>: <value>      (for existing settings)
          Reset all settings  (action item)
subtitle: Current setting     (for existing settings)
          wf config reset ... (for action item)
valid:    false (settings display)
          true  (reset action)
```

### `help` items

```
title:    wf <command> <args>
subtitle: <command description>
valid:    false
autocomplete: <command trigger string>
```
