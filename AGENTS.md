# AI Assistant Guidelines

Please read CONTRIBUTING.md first.

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Common Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools
- `ci`: Changes to CI configuration files and scripts

### Examples

```
feat: add GitHub token management via keyring
fix: handle empty configuration file correctly
docs: add function documentation to controller package
chore(deps): update dependency aquaproj/aqua-registry to v4.403.0
```

## Code Validation

After making code changes, **always run** the following commands to validate and test:

### Validation

```sh
go vet ./...
```

### Testing

```sh
go test ./... -race -covermode=atomic
```

### Lint

```sh
golangci-lint run
```

### Auto fix lint errors

Note that only a few errors can be fixed by this command.
Many lint errors needs to be fixed manually.

```sh
golangci-lint fmt
```

## Test Framework Guidelines

- **DO NOT** use `testify` for writing tests
- **DO** use `google/go-cmp` for comparing expected and actual values
- Use standard Go testing package (`testing`) for all tests

## Dependencies

This project uses:

- [aqua](https://aquaproj.github.io/) for tool version management
- [cmdx](https://github.com/suzuki-shunsuke/cmdx) for task runner
- [golangci-lint](https://github.com/golangci/golangci-lint) for lint
- [goreleaser](https://goreleaser.com/) for releases

## Run docfresh locally

Show help:

```sh
go run ./cmd/docfresh help-all
```

Update documents by docfresh:

```sh
go run ./cmd/docfresh run <file path> [<file path> ...]
```

e.g.

```sh
go run ./cmd/docfresh run README.md examples/*.md
```

## Debugging

Enable debug logging:

```sh
export DOCFRESH_LOG_LEVEL=debug
```

## File Naming Conventions

- Internal test files: append `_internal_test.go` for internal testing

## Error Handling

Use slog and [slog-error](github.com/suzuki-shunsuke/slog-error).
