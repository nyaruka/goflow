# CLAUDE.md

## Project Overview

Goflow is a Go implementation of the RapidPro/Textit flow engine specification. It orchestrates messaging flows involving contact management, message routing, webhook calls, and event handling.

- **Module:** `github.com/nyaruka/goflow`
- **Go version:** 1.25+
- **License:** Apache 2.0
- **Spec:** https://textit.com/mr/docs/

## Build & Test Commands

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./flows/actions/

# Update test fixtures (golden files)
go test github.com/nyaruka/goflow/test -update

# Regenerate ANTLR parsers
go generate ./...

# Build CLI tools
go install github.com/nyaruka/goflow/cmd/flowrunner
go install github.com/nyaruka/goflow/cmd/docgen
```

System dependencies for full builds: `pandoc`, `gettext`.

## Project Structure

| Directory | Purpose |
|-----------|---------|
| `flows/` | Core flow engine: actions, events, routers, triggers, resumes, modifiers |
| `excellent/` | Expression language engine (types, operators, 100+ built-in functions) |
| `contactql/` | Contact query language (ANTLR-based parser) |
| `assets/` | Asset management and static JSON loading |
| `envs/` | Environment config (timezone, language, locale, number formatting) |
| `services/` | External integrations (email, webhooks, classifiers) |
| `cmd/` | CLI tools: `flowrunner`, `flowmigrate`, `docgen`, `flowxgettext` |
| `test/` | Test infrastructure, helpers, and JSON fixtures |
| `antlr/` | ANTLR grammars (.g4) and generated parsers |
| `locale/` | i18n translations (.po files) |
| `utils/` | Shared utilities (validation, text, attachments) |

## Code Conventions

### Style
- Standard `gofmt` formatting
- Interface-based design with small, focused interfaces (1-5 methods)
- Embedded base structs for shared behavior (`baseAction`, `BaseEvent`)
- Distinct UUID types per domain: `ContactUUID`, `NodeUUID`, `ActionUUID`, etc.
- Type suffixes: `*Reference` for asset pointers, `*List` for collections

### Struct Tags
- JSON tags with validation: `json:"field" validate:"required,uuid"`
- Engine tags for special handling: `engine:"localized,evaluated"`

### Error Handling
- Explicit `(result, error)` returns
- Error wrapping with context: `fmt.Errorf("message: %w", err)`
- Validation via `utils.Validate()` using go-playground/validator

### Design Patterns
- **Type Registry:** Actions and events registered by type string in `init()`
- **Event Sourcing:** Flows produce events that are persisted and inspected
- **Visitor Pattern:** Node/action inspection via callback functions
- **Builder Pattern:** Environment and engine configuration

### Documentation Comments
- Doc comments on all exported identifiers
- Special docgen tags: `@action`, `@event`, `@operator`, `@function`
- Expression examples: `@("hello" & " world") -> hello world`

## Testing Conventions

- Test files use `package_name_test` (black-box testing)
- `testify/require` for setup assertions, `testify/assert` for validation
- `test.MockUniverse()` freezes time and UUIDs for deterministic tests
- JSON fixtures in `test/testdata/runner/` (`*.test.json` for expected output)
- Key helpers: `test.AssertXEqual()`, `test.AssertEqualJSON()`, `test.NewHTTPServer()`

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `nyaruka/gocommon` | Shared utils (dates, i18n, URNs, UUIDs) |
| `antlr4-go/antlr/v4` | ANTLR4 runtime for expression/query parsing |
| `go-playground/validator` | Struct field validation |
| `shopspring/decimal` | Precise decimal arithmetic |
| `buger/jsonparser` | High-performance JSON parsing |
| `stretchr/testify` | Test assertions |
