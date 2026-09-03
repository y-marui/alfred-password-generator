# Architecture

## Overview

An Alfred Workflow (Go): `cmd/password-generator-alfred` is the single
universal (amd64+arm64) binary `workflow/info.plist` invokes. Its Script
Filter node passes the query following the `passgen` keyword as `$1`; the
binary parses it, generates one or more passwords via `internal/passgen`,
and prints Alfred Script Filter JSON via `internal/scriptfilter`. Selecting
a result copies its password to the clipboard through Alfred's own native
Conditional/Clipboard Output nodes — no script is involved in that step; see
[docs/specification.md](specification.md) for the full data flow.
`scripts/build-workflow.sh` packages the binary with `workflow/info.plist`
and `workflow/icon.png` into a `.alfredworkflow`.

This structure — a thin `cmd/` entry point over independently testable
`internal/` packages, no generic command-router abstraction, Script Filter
JSON via a small `scriptfilter` package — deliberately matches
[y-marui/alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text)
and [y-marui/alfred-markdown-ref](https://github.com/y-marui/alfred-markdown-ref),
this author's other Alfred Workflows already implemented in Go. This workflow
itself was originally a Python implementation
([`src/alfred`/`src/app`](https://github.com/y-marui/alfred-password-generator/tree/v1.0.0/src));
see `CHANGELOG.md`'s `[Unreleased]` entry for what changed and why in that
rewrite.

## Entry Points

- `cmd/password-generator-alfred` — a single command, no subcommands. The
  query it receives (e.g. `"panc split 18 6"`, `"24"`, `"help"`) determines
  behavior; see [docs/specification.md](specification.md#commands).

One Alfred trigger reaches it: the `passgen` keyword, wired in
`workflow/info.plist`.

## Directory Structure

| Directory | Role |
|---|---|
| `cmd/password-generator-alfred/` | The binary Alfred invokes; recovers panics into a Script Filter error item and writes the response |
| `internal/passgencmd/` | Query dispatch, argument parsing, overview/help rendering — the Alfred result rows |
| `internal/passgen/` | Pattern expansion and password generation, unit tested independently of Alfred |
| `internal/scriptfilter/` | Alfred Script Filter JSON response types |
| `workflow/` | `info.plist` (the Alfred object graph), `icon.png` |
| `scripts/build-workflow.sh` | Builds the universal binary and packages `workflow/` into `dist/*.alfredworkflow` |
| `scripts/extract-changelog.sh` | Extracts one version's notes from `CHANGELOG.md` for GitHub Releases |
| `docs/` | Specification, ADRs |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`) |

## Key Dependencies

None. `internal/passgen` and `internal/passgencmd` use only the Go standard
library (`crypto/rand`, `strconv`, `strings`).

## Alfred Configuration Builder (`userconfigurationconfig`)

Alfred 5 の Configuration Builder は `info.plist` の `userconfigurationconfig` キーで定義する。
利用可能な全型・各キーの詳細は [`docs/configuration-builder.md`](configuration-builder.md) を参照。

This workflow declares two variables: `history` (checkbox) — it controls
which of the two Clipboard Output nodes the Conditional node routes to
(`workflow/info.plist`), and is never read by the Go binary itself — and
`max_attempts` (textfield) — the character-class diversity retry budget
`internal/passgencmd.maxAttempts()` reads via `os.Getenv`
(`docs/specification.md#character-class-diversity`).

### Passing variables

Alfred はスクリプト実行時に各 `variable` を環境変数として渡す。
インストール直後は `prefs.plist` が存在しないため変数は未セットになる場合がある。
この Go バイナリはどの Config Builder 変数も読まないため該当しないが、将来スクリプト側で
変数を読む場合は常にデフォルト値を持たせること。

~~~go
// Go
value := os.Getenv("my_variable")
if value == "" {
	value = "fallback"
}
~~~

**注意:** `checkbox` 型の unchecked 値は `"0"` ではなく空文字 `""` になる。
`[ "$var" = "1" ]` で判定し、`"0"` との比較は避けること。

### Relationship between `variables` / `prefs.plist` / `default`

| 場所 | 役割 |
|---|---|
| `userconfigurationconfig[].config.default` | Configuration Builder UI の初期表示のみ。変数への書き込みは行わない。 |
| `prefs.plist`（同ディレクトリ） | ユーザーが Configuration Builder で保存した値。Alfred が自動生成・更新する。 |
| `info.plist` の `variables` | スクリプトに常に渡したい固定の環境変数。Configuration Builder で管理する変数はここに入れない。 |
