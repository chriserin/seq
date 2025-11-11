# sq Project Guidelines

## Build & Test Commands

- **IMPORTANT**: Run every go command with the `-tags lua54` option
- Run all tests: `go test -tags lua54 ./...`
- Run specific package tests: `go test -tags lua54 ./internal/packageName`
- Run single test: `go test -tags lua54 -run TestName ./package/path`
- Run specific test function: `go test -tags lua54 -run TestName/SubTest ./package/path`
- Run with verbose output: `go test -tags lua54 -v ./...`
- Build: `go build -tags lua54 -o sq`
- Format code: `go fmt ./...`

## Code Style

- Go version: 1.23.3
- Use `stretchr/testify/assert` for test assertions
- Organize imports: standard library first, then external packages, then internal packages
- Follow Go naming conventions: CamelCase for exported, camelCase for unexported
- Error handling: check errors immediately, don't use panics in production code
- Internal packages: use `internal/` directory for code not meant for external import
- UI components: use charmbracelet/bubbletea framework patterns
- Tests: use table-driven tests with descriptive test cases
- Package naming: use lowercase package names (e.g., 'overlaykey')

## Project Structure

- Main application: root directory
- Internal libraries: `/internal` directory
- Configuration: Lua files in `/config`
