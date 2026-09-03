# Developing

This document covers the development workflow, conventions, and guidelines for contributors to this project.

## Development Setup

```bash
git clone https://github.com/y-marui/alfred-password-generator
cd alfred-password-generator
go build ./...
```

**Prerequisites:**
- macOS (required for Alfred)
- Go (see `go.mod` for the toolchain version)
- Alfred 5 with Powerpack
- `jq` (optional, for pretty-printed dev output): `brew install jq`
- `gh` CLI (required for releases): `brew install gh`

## Development Workflow

### Daily commands

```bash
go run ./cmd/password-generator-alfred "18 A-Z"   # simulate Alfred locally
go run ./cmd/password-generator-alfred "panc split"
go run ./cmd/password-generator-alfred "help"
go run ./cmd/password-generator-alfred ""
make test                 # run test suite
make lint                 # gofmt -l + go vet
make fmt                  # gofmt -w (auto-fix)
make build                # go build ./...
make build-workflow       # build dist/*.alfredworkflow
```

Pipe through `jq` for pretty-printed JSON:

```bash
go run ./cmd/password-generator-alfred "panc" | jq .
```

### Testing in Alfred

1. `make build-workflow` — generates `dist/*.alfredworkflow`
2. Double-click the `.alfredworkflow` file to install in Alfred
3. Open Alfred and type your keyword to verify behavior

During rapid iteration you can symlink `workflow/` to Alfred's workflow directory,
but `go run ./cmd/password-generator-alfred "query"` is usually faster for logic changes.

## Adding a New Command

1. Add a `handleFoo(args string) scriptfilter.Response` function to
   `internal/passgencmd/passgencmd.go`:

```go
func handleFoo(args string) scriptfilter.Response {
	return scriptfilter.Response{
		Items: []scriptfilter.Item{
			{Title: "My command", Subtitle: "Args: " + args, Arg: args, Valid: scriptfilter.BoolPtr(true)},
		},
	}
}
```

2. Register it in `Dispatch`'s switch statement.
3. Add tests in `internal/passgencmd/passgencmd_test.go`.
4. Update `README.md` Usage section and `workflow/info.plist`'s Script Filter `subtext`.

## Naming Conventions

| Scope | Convention | Example |
|---|---|---|
| Go files | `snake_case` avoided — one file per package concern | `passgen.go`, `passgencmd.go` |
| Go packages | short, lowercase, no underscores | `passgen`, `passgencmd`, `scriptfilter` |
| Exported functions / types | `PascalCase` | `Generate`, `Response`, `Item` |
| Unexported functions / variables | `camelCase` | `handleBasic`, `defaultLength` |
| Alfred command names | lowercase | `"panc"`, `"split"`, `"help"` |
| Alfred variable names | `lowercase_with_underscores` | `history` |
| Commit messages | Conventional Commits | `feat:`, `fix:`, `docs:`, `chore:` |
| Branch names | `feat/`, `fix/`, `docs/`, `chore/` | `feat/add-open-browser` |

## Code Style

- **Formatter:** `gofmt`. CI enforces this (`make lint`).
- **Linter:** `go vet`.
- **Comments:** Write *why*, not *what*. Do not comment self-evident code.
- **No debug prints:** Remove all stray `fmt.Print*` statements before committing;
  the only writer to stdout is `scriptfilter.Response.Write`.
- **No third-party dependencies** unless clearly justified — keep `go.mod` dependency-free.

## Commit Guidelines

- Commit per **feature unit**, after confirming it works.
- **No WIP commits** — do not commit code that does not run.
- **No `--no-verify`** — never skip pre-commit hooks.

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add clipboard copy action
fix: cache miss on special characters in query
chore: update Go toolchain to 1.28
docs: update README usage section
refactor: simplify passgencmd dispatch logic
```

## Pull Request Checklist

- [ ] `make lint` passes
- [ ] `make test` passes
- [ ] `make build-workflow` succeeds
- [ ] New commands have tests
- [ ] `README.md` updated if user-facing changes
- [ ] `CHANGELOG.md` entry added under `[Unreleased]`

## Code Review Guidelines

**Reviewers check for:**
- Architectural constraints respected (no business logic in `cmd/password-generator-alfred`)
- No hardcoded absolute paths (use `$HOME` / env vars)
- No debug prints in production code
- No Unicode emoji in Alfred result item `title` / `subtitle`
- Tests cover the new or changed behavior
- Alfred env variables managed via Config Builder, not `variables` key

**Security-sensitive changes** (auth, encryption, data access) require explicit
security review before merge.

**Self-review:** Individual contributors open a PR and self-review before merging
to `main`.
