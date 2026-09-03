# File Map

> File-level dependency map for alfred-password-generator.
> Update this as the codebase evolves.

## Entry Points

| File | Role |
|---|---|
| `cmd/password-generator-alfred/main.go` | Alfred executes this binary — the sole entry point |

## Call Flow

```
cmd/password-generator-alfred/main.go
  └─ dispatch(query)                                [recovers panics into an error item]
       └─ internal/passgencmd.Dispatch(query)
            ├─ handleBasic(args)                     [default]
            │    └─ internal/passgen.Generate(pattern, length)
            ├─ handlePanc(args)
            │    └─ internal/passgen.Generate / GenerateSplit
            ├─ handleSplit(args)
            │    └─ internal/passgen.GenerateSplit(pattern, length, by)
            └─ handleHelp()
```

## Package Dependency Table

| Package | Imports from | Notes |
|---|---|---|
| `internal/scriptfilter` | stdlib only | Alfred Script Filter JSON types (`Item`, `Response`) and `Write` |
| `internal/passgen` | stdlib only (`crypto/rand`) | Core password generation logic — pattern expansion, `Generate`, `GenerateSplit` |
| `internal/passgencmd` | `internal/passgen`, `internal/scriptfilter` | Query dispatch, argument parsing, overview/help rendering |
| `cmd/password-generator-alfred` | `internal/passgencmd`, `internal/scriptfilter` | Reads `os.Args[1]`, recovers panics, writes JSON to stdout |

## Tests

| File | Tests |
|---|---|
| `internal/passgen/passgen_test.go` | `Generate`, `GenerateSplit`, pattern expansion (valid/invalid ranges) |
| `internal/passgencmd/passgencmd_test.go` | Command dispatch (passgen, panc, split, help), overview mode, error items |
