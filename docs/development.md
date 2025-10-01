# Development Guide

This guide covers development practices, setup, and contribution guidelines for TykCtl-Go.

## 🛠️ Development Setup

### Prerequisites

- **Go 1.21+** - Required for development
- **Git** - Version control
- **Make** - Build automation (optional)

### Environment Setup

```bash
# Clone the repository
git clone https://github.com/edsonmichaque/tykctl-go.git
cd tykctl-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build ./...
```

### Development Dependencies

```bash
# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest  # Hot reload for development
```

## 📝 Code Style Guidelines

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` and `goimports` for formatting
- Follow the project's existing patterns

### Documentation

- All public functions must have Go doc comments
- Use `// Package` comments for package documentation
- Include examples in documentation when helpful

### Error Handling

```go
// Good: Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// Good: Use sentinel errors for specific cases
var ErrConfigNotFound = errors.New("config file not found")

// Good: Check error types when needed
if errors.Is(err, ErrConfigNotFound) {
    // Handle specific error
}
```

### Testing

```go
// Good: Table-driven tests
func TestParseTimeout(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected time.Duration
        wantErr  bool
    }{
        {"valid seconds", "30s", 30 * time.Second, false},
        {"valid minutes", "5m", 5 * time.Minute, false},
        {"invalid format", "invalid", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := time.ParseDuration(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseTimeout() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.expected {
                t.Errorf("ParseTimeout() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./plugin/...

# Run tests with verbose output
go test -v ./...
```

### Test Structure

- Use table-driven tests for multiple scenarios
- Test both success and failure cases
- Use `testify/assert` for assertions when helpful
- Mock external dependencies

### Integration Tests

```bash
# Run integration tests (if any)
go test -tags=integration ./...

# Run tests with specific environment
TYKCTL_DEBUG=true go test ./...
```

## 🔍 Code Quality

### Linting

```bash
# Run linter
golangci-lint run

# Run specific linters
golangci-lint run --enable=gofmt,goimports,vet

# Fix auto-fixable issues
golangci-lint run --fix
```

### Common Linting Rules

- **gofmt**: Code formatting
- **goimports**: Import organization
- **govet**: Go vet checks
- **golint**: Style and convention checks
- **gosec**: Security checks
- **ineffassign**: Ineffective assignments

## 📦 Building

### Local Build

```bash
# Build all packages
go build ./...

# Build specific package
go build ./plugin

# Build with specific tags
go build -tags=debug ./...
```

### Cross-Platform Build

```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build ./...
GOOS=windows GOARCH=amd64 go build ./...
GOOS=darwin GOARCH=amd64 go build ./...
```

## 🚀 Release Process

### Version Management

- Use semantic versioning (v1.2.3)
- Update version in `go.mod`
- Tag releases in Git

### Release Checklist

- [ ] All tests pass
- [ ] Linting passes
- [ ] Documentation is updated
- [ ] Version is bumped
- [ ] CHANGELOG is updated
- [ ] Release notes are prepared

## 🤝 Contributing

### Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Pull Request Guidelines

- **Clear Description**: Explain what the PR does
- **Tests**: Include tests for new functionality
- **Documentation**: Update documentation as needed
- **Breaking Changes**: Clearly mark breaking changes
- **Small PRs**: Keep PRs focused and reasonably sized

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tooling changes

**Examples:**
```
feat(plugin): add timeout configuration support

fix(config): resolve environment variable parsing issue

docs(api): update client documentation
```

## 🐛 Debugging

### Debug Mode

```bash
# Enable debug logging
TYKCTL_DEBUG=true go run main.go

# Verbose output
TYKCTL_VERBOSE=true go run main.go
```

### Common Debug Issues

1. **Configuration Issues**: Check environment variables and config files
2. **Plugin Issues**: Verify plugin paths and permissions
3. **Network Issues**: Check API endpoints and connectivity
4. **Permission Issues**: Verify file/directory permissions

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Testing in Go](https://golang.org/doc/tutorial/add-a-test)
- [Go Modules](https://golang.org/doc/modules/)

## 🆘 Getting Help

- **Issues**: [GitHub Issues](https://github.com/edsonmichaque/tykctl-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/edsonmichaque/tykctl-go/discussions)
- **Documentation**: Browse the docs directory