# Contributing to go-chatbot

Thank you for your interest in contributing to go-chatbot! We welcome contributions from the community and are pleased to have you join us.

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to contact@rumenx.com.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title**: Summarize the problem in the title
- **Detailed description**: Explain what you expected vs. what actually happened
- **Steps to reproduce**: List the steps to reproduce the behavior
- **Environment**: Include Go version, OS, and package version
- **Code samples**: Include relevant code snippets or configuration

### Suggesting Features

Feature requests are welcome! Please:

- Check existing issues for similar requests
- Explain the use case and why it would be beneficial
- Provide examples of how the feature would work
- Consider how it fits with the project's goals

### Development Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/go-chatbot.git
   cd go-chatbot
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

### Running Tests

Before submitting changes, ensure all tests pass:

```bash
# Run the test suite
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run static analysis
golangci-lint run

# Format code
go fmt ./...

# Run imports formatting
goimports -w .

# Run Go vet
go vet ./...
```

### Code Standards

We follow strict code quality standards:

- **Go Coding Standards**: Follow standard Go conventions (gofmt, goimports, go vet)
- **golangci-lint**: Static analysis for code quality and potential issues
- **High Test Coverage**: All new code should include tests (aim for 80%+ coverage)
- **Documentation**: Public functions, types, and packages must be documented with Go doc comments
- **Context Awareness**: Use context.Context for cancellation and timeouts where appropriate

### Adding New AI Providers/Models

To add a new AI provider/model:

1. Implement the `Model` interface in a new file in `models/`.
2. Add config options for the new provider in `config/`.
3. Update the model factory to support the new model.
4. Add tests for your new model in `models/`.
5. Document usage in the README and examples.

### Writing Tests

We use Go's standard testing package with testify for assertions. When adding new features:

1. **Write tests first** (Test-Driven Development)
2. **Cover edge cases** and error conditions
3. **Use descriptive test names** that explain what is being tested
4. **Mock external dependencies** (API calls, network requests, etc.)
5. **Use table-driven tests** for multiple scenarios

Example test structure:

```go
func TestChatMessageFilter_Handle(t *testing.T) {
    tests := []struct {
        name           string
        config         config.MessageFilteringConfig
        inputMessage   string
        expectedOutput string
        expectError    bool
    }{
        {
            name: "filters profanity from user messages",
            config: config.MessageFilteringConfig{
                Profanities: []string{"badword"},
            },
            inputMessage:   "This contains badword content",
            expectedOutput: "This contains *** content",
            expectError:    false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            filter := middleware.NewChatMessageFilter(tt.config)
            result, err := filter.Handle(context.Background(), tt.inputMessage)

            if tt.expectError {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.expectedOutput, result.Message)
        })
    }
}
```

### Submitting Changes

1. **Commit your changes** with clear, descriptive commit messages
2. **Push to your fork** and submit a pull request
3. **Ensure CI passes** - all tests, linting, and checks must pass
4. **Respond to feedback** from maintainers promptly

### Pull Request Guidelines

- **Clear title and description**: Explain what the PR does and why
- **Link related issues**: Reference any issues the PR addresses
- **Keep PRs focused**: One feature or fix per PR
- **Include tests**: All new functionality should include tests
- **Update documentation**: Update README or other docs if needed
- **Follow Go conventions**: Ensure code follows Go best practices

### Commit Message Format

Use clear, descriptive commit messages following conventional commit format:

```
type(scope): description

Examples:
feat(models): add support for Gemini Pro model
fix(middleware): handle empty message filtering rules
docs(readme): update configuration examples
test(models): add edge case tests for rate limiting
refactor(config): simplify configuration loading
style(lint): fix golangci-lint warnings
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `style`, `chore`

### Documentation

- Update relevant documentation for any changes
- Include code examples for new features
- Update the CHANGELOG.md for notable changes
- Ensure README.md stays current
- Write Go doc comments for all public APIs

## Getting Help

- **Issues**: For bugs and feature requests
- **Discussions**: For questions and general discussion
- **Email**: contact@rumenx.com for private inquiries

## Recognition

Contributors are recognized in:

- CHANGELOG.md for significant contributions
- GitHub contributors page
- Release notes for major features

## Development Philosophy

go-chatbot aims to be:

- **Framework-agnostic**: Works with any Go web framework
- **Secure by default**: Built-in security features
- **Easy to extend**: Clean architecture and interfaces
- **Well-tested**: High test coverage and quality
- **Performance-focused**: Efficient and scalable
- **Idiomatic Go**: Follows Go conventions and best practices
- **Context-aware**: Proper use of context.Context throughout

Thank you for contributing to go-chatbot! 🚀
