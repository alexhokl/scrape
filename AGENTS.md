# AGENTS.md

This file provides guidance for AI coding agents working in this repository.

## Project Overview

`scrape` is a Go CLI tool that retrieves links from websites and scrapes article
content to produce markdown output. It uses the [cobra](https://github.com/spf13/cobra)
CLI framework and [colly](https://github.com/gocolly/colly) for web scraping.

### Structure

```
main.go              # Entrypoint, calls cmd.Execute()
cmd/                 # Cobra CLI commands (root, article, links, title, filename)
scraper/             # Scraper interfaces, implementations, and utilities
  scraper.go         # LinkScraper and ArticleScraper interfaces
  factory.go         # Factory functions to create scrapers by source type
  guardian.go        # GuardianScraper implementation
  go_doc.go          # GoDocScraper implementation
  microsoft_learn.go # MicrosoftLearnScraper implementation
  new_york_times.go  # NewYorkTimesScraper implementation
  tailscale.go       # TailscaleScraper implementation
  tofugu.go          # TofuguScraper implementation
  util.go            # Shared utility functions
```

## Build / Lint / Test Commands

This project uses [Task](https://taskfile.dev/) as its task runner (see `Taskfile.yml`).

```bash
# Build (compiles to /dev/null to check for errors)
task build
# or: go build -o /dev/null

# Install the binary
task install
# or: go install

# Run all tests
task test
# or: go test ./...

# Run tests in a single package
go test ./scraper/
go test ./cmd/

# Run a single test by name
go test ./scraper/ -run TestRemoveExtraSpaces
go test ./cmd/ -run TestValidateArticleOptions_ValidGuardianSource

# Run tests with subtests by name
go test ./scraper/ -run TestCreateArticleScraper_CaseSensitive/Guardian

# Run tests with verbose output
go test -v ./...

# Test coverage
task coverage
# or: go test --cover ./...

# Generate HTML coverage report
task coverage-html

# Lint (requires golangci-lint)
task lint
# or: golangci-lint run

# Security check (requires gosec)
task sec
# or: gosec ./...
```

## Code Style Guidelines

### General Formatting

- Use `gofmt` / `goimports` standard formatting (tabs for indentation).
- No custom formatting rules beyond standard Go conventions.

### Package & File Organization

- Package `cmd` contains all CLI commands. Each subcommand lives in its own file
  (e.g., `article.go`, `links.go`, `title.go`, `filename.go`).
- Package `scraper` contains interfaces, implementations, factory, and utilities.
- Each scraper implementation gets its own file named after the source
  (e.g., `guardian.go`, `go_doc.go`, `microsoft_learn.go`).
- Test files use the `_test.go` suffix and are in the same package (not `_test`
  external test packages). This allows testing unexported functions.

### Imports

- Group imports into standard library, then project-internal, then third-party,
  separated by blank lines. Example from `cmd/article.go`:
  ```go
  import (
      "fmt"
      "os"

      "github.com/alexhokl/scrape/scraper"
      "github.com/spf13/cobra"
  )
  ```
- Use the module path `github.com/alexhokl/scrape` for internal imports.

### Naming Conventions

- **Interfaces**: PascalCase with descriptive names (`LinkScraper`, `ArticleScraper`).
- **Structs**: PascalCase, named after the source (`GuardianScraper`, `GoDocScraper`).
- **Exported functions**: PascalCase (`CreateLinkScraper`, `CreateArticleScraper`).
- **Unexported functions**: camelCase (`removeExtraSpaces`, `trimSpacesAndLineBreaks`,
  `generateFileNameFromTitle`, `parseGoDocParagraph`).
- **Options structs**: camelCase unexported structs with `Options` suffix
  (`articleOptions`, `linksOptions`, `titleOptions`, `filenameOptions`).
- **Package-level vars for options**: camelCase with `Opts` suffix
  (`articleOpts`, `linksOpts`, `titleOpts`, `filenameOpts`).
- **Cobra commands**: camelCase with `Cmd` suffix (`articleCmd`, `linksCmd`).
- **Source type strings**: all lowercase (`"guardian"`, `"microsoft"`, `"go"`,
  `"tofugu"`, `"newyorktimes"`, `"tailscale"`). Case-sensitive by design.

### Error Handling

- Return errors using `fmt.Errorf` with descriptive messages.
- Wrap errors with `%w` verb when adding context:
  `fmt.Errorf("error creating scraper: %w", err)`.
- Use bare `fmt.Errorf` (no wrapping) for validation errors:
  `fmt.Errorf("invalid source: %s", opts.source)`.
- Factory functions return `nil, fmt.Errorf(...)` for unsupported sources.
- Scraper methods return zero value + error on failure (e.g., `"", err`).
- Cobra commands use `RunE` / `PersistentPreRunE` to return errors rather than
  calling `os.Exit` or `log.Fatal`.

### Interfaces & Patterns

- Interfaces are defined in `scraper/scraper.go` and kept minimal.
- `LinkScraper` has one method: `ScrapeLinks(url string) (map[string]string, error)`.
- `ArticleScraper` has three methods: `ScrapeArticle`, `ScrapeTitle`, `ScrapeFilename`.
- Factory pattern in `scraper/factory.go` maps source type strings to concrete
  implementations via switch statements.
- Each scraper struct is a pointer receiver with no fields (stateless).
- `ScrapeFilename` implementations typically delegate to `ScrapeTitle` then call
  `generateFileNameFromTitle`.

### CLI Command Pattern

Each command file follows this pattern:
1. Define an unexported options struct.
2. Declare a package-level var for the options.
3. Define the cobra.Command with `Use`, `Short`, `PersistentPreRunE`, and `RunE`.
4. In `init()`, add the command to `rootCmd` and register flags.
5. Implement a validation function and a run function, both matching
   `func(*cobra.Command, []string) error`.
6. Validation uses switch statements and returns errors for invalid inputs.
7. Output goes to `os.Stdout` via `fmt.Fprintln` or `fmt.Fprintf`.

### Testing Conventions

- Use the standard `testing` package only; no third-party test frameworks.
- Test function names: `TestFunctionName_Scenario` (e.g.,
  `TestCreateArticleScraper_UnsupportedSource`).
- Use table-driven tests with `[]struct{ name, input, expected }` slices for
  functions with many input variations (see `util_test.go`).
- Use `t.Run(name, func(t *testing.T) {...})` for subtests.
- Use `t.Fatalf` for fatal setup errors, `t.Errorf` for assertion failures.
- Mock implementations are defined in test files with compile-time interface
  checks: `var _ LinkScraper = (*mockLinkScraper)(nil)`.
- When testing commands that use package-level option vars, save and restore
  the original value with `defer`.

### Adding a New Scraper

1. Create `scraper/<source_name>.go` with a struct implementing `ArticleScraper`
   and/or `LinkScraper`.
2. Register it in `scraper/factory.go` by adding a case to the appropriate
   `Create*Scraper` switch.
3. Add the source string to validation switch statements in the relevant
   `cmd/*.go` files.
4. Create `scraper/<source_name>_test.go` with tests.
5. Source type strings must be all lowercase, single-word or concatenated
   (e.g., `"newyorktimes"`, not `"new-york-times"`).
