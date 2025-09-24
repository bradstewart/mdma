# Repository Guidelines

## Project Structure & Module Organization
- `cmd/mdma` holds the TUI entry point; start here when wiring top-level commands or new flows.
- `internal/hello` currently houses sample domain logic; follow this pattern for internal packages that should not be exported.
- `docs` collects research notes (e.g., library evaluations) that inform future feature work; keep architectural research there.
- `prompts` contains strategic guidance; review updates there before making large changes.

## Build, Test, and Development Commands
- `make run` (or `go run cmd/mdma/main.go`) launches the app for manual smoke tests.
- `go test ./...` executes the Go test suite; run it locally before submitting changes, even when no new tests were added.
- `go fmt ./...` and `go vet ./...` keep formatting consistent and catch static analysis issues prior to review.
- Optional: `nix develop` (via `devenv.*`) drops you into the configured development environment if you use Nix.

## Coding Style & Naming Conventions
- Follow standard Go conventions: tabs for indentation, CamelCase for exported identifiers, and `snake_case` for file names when packages span multiple files.
- Format all Go code with `gofmt`; commit the formatted output, not manual tweaks.
- Keep packages focused; avoid cyclic dependencies by placing shared helpers under `internal/<feature>`.
- Use concise doc comments (`// Name ...`) for exported types and functions when behaviour is non-obvious.

## Testing Guidelines
- Place tests beside implementation files using the `_test.go` suffix and `package name_test` when black-box testing.
- Prefer table-driven tests for combinations of inputs, and assert user-visible strings stay stable.
- When introducing audio or UI logic, add integration-style tests that validate wiring via thin seams (e.g., interface mocks) to keep `go test ./...` fast.

## Commit & Pull Request Guidelines
- Base commit messages on recent history: a short (<60 char) sentence-case summary such as `Update README.md`, followed by optional descriptive body text.
- Group logically related changes per commit; avoid bundling refactors with feature work without clear justification.
- Pull requests should describe the user-facing impact, list testing performed (commands from above), and mention related prompts or docs when relevant.
- Include screenshots or terminal recordings when UI output changes to help reviewers understand the effect quickly.
